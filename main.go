package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"

	"github.com/bwmarrin/discordgo"
	"github.com/wcharczuk/go-chart"
)

const (
	region string = "ue-east-1"
)

type Credential struct {
	BotToken  string
	ChannelID string
}

type Record struct {
	Name  string
	Cost  float64
	Ratio float64
}

type Payload struct {
	Content string
	File    *bytes.Buffer
}

func getPeriod() (string, string) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	today := time.Now().UTC().In(jst)
	start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, jst).Format("2006-01-02")
	end := today.Format("2006-01-02")

	return start, end
}

func getCost(ctx context.Context) (*costexplorer.GetCostAndUsageOutput, error) {
	start, end := getPeriod()
	period := types.DateInterval{
		Start: aws.String(start),
		End:   aws.String(end),
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	client := costexplorer.NewFromConfig(cfg)

	output, err := client.GetCostAndUsage(ctx, &costexplorer.GetCostAndUsageInput{
		TimePeriod:  &period,
		Granularity: "MONTHLY",
		GroupBy: []types.GroupDefinition{
			{
				Key:  aws.String("SERVICE"),
				Type: "DIMENSION",
			},
		},
		Metrics: []string{"BlendedCost"},
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return output, nil
}

func createList(cost *costexplorer.GetCostAndUsageOutput) []Record {
	list := make([]Record, 0, 0)
	totalCost := 0.0

	for _, g := range cost.ResultsByTime[0].Groups {
		cost, _ := strconv.ParseFloat(*g.Metrics["BlendedCost"].Amount, 64)
		totalCost = totalCost + cost
	}

	for _, g := range cost.ResultsByTime[0].Groups {
		name := g.Keys[0]
		cost, _ := strconv.ParseFloat(*g.Metrics["BlendedCost"].Amount, 64)
		list = append(list, Record{
			Name:  name,
			Cost:  cost,
			Ratio: cost / totalCost * 100,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Cost > list[j].Cost
	})

	return list
}

func drawChart(list []Record) (*bytes.Buffer, error) {
	var values []chart.Value
	others := chart.Value{
		Label: "Others",
		Value: 0.0,
	}

	for _, i := range list {
		if i.Ratio < 1.0 {
			others.Value += i.Cost
			continue
		}

		values = append(values, chart.Value{
			Value: i.Cost,
			Label: i.Name,
		})
	}

	values = append(values, others)

	c := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
	}

	buffer := bytes.NewBuffer([]byte{})
	err := c.Render(chart.PNG, buffer)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return buffer, nil
}

func createPayload(list []Record, chart *bytes.Buffer) (Payload, error) {
	start, end := getPeriod()

	totalCost := 0.0
	for _, r := range list {
		totalCost += r.Cost
	}

	content := fmt.Sprintf("__Daily Report__\n\nPeriod: %v - %v\nTotal: $%.2f\n\n", start, end, totalCost)
	for _, r := range list {
		content += fmt.Sprintf("- %v: $%.2f (%.1f%%)\n", r.Name, r.Cost, r.Ratio)
	}

	payload := Payload{
		Content: content,
		File:    chart,
	}

	return payload, nil
}

func postCost(payload Payload, credential Credential) (*discordgo.Message, error) {
	client, err := discordgo.New("Bot " + credential.BotToken)
	if err != nil {
		return nil, fmt.Errorf("cannot create Discord session, %v", err)
	}

	webhook, err := client.WebhookCreate(credential.ChannelID, "aws", "aws")
	if err != nil {
		return nil, fmt.Errorf("cannot create Webhook, %v", err)
	}

	files := []*discordgo.File{
		{
			Name:        "aws.png",
			ContentType: "image/png",
			Reader:      payload.File,
		},
	}

	params := &discordgo.WebhookParams{
		Username: webhook.Name,
		Content:  payload.Content,
		Files:    files,
	}

	message, err := client.WebhookExecute(webhook.ID, webhook.Token, true, params)
	if err != nil {
		return nil, fmt.Errorf("error execute Webhook, %v", err)
	}

	err = client.WebhookDelete(webhook.ID)
	if err != nil {
		return nil, fmt.Errorf("error delete Webhook, %v", err)
	}

	return message, nil
}

func handler(ctx context.Context) error {
	// TODO: 環境変数を読み込む
	credential := Credential{
		BotToken:  os.Getenv("BOT_TOKEN"),
		ChannelID: os.Getenv("CHANNEL_ID"),
	}

	// TODO: 料金を計算する
	cost, err := getCost(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	list := createList(cost)

	// TODO: 料金をもとに円グラフを作成する
	chart, err := drawChart(list)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// TODO: Discord に料金とその円グラフを投稿する
	payload, err := createPayload(list, chart)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	message, err := postCost(payload, credential)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(message)
	return nil
}

func main() {
	lambda.Start(handler)
}

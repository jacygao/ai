package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable.")
		return
	}
	client := openai.NewClient(option.WithAPIKey(apiKey))

	rankeeHost := os.Getenv("RANKEE_HOST")
	if rankeeHost == "" {
		fmt.Println("Please set the RANKEE_HOST environment variable.")
		return
	}
	tools := &Tools{
		RankeeClient: NewRankeeClient(rankeeHost),
	}

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(&client, tools, getUserMessage)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

type Tools struct {
	RankeeClient *RankeeClient
}

func NewAgent(client *openai.Client, tools *Tools, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

type Agent struct {
	client         *openai.Client
	getUserMessage func() (string, bool)
	tools          *Tools
}

func (a *Agent) Run(ctx context.Context) error {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{},
		Seed:     openai.Int(0),
		Model:    openai.ChatModelGPT4o,
		Tools: []openai.ChatCompletionToolParam{
			{
				Function: openai.FunctionDefinitionParam{
					Name:        "get_weather",
					Description: openai.String("Get weather at the given location"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"location"},
					},
				},
			},
			{
				Function: openai.FunctionDefinitionParam{
					Name:        "run_rankee_evaluation",
					Description: openai.String("Run a Rankee evaluation with given parameters"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"AppId": map[string]string{
								"type": "string",
							},
							"Index": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"AppId", "Index"},
					},
				},
			},
		},
	}

	fmt.Println("Chat with GPT (use 'ctrl-c' to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		params.Messages = append(params.Messages, openai.UserMessage(userInput))
		completion, err := a.client.Chat.Completions.New(ctx, params)
		if err != nil {
			return fmt.Errorf("error generating completion: %w", err)
		}

		toolCalls := completion.Choices[0].Message.ToolCalls

		// Return early if there are no tool calls
		if len(toolCalls) > 0 {
			params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
			for _, toolCall := range toolCalls {
				if toolCall.Function.Name == "get_weather" {
					// Extract the location from the function call arguments
					var args map[string]interface{}
					err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
					if err != nil {
						panic(err)
					}
					location := args["location"].(string)

					// Simulate getting weather data
					weatherData := getWeather(location)

					// Print the weather data
					fmt.Printf("Getting the weather for %s: %s\n", location, weatherData)

					params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
				}
				if toolCall.Function.Name == "run_rankee_evaluation" {
					// Extract the location from the function call arguments
					var args map[string]interface{}
					err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
					if err != nil {
						panic(err)
					}
					AppId := args["AppId"].(string)
					Index := args["Index"].(string)
					// Simulate getting weather data
					result := runRankeeEvaluation(AppId, Index)

					params.Messages = append(params.Messages, openai.ToolMessage(result, toolCall.ID))
				}
			}
			completion, err = a.client.Chat.Completions.New(ctx, params)
			if err != nil {
				return fmt.Errorf("error generating completion after tool call: %w", err)
			}
		}

		params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
		fmt.Printf("\u001b[93mGPT\u001b[0m: %s\n", completion.Choices[0].Message.Content)
	}

	return nil
}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}

func runRankeeEvaluation(appID string, index string) string {
	// This function would contain the logic to run the Rankee evaluation
	// using the Client struct and its methods as described in the prompt.
	// For now, it's left empty as a placeholder.
	fmt.Println("Running Rankee evaluation...")
	return "running"
}

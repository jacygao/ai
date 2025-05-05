package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var promptHeading = "You are an AI Assistant designed to answer user questions based on the provided context.\n" +
	"Answer the question within the <answer> tags using the context within the <context> tags below:\n"
var propmtFooter = "Only generate answers strictly based on the context. Do not infer or assume additional details beyond what is provided.\n" +
	"For example, given <question>Who is Charlotte</question> <context>[Christiane is Charlotte's mum. Jacy is a software engineer.]</context> " +
	"the output should be Charlotte is Christinane's daughter - since that is explicitly mentioned in the context.\n" +
	"If the context does not contain sufficient information, respond exactly with: I do not have any knowledge to answer that question."

func GetPrompt(question string, context []string) string {
	return fmt.Sprintf(promptHeading+"<question>\n%s\n</question>\n"+"<context>\n%v\n</context>\n"+propmtFooter, question, context)
}

type RequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ResponseBody struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func Chat(query string, context []string) {
	url := "http://localhost:11434/api/generate"
	requestData := RequestBody{
		Model:  "llama3.2",
		Prompt: GetPrompt(query, context),
	}
	fmt.Println("prompt: ", requestData.Prompt)
	jsonData, _ := json.Marshal(requestData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read the full response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(body)))

	for scanner.Scan() {
		line := scanner.Text()
		var data ResponseBody
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			fmt.Println("Error unmarshalling response:", err)
			continue
		}
		fmt.Print(data.Response)
		time.Sleep(time.Millisecond * 300)

		if data.Done {
			break
		}
	}
}

type EmbedResponseBody struct {
	Model string      `json:"model"`
	Data  [][]float64 `json:"embeddings"`
}

func Embed(query string) ([]float64, error) {
	// Define the API endpoint and request payload
	url := "http://localhost:11434/api/embed"
	payload := map[string]any{
		"model": "all-minilm:l6-v2",
		"input": query,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return nil, err
	}

	// Create an HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	vector := &EmbedResponseBody{}
	if err := json.Unmarshal(body, vector); err != nil {
		fmt.Println(err)
	}

	// Print the response
	return vector.Data[0], nil
}

func Rerank() {
	url := "http://localhost:11434/api/generate"
	payload := map[string]any{
		"model": "all-minilm:l6-v2",
		"input": query,
	}
}

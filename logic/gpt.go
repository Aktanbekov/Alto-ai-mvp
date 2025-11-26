package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GetGPTResponse sends messages to OpenAI and returns the response
func GetGPTResponse(messages []Message) (string, error) {
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		return "", errors.New("GPT_API_KEY not set")
	}

	body := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": messages,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(data, &openAIResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		if openAIResp.Error.Message != "" {
			return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
		}
		return "", fmt.Errorf("OpenAI API returned status %d: %s", res.StatusCode, string(data))
	}

	if len(openAIResp.Choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func LogicGpt() { // Export the function for testing
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GPT_API_KEY not set")
		return
	}

	messages := []Message{
		{
			Role:    "user",
			Content: "Hello, how are you?",
		},
	}

	response, err := GetGPTResponse(messages)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Response:", response)
}

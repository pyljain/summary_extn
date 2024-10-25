package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type openai struct {
	apiKey  string
	baseURL string
	stream  bool
}

type openAIReqBody struct {
	MaxTokens int             `json:"max_completion_tokens"`
	Messages  []OpenAIMessage `json:"messages"`
	Model     string          `json:"model"`
	Stream    bool            `json:"stream,omitempty"`
}

type OpenAIMessage struct {
	Role    string          `json:"role"`
	Content []OpenAIContent `json:"content"`
}

type OpenAIContent struct {
	Type     string          `json:"type,omitempty"`
	Text     string          `json:"text,omitempty"`
	ImageUrl *OpenAIImageUrl `json:"image_url,omitempty"`
}

type OpenAIImageUrl struct {
	Url string `json:"url"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
	Id      string         `json:"id"`
	Model   string         `json:"model"`
}

type openAIStreamingResponse struct {
	Id               string                  `json:"id"`
	Model            string                  `json:"model"`
	StreamingChoices []streamingOpenAIChoice `json:"choices"`
}

type streamingOpenAIChoice struct {
	Index int     `json:"index"`
	Delta Content `json:"delta"`
}
type openAIChoice struct {
	Message Content `json:"message"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Role        string                 `json:"role,omitempty"`
	Text        string                 `json:"text,omitempty"`
	ContentType string                 `json:"type,omitempty"`
	Id          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Content     string                 `json:"content,omitempty"`
	ToolUseId   string                 `json:"tool_use_id,omitempty"`
	PartialJson *string                `json:"partial_json,omitempty"`
}

func CallOpenAI(prompt string) (string, error) {
	rb := openAIReqBody{
		Messages: []OpenAIMessage{
			{
				Role: "user",
				Content: []OpenAIContent{
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
		MaxTokens: 4096,
		Model:     "gpt-4o",
		Stream:    false,
	}

	rbBytes, err := json.Marshal(rb)
	if err != nil {
		return "", err
	}

	log.Printf("Request: %s", string(rbBytes))

	bufferedReq := bytes.NewBuffer(rbBytes)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/chat/completions", "https://api.openai.com/v1"), bufferedReq)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	req.Header.Add("content-type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to OpenAI: %s", err)
		return "", err
	}

	if resp.StatusCode >= 400 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading contents: %s", err)
			return "", errors.New(resp.Status)
		}
		log.Printf("Error calling OpenAI %s", string(respBytes))
		return "", errors.New(string(respBytes))
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading contents: %s", err)
		return "", err
	}

	var response openAIResponse
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		log.Printf("Error unmarshalling response: %s", err)
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

package llm

import (
	"fmt"
	"net/http"
	"time"
	globalconstant "voice-hack-backend/globalConstant"
	"voice-hack-backend/utilities/httpRequest"
)

// Function using MakeHttpCall
func GenerateInsightsViaLLM(userQuery string, systemPrompt string) (string, error) {
	// Prepare request body
	reqBody := map[string]any{
		"model": globalconstant.LLM_MODEL,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": userQuery,
			},
			{
				"role":    "system",
				"content": systemPrompt,
			},
		},
		"temperature": 0.7,
		// "max_tokens":  500,
	}

	// Prepare HttpRequest
	req := httpRequest.HttpRequest{
		Method: http.MethodPost,
		URL:    globalconstant.LLM_API_URL,
		Headers: map[string]any{
			"Authorization": "Bearer " + globalconstant.LLM_API_KEY,
			"Content-Type":  "application/json",
		},
		Body:    reqBody,
		Timeout: 30 * time.Second,
	}

	// Call LLM Gateway via MakeHttpCall
	resp := httpRequest.MakeHttpCall(req)
	if resp.Err != nil {
		return "", fmt.Errorf("LLM Gateway call failed: %v", resp.Err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM Gateway returned status %d", resp.StatusCode)
	}

	// Extract response text
	choices, ok := resp.Body["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices returned from LLM")
	}

	firstChoice, ok := choices[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid choice structure")
	}

	message, ok := firstChoice["message"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid message structure")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}

	return content, nil
}

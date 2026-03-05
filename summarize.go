package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type apiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type apiRequest struct {
	Model    string       `json:"model"`
	MaxToks  int          `json:"max_tokens"`
	System   string       `json:"system"`
	Messages []apiMessage `json:"messages"`
}

type apiContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type apiResponse struct {
	Content []apiContentBlock `json:"content"`
}

func formatCommits(commits []Commit) string {
	grouped := make(map[string][]string)
	var order []string
	for _, c := range commits {
		if _, exists := grouped[c.Repo]; !exists {
			order = append(order, c.Repo)
		}
		grouped[c.Repo] = append(grouped[c.Repo], c.Message)
	}

	var sb strings.Builder
	for i, repo := range order {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(repo + ":\n")
		for _, msg := range grouped[repo] {
			sb.WriteString("- " + msg + "\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func Summarize(commits []Commit, apiKey string, model string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("API key not found. Run: clog config --api-key \"your-key\"")
	}

	reqBody := apiRequest{
		Model:   model,
		MaxToks: 300,
		System:  "You are a developer assistant. Summarize the following git commits into 2-4 sentences written in first person, past tense, suitable for a daily standup or async team update. Be concise. Focus on what changed and why it matters, not technical details. Do not use bullet points. Output plain text only.",
		Messages: []apiMessage{
			{Role: "user", Content: formatCommits(commits)},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return apiResp.Content[0].Text, nil
}

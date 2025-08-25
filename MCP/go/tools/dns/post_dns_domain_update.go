package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/solar-vps/mcp-server/config"
	"github.com/solar-vps/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Post_dns_domain_updateHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		domainVal, ok := args["domain"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: domain"), nil
		}
		domain, ok := domainVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: domain"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["id"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("id=%v", val))
		}
		if val, ok := args["name"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("name=%v", val))
		}
		if val, ok := args["type"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("type=%v", val))
		}
		if val, ok := args["content"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("content=%v", val))
		}
		if val, ok := args["ttl"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("ttl=%v", val))
		}
		if val, ok := args["prio"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("prio=%v", val))
		}
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("api_key=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/dns/%s/update%s", cfg.BaseURL, domain, queryString)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			// API key already added to query string
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreatePost_dns_domain_updateTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_dns_domain_update",
		mcp.WithDescription("Update dns record for a given domain"),
		mcp.WithString("domain", mcp.Required(), mcp.Description("Domain name to add record under")),
		mcp.WithString("id", mcp.Required(), mcp.Description("Id of DNS record")),
		mcp.WithString("name", mcp.Description("Fully qualified name for the DNS record")),
		mcp.WithString("type", mcp.Description("Type for DNS record")),
		mcp.WithString("content", mcp.Description("Content for the DNS Record")),
		mcp.WithString("ttl", mcp.Description("Time To Live for DNS record")),
		mcp.WithString("prio", mcp.Description("Priority of the DNS record")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Post_dns_domain_updateHandler(cfg),
	}
}

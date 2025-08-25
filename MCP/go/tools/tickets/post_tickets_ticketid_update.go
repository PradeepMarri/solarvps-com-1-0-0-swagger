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

func Post_tickets_ticketid_updateHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		ticketidVal, ok := args["ticketid"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: ticketid"), nil
		}
		ticketid, ok := ticketidVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: ticketid"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["contents"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("contents=%v", val))
		}
		// Fallback to single auth parameter
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("api_key=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/tickets/%s/update%s", cfg.BaseURL, ticketid, queryString)
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

func CreatePost_tickets_ticketid_updateTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_tickets_ticketid_update",
		mcp.WithDescription("Post a reply to a ticket"),
		mcp.WithString("ticketid", mcp.Required(), mcp.Description("TicketId of the ticket you want to post an update to")),
		mcp.WithString("contents", mcp.Required(), mcp.Description("Message reply being sent")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Post_tickets_ticketid_updateHandler(cfg),
	}
}

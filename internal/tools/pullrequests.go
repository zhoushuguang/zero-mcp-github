package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zhoushuguang/zero-mcp-github/internal/svc"

	"github.com/google/go-github/v69/github"
	"github.com/zeromicro/go-zero/mcp"
)

func listPullRequests(svcCtx *svc.ServiceContext) mcp.Tool {
	var listPullRequestsHandler = func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Owner     string  `json:"owner"`
			Repo      string  `json:"repo"`
			State     string  `json:"state,optional"`
			Head      string  `json:"head,optional"`
			Base      string  `json:"base,optional"`
			Sort      string  `json:"sort,optional"`
			Direction string  `json:"direction,optional"`
			Page      float64 `json:"page,optional"`
			PerPage   float64 `json:"perPage,optional"`
		}
		err := mcp.ParseArguments(params, &req)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		opts := &github.PullRequestListOptions{
			State:     req.State,
			Head:      req.Head,
			Base:      req.Base,
			Sort:      req.Sort,
			Direction: req.Direction,
			ListOptions: github.ListOptions{
				PerPage: int(req.PerPage),
				Page:    int(req.Page),
			},
		}

		prs, resp, err := svcCtx.GithubClient.PullRequests.List(ctx, req.Owner, req.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}
			return mcp.CallToolResult{
				Content: []any{fmt.Sprintf("Error: %s", string(body))},
				IsError: true,
			}, nil
		}
		r, err := json.Marshal(prs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issue: %w", err)
		}
		return mcp.CallToolResult{
			Content: []any{r},
			IsError: false,
		}, nil
	}

	return mcp.Tool{
		Name:        "list_pull_requests",
		Description: "List pull requests in a GitHub repository.",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"owner": map[string]any{
					"type":        "string",
					"description": "Repository owner",
				},
				"repo": map[string]any{
					"type":        "string",
					"description": "Repository name",
				},
				"state": map[string]any{
					"type":        "string",
					"description": "Filter by state",
					"enum":        []string{"open", "closed", "all"},
				},
				"head": map[string]any{
					"type":        "string",
					"description": "Filter by head user/org and branch",
				},
				"base": map[string]any{
					"type":        "string",
					"description": "Filter by base branch",
				},
				"sort": map[string]any{
					"type":        "string",
					"description": "Sort by",
					"enum":        []string{"created", "updated", "popularity", "long-running"},
				},
				"direction": map[string]any{
					"type":        "string",
					"description": "Sort direction",
					"enum":        []string{"asc", "desc"},
				},
				"page": map[string]any{
					"type":        "number",
					"description": "Page number for pagination (min 1)",
					"minimum":     1,
				},
				"perPage": map[string]any{
					"type":        "number",
					"description": "Results per page for pagination (min 1, max 100)",
					"minimum":     1,
				},
			},
			Required: []string{"owner", "repo"},
		},
		Handler: listPullRequestsHandler,
	}
}

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zhoushuguang/zero-mcp-github/internal/svc"

	"github.com/google/go-github/v69/github"
	"github.com/zeromicro/go-zero/mcp"
)

func listIssuesTool(svcCtx *svc.ServiceContext) mcp.Tool {
	var listIssuesTool = func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Owner     string   `json:"owner"`
			Repo      string   `json:"repo"`
			State     string   `json:"state,optional"`
			Labels    []string `json:"labels,optional"`
			Sort      string   `json:"sort,optional"`
			Direction string   `json:"direction,optional"`
			Since     string   `json:"since,optional"`
			Page      float64  `json:"page,optional"`
			PerPage   float64  `json:"perPage,optional"`
		}
		err := mcp.ParseArguments(params, &req)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}

		var sinceTime time.Time
		if req.Since != "" {
			sinceTime, err = parseISOTimestamp(req.Since)
			if err != nil {
				return nil, fmt.Errorf("failed to parse timestamp: %w", err)
			}
		}
		opts := &github.IssueListByRepoOptions{
			State:     req.State,
			Labels:    req.Labels,
			Sort:      req.Sort,
			Direction: req.Direction,
			Since:     sinceTime,
			ListOptions: github.ListOptions{
				Page:    int(req.Page),
				PerPage: int(req.PerPage),
			},
		}

		issues, resp, err := svcCtx.GithubClient.Issues.ListByRepo(ctx, req.Owner, req.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
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

		r, err := json.Marshal(issues)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issue: %w", err)
		}
		return mcp.CallToolResult{
			Content: []any{r},
			IsError: false,
		}, nil
	}

	return mcp.Tool{
		Name:        "list_issues",
		Description: "List issues in a GitHub repository.",
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
					"description": "TFilter by state",
					"enum":        []string{"open", "closed", "all"},
				},
				"labels": map[string]any{
					"type":        "array",
					"description": "Filter by labels",
					"items": map[string]any{
						"type": "string",
					},
				},
				"sort": map[string]any{
					"type":        "string",
					"description": "Sort order",
					"enum":        []string{"created", "updated", "comments"},
				},
				"direction": map[string]any{
					"type":        "string",
					"description": "Sort direction",
					"enum":        []string{"asc", "desc"},
				},
				"since": map[string]any{
					"type":        "string",
					"description": "Filter by date (ISO 8601 timestamp)",
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
					"maximum":     100,
				},
			},
			Required: []string{"owner", "repo"},
		},
		Handler: listIssuesTool,
	}
}

func getIssueTool(svcCtx *svc.ServiceContext) mcp.Tool {
	var getIssueHandler = func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Owner       string `json:"owner"`
			Repo        string `json:"repo"`
			IssueNumber int    `json:"issue_number"`
		}
		err := mcp.ParseArguments(params, &req)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		issue, resp, err := svcCtx.GithubClient.Issues.Get(ctx, req.Owner, req.Repo, req.IssueNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get issue: %w", err)
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

		r, err := json.Marshal(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issue: %w", err)
		}
		return mcp.CallToolResult{
			Content: []any{r},
			IsError: false,
		}, nil
	}

	return mcp.Tool{
		Name:        "get_issue",
		Description: "Get details of a specific issue in a GitHub repository.",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"owner": map[string]any{
					"type":        "string",
					"description": "The owner of the repository",
				},
				"repo": map[string]any{
					"type":        "string",
					"description": "The name of the repository",
				},
				"issue_number": map[string]any{
					"type":        "number",
					"description": "The number of the issue",
				},
			},
			Required: []string{"owner", "repo", "issue_number"},
		},
		Handler: getIssueHandler,
	}
}

func createIssueTool(svcCtx *svc.ServiceContext) mcp.Tool {
	var createIssueHandler = func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Owner     string   `json:"owner"`
			Repo      string   `json:"repo"`
			Title     string   `json:"title"`
			Body      string   `json:"body,optional"`
			Assignees []string `json:"assignees,optional"`
			Labels    []string `json:"labels,optional"`
			Milestone float64  `json:"milestone,optional"`
		}
		err := mcp.ParseArguments(params, &req)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		var milestoneNum *int
		if req.Milestone != 0 {
			v := int(req.Milestone)
			milestoneNum = &v
		}
		issueRequest := github.IssueRequest{
			Title:     github.Ptr(req.Title),
			Body:      github.Ptr(req.Body),
			Milestone: milestoneNum,
		}
		if req.Assignees != nil {
			issueRequest.Assignees = &req.Assignees
		}
		if req.Labels != nil {
			issueRequest.Labels = &req.Labels
		}
		issue, resp, err := svcCtx.GithubClient.Issues.Create(ctx, req.Owner, req.Repo, &issueRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to create issue: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusCreated {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}
			return mcp.CallToolResult{
				Content: []any{fmt.Sprintf("Error: %s", string(body))},
				IsError: true,
			}, nil
		}

		r, err := json.Marshal(issue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal issue: %w", err)
		}
		return mcp.CallToolResult{
			Content: []any{r},
			IsError: false,
		}, nil
	}

	return mcp.Tool{
		Name:        "create_issue",
		Description: "Create a new issue in a GitHub repository.",
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
				"title": map[string]any{
					"type":        "string",
					"description": "Issue title",
				},
				"body": map[string]any{
					"type":        "string",
					"description": "Issue body content",
				},
				"assignees": map[string]any{
					"type":        "array",
					"description": "usernames to assign to the issue",
					"items": map[string]any{
						"type": "string",
					},
				},
				"labels": map[string]any{
					"type":        "array",
					"description": "Labels to apply to the issue",
					"items": map[string]any{
						"type": "string",
					},
				},
				"milestone": map[string]any{
					"type":        "number",
					"description": "Milestone number",
				},
			},
			Required: []string{"owner", "repo", "title"},
		},
		Handler: createIssueHandler,
	}
}

func parseISOTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("2006-01-02", timestamp)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid ISO 8601 timestamp: %s (supported formats: YYYY-MM-DDThh:mm:ssZ or YYYY-MM-DD)", timestamp)
}

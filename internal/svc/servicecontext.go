package svc

import (
	"github.com/google/go-github/v69/github"
	"github.com/zhoushuguang/zero-mcp-github/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	GithubClient *github.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		GithubClient: github.NewClient(nil).WithAuthToken(c.Github.Token),
	}
}

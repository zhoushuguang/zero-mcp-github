package config

import (
	"github.com/zeromicro/go-zero/mcp"
)

type Config struct {
	mcp.McpConf
	Github struct {
		Token string `json:"token"`
	}
}

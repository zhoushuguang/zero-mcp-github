package main

import (
	"flag"
	"fmt"
	
	"github.com/zhoushuguang/zero-mcp-github/internal/svc"
	"github.com/zhoushuguang/zero-mcp-github/internal/tools"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/mcp"
	"github.com/zhoushuguang/zero-mcp-github/internal/config"
)

var configFile = flag.String("f", "etc/zero-mcp-github.yaml", "the config file")

func main() {
	flag.Parse()

	logx.DisableStat()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	mcpSrv := mcp.NewMcpServer(c.McpConf)
	defer mcpSrv.Stop()

	svcCtx := svc.NewServiceContext(c)
	tools.NewToolset(mcpSrv, svcCtx)

	fmt.Printf("Starting MCP Server on %s:%d\n", c.McpConf.Host, c.McpConf.Port)
	mcpSrv.Start()
}

package tools

import (
	"github.com/zhoushuguang/zero-mcp-github/internal/svc"

	"github.com/zeromicro/go-zero/mcp"
)

type Toolset struct {
	svcCtx *svc.ServiceContext
	mcpSrv mcp.McpServer
}

func NewToolset(mcpSrv mcp.McpServer, svcCtx *svc.ServiceContext) {
	toolSet := &Toolset{
		svcCtx: svcCtx,
		mcpSrv: mcpSrv,
	}
	toolSet.addTools()
}

func (t *Toolset) addTools() {
	if err := t.mcpSrv.RegisterTool(listIssuesTool(t.svcCtx)); err != nil {
		panic(err)
	}
	if err := t.mcpSrv.RegisterTool(getIssueTool(t.svcCtx)); err != nil {
		panic(err)
	}
	if err := t.mcpSrv.RegisterTool(createIssueTool(t.svcCtx)); err != nil {
		panic(err)
	}
	if err := t.mcpSrv.RegisterTool(listPullRequests(t.svcCtx)); err != nil {
		panic(err)
	}
}

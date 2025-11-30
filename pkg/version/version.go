package version

// Version 服务端版本号（通过 -ldflags 注入）
var Version = "dev"

// AgentVersion Agent 版本号（通过 -ldflags 注入）
var AgentVersion = "dev"

// GetVersion 获取服务端版本号
func GetVersion() string {
	if Version == "" {
		return "dev"
	}
	return Version
}

// GetAgentVersion 获取 Agent 版本号
func GetAgentVersion() string {
	if AgentVersion == "" {
		return "dev"
	}
	return AgentVersion
}

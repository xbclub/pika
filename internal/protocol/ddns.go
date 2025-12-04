package protocol

// DDNSConfigData DDNS 配置数据（服务端下发给客户端）
type DDNSConfigData struct {
	Enabled bool `json:"enabled"` // 是否启用 DDNS

	// IP 获取配置
	EnableIPv4    bool   `json:"enableIpv4"`              // 是否启用 IPv4
	EnableIPv6    bool   `json:"enableIpv6"`              // 是否启用 IPv6
	IPv4GetMethod string `json:"ipv4GetMethod,omitempty"` // IPv4 获取方式: api, interface, command
	IPv6GetMethod string `json:"ipv6GetMethod,omitempty"` // IPv6 获取方式: api, interface, command
	IPv4GetValue  string `json:"ipv4GetValue,omitempty"`  // IPv4 获取配置值（接口名/API URL/命令）
	IPv6GetValue  string `json:"ipv6GetValue,omitempty"`  // IPv6 获取配置值（接口名/API URL/命令）
}

// DDNSIPReportData DDNS IP 上报数据（客户端发送）
type DDNSIPReportData struct {
	IPv4 string `json:"ipv4,omitempty"` // IPv4 地址
	IPv6 string `json:"ipv6,omitempty"` // IPv6 地址
}

package ddns

import (
	"context"
	"fmt"
	"strings"
)

// Provider DNS 服务商接口
type Provider interface {
	// UpdateRecord 更新 DNS 记录
	// domain: 完整域名，如 ddns.example.com
	// recordType: 记录类型，A 或 AAAA
	// ip: IP 地址
	UpdateRecord(ctx context.Context, domain, recordType, ip string) error

	// GetRecord 获取 DNS 记录
	// domain: 完整域名
	// recordType: 记录类型，A 或 AAAA
	// 返回: 当前记录的 IP 地址
	GetRecord(ctx context.Context, domain, recordType string) (string, error)
}

// RecordType DNS 记录类型
const (
	RecordTypeA    = "A"
	RecordTypeAAAA = "AAAA"
)

// parseDomain 解析域名，提取主域名和子域名
// 例如: ddns.example.com -> example.com, ddns
// 例如: example.com -> example.com, @
func parseDomain(fullDomain string) (zone, name string, err error) {
	parts := strings.Split(fullDomain, ".")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("无效的域名格式: %s", fullDomain)
	}

	if len(parts) == 2 {
		// 例如: example.com，主机记录是 @
		return fullDomain, "@", nil
	}

	// 例如: ddns.example.com -> zone: example.com, name: ddns
	zone = strings.Join(parts[len(parts)-2:], ".")
	name = strings.Join(parts[:len(parts)-2], ".")

	return zone, name, nil
}

package ddns

import (
	"fmt"

	"github.com/libdns/alidns"
	"github.com/libdns/cloudflare"
	"github.com/libdns/huaweicloud"
	"github.com/libdns/tencentcloud"
)

// NewProvider 创建 DNS 提供商
func NewProvider(providerType string, config map[string]string) (Provider, error) {
	var libdnsProvider interface{}

	switch providerType {
	case "aliyun":
		accessKeyID, ok := config["accessKeyId"]
		if !ok || accessKeyID == "" {
			return nil, fmt.Errorf("阿里云 AccessKeyId 不能为空")
		}
		accessKeySecret, ok := config["accessKeySecret"]
		if !ok || accessKeySecret == "" {
			return nil, fmt.Errorf("阿里云 AccessKeySecret 不能为空")
		}

		libdnsProvider = &alidns.Provider{
			CredentialInfo: alidns.CredentialInfo{
				AccessKeyID:     accessKeyID,
				AccessKeySecret: accessKeySecret,
			},
		}

	case "tencentcloud":
		secretID, ok := config["secretId"]
		if !ok || secretID == "" {
			return nil, fmt.Errorf("腾讯云 SecretId 不能为空")
		}
		secretKey, ok := config["secretKey"]
		if !ok || secretKey == "" {
			return nil, fmt.Errorf("腾讯云 SecretKey 不能为空")
		}

		libdnsProvider = &tencentcloud.Provider{
			SecretId:  secretID,
			SecretKey: secretKey,
		}

	case "cloudflare":
		apiToken, ok := config["apiToken"]
		if !ok || apiToken == "" {
			return nil, fmt.Errorf("Cloudflare API Token 不能为空")
		}

		libdnsProvider = &cloudflare.Provider{
			APIToken: apiToken,
		}

	case "huaweicloud":
		accessKeyID, ok := config["accessKeyId"]
		if !ok || accessKeyID == "" {
			return nil, fmt.Errorf("华为云 AccessKeyId 不能为空")
		}
		secretAccessKey, ok := config["secretAccessKey"]
		if !ok || secretAccessKey == "" {
			return nil, fmt.Errorf("华为云 SecretAccessKey 不能为空")
		}
		region, ok := config["region"]
		if !ok || region == "" {
			region = "cn-south-1" // 默认区域
		}

		libdnsProvider = &huaweicloud.Provider{
			AccessKeyId:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			RegionId:        region,
		}

	default:
		return nil, fmt.Errorf("不支持的 DNS 服务商: %s", providerType)
	}

	return NewLibDNSProvider(libdnsProvider)
}

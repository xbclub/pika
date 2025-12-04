package ddns

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/libdns/libdns"
)

// LibDNSProvider 基于 libdns 的通用 DNS 提供商
type LibDNSProvider struct {
	getter libdns.RecordGetter
	setter libdns.RecordSetter
}

// NewLibDNSProvider 创建基于 libdns 的提供商
func NewLibDNSProvider(provider interface{}) (*LibDNSProvider, error) {
	getter, okGetter := provider.(libdns.RecordGetter)
	setter, okSetter := provider.(libdns.RecordSetter)

	if !okGetter || !okSetter {
		return nil, fmt.Errorf("提供商必须实现 RecordGetter 和 RecordSetter 接口")
	}

	return &LibDNSProvider{
		getter: getter,
		setter: setter,
	}, nil
}

// UpdateRecord 更新 DNS 记录
func (p *LibDNSProvider) UpdateRecord(ctx context.Context, domain, recordType, ip string) error {
	zone, name, err := parseDomain(domain)
	if err != nil {
		return err
	}

	// 解析 IP 地址
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("无效的 IP 地址: %w", err)
	}

	// 获取现有记录
	records, err := p.getter.GetRecords(ctx, zone)
	if err != nil {
		return fmt.Errorf("获取 DNS 记录失败: %w", err)
	}

	// 查找匹配的记录
	existingRecord := findAddressRecord(records, name, recordType)

	// 如果记录存在且 IP 相同，无需更新
	if existingRecord != nil {
		existingRR := existingRecord.RR()
		existingAddr, ok := existingRecord.(libdns.Address)
		if ok && existingAddr.IP == addr {
			return nil
		}
		// 检查是否是 RR 类型且 Data 相同
		if existingRR.Data == ip {
			return nil
		}
	}

	// 构建新记录
	newRecord := libdns.Address{
		Name: name,
		IP:   addr,
		TTL:  10 * time.Minute, // 默认 TTL 10 分钟
	}

	// 更新记录
	_, err = p.setter.SetRecords(ctx, zone, []libdns.Record{newRecord})
	if err != nil {
		return fmt.Errorf("更新 DNS 记录失败: %w", err)
	}

	return nil
}

// GetRecord 获取 DNS 记录
func (p *LibDNSProvider) GetRecord(ctx context.Context, domain, recordType string) (string, error) {
	zone, name, err := parseDomain(domain)
	if err != nil {
		return "", err
	}

	// 获取所有记录
	records, err := p.getter.GetRecords(ctx, zone)
	if err != nil {
		return "", fmt.Errorf("获取 DNS 记录失败: %w", err)
	}

	// 查找匹配的记录
	record := findAddressRecord(records, name, recordType)
	if record == nil {
		return "", fmt.Errorf("未找到 DNS 记录")
	}

	// 尝试从 Address 类型获取 IP
	if addr, ok := record.(libdns.Address); ok {
		return addr.IP.String(), nil
	}

	// 降级使用 RR.Data
	rr := record.RR()
	return rr.Data, nil
}

// findAddressRecord 在记录列表中查找匹配的地址记录
func findAddressRecord(records []libdns.Record, name, recordType string) libdns.Record {
	for _, record := range records {
		rr := record.RR()
		// 匹配名称和类型
		if rr.Name == name && rr.Type == recordType {
			return record
		}
	}
	return nil
}

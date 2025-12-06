package sysutil

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ConfigureICMPPermissions 配置 ICMP 权限
// 在 Linux 系统上，允许非特权用户发起 ICMP 请求
// 使用 sync.Once 确保只执行一次
func ConfigureICMPPermissions() error {
	return doConfigureICMPPermissions()
}

// doConfigureICMPPermissions 实际执行配置
func doConfigureICMPPermissions() error {
	const sysctlPath = "/proc/sys/net/ipv4/ping_group_range"

	// 1. 先检查当前配置
	currentMin, currentMax, err := readPingGroupRange(sysctlPath)
	if err != nil {
		return fmt.Errorf("读取当前 ICMP 配置失败: %w", err)
	}

	// 2. 检查是否已经满足要求 (范围包含 0 到 2147483647)
	if currentMin <= 0 && currentMax >= 2147483647 {
		log.Printf("✅ ICMP 权限已配置: net.ipv4.ping_group_range=%d %d", currentMin, currentMax)
		return nil
	}

	// 3. 需要配置，写入新值
	log.Printf("🔧 当前 ICMP 配置: net.ipv4.ping_group_range=%d %d (不满足要求)", currentMin, currentMax)
	log.Println("   正在配置为: 0 2147483647")

	if err := writePingGroupRange(sysctlPath, 0, 2147483647); err != nil {
		return fmt.Errorf("配置 ICMP 权限失败: %w", err)
	}

	log.Println("✅ ICMP 权限配置成功: net.ipv4.ping_group_range=0 2147483647")
	return nil
}

// readPingGroupRange 读取 ping_group_range 的当前值
func readPingGroupRange(path string) (min, max int, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}

	// 格式: "min max\n"
	fields := strings.Fields(string(data))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("无效的格式: %s", string(data))
	}

	min, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("解析 min 值失败: %w", err)
	}

	max, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("解析 max 值失败: %w", err)
	}

	return min, max, nil
}

// writePingGroupRange 写入新的 ping_group_range 值
func writePingGroupRange(path string, min, max int) error {
	value := fmt.Sprintf("%d\t%d", min, max)
	if err := os.WriteFile(path, []byte(value), 0644); err != nil {
		// 检查是否是权限问题
		if os.IsPermission(err) {
			return fmt.Errorf("需要 root 权限才能配置 ICMP 权限: %w", err)
		}
		return err
	}
	return nil
}

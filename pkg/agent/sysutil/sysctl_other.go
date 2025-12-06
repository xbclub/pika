//go:build !linux
// +build !linux

package sysutil

// ConfigureICMPPermissions 在非 Linux 平台上不需要配置
func ConfigureICMPPermissions() error {
	// Windows 和 macOS 不需要特殊配置
	return nil
}

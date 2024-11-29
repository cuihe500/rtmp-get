package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 最低要求的 Npcap 版本
const minRequiredVersion = "1.60"

// 添加 checkNpcapDependency 函数
func checkNpcapDependency() error {
	Debug("开始检查 Npcap 依赖...")

	// 获取版本（如果未安装会自动安装）
	version, err := getNpcapVersion()
	if err != nil {
		return fmt.Errorf("Npcap 检查失败: %v", err)
	}

	Debug("检查版本兼容性...")
	// 检查最低版本要求
	if !isVersionCompatible(version) {
		return fmt.Errorf("Npcap 版本过低，请升级到最新版本")
	}

	Info("Npcap 依赖检查通过")
	return nil
}

func getNpcapVersion() (string, error) {
	Debug("开始检查 Npcap 版本...")

	// 首先确保已安装
	if !isNpcapInstalled() {
		Debug("Npcap 未安装，开始安装...")
		if err := ensureNpcapInstalled(); err != nil {
			return "", fmt.Errorf("安装 Npcap 失败: %v", err)
		}
		// 安装后等待系统初始化
		Debug("等待系统初始化...")
		time.Sleep(15 * time.Second)
	}

	// 首先尝试使用注册表查询版本
	Debug("尝试从注册表获取版本...")
	cmd := exec.Command("reg", "query", "HKLM\\SOFTWARE\\Npcap", "/v", "Version")
	output, err := cmd.Output()
	if err == nil {
		re := regexp.MustCompile(`\d+\.\d+`)
		version := re.FindString(string(output))
		if version != "" {
			Debug("从注册表获取到版本: %s", version)
			return version, nil
		}
	}

	// 如果注册表查询失败，尝试从 DLL 文件获取版本
	Debug("尝试从 DLL 获取版本...")
	dllPaths := []string{
		"C:\\Windows\\System32\\Npcap\\NPFInstall.exe",
		"C:\\Windows\\System32\\Npcap\\npcap.dll",
		"C:\\Program Files\\Npcap\\NPFInstall.exe",
		"C:\\Program Files (x86)\\Npcap\\NPFInstall.exe",
	}

	for _, path := range dllPaths {
		Debug("检查文件: %s", path)
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("(Get-Item '%s').VersionInfo.FileVersion", path))
		output, err := cmd.Output()
		if err == nil {
			version := strings.TrimSpace(string(output))
			re := regexp.MustCompile(`\d+\.\d+`)
			version = re.FindString(version)
			if version != "" {
				Debug("从文件获取到版本: %s", version)
				return version, nil
			}
		}
	}

	// 如果还是无法获取版本，但确认已安装，返回默认版本
	if isNpcapInstalled() {
		Warn("无法获取具体版本，但确认已安装，使用默认版本 1.60")
		return "1.60", nil
	}

	return "", fmt.Errorf("无法获取 Npcap 版本信息")
}

func isVersionCompatible(version string) bool {
	Debug("检查版本兼容性: 当前版本 %s，最低要求版本 %s", version, minRequiredVersion)

	current := strings.Split(version, ".")
	required := strings.Split(minRequiredVersion, ".")

	if len(current) < 2 || len(required) < 2 {
		Debug("版本号格式无效")
		return false
	}

	currentMajor, _ := strconv.Atoi(current[0])
	currentMinor, _ := strconv.Atoi(current[1])
	requiredMajor, _ := strconv.Atoi(required[0])
	requiredMinor, _ := strconv.Atoi(required[1])

	if currentMajor > requiredMajor {
		Debug("版本兼容：主版本号高于要求")
		return true
	}
	if currentMajor == requiredMajor {
		compatible := currentMinor >= requiredMinor
		Debug("版本兼容性检查结果: %v", compatible)
		return compatible
	}
	Debug("版本不兼容：版本过低")
	return false
}

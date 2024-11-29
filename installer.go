package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed resources/npcap-installer.exe
var npcapInstaller embed.FS

func ensureNpcapInstalled() error {
	Debug("检查 Npcap 是否已安装...")
	if isNpcapInstalled() {
		Debug("检测到 Npcap 已安装")
		return nil
	}

	Debug("未检测到 Npcap，准备安装...")

	Debug("正在提取 Npcap 安装程序...")
	installerData, err := npcapInstaller.ReadFile("resources/npcap-installer.exe")
	if err != nil {
		return fmt.Errorf("读取安装程序失败: %v", err)
	}

	tempFile := filepath.Join(os.TempDir(), "npcap-installer.exe")
	Debug("创建临时安装文件: %s", tempFile)
	if err := os.WriteFile(tempFile, installerData, 0755); err != nil {
		return fmt.Errorf("创建临时安装文件失败: %v", err)
	}
	defer os.Remove(tempFile)

	Debug("正在启动 Npcap 安装程序...")

	cmd := exec.Command(tempFile, "/S", "/winpcap_mode=yes")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}

	Debug("等待安装完成...")
	time.Sleep(10 * time.Second)

	Debug("验证安装结果...")
	if !isNpcapInstalled() {
		return fmt.Errorf("安装似乎失败了，请手动安装 Npcap")
	}

	Debug("Npcap 安装成功！")
	return nil
}

func isNpcapInstalled() bool {
	paths := []string{
		"C:\\Windows\\System32\\Npcap",
		"C:\\Program Files\\Npcap",
		"C:\\Program Files (x86)\\Npcap",
	}

	for _, path := range paths {
		Debug("检查路径: %s", path)
		if _, err := os.Stat(path); err == nil {
			Debug("在路径 %s 找到 Npcap", path)
			return true
		}
	}
	Debug("未找到 Npcap 安装")
	return false
}

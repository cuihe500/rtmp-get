package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Config 配置结构体
type Config struct {
	Interface     string `json:"interface"`      // 网络接口名称
	DisplayFilter string `json:"display_filter"` // 过滤规则
	TsharkPath    string `json:"tshark_path"`    // tshark路径（在Go版本中不需要）
}

// PacketSniffer 数据包嗅探器结构体
type PacketSniffer struct {
	interface_    string
	displayFilter string
	wg            sync.WaitGroup
	stopChan      chan struct{}
}

func NewPacketSniffer(interface_ string, displayFilter string) *PacketSniffer {
	return &PacketSniffer{
		interface_:    interface_,
		displayFilter: displayFilter,
		stopChan:      make(chan struct{}),
	}
}

func (ps *PacketSniffer) Start() error {
	handle, err := pcap.OpenLive(ps.interface_, 1600, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("打开网络接口错误: %v", err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(ps.displayFilter); err != nil {
		return fmt.Errorf("设置过滤器错误: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	ps.wg.Add(1)
	go func() {
		defer ps.wg.Done()
		for {
			select {
			case <-ps.stopChan:
				Info("正在停止数据包捕获...")
				return
			default:
				packet, err := packetSource.NextPacket()
				if err != nil {
					continue
				}
				packetData := packet.String()

				if server := filterStrings(packetData, "rtmp://"); server != "" {
					Info("服务器: %s", server)
				}

				if code := filterStrings(packetData, "stream-"); code != "" {
					code = strings.Trim(code, "\"")
					Info("推流码: %s", code)
				}
			}
		}
	}()

	return nil
}

func (ps *PacketSniffer) Stop() {
	close(ps.stopChan)
}

func (ps *PacketSniffer) Wait() {
	ps.wg.Wait()
}

func filterStrings(inputStr, targetStr string) string {
	words := strings.Fields(inputStr)
	for _, word := range words {
		if strings.Contains(word, targetStr) {
			return word
		}
	}
	return ""
}

func listInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		Fatal("获取网络接口列表失败: %v", err)
	}

	if IsDebugMode() {
		Debug("可用网络接口列表:")
		for _, device := range devices {
			Debug("----------------------------------------")
			Debug("接口名称: %s", device.Name)
			Debug("接口描述: %s", device.Description)
			if len(device.Addresses) > 0 {
				Debug("IP地址列表:")
				for _, address := range device.Addresses {
					Debug("  - IP地址: %s", address.IP)
					Debug("    子网掩码: %s", address.Netmask)
					if address.Broadaddr != nil {
						Debug("    广播地址: %s", address.Broadaddr)
					}
				}
			}
			Debug("----------------------------------------")
		}
	}
}

func main() {
	// 解析命令行参数
	debug := flag.Bool("debug", false, "启用调试模式")
	logLevel := flag.String("log-level", "INFO", "日志级别 (TRACE/DEBUG/INFO/WARN/ERROR)")
	flag.Parse()

	// 设置日志级别和调试模式
	SetDebugMode(*debug)
	SetLogLevel(*logLevel)

	// 添加测试输出
	Debug("测试调试输出")
	Info("测试信息输出")

	Info("程序启动...")
	Debug("调试模式已启用")

	// 确保 Npcap 已安装
	if err := checkNpcapDependency(); err != nil {
		Fatal("检查 Npcap 依赖失败: %v", err)
	}

	// 列出所有网络接口
	listInterfaces()

	// 读取配置文件
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		Fatal("读取配置文件错误: %v", err)
	}

	// 解析配置
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		Fatal("解析配置文件错误: %v", err)
	}

	Debug("配置信息:")
	Debug("- 网络接口: %s", config.Interface)
	Debug("- 过滤规则: %s", config.DisplayFilter)

	// 创建并启动嗅探器
	sniffer := NewPacketSniffer(config.Interface, config.DisplayFilter)
	if err := sniffer.Start(); err != nil {
		Fatal("启动嗅探器错误: %v", err)
	}

	Info("开始捕获数据包...")
	Info("按 Ctrl+C 停止捕获...")

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待中断信号
	<-sigChan

	// 优雅停止
	Info("正在停止程序...")
	sniffer.Stop()
	sniffer.Wait()

	Info("程序已停止")
	fmt.Println("\n按回车键退出...")
	fmt.Scanln()
}

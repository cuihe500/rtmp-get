# RTMP-Get

一个用于捕获RTMP推流地址和推流码的工具。

## 功能特性

- 自动捕获网络数据包中的RTMP服务器地址
- 自动提取推流码
- 支持自定义网络接口和过滤规则
- 自动安装Npcap依赖（Windows）

## 环境要求

### Windows
- Go 1.16 或更高版本
- Npcap（程序会自动安装）

### Linux（用于交叉编译）
- Go 1.16 或更高版本
- MinGW-w64 交叉编译工具链

## 安装依赖

### CentOS/RHEL 

#### 安装交叉编译工具链

```bash
dnf --enablerepo=crb install mingw64-gcc
```

### Ubuntu/Debian

#### 安装交叉编译工具链

```bash
sudo apt-get update
sudo apt-get install gcc-mingw-w64
```

## 编译

### Windows下编译

```bash
go build -o rtmp-get.exe main.go
```

### Linux下交叉编译Windows版本

#### 设置交叉编译环境变量

```bash
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/lib/pkgconfig
```

## 配置

创建 `config.json` 文件：

```json
{
"interface": "\\Device\\NPF_{XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}",
"display_filter": "tcp"
}
```

- `interface`: 网络接口名称（运行程序时会显示可用接口列表）
- `display_filter`: 数据包过滤规则

## 使用方法

1. 运行程序（需要管理员权限）：

```bash
rtmp-get.exe
```

2. 程序会显示可用的网络接口列表

3. 将需要监听的网络接口名称复制到 `config.json` 中

4. 重新运行程序开始捕获

## 注意事项

1. 程序需要管理员权限运行

2. 首次运行会自动安装Npcap（如果未安装）

3. 确保配置文件中的网络接口名称正确

4. 如果手动安装Npcap，建议使用最新版本

## 常见问题

1. 找不到网络接口
   - 确保以管理员权限运行

   - 检查Npcap是否正确安装

   - 检查配置文件中的接口名称


2. 编译错误
   - 确保已安装所有必要的依赖

   - 检查环境变量是否正确设置

   - 尝试清理并重新下载依赖：

     ```bash
     go clean -modcache
     go mod tidy
     ```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！
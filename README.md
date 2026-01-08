# Aliyun Security Group Manager

阿里云安全组规则管理器 - 一个用于自动化管理阿里云 ECS 安全组规则的工具。

## 功能特性

- 📝 **文本配置管理**: 通过简单的配置文件管理安全组规则
- 🔄 **自动同步**: 实时监控配置文件变化并自动同步到阿里云
- ⏰ **规则过期管理**: 支持为规则设置过期时间，自动清理过期规则
- 🚀 **易于部署**: 支持 Worker 模式持续运行或 CLI 模式手动执行

## 项目结构

```
aliyun-security-group-mgr/
├── cmd/
│   ├── cli/          # CLI 命令行工具（待实现）
│   └── worker/       # Worker 后台服务
├── internal/
│   ├── conf/         # 配置管理
│   ├── ecs/          # 阿里云 ECS 接口封装
│   ├── reloader/     # 文件监控和规则解析
│   ├── service/      # 业务逻辑服务层
│   └── utils/        # 工具函数
└── sgmgr_rules.conf  # 规则配置文件示例
```

## 快速开始

### 前置要求

- Go 1.22.6 或更高版本
- 阿里云账号及 AccessKey
- 目标安全组的管理权限

### 安装

```bash
# 克隆项目
git clone <repository-url>
cd aliyun-security-group-mgr

# 编译 Worker
go build -o worker ./cmd/worker
```

### 配置

创建 `.env` 配置文件（或使用环境变量）：

```bash
# 阿里云凭证配置
ALIYUN_SGMGR_CREDENTIAL_TYPE=access_key
ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_ID=your_access_key_id
ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_SECRET=your_access_key_secret

# ECS 配置
ALIYUN_SGMGR_ECS_REGION_ID=cn-hangzhou
ALIYUN_SGMGR_ECS_ENDPOINT=ecs.aliyuncs.com

# 安全组配置
ALIYUN_SGMGR_SECURITY_GROUP_ID=sg-xxxxxxxxxxxxx

# 文件监控配置
ALIYUN_SGMGR_RELOADER_ENABLED=true
ALIYUN_SGMGR_RELOADER_INTERVAL=5
ALIYUN_SGMGR_RELOADER_WATCH_PATH=./sgmgr_rules.conf

# 调试模式（可选）
DEBUG=false
```

### 规则配置文件格式

在 `sgmgr_rules.conf` 文件中定义安全组规则，格式如下：

```
<policy> <direction> <protocol> <port_range> from <cidr_ip> priority <priority> until <expire_time> # <description>
```

第一次时，可以不创建该文件，Worker会自动从阿里云拉取现有规则并生成初始配置文件。

**参数说明**：
- `policy`: 授权策略，可选值 `accept` 或 `drop`
- `direction`: 方向，`ingress`（入方向）或 `egress`（出方向）
- `protocol`: 协议类型，如 `tcp`、`udp`、`icmp` 等
- `port_range`: 端口范围，格式 `起始端口/结束端口`，如 `80/80` 或 `1000/2000`
- `cidr_ip`: 授权的 IP 地址范围，如 `0.0.0.0/0` 或 `192.168.1.0/24`
- `priority`: 优先级，取值范围 1-100，数字越小优先级越高
- `expire_time`: 规则过期时间，RFC3339 格式，如 `2026-01-01T00:00:00Z`
- `description`: 规则描述（注释部分）

**示例**：

```conf
# 允许来自任何 IP 的 HTTP 访问
accept ingress tcp 80/80 from 0.0.0.0/0 priority 1 until 2100-01-01T00:00:00Z # HTTP Server

# 限制特定 IP 段访问 SSH
accept ingress tcp 22/22 from 192.168.1.0/24 priority 1 until 2025-12-31T23:59:59Z # SSH for office

# 允许特定端口范围
accept ingress tcp 8000/8100 from 10.0.0.0/8 priority 10 until 2100-01-01T00:00:00Z # Internal services
```

### 运行

```bash
# 使用默认配置文件 (.env)
./worker

# 指定配置文件
./worker -config /path/to/your/config.env
```

Worker 将会：
1. 加载配置并连接到阿里云
2. 读取规则配置文件
3. 同步规则到指定的安全组
4. 持续监控配置文件变化并自动更新

## 工作原理

1. **规则解析**: 读取并解析 `sgmgr_rules.conf` 配置文件
2. **规则比对**: 获取当前安全组的所有规则，与配置文件进行比对
3. **增量同步**: 
   - 添加配置文件中存在但安全组中不存在的规则
   - 删除安全组中存在但配置文件中不存在的规则
   - 删除已过期的规则
4. **文件监控**: 定期检查配置文件的修改时间，发现变化时自动重新同步

## 环境变量配置说明

| 环境变量 | 说明 | 必填 | 默认值 |
|---------|------|------|--------|
| `ALIYUN_SGMGR_CREDENTIAL_TYPE` | 凭证类型 | 是 | - |
| `ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_ID` | AccessKey ID | 是 | - |
| `ALIYUN_SGMGR_CREDENTIAL_ACCESS_KEY_SECRET` | AccessKey Secret | 是 | - |
| `ALIYUN_SGMGR_ECS_REGION_ID` | 地域 ID | 是 | - |
| `ALIYUN_SGMGR_ECS_ENDPOINT` | ECS API 端点 | 否 | ecs.aliyuncs.com |
| `ALIYUN_SGMGR_SECURITY_GROUP_ID` | 安全组 ID | 是 | - |
| `ALIYUN_SGMGR_RELOADER_ENABLED` | 是否启用自动重载 | 否 | true |
| `ALIYUN_SGMGR_RELOADER_INTERVAL` | 检查间隔（秒） | 否 | 60 |
| `ALIYUN_SGMGR_RELOADER_WATCH_PATH` | 监控的配置文件路径 | 是 | - |
| `ALIYUN_SGMGR_DEBUG` | 调试模式 | 否 | false |

## 开发

### 运行测试

```bash
go test ./...
```

### 构建

```bash
# 构建 Worker
go build -o worker ./cmd/worker

# 构建 CLI（待实现）
go build -o cli ./cmd/cli
```

## 注意事项

⚠️ **重要提示**：

1. **权限要求**: 确保使用的 AccessKey 具有对目标安全组的读写权限
2. **规则限制**: 阿里云安全组规则有数量限制，请注意不要超过配额
3. **过期时间**: 设置合理的过期时间，避免重要规则意外过期
4. **配置备份**: 建议定期备份 `sgmgr_rules.conf` 配置文件
5. **测试环境**: 首次使用建议先在测试环境的安全组上进行验证

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题或建议，请通过 Issue 反馈。

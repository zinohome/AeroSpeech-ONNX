# AeroSpeech-ONNX

基于 sherpa-onnx 的语音识别（STT）和文本转语音（TTS）服务器系统，使用 Go 语言开发。

## 特性

### 核心功能
- 🎤 **实时语音识别**：支持 WebSocket 流式识别和 REST API 文件识别
- 🔊 **文本转语音**：支持103个说话人、多语言、语速和音调调节（Kokoro v1.1）
- ⚡ **高性能**：支持 CPU/GPU 可配置，GPU 模式性能提升 3-5 倍
- 🌐 **多语言支持**：支持中文、英文（STT支持多语言，TTS支持中英文）
- 🎭 **丰富音色**：103个专业说话人（55个中文女声、45个中文男声、3个英语女声）
- 📊 **监控面板**：实时监控系统状态和性能指标
- 🐳 **容器化部署**：支持 Docker 和 Docker Compose 部署

### v2.0 新增功能 ⭐
- ✨ **限流中间件** - 保护系统免受过载攻击，支持IP级别限流
- ✨ **配置热重载** - 无需重启即可应用新配置（2秒生效）
- ✨ **Bootstrap启动** - 统一初始化流程和依赖注入
- ✨ **会话增强** - 错误计数、原子操作、增强统计
- ✨ **资源池优化** - 并行初始化（提升50%启动速度）、超时处理
- ✨ **VAD池管理** - 支持Voice Activity Detection资源池
- ✨ **增强统计** - 多维度性能监控（P95/P99延迟、QPS等）
- ✨ **WebSocket优化** - 更稳定的连接处理和错误恢复
- 📖 **STT配置优化指南** - 详细的语音识别优化方案（模型选择、参数调优、场景配置）
- 📖 **TTS配置优化指南** - 全面的语音合成优化方案（音色选择、语音参数、音质提升）
- 📖 **TTS说话人使用指南** - 103个说话人完整列表和使用方法 ⭐ 最新

📖 详细说明请查看：
- [新增功能文档](docs/12-新增功能说明.md)
- [STT配置优化指南](docs/13-STT配置优化指南.md)
- [TTS配置优化指南](docs/14-TTS配置优化指南.md)
- [TTS说话人使用指南](docs/15-TTS说话人使用指南.md) ⭐ 最新

## 技术栈

- **语言**: Go 1.21+
- **语音引擎**: sherpa-onnx (C API)
- **Go 绑定**: github.com/k2-fsa/sherpa-onnx-go
- **执行后端**: ONNX Runtime (支持 CPU/CUDA)
- **Web 框架**: Gin
- **WebSocket**: Gorilla WebSocket
- **配置管理**: Viper
- **日志**: Logrus / Zap

## 快速开始

### 前置要求

- Go 1.21+
- sherpa-onnx Go 绑定库
- ONNX Runtime (CPU 或 GPU 版本)
- 模型文件（ASR 和 TTS 模型）

### 安装依赖

```bash
go mod tidy
```

### 编译二进制文件

#### 方式1: 使用编译脚本（推荐）

```bash
# 编译所有服务到 bin/ 目录
./scripts/build.sh
```

编译后的二进制文件位于 `bin/` 目录：
- `bin/stt-server` - STT服务
- `bin/tts-server` - TTS服务
- `bin/speech-server` - 统一服务

#### 方式2: 手动编译

```bash
# 创建 bin 目录
mkdir -p bin

# 编译 STT 服务
go build -o bin/stt-server ./cmd/stt-server

# 编译 TTS 服务
go build -o bin/tts-server ./cmd/tts-server

# 编译统一服务
go build -o bin/speech-server ./cmd/speech-server
```

#### 运行编译后的二进制文件

```bash
# STT 服务
./bin/stt-server

# TTS 服务
./bin/tts-server

# 统一服务
./bin/speech-server
```

### 配置

#### 统一服务配置

```bash
cp configs/speech-config.example.json configs/speech-config.json
```

编辑 `configs/speech-config.json`，设置：
- `mode`: `"unified"` 或 `"separated"`
- `stt`: STT模型配置
- `tts`: TTS模型配置
- `server`: 服务器配置

#### 分离服务配置

```bash
cp configs/stt-config.example.json configs/stt-config.json
cp configs/tts-config.example.json configs/tts-config.json
```

编辑配置文件，设置模型路径和 Provider（CPU/CUDA/Auto）。

### 运行服务

#### 方式1: 统一服务（推荐用于开发/测试环境）

统一服务将STT和TTS合并到一个端口，简化部署：

```bash
# 复制统一配置文件
cp configs/speech-config.example.json configs/speech-config.json

# 编辑配置文件，设置模型路径和Provider

# 运行统一服务
go run cmd/speech-server/main.go
```

**统一服务特性**:
- 单一端口（默认8080）
- 同时支持STT和TTS
- 路由区分：`/api/v1/stt/*` 和 `/api/v1/tts/*`
- WebSocket：`/ws/stt` 和 `/ws/tts`

**访问地址**:
- **Web界面**:
  - STT 测试页面: http://localhost:8080/stt
  - TTS 测试页面: http://localhost:8080/tts
  - 监控面板: http://localhost:8080/monitor
- **REST API**:
  - STT API: http://localhost:8080/api/v1/stt/*
  - TTS API: http://localhost:8080/api/v1/tts/*
  - 健康检查: http://localhost:8080/api/v1/health
- **WebSocket**:
  - STT WebSocket: ws://localhost:8080/ws/stt
  - TTS WebSocket: ws://localhost:8080/ws/tts
  - 兼容路由: ws://localhost:8080/ws?type=stt 或 ws://localhost:8080/ws?type=tts

#### 方式2: 分离服务（推荐用于生产环境）

分离服务将STT和TTS分别部署，提供更好的故障隔离和扩展性：

```bash
# STT 服务
go run cmd/stt-server/main.go

# TTS 服务（另一个终端）
go run cmd/tts-server/main.go
```

**分离服务特性**:
- 独立端口（STT: 8080, TTS: 8081）
- 故障隔离
- 独立扩展
- 资源优化

**访问地址**:
- STT 测试页面: http://localhost:8080
- TTS 测试页面: http://localhost:8081
- 监控面板: http://localhost:8080/monitor

## 项目结构

```
AeroSpeech-ONNX/
├── cmd/                    # 服务入口
│   ├── speech-server/      # 统一服务入口（支持两种模式）
│   ├── stt-server/         # STT服务入口（分离模式）
│   └── tts-server/         # TTS服务入口（分离模式）
├── internal/               # 内部代码
│   ├── common/             # 共享代码
│   │   ├── config/         # 配置管理
│   │   ├── handlers/       # HTTP处理器
│   │   ├── logger/         # 日志
│   │   ├── router/         # 路由
│   │   ├── session/        # 会话管理
│   │   └── ws/             # WebSocket处理
│   ├── asr/                # ASR模块
│   └── tts/                # TTS模块
├── pkg/                    # 可复用包
│   └── utils/              # 工具函数
├── web/                    # Web前端
│   ├── static/             # 静态资源
│   └── templates/          # HTML模板
├── configs/                # 配置文件
├── docs/                   # 文档
└── test/                   # 测试文件
```

## API 文档

### Swagger 文档

项目已集成 Swagger API 文档，启动服务后可通过以下地址访问：

- **Swagger UI**: http://localhost:8080/swagger/index.html

Swagger 文档提供了：
- 完整的 REST API 接口列表（批量处理）
- 请求/响应参数说明
- 在线测试功能
- 交互式 API 文档

> **注意**: Swagger 文档主要展示 REST API（批量处理）。**流式处理**通过 WebSocket 接口实现，请参考 [WebSocket 接口文档](docs/03-websocket接口设计.md)

### 其他文档

详细的 API 文档请参考 [docs/04-API设计.md](docs/04-API设计.md)

## WebSocket 接口（流式处理）

系统支持通过 WebSocket 进行**流式处理**：

- **STT 流式识别**: `ws://localhost:8080/ws/stt` - 实时音频流识别
- **TTS 流式合成**: `ws://localhost:8080/ws/tts` - 实时文本转语音

详细的 WebSocket 接口文档请参考 [docs/03-websocket接口设计.md](docs/03-websocket接口设计.md)

### 接口对比

| 接口类型 | 用途 | 延迟 | 适用场景 |
|---------|------|------|---------|
| **REST API** | 批量处理 | 较高 | 文件上传、批量处理 |
| **WebSocket** | 流式处理 | 低 | 实时对话、流式播放 |

## 配置说明

### Provider 配置

- `"cpu"`: CPU 模式（默认）
- `"cuda"`: GPU 模式（需要 NVIDIA GPU）
- `"auto"`: 自动选择（优先 GPU，失败回退 CPU）

详细配置示例请参考 [docs/11-GPU配置示例.md](docs/11-GPU配置示例.md)

## 性能指标

### CPU 模式
- STT 延迟: 200-400ms（1秒音频）
- TTS 生成时间: 2-5秒（5秒音频）

### GPU 模式
- STT 延迟: 50-100ms（1秒音频）
- TTS 生成时间: 0.5-1秒（5秒音频）
- 性能提升: 3-5倍

## 开发

### 开发规范

请遵循 [docs/06-开发规范.md](docs/06-开发规范.md) 中的开发规范。

### 测试

```bash
# 运行所有测试
go test ./...

# 查看测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**测试覆盖率要求**: ≥80%

**当前测试覆盖率**:
- WebSocket模块: 82.4% ✅
- Handlers模块: 85.0% ✅
- 整体覆盖率: 持续提升中

## 部署

### Docker 部署

#### 方式1: 统一服务（推荐用于开发/测试环境）

```bash
# 使用 docker-compose 启动统一服务（默认）
docker-compose up -d speech-server

# 或者手动构建
docker build -f Dockerfile.speech -t aerospeech-unified .
docker run -d -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  -e SPEECH_CONFIG_PATH=/app/configs/speech-config.json \
  aerospeech-unified
```

#### 方式2: 分离服务（推荐用于生产环境）

```bash
# 使用 docker-compose 启动分离服务
docker-compose --profile separated up -d

# 或者手动构建
docker build -f Dockerfile.stt -t aerospeech-stt .
docker build -f Dockerfile.tts -t aerospeech-tts .

docker run -d -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  aerospeech-stt

docker run -d -p 8081:8081 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  aerospeech-tts
```

详细部署文档请参考 [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)

## 文档

完整的项目文档位于 `docs/` 目录：

### 核心文档
- [项目概述](docs/01-项目概述.md)
- [架构设计](docs/02-架构设计.md)
- [WebSocket接口设计](docs/03-websocket接口设计.md)
- [API设计](docs/04-API设计.md)
- [部署文档](docs/DEPLOYMENT.md)
- [统一服务使用说明](docs/统一服务使用说明.md)

### 开发文档
- [实施路线图](docs/05-实施路线图.md)
- [开发规范](docs/06-开发规范.md)
- [测试规范](docs/07-测试规范.md)
- [开发流程管控](docs/08-开发流程管控.md)
- [测试计划](docs/09-测试计划.md)

### 技术文档
- [Sherpa-ONNX技术分析](docs/10-sherpa-onnx技术分析.md)
- [GPU配置示例](docs/11-GPU配置示例.md)
- [新增功能说明](docs/12-新增功能说明.md) ⭐ v2.0新增
- [STT配置优化指南](docs/13-STT配置优化指南.md) ⭐ v2.0新增
- [架构设计分析-端口合并可行性](docs/架构设计分析-端口合并可行性.md)

### 项目状态
- [项目完整度分析](docs/项目完整度分析-最新.md)
- [更新日志](docs/CHANGELOG.md) ⭐ v2.0新增

## 许可证

Apache-2.0 License

## 参考项目

- [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx) - 核心语音引擎
- [achatbot-go](https://github.com/weedge/achatbot-go) - 多模态聊天机器人参考
- [asr_server](https://github.com/bbeyondllove/asr_server) - ASR 服务器参考


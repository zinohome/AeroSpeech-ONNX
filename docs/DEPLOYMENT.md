# 部署文档

## 1. 前置要求

### 1.1 系统要求
- Linux/macOS/Windows
- Go 1.21+ (开发环境)
- Docker 和 Docker Compose (容器化部署)
- sherpa-onnx 模型文件

### 1.2 模型文件
- ASR模型: SenseVoice/Paraformer等
- TTS模型: Kokoro/VITS等
- 模型文件需要放在 `models/` 目录下

## 2. 本地部署

### 2.1 安装依赖

```bash
go mod download
```

### 2.2 配置

复制配置文件模板：

```bash
cp configs/stt-config.example.json configs/stt-config.json
cp configs/tts-config.example.json configs/tts-config.json
```

编辑配置文件，设置模型路径和Provider。

### 2.3 运行服务

#### STT服务
```bash
go run cmd/stt-server/main.go
```

#### TTS服务
```bash
go run cmd/tts-server/main.go
```

## 3. Docker部署

### 3.1 统一服务部署（推荐用于开发/测试环境）

#### 3.1.1 使用Docker Compose（推荐）

```bash
# 启动统一服务（默认）
docker-compose up -d speech-server

# 查看日志
docker-compose logs -f speech-server

# 停止服务
docker-compose stop speech-server
```

#### 3.1.2 手动构建和运行

```bash
# 构建镜像
docker build -f Dockerfile.speech -t aerospeech-unified .

# 运行容器
docker run -d \
  --name aerospeech-unified \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  -e SPEECH_CONFIG_PATH=/app/configs/speech-config.json \
  aerospeech-unified

# 查看日志
docker logs -f aerospeech-unified
```

**统一服务特性**:
- 单一端口（8080）
- 同时支持STT和TTS
- 简化部署和管理

### 3.2 分离服务部署（推荐用于生产环境）

#### 3.2.1 使用Docker Compose（推荐）

```bash
# 启动分离服务（STT和TTS）
docker-compose --profile separated up -d

# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f stt-server
docker-compose logs -f tts-server

# 停止服务
docker-compose --profile separated down
```

#### 3.2.2 手动构建和运行

```bash
# 构建STT服务镜像
docker build -f Dockerfile.stt -t aerospeech-stt .

# 构建TTS服务镜像
docker build -f Dockerfile.tts -t aerospeech-tts .

# 运行STT服务
docker run -d \
  --name aerospeech-stt \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  -e STT_CONFIG_PATH=/app/configs/stt-config.json \
  aerospeech-stt

# 运行TTS服务
docker run -d \
  --name aerospeech-tts \
  -p 8081:8081 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/models:/app/models \
  -v $(pwd)/logs:/app/logs \
  -e TTS_CONFIG_PATH=/app/configs/tts-config.json \
  aerospeech-tts
```

**分离服务特性**:
- 独立端口（STT: 8080, TTS: 8081）
- 故障隔离
- 独立扩展
- 资源优化

### 3.3 查看日志

```bash
# Docker Compose方式
docker-compose logs -f [service-name]

# 手动运行方式
docker logs -f [container-name]
```

### 3.4 服务管理

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose stop

# 重启服务
docker-compose restart

# 删除服务
docker-compose down

# 查看服务状态
docker-compose ps
```

## 4. GPU部署

### 4.1 安装NVIDIA Container Toolkit

参考 [NVIDIA Container Toolkit文档](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

### 4.2 配置GPU支持

在配置文件中设置：
```json
{
  "provider": {
    "provider": "cuda",
    "device_id": 0,
    "num_threads": 1
  }
}
```

### 4.3 Docker Compose GPU配置

```yaml
services:
  stt-server:
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

## 5. 访问服务

### 统一服务模式
- STT测试页面: http://localhost:8080/stt
- TTS测试页面: http://localhost:8080/tts
- 监控面板: http://localhost:8080/monitor
- STT API: http://localhost:8080/api/v1/stt/*
- TTS API: http://localhost:8080/api/v1/tts/*
- STT WebSocket: ws://localhost:8080/ws/stt
- TTS WebSocket: ws://localhost:8080/ws/tts

### 分离服务模式
- STT服务: http://localhost:8080
- TTS服务: http://localhost:8081
- STT监控面板: http://localhost:8080/monitor
- TTS监控面板: http://localhost:8081/monitor

## 6. 健康检查

### 统一服务模式
```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# STT健康检查
curl http://localhost:8080/api/v1/stt/config

# TTS健康检查
curl http://localhost:8080/api/v1/tts/config
```

### 分离服务模式
```bash
# STT服务健康检查
curl http://localhost:8080/api/v1/health

# TTS服务健康检查
curl http://localhost:8081/api/v1/health
```

## 7. 配置说明

### 7.1 统一服务配置

创建统一配置文件：
```bash
cp configs/speech-config.example.json configs/speech-config.json
```

编辑 `configs/speech-config.json`，设置：
- `mode`: `"unified"` - 统一模式
- `stt`: STT模型配置（可选）
- `tts`: TTS模型配置（可选）
- `server`: 服务器配置

### 7.2 分离服务配置

创建分离服务配置文件：
```bash
cp configs/stt-config.example.json configs/stt-config.json
cp configs/tts-config.example.json configs/tts-config.json
```

编辑配置文件，设置模型路径和Provider（CPU/CUDA/Auto）。


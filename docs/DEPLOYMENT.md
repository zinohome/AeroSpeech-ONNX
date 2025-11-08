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

### 3.1 构建镜像

```bash
docker build -f Dockerfile.stt -t aerospeech-stt .
docker build -f Dockerfile.tts -t aerospeech-tts .
```

### 3.2 使用Docker Compose

```bash
docker-compose up -d
```

### 3.3 查看日志

```bash
docker-compose logs -f
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

- STT服务: http://localhost:8080
- TTS服务: http://localhost:8081
- 监控面板: http://localhost:8080/monitor

## 6. 健康检查

```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8081/api/v1/health
```


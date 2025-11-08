# GPU配置示例

## 1. 配置文件示例

### 1.1 STT服务配置（CPU模式）

```json
{
  "recognition": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "use_inverse_text_normalization": false,
    "num_threads": 16,
    "provider": "cpu",
    "debug": false
  }
}
```

### 1.2 STT服务配置（GPU模式）

```json
{
  "recognition": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "use_inverse_text_normalization": false,
    "num_threads": 1,
    "provider": "cuda",
    "device_id": 0,
    "debug": false
  }
}
```

### 1.3 STT服务配置（自动模式）

```json
{
  "recognition": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "use_inverse_text_normalization": false,
    "num_threads": 8,
    "provider": "auto",
    "device_id": 0,
    "debug": false
  }
}
```

### 1.4 TTS服务配置（CPU模式）

```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_0/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_0/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_0/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_0/espeak-ng-data",
    "num_threads": 4,
    "provider": "cpu",
    "debug": false
  }
}
```

### 1.5 TTS服务配置（GPU模式）

```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_0/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_0/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_0/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_0/espeak-ng-data",
    "num_threads": 1,
    "provider": "cuda",
    "device_id": 0,
    "debug": false
  }
}
```

## 2. 环境变量配置

### 2.1 CPU模式
```bash
export ASR_PROVIDER=cpu
export ASR_NUM_THREADS=16
export TTS_PROVIDER=cpu
export TTS_NUM_THREADS=4
```

### 2.2 GPU模式
```bash
export ASR_PROVIDER=cuda
export ASR_DEVICE_ID=0
export ASR_NUM_THREADS=1
export TTS_PROVIDER=cuda
export TTS_DEVICE_ID=0
export TTS_NUM_THREADS=1
```

### 2.3 自动模式
```bash
export ASR_PROVIDER=auto
export ASR_DEVICE_ID=0
export TTS_PROVIDER=auto
export TTS_DEVICE_ID=0
```

## 3. Go代码配置示例

### 3.1 ASR Provider配置

```go
// CPU配置
config := sherpa.OfflineRecognizerConfig{
    ModelConfig: sherpa.OfflineModelConfig{
        Provider:   "cpu",
        NumThreads: 16,
    },
}

// GPU配置
config := sherpa.OfflineRecognizerConfig{
    ModelConfig: sherpa.OfflineModelConfig{
        Provider:   "cuda",
        NumThreads: 1,
        DeviceID:   0,
    },
}

// 自动配置（优先GPU，失败回退CPU）
config := sherpa.OfflineRecognizerConfig{
    ModelConfig: sherpa.OfflineModelConfig{
        Provider:   "auto",
        NumThreads: 8,
        DeviceID:   0,
    },
}
```

### 3.2 TTS Provider配置

```go
// CPU配置
config := sherpa.OfflineTtsConfig{
    Model: sherpa.OfflineTtsModelConfig{
        Provider:   "cpu",
        NumThreads: 4,
    },
}

// GPU配置
config := sherpa.OfflineTtsConfig{
    Model: sherpa.OfflineTtsModelConfig{
        Provider:   "cuda",
        NumThreads: 1,
        DeviceID:   0,
    },
}
```

### 3.3 Provider自动选择

```go
func NewProvider(config ProviderConfig) (Provider, error) {
    switch config.Provider {
    case "auto":
        if isGPUAvailable() {
            config.Provider = "cuda"
            log.Info("GPU available, using CUDA provider")
        } else {
            config.Provider = "cpu"
            log.Info("GPU not available, using CPU provider")
        }
    case "cuda":
        if !isGPUAvailable() {
            log.Warn("GPU not available, falling back to CPU")
            config.Provider = "cpu"
        }
    }
    
    // 创建Provider
    return createProvider(config)
}
```

## 4. Docker配置示例

### 4.1 CPU Dockerfile

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o stt-server ./cmd/stt-server

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
    libonnxruntime.so \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/stt-server /usr/local/bin/
CMD ["stt-server"]
```

### 4.2 GPU Dockerfile

```dockerfile
FROM nvidia/cuda:11.8.0-runtime-ubuntu22.04 AS base
RUN apt-get update && apt-get install -y \
    libonnxruntime-gpu.so \
    && rm -rf /var/lib/apt/lists/*

FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o stt-server ./cmd/stt-server

FROM base
COPY --from=builder /app/stt-server /usr/local/bin/
CMD ["stt-server"]
```

### 4.3 docker-compose.yml（CPU）

```yaml
version: '3.8'
services:
  stt-server:
    build:
      context: .
      dockerfile: Dockerfile.cpu
    ports:
      - "8080:8080"
    volumes:
      - ./models:/app/models
      - ./configs/stt-config.json:/app/config.json
```

### 4.4 docker-compose.yml（GPU）

```yaml
version: '3.8'
services:
  stt-server:
    build:
      context: .
      dockerfile: Dockerfile.gpu
    ports:
      - "8080:8080"
    volumes:
      - ./models:/app/models
      - ./configs/stt-config.json:/app/config.json
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
    environment:
      - NVIDIA_VISIBLE_DEVICES=0
```

## 5. 混合部署配置

### 5.1 STT使用GPU，TTS使用CPU

```json
{
  "stt": {
    "recognition": {
      "provider": "cuda",
      "device_id": 0,
      "num_threads": 1
    }
  },
  "tts": {
    "tts": {
      "provider": "cpu",
      "num_threads": 4
    }
  }
}
```

### 5.2 多GPU配置

```json
{
  "stt": {
    "recognition": {
      "provider": "cuda",
      "device_id": 0,
      "num_threads": 1
    }
  },
  "tts": {
    "tts": {
      "provider": "cuda",
      "device_id": 1,
      "num_threads": 1
    }
  }
}
```

## 6. 性能调优建议

### 6.1 CPU模式调优
- **NumThreads**: 设置为CPU核心数（8-16）
- **批处理**: 使用批处理提高吞吐量
- **模型量化**: 使用int8量化模型

### 6.2 GPU模式调优
- **NumThreads**: 通常设置为1
- **批处理**: 批处理显著提高GPU利用率
- **显存管理**: 监控显存使用，避免溢出
- **多GPU**: 使用多GPU负载均衡

### 6.3 自动模式调优
- **检测机制**: 启动时检测GPU可用性
- **回退策略**: GPU不可用时自动回退CPU
- **日志记录**: 记录Provider选择过程

## 7. 故障排查

### 7.1 GPU不可用
```bash
# 检查NVIDIA驱动
nvidia-smi

# 检查CUDA版本
nvcc --version

# 检查ONNX Runtime GPU版本
ldd libonnxruntime.so | grep cuda
```

### 7.2 显存不足
- 减少批大小
- 使用量化模型
- 限制并发数
- 使用多GPU

### 7.3 性能不达标
- 确认使用GPU Provider
- 启用批处理
- 检查GPU利用率
- 使用量化模型


# API文档

## 1. STT API

### 1.1 文件上传识别

**POST** `/api/v1/stt/recognize`

**请求**: multipart/form-data
- `audio`: 音频文件

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "text": "识别结果",
    "timestamp": 1234567890
  }
}
```

### 1.2 批量识别

**POST** `/api/v1/stt/batch`

**请求**:
```json
{
  "files": ["/path/to/file1.wav", "/path/to/file2.wav"]
}
```

### 1.3 获取配置

**GET** `/api/v1/stt/config`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "sample_rate": 16000,
    "provider": "cpu",
    "gpu_available": false
  }
}
```

### 1.4 获取统计信息

**GET** `/api/v1/stt/stats`

## 2. TTS API

### 2.1 文本合成

**POST** `/api/v1/tts/synthesize`

**请求**:
```json
{
  "text": "要合成的文本",
  "speaker_id": 0,
  "speed": 1.0
}
```

**响应**: audio/wav 音频数据

### 2.2 批量合成

**POST** `/api/v1/tts/batch`

### 2.3 获取说话人列表

**GET** `/api/v1/tts/speakers`

### 2.4 获取配置

**GET** `/api/v1/tts/config`

## 3. 监控API

### 3.1 综合统计信息

**GET** `/api/v1/stats`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "timestamp": "2025-11-16T10:00:00Z",
    "asr": {
      "total_requests": 1000,
      "successful_requests": 980,
      "failed_requests": 20,
      "avg_latency_ms": 150.5,
      "p95_latency_ms": 300,
      "p99_latency_ms": 500,
      "requests_per_second": 10.5
    },
    "tts": {
      "total_requests": 800,
      "successful_requests": 790,
      "failed_requests": 10,
      "avg_latency_ms": 200.3,
      "p95_latency_ms": 400,
      "p99_latency_ms": 600,
      "requests_per_second": 8.2
    },
    "sessions": {
      "active": 42,
      "total": 150,
      "total_sessions": 5000,
      "active_sessions": 42,
      "total_messages": 150000
    },
    "resources": {
      "cpu_usage_percent": 45.2,
      "memory_usage_mb": 1024,
      "goroutine_count": 150,
      "pool_usage_percent": 65.5
    }
  }
}
```

### 3.2 限流器统计信息 ⭐ 新增

**GET** `/api/v1/rate-limit/stats`

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "enabled": true,
    "active_limiters": 15,
    "current_connections": 42,
    "max_connections": 2000,
    "requests_per_second": 1000,
    "burst_size": 2000
  }
}
```

### 3.3 实时监控数据

**GET** `/api/v1/monitor`

### 3.4 会话列表

**GET** `/api/v1/sessions`

### 3.5 会话详情

**GET** `/api/v1/sessions/{session_id}`

## 4. WebSocket接口

### 4.1 STT WebSocket

**连接**: `ws://host:8080/ws`

**消息格式**:
- 发送: 二进制音频数据 (PCM 16-bit)
- 接收: JSON格式识别结果

### 4.2 TTS WebSocket

**连接**: `ws://host:8081/ws`

**消息格式**:
- 发送: JSON格式合成请求
- 接收: 二进制音频数据 (PCM 16-bit)


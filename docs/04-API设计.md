# API 设计

## 1. 概述

本系统提供 RESTful API 和 WebSocket API 两种接口方式。本文档主要描述 RESTful API 设计。

### 1.1 基础信息
- **API版本**: v1
- **基础路径**: `/api/v1`
- **内容类型**: `application/json`
- **字符编码**: UTF-8

### 1.2 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { ... },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "error message",
  "error": {
    "type": "INVALID_PARAMS",
    "details": "具体错误信息"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 1.3 HTTP状态码
- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 无权限
- `404 Not Found`: 资源不存在
- `429 Too Many Requests`: 请求频率超限
- `500 Internal Server Error`: 服务器错误
- `503 Service Unavailable`: 服务不可用

## 2. STT API

### 2.1 健康检查

#### GET /health
检查服务健康状态

**请求**:
```http
GET /health HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "version": "1.0.0",
    "uptime_seconds": 3600,
    "components": {
      "asr": "ready",
      "vad": "ready"
    },
    "provider": {
      "asr": "cuda",
      "gpu_available": true,
      "gpu_device_id": 0
    }
  }
}
```

### 2.2 统计信息

#### GET /stats
获取服务统计信息

**请求**:
```http
GET /stats HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "timestamp": "2024-01-01T12:00:00Z",
    "asr": {
      "total_requests": 1000,
      "successful_requests": 980,
      "failed_requests": 20,
      "avg_latency_ms": 250,
      "p95_latency_ms": 400,
      "p99_latency_ms": 600,
      "requests_per_second": 10.5
    },
    "sessions": {
      "active": 5,
      "total": 100,
      "avg_duration_seconds": 30
    },
    "resources": {
      "cpu_usage_percent": 45.5,
      "memory_usage_mb": 512,
      "pool_usage_percent": 60.0
    }
  }
}
```

### 2.3 音频文件识别

#### POST /stt/recognize
上传音频文件进行识别

**请求**:
```http
POST /api/v1/stt/recognize HTTP/1.1
Host: host:8080
Content-Type: multipart/form-data

------WebKitFormBoundary
Content-Disposition: form-data; name="audio"; filename="test.wav"
Content-Type: audio/wav

[音频文件二进制数据]
------WebKitFormBoundary
Content-Disposition: form-data; name="language"

auto
------WebKitFormBoundary
Content-Disposition: form-data; name="punctuation"

true
------WebKitFormBoundary--
```

**参数**:
- `audio` (file, required): 音频文件（WAV/MP3/FLAC）
- `language` (string, optional): 语言代码（auto/zh/en/ja/ko）
- `punctuation` (boolean, optional): 是否添加标点（默认true）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "text": "识别结果文本",
    "confidence": 0.95,
    "language": "zh",
    "duration_ms": 1500,
    "processing_time_ms": 250
  }
}
```

### 2.4 批量识别

#### POST /stt/batch
批量识别多个音频文件

**请求**:
```http
POST /api/v1/stt/batch HTTP/1.1
Host: host:8080
Content-Type: multipart/form-data

------WebKitFormBoundary
Content-Disposition: form-data; name="files"; filename="audio1.wav"
Content-Type: audio/wav

[音频文件1]
------WebKitFormBoundary
Content-Disposition: form-data; name="files"; filename="audio2.wav"
Content-Type: audio/wav

[音频文件2]
------WebKitFormBoundary--
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "results": [
      {
        "file": "audio1.wav",
        "text": "识别结果1",
        "confidence": 0.95,
        "duration_ms": 1500
      },
      {
        "file": "audio2.wav",
        "text": "识别结果2",
        "confidence": 0.92,
        "duration_ms": 2000
      }
    ],
    "total": 2,
    "successful": 2,
    "failed": 0
  }
}
```

### 2.5 获取配置

#### GET /config
获取服务配置信息

**请求**:
```http
GET /config HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "audio": {
      "sample_rate": 16000,
      "chunk_size": 4096,
      "format": "pcm_s16le"
    },
    "recognition": {
      "model": "sense-voice",
      "language": "auto",
      "num_threads": 8,
      "provider": "cuda",
      "gpu_device_id": 0
    }
  }
}
```

## 3. TTS API

### 3.1 健康检查

同STT服务的健康检查接口，但检查TTS组件状态。

### 3.2 文本合成

#### POST /tts/synthesize
文本转语音合成

**请求**:
```http
POST /api/v1/tts/synthesize HTTP/1.1
Host: host:8081
Content-Type: application/json

{
  "text": "要合成的文本",
  "speaker_id": 0,
  "speed": 1.0,
  "pitch": 1.0,
  "format": "wav"
}
```

**参数**:
- `text` (string, required): 要合成的文本（最大1000字符）
- `speaker_id` (integer, optional): 说话人ID（默认0）
- `speed` (float, optional): 语速（0.5-2.0，默认1.0）
- `pitch` (float, optional): 音调（0.5-2.0，默认1.0）
- `format` (string, optional): 输出格式（wav/mp3，默认wav）

**响应**:
```http
HTTP/1.1 200 OK
Content-Type: audio/wav
Content-Length: 123456
X-Text-Length: 10
X-Duration-Ms: 5000
X-Processing-Time-Ms: 2000

[音频文件二进制数据]
```

**JSON响应（可选）**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "audio_url": "https://host:8081/api/v1/tts/audio/abc123",
    "duration_ms": 5000,
    "processing_time_ms": 2000,
    "text_length": 10
  }
}
```

### 3.3 批量合成

#### POST /tts/batch
批量文本合成

**请求**:
```http
POST /api/v1/tts/batch HTTP/1.1
Host: host:8081
Content-Type: application/json

{
  "texts": [
    "文本1",
    "文本2",
    "文本3"
  ],
  "options": {
    "speaker_id": 0,
    "speed": 1.0
  }
}
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "results": [
      {
        "index": 0,
        "text": "文本1",
        "audio_url": "https://host:8081/api/v1/tts/audio/abc123",
        "duration_ms": 2000
      },
      {
        "index": 1,
        "text": "文本2",
        "audio_url": "https://host:8081/api/v1/tts/audio/def456",
        "duration_ms": 2500
      }
    ],
    "total": 3,
    "successful": 3,
    "failed": 0
  }
}
```

### 3.4 获取可用说话人

#### GET /tts/speakers
获取可用的说话人列表

**请求**:
```http
GET /api/v1/tts/speakers HTTP/1.1
Host: host:8081
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "speakers": [
      {
        "id": 0,
        "name": "默认女声",
        "gender": "female",
        "language": "zh"
      },
      {
        "id": 1,
        "name": "默认男声",
        "gender": "male",
        "language": "zh"
      }
    ]
  }
}
```

### 3.5 获取音频文件

#### GET /tts/audio/{audio_id}
获取已生成的音频文件

**请求**:
```http
GET /api/v1/tts/audio/abc123 HTTP/1.1
Host: host:8081
```

**响应**:
```http
HTTP/1.1 200 OK
Content-Type: audio/wav
Content-Length: 123456
Cache-Control: public, max-age=3600

[音频文件二进制数据]
```

## 4. 监控API

### 4.1 实时监控

#### GET /monitor
获取实时监控数据

**请求**:
```http
GET /api/v1/monitor HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "timestamp": "2024-01-01T12:00:00Z",
    "metrics": {
      "active_connections": 10,
      "requests_per_second": 15.5,
      "avg_latency_ms": 250,
      "p95_latency_ms": 400,
      "error_rate": 0.02
    },
    "resources": {
      "cpu_usage_percent": 45.5,
      "memory_usage_mb": 512,
      "thread_count": 16,
      "goroutine_count": 50
    },
    "performance": {
      "asr_pool_usage_percent": 60.0,
      "tts_pool_usage_percent": 40.0,
      "queue_length": 5
    }
  }
}
```

### 4.2 历史统计

#### GET /monitor/history
获取历史统计数据

**请求**:
```http
GET /api/v1/monitor/history?duration=1h&interval=1m HTTP/1.1
Host: host:8080
```

**参数**:
- `duration` (string, optional): 时间范围（1h/6h/24h，默认1h）
- `interval` (string, optional): 时间间隔（1m/5m/15m，默认1m）

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "duration": "1h",
    "interval": "1m",
    "data_points": [
      {
        "timestamp": "2024-01-01T12:00:00Z",
        "requests_per_second": 10.5,
        "avg_latency_ms": 250,
        "error_rate": 0.02
      },
      {
        "timestamp": "2024-01-01T12:01:00Z",
        "requests_per_second": 12.3,
        "avg_latency_ms": 230,
        "error_rate": 0.01
      }
    ]
  }
}
```

## 5. 会话管理API

### 5.1 获取活跃会话

#### GET /sessions
获取当前活跃会话列表

**请求**:
```http
GET /api/v1/sessions HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "sessions": [
      {
        "session_id": "uuid-string",
        "type": "stt",
        "status": "active",
        "created_at": "2024-01-01T12:00:00Z",
        "duration_seconds": 30,
        "audio_bytes_received": 102400,
        "text_bytes_sent": 2048
      }
    ],
    "total": 5,
    "active": 3,
    "idle": 2
  }
}
```

### 5.2 获取会话详情

#### GET /sessions/{session_id}
获取指定会话的详细信息

**请求**:
```http
GET /api/v1/sessions/abc123 HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "session_id": "abc123",
    "type": "stt",
    "status": "active",
    "created_at": "2024-01-01T12:00:00Z",
    "last_activity": "2024-01-01T12:00:30Z",
    "duration_seconds": 30,
    "statistics": {
      "audio_bytes_received": 102400,
      "text_bytes_sent": 2048,
      "requests_count": 10,
      "avg_latency_ms": 250
    }
  }
}
```

### 5.3 关闭会话

#### DELETE /sessions/{session_id}
关闭指定会话

**请求**:
```http
DELETE /api/v1/sessions/abc123 HTTP/1.1
Host: host:8080
```

**响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "session_id": "abc123",
    "closed_at": "2024-01-01T12:00:35Z"
  }
}
```

## 6. 认证和授权

### 6.1 API Key认证

#### 请求头
```http
X-API-Key: your-api-key-here
```

### 6.2 Token认证（可选）

#### 请求头
```http
Authorization: Bearer your-token-here
```

### 6.3 认证错误响应
```json
{
  "code": 401,
  "message": "Unauthorized",
  "error": {
    "type": "AUTH_FAILED",
    "details": "Invalid API key"
  }
}
```

## 7. 限流

### 7.1 限流响应头
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1609459200
```

### 7.2 限流错误响应
```json
{
  "code": 429,
  "message": "Too Many Requests",
  "error": {
    "type": "RATE_LIMIT_EXCEEDED",
    "details": "Rate limit exceeded. Please retry after 60 seconds."
  },
  "retry_after": 60
}
```

## 8. 错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| SUCCESS | 200 | 请求成功 |
| INVALID_PARAMS | 400 | 参数错误 |
| AUTH_FAILED | 401 | 认证失败 |
| FORBIDDEN | 403 | 无权限 |
| NOT_FOUND | 404 | 资源不存在 |
| RATE_LIMIT_EXCEEDED | 429 | 请求频率超限 |
| INTERNAL_ERROR | 500 | 服务器内部错误 |
| SERVICE_UNAVAILABLE | 503 | 服务不可用 |
| AUDIO_FORMAT_ERROR | 400 | 音频格式错误 |
| MODEL_LOAD_ERROR | 500 | 模型加载失败 |
| RECOGNITION_ERROR | 500 | 识别错误 |
| SYNTHESIS_ERROR | 500 | 合成错误 |
| RESOURCE_EXHAUSTED | 503 | 资源耗尽 |

## 9. 版本控制

### 9.1 URL版本控制
```
/api/v1/...
/api/v2/...
```

### 9.2 版本兼容性
- 新版本保持向后兼容
- 废弃的API提前通知
- 至少保留一个旧版本

## 10. 测试工具

### 10.1 cURL示例

#### STT识别
```bash
curl -X POST http://host:8080/api/v1/stt/recognize \
  -F "audio=@test.wav" \
  -F "language=auto"
```

#### TTS合成
```bash
curl -X POST http://host:8081/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{
    "text": "测试文本",
    "speaker_id": 0,
    "speed": 1.0
  }' \
  --output output.wav
```

### 10.2 Postman集合
提供Postman API集合文件，包含所有API的测试用例。


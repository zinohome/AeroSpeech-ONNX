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

### 3.1 实时监控数据

**GET** `/api/v1/monitor`

### 3.2 会话列表

**GET** `/api/v1/sessions`

### 3.3 会话详情

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


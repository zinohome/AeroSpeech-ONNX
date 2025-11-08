# WebSocket 接口设计

## 1. 概述

WebSocket 接口用于实时双向通信，支持流式音频传输和实时结果返回。

### 1.1 连接地址
- **STT服务**: `ws://host:8080/ws` 或 `wss://host:8080/ws`
- **TTS服务**: `ws://host:8081/ws` 或 `wss://host:8081/ws`

### 1.2 协议版本
- WebSocket 协议版本: RFC 6455
- 子协议: 无（使用二进制和JSON消息）

## 2. STT WebSocket 接口

### 2.1 连接建立

#### 2.1.1 连接请求
```
GET /ws HTTP/1.1
Host: host:8080
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: <base64-encoded-key>
Sec-WebSocket-Version: 13
```

#### 2.1.2 连接响应
```http
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: <calculated-accept>
```

#### 2.1.3 连接确认消息
服务器连接成功后立即发送：
```json
{
  "type": "connection",
  "status": "connected",
  "session_id": "uuid-string",
  "message": "WebSocket connected, ready for audio",
  "config": {
    "sample_rate": 16000,
    "chunk_size": 4096,
    "format": "pcm_s16le",
    "provider": "cuda",
    "gpu_available": true,
    "gpu_device_id": 0
  }
}
```

### 2.2 消息格式

#### 2.2.1 客户端 → 服务器（音频数据）
**消息类型**: 二进制消息

**格式**: PCM 16-bit 小端序，单声道，16kHz采样率

**数据包大小**: 建议 4096 字节（约 128ms 音频）

**示例**:
```javascript
// JavaScript示例
const audioData = new Int16Array(2048); // 2048 samples = 4096 bytes
// ... 填充音频数据 ...
const buffer = new ArrayBuffer(audioData.length * 2);
const view = new DataView(buffer);
for (let i = 0; i < audioData.length; i++) {
    view.setInt16(i * 2, audioData[i], true); // little-endian
}
ws.send(buffer);
```

#### 2.2.2 服务器 → 客户端（识别结果）
**消息类型**: 文本消息（JSON）

**格式**:
```json
{
  "type": "partial",
  "text": "部分识别结果",
  "is_final": false,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

```json
{
  "type": "final",
  "text": "最终识别结果",
  "is_final": true,
  "confidence": 0.95,
  "timestamp": "2024-01-01T12:00:00Z",
  "duration_ms": 1500
}
```

#### 2.2.3 错误消息
```json
{
  "type": "error",
  "code": "ERROR_CODE",
  "message": "错误描述",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**错误代码**:
- `AUDIO_FORMAT_ERROR`: 音频格式错误
- `MODEL_LOAD_ERROR`: 模型加载失败
- `RECOGNITION_ERROR`: 识别过程错误
- `SESSION_TIMEOUT`: 会话超时
- `RATE_LIMIT_EXCEEDED`: 请求频率超限

### 2.3 控制消息

#### 2.3.1 开始识别
客户端发送（可选，用于显式控制）:
```json
{
  "type": "start",
  "config": {
    "language": "auto",
    "punctuation": true
  }
}
```

#### 2.3.2 结束识别
客户端发送:
```json
{
  "type": "end"
}
```

服务器响应:
```json
{
  "type": "end_ack",
  "final_text": "完整识别文本",
  "duration_ms": 5000
}
```

#### 2.3.3 重置会话
客户端发送:
```json
{
  "type": "reset"
}
```

服务器响应:
```json
{
  "type": "reset_ack"
}
```

### 2.4 心跳机制

#### 2.4.1 Ping消息
服务器定期发送（每30秒）:
```json
{
  "type": "ping",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 2.4.2 Pong响应
客户端响应:
```json
{
  "type": "pong",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 2.5 状态消息

#### 2.5.1 VAD状态
```json
{
  "type": "vad_status",
  "is_speech": true,
  "confidence": 0.85,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 3. TTS WebSocket 接口

### 3.1 连接建立

连接流程同STT，但配置不同：
```json
{
  "type": "connection",
  "status": "connected",
  "session_id": "uuid-string",
  "message": "WebSocket connected, ready for text",
  "config": {
    "sample_rate": 24000,
    "format": "pcm_s16le",
    "available_speakers": [0, 1, 2, 3],
    "provider": "cuda",
    "gpu_available": true,
    "gpu_device_id": 0
  }
}
```

### 3.2 消息格式

#### 3.2.1 客户端 → 服务器（文本数据）
**消息类型**: 文本消息（JSON）

**格式**:
```json
{
  "type": "synthesize",
  "text": "要合成的文本",
  "options": {
    "speaker_id": 0,
    "speed": 1.0,
    "pitch": 1.0
  }
}
```

#### 3.2.2 服务器 → 客户端（音频数据）
**消息类型**: 二进制消息

**格式**: PCM 16-bit 小端序，单声道，24kHz采样率

**分块发送**: 音频数据分块发送，每块约 4096 字节

**开始标记**:
```json
{
  "type": "audio_start",
  "text": "要合成的文本",
  "total_samples": 120000,
  "sample_rate": 24000
}
```

**音频数据**: 二进制消息（PCM格式）

**结束标记**:
```json
{
  "type": "audio_end",
  "duration_ms": 5000,
  "samples_sent": 120000
}
```

#### 3.2.3 错误消息
同STT错误消息格式

### 3.3 控制消息

#### 3.3.1 停止合成
```json
{
  "type": "stop"
}
```

#### 3.3.2 暂停/恢复
```json
{
  "type": "pause"
}
```

```json
{
  "type": "resume"
}
```

## 4. 会话管理

### 4.1 会话生命周期

```
连接建立 → 会话创建 → 活跃状态 → 空闲超时 → 会话关闭
```

### 4.2 会话超时

- **STT会话**: 60秒无音频数据自动关闭
- **TTS会话**: 30秒无文本数据自动关闭
- **心跳超时**: 90秒无心跳响应自动关闭

### 4.3 会话状态

```json
{
  "type": "session_status",
  "session_id": "uuid-string",
  "status": "active|idle|timeout",
  "duration_ms": 5000,
  "audio_bytes_received": 102400,
  "text_bytes_sent": 2048
}
```

## 5. 错误处理

### 5.1 连接错误

#### 5.1.1 连接失败
```json
{
  "type": "error",
  "code": "CONNECTION_FAILED",
  "message": "Failed to establish WebSocket connection"
}
```

#### 5.1.2 认证失败
```json
{
  "type": "error",
  "code": "AUTH_FAILED",
  "message": "Authentication failed"
}
```

### 5.2 业务错误

#### 5.2.1 参数错误
```json
{
  "type": "error",
  "code": "INVALID_PARAMS",
  "message": "Invalid parameters",
  "details": {
    "field": "text",
    "reason": "text is empty"
  }
}
```

#### 5.2.2 资源不足
```json
{
  "type": "error",
  "code": "RESOURCE_EXHAUSTED",
  "message": "No available resources, please retry later"
}
```

### 5.3 重连机制

客户端应实现自动重连：
1. 检测连接断开
2. 等待1秒后重连
3. 重连失败，指数退避（1s, 2s, 4s, 8s）
4. 最多重试5次

## 6. 性能优化

### 6.1 消息大小
- **音频数据**: 建议 4096 字节/包（约128ms）
- **文本数据**: 建议 ≤ 1024 字符/消息
- **控制消息**: 建议 ≤ 512 字节

### 6.2 缓冲策略
- **客户端**: 音频数据缓冲后批量发送
- **服务器**: 识别结果缓冲后批量返回

### 6.3 压缩
- **文本消息**: 可选的gzip压缩
- **二进制消息**: 不压缩（已为PCM格式）

## 7. 安全考虑

### 7.1 认证
```http
GET /ws?token=<auth-token> HTTP/1.1
```

### 7.2 限流
- **连接数限制**: 单IP最多10个连接
- **消息频率**: 最多100条消息/秒
- **数据大小**: 单条消息最大1MB

### 7.3 超时
- **连接超时**: 10秒
- **消息超时**: 30秒
- **会话超时**: 见4.2节

## 8. 客户端示例

### 8.1 JavaScript STT客户端

```javascript
class STTWebSocketClient {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.sessionId = null;
    }
    
    connect() {
        this.ws = new WebSocket(this.url);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected');
        };
        
        this.ws.onmessage = (event) => {
            if (typeof event.data === 'string') {
                const msg = JSON.parse(event.data);
                this.handleMessage(msg);
            }
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket closed');
            // 自动重连
            setTimeout(() => this.connect(), 1000);
        };
    }
    
    sendAudio(audioData) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(audioData);
        }
    }
    
    handleMessage(msg) {
        switch(msg.type) {
            case 'connection':
                this.sessionId = msg.session_id;
                break;
            case 'partial':
            case 'final':
                this.onResult(msg.text, msg.is_final);
                break;
            case 'error':
                this.onError(msg);
                break;
        }
    }
    
    onResult(text, isFinal) {
        // 处理识别结果
    }
    
    onError(error) {
        // 处理错误
    }
    
    close() {
        if (this.ws) {
            this.ws.close();
        }
    }
}
```

### 8.2 Go TTS客户端

```go
type TTSWebSocketClient struct {
    conn *websocket.Conn
    url  string
}

func (c *TTSWebSocketClient) Connect() error {
    conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
    if err != nil {
        return err
    }
    c.conn = conn
    
    go c.readMessages()
    return nil
}

func (c *TTSWebSocketClient) Synthesize(text string, opts TTSOptions) error {
    msg := map[string]interface{}{
        "type": "synthesize",
        "text": text,
        "options": opts,
    }
    return c.conn.WriteJSON(msg)
}

func (c *TTSWebSocketClient) readMessages() {
    for {
        var msg map[string]interface{}
        if err := c.conn.ReadJSON(&msg); err != nil {
            break
        }
        c.handleMessage(msg)
    }
}
```

## 9. 测试用例

### 9.1 连接测试
- 正常连接
- 认证失败
- 连接超时

### 9.2 消息测试
- 正常音频传输
- 音频格式错误
- 消息过大
- 消息频率超限

### 9.3 异常测试
- 网络断开
- 服务器重启
- 会话超时
- 资源耗尽

## 10. 监控指标

### 10.1 连接指标
- 活跃连接数
- 连接建立速率
- 连接断开速率
- 平均连接时长

### 10.2 消息指标
- 消息发送速率
- 消息接收速率
- 消息延迟
- 消息错误率

### 10.3 业务指标
- STT识别延迟
- TTS生成延迟
- 会话成功率
- 错误率


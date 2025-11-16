# STT配置优化指南

> 版本: v2.0  
> 更新时间: 2025-11-16  
> 目的: 提高语音识别准确率和性能

## 概述

本文档提供全面的STT（语音识别）配置优化方案，涵盖模型选择、参数调优、硬件加速等多个方面，帮助您在不同场景下获得最佳的识别效果。

---

## 1. 核心配置参数

### 1.1 当前ASR配置结构

```go
type ASRConfig struct {
    ModelPath  string         // 模型文件路径
    TokensPath string         // 词表文件路径
    Language   string         // 语言代码
    Provider   ProviderConfig // 计算Provider配置
    Debug      bool           // 调试模式
}
```

**建议扩展**：添加以下高级参数以提高控制精度

```go
type ASRConfig struct {
    ModelPath       string
    TokensPath      string
    Language        string
    Provider        ProviderConfig
    Debug           bool
    
    // 新增高级参数
    MaxActivePaths  int     // 最大活跃路径数（默认4）
    NumThreads      int     // 解码线程数
    DecodingMethod  string  // 解码方法：greedy/modified_beam_search
    EnableEndpoint  bool    // 启用端点检测
    Rule1MinTrailingS float32 // 端点检测规则1
    Rule2MinTrailingS float32 // 端点检测规则2
    Rule3MinUtteranceS float32 // 端点检测规则3
}
```

---

## 2. 模型选择与优化

### 2.1 推荐模型列表

#### 中文场景

**1. SenseVoice（推荐）⭐**
```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto"
  }
}
```
- **优点**: 多语言支持、高准确率、支持语言自动识别
- **适用**: 中英混合、多语言场景
- **识别率**: ★★★★★

**2. Paraformer（大模型）**
```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-paraformer-zh-2023-09-14/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-paraformer-zh-2023-09-14/tokens.txt",
    "language": "zh"
  }
}
```
- **优点**: 准确率极高、标点自动添加
- **缺点**: 模型较大、速度较慢
- **适用**: 对准确率要求高、实时性要求不高的场景
- **识别率**: ★★★★★

**3. Paraformer（小模型）**
```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/tokens.txt",
    "language": "zh"
  }
}
```
- **优点**: 速度快、资源占用小
- **适用**: 实时性要求高、硬件资源受限
- **识别率**: ★★★★☆

#### 英文场景

**Whisper模型**
```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-whisper-base.en/model.onnx",
    "tokens_path": "models/asr/sherpa-onnx-whisper-base.en/tokens.txt",
    "language": "en"
  }
}
```
- **优点**: 英文识别准确率高、抗噪性强
- **适用**: 纯英文场景
- **识别率**: ★★★★★

### 2.2 模型量化对比

| 量化方式 | 文件后缀 | 精度 | 速度 | 模型大小 | 推荐场景 |
|---------|---------|------|------|---------|---------|
| FP32 | .onnx | ★★★★★ | ★★☆☆☆ | 大 | GPU推理 |
| FP16 | .fp16.onnx | ★★★★☆ | ★★★☆☆ | 中 | GPU推理 |
| INT8 | .int8.onnx | ★★★★☆ | ★★★★★ | 小 | CPU推理（推荐）|

**建议**：
- CPU环境：使用INT8量化模型
- GPU环境：使用FP16或FP32模型
- 生产环境：INT8模型可满足大多数场景需求

---

## 3. Provider配置优化

### 3.1 CPU配置

#### 基础配置
```json
{
  "provider": {
    "provider": "cpu",
    "num_threads": 4
  }
}
```

#### 优化建议

**1. 线程数设置**
```python
# 推荐公式
num_threads = min(物理核心数, 8)

# 示例
# 4核CPU: num_threads = 4
# 8核CPU: num_threads = 8
# 16核CPU: num_threads = 8 (超过8个线程收益递减)
```

**2. 高性能配置**
```json
{
  "provider": {
    "provider": "cpu",
    "num_threads": 8
  },
  "audio": {
    "chunk_size": 4096
  }
}
```
- 适用于8核以上CPU
- chunk_size保持4096获得最佳吞吐量

**3. 低延迟配置**
```json
{
  "provider": {
    "provider": "cpu",
    "num_threads": 4
  },
  "audio": {
    "chunk_size": 2048
  }
}
```
- chunk_size减小可降低延迟
- 但会增加CPU使用率

### 3.2 GPU配置

#### CUDA配置
```json
{
  "provider": {
    "provider": "cuda",
    "device_id": 0,
    "num_threads": 2
  }
}
```

**优化要点**：
1. **设备选择**：多GPU环境使用`device_id`指定
2. **线程数**：GPU模式下建议2-4个线程
3. **批处理**：GPU更适合批量处理场景

**性能对比**：
```
场景：SenseVoice INT8模型
- CPU (8核): ~150ms/秒音频
- GPU (RTX 3060): ~30ms/秒音频
- 性能提升：约5倍
```

#### 自动选择
```json
{
  "provider": {
    "provider": "auto"
  }
}
```
- 自动检测并使用可用的最佳Provider
- 优先级：CUDA > CPU

---

## 4. 音频参数优化

### 4.1 采样率配置

```json
{
  "audio": {
    "sample_rate": 16000
  }
}
```

**采样率选择指南**：

| 采样率 | 音质 | 识别效果 | 推荐场景 |
|--------|------|---------|---------|
| 8000 Hz | 低 | 一般 | 电话场景 |
| 16000 Hz | 中 | 好 | **通用场景（推荐）** |
| 44100 Hz | 高 | 好 | 高音质录音 |
| 48000 Hz | 很高 | 好 | 专业录音 |

**注意事项**：
- 大多数ASR模型训练时使用16kHz
- 使用其他采样率会自动重采样
- 重采样可能影响识别准确率

### 4.2 音频块大小（Chunk Size）

```json
{
  "audio": {
    "chunk_size": 4096
  }
}
```

**chunk_size影响分析**：

| Chunk Size | 延迟 | 吞吐量 | CPU使用 | 推荐场景 |
|------------|------|--------|---------|---------|
| 1024 | 64ms | 低 | 高 | 极低延迟场景 |
| 2048 | 128ms | 中 | 中高 | 实时对话 |
| 4096 | 256ms | 高 | 中 | **平衡场景（推荐）** |
| 8192 | 512ms | 很高 | 低 | 批处理/离线识别 |

**计算公式**：
```
延迟(ms) = (chunk_size / sample_rate) * 1000
例如: (4096 / 16000) * 1000 = 256ms
```

### 4.3 特征维度

```json
{
  "audio": {
    "feature_dim": 80
  }
}
```

- **80维**：标准配置，适用于大多数模型
- **不建议修改**：除非模型文档明确要求

### 4.4 归一化因子

```json
{
  "audio": {
    "normalize_factor": 32768.0
  }
}
```

- **32768.0**：16位PCM音频标准值
- **作用**：将音频数据归一化到[-1, 1]区间
- **不建议修改**

---

## 5. VAD集成优化

### 5.1 VAD的作用

Voice Activity Detection（语音活动检测）可以：
- **过滤静音**：减少无效音频处理
- **智能分段**：在语音停顿处自动分段
- **提高准确率**：避免噪音干扰
- **降低成本**：只处理有效音频

### 5.2 VAD配置

```json
{
  "vad": {
    "enabled": true,
    "provider": "silero",
    "threshold": 0.5,
    "pool_size": 200
  }
}
```

**参数说明**：

**threshold（阈值）**：
- `0.3`：宽松，捕获更多音频（可能包含噪音）
- `0.5`：**平衡（推荐）**
- `0.7`：严格，只保留清晰语音（可能丢失轻声）

**不同场景推荐**：

```json
// 安静环境（办公室、录音棚）
{
  "vad": {
    "enabled": true,
    "threshold": 0.6
  }
}

// 一般环境（家庭、会议室）
{
  "vad": {
    "enabled": true,
    "threshold": 0.5
  }
}

// 嘈杂环境（街道、工厂）
{
  "vad": {
    "enabled": true,
    "threshold": 0.4
  }
}
```

### 5.3 VAD性能影响

**启用VAD后的提升**：
- 识别准确率：+5-10%
- 处理速度：+20-30%（过滤无效音频）
- CPU使用率：-15-20%

---

## 6. 场景化配置方案

### 6.1 实时对话场景

**特点**：低延迟、高并发

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 2048,
    "feature_dim": 80,
    "normalize_factor": 32768.0
  },
  "vad": {
    "enabled": true,
    "threshold": 0.5
  },
  "websocket": {
    "read_timeout": 20,
    "max_message_size": 1048576
  }
}
```

**关键点**：
- chunk_size=2048，延迟约128ms
- 启用VAD过滤静音
- INT8模型平衡速度和准确率

### 6.2 会议记录场景

**特点**：高准确率、支持标点、多说话人

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-paraformer-zh-2023-09-14/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-paraformer-zh-2023-09-14/tokens.txt",
    "language": "zh",
    "provider": {
      "provider": "cuda",
      "device_id": 0
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 8192,
    "feature_dim": 80,
    "normalize_factor": 32768.0
  },
  "vad": {
    "enabled": true,
    "threshold": 0.4
  }
}
```

**关键点**：
- 使用Paraformer大模型，准确率最高
- GPU加速提升处理速度
- chunk_size较大，优先准确率
- VAD阈值适当降低，避免漏掉轻声

### 6.3 客服质检场景

**特点**：批量处理、高准确率、支持情感分析

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "provider": {
      "provider": "cuda",
      "device_id": 0,
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 8192,
    "feature_dim": 80,
    "normalize_factor": 32768.0
  },
  "vad": {
    "enabled": true,
    "threshold": 0.5
  }
}
```

**关键点**：
- 使用FP32模型（非INT8）获得最高准确率
- GPU批量处理提高吞吐量
- SenseVoice支持情感识别

### 6.4 移动端/边缘设备场景

**特点**：资源受限、低功耗

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/tokens.txt",
    "language": "zh",
    "provider": {
      "provider": "cpu",
      "num_threads": 2
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 4096,
    "feature_dim": 80,
    "normalize_factor": 32768.0
  },
  "vad": {
    "enabled": true,
    "threshold": 0.6
  }
}
```

**关键点**：
- 使用小模型减少内存占用
- INT8量化降低计算量
- 线程数=2，降低功耗
- VAD阈值适当提高，减少误识别

### 6.5 电话语音识别

**特点**：8kHz采样率、信道噪音

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-telecom-model/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-telecom-model/tokens.txt",
    "language": "zh",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 8000,
    "chunk_size": 2048,
    "feature_dim": 80,
    "normalize_factor": 32768.0
  },
  "vad": {
    "enabled": true,
    "threshold": 0.4
  }
}
```

**关键点**：
- 采样率=8000Hz（电话标准）
- 使用专门针对电话场景训练的模型
- VAD阈值降低，应对信道噪音

---

## 7. 高级优化技巧

### 7.1 端点检测优化

端点检测（Endpoint Detection）可以自动判断语音结束点，实现自动分段。

**推荐配置**（需扩展配置结构）：
```json
{
  "stt": {
    "enable_endpoint": true,
    "rule1_min_trailing_silence": 2.4,
    "rule2_min_trailing_silence": 1.2,
    "rule3_min_utterance_length": 0.0
  }
}
```

**参数说明**：
- `rule1_min_trailing_silence`：检测到静音的最小持续时间（秒）
- `rule2_min_trailing_silence`：第二条规则的静音时间
- `rule3_min_utterance_length`：最小语音长度

### 7.2 解码方法选择

**Greedy Search（贪心搜索）**：
```json
{
  "stt": {
    "decoding_method": "greedy_search"
  }
}
```
- 速度快
- 适用于实时场景

**Modified Beam Search（改进束搜索）**：
```json
{
  "stt": {
    "decoding_method": "modified_beam_search",
    "max_active_paths": 4
  }
}
```
- 准确率高5-10%
- 速度稍慢
- 适用于对准确率要求高的场景

### 7.3 多路径解码

```json
{
  "stt": {
    "max_active_paths": 4
  }
}
```

**max_active_paths影响**：

| 路径数 | 准确率 | 速度 | 内存占用 |
|--------|--------|------|---------|
| 1 | 基准 | 最快 | 最小 |
| 4 | +3-5% | 中等 | **平衡（推荐）** |
| 8 | +5-7% | 慢 | 大 |
| 16 | +6-8% | 很慢 | 很大 |

**建议**：
- 实时场景：1-4路径
- 离线场景：4-8路径
- 高精度需求：8-16路径

---

## 8. 性能监控与调优

### 8.1 关键指标

**识别准确率（WER - Word Error Rate）**：
```
WER = (替换词数 + 删除词数 + 插入词数) / 总词数
```
- 目标：WER < 5%（优秀）
- 可接受：WER < 10%

**实时率（RTF - Real Time Factor）**：
```
RTF = 处理时间 / 音频时长
```
- RTF < 1：实时处理
- RTF = 0.1：10秒音频1秒处理完
- 目标：RTF < 0.3

**延迟**：
```
总延迟 = 网络延迟 + chunk延迟 + 模型推理延迟
```
- 实时对话：< 500ms
- 一般场景：< 1000ms

### 8.2 监控API

```bash
# 获取STT统计信息
curl http://localhost:8080/api/v1/stt/stats

# 响应示例
{
  "total_requests": 1000,
  "successful_requests": 985,
  "failed_requests": 15,
  "avg_latency_ms": 200,
  "p95_latency_ms": 450,
  "p99_latency_ms": 800,
  "avg_rtf": 0.25,
  "avg_audio_length_s": 5.2
}
```

### 8.3 性能优化步骤

**步骤1：建立基准**
```bash
# 使用标准配置测试
# 记录：WER、RTF、延迟
```

**步骤2：逐项优化**
```bash
# 2.1 调整chunk_size
# 测试：1024, 2048, 4096, 8192
# 记录每种配置的指标

# 2.2 调整num_threads
# 测试：2, 4, 8
# 记录CPU使用率和RTF

# 2.3 启用VAD
# 对比启用前后的准确率提升
```

**步骤3：A/B测试**
```bash
# 对比不同模型的效果
# SenseVoice vs Paraformer
# INT8 vs FP16
```

---

## 9. 常见问题与解决方案

### 9.1 识别准确率低

**症状**：识别结果错误率高

**可能原因与解决方案**：

1. **音频质量差**
   ```json
   // 解决：提高音频质量要求
   {
     "audio": {
       "sample_rate": 16000,  // 确保使用16kHz
       "normalize_factor": 32768.0  // 正确归一化
     }
   }
   ```

2. **模型不匹配**
   ```bash
   # 解决：更换适合场景的模型
   # 中英混合 -> SenseVoice
   # 纯中文高准确率 -> Paraformer大模型
   ```

3. **VAD过滤过度**
   ```json
   // 解决：降低VAD阈值
   {
     "vad": {
       "threshold": 0.4  // 从0.6降到0.4
     }
   }
   ```

4. **音频分段不当**
   ```json
   // 解决：调整chunk_size
   {
     "audio": {
       "chunk_size": 4096  // 增大chunk提高上下文
     }
   }
   ```

### 9.2 实时性差（延迟高）

**症状**：用户说完话后很久才得到结果

**解决方案**：

1. **减小chunk_size**
   ```json
   {
     "audio": {
       "chunk_size": 2048  // 从4096减到2048
     }
   }
   ```

2. **使用GPU加速**
   ```json
   {
     "provider": {
       "provider": "cuda",
       "device_id": 0
     }
   }
   ```

3. **使用小模型**
   ```bash
   # 从大模型切换到小模型
   paraformer-zh-2023-09-14 -> paraformer-zh-small-2024-03-09
   ```

4. **增加处理线程**
   ```json
   {
     "provider": {
       "num_threads": 8  // 从4增到8
     }
   }
   ```

### 9.3 CPU/GPU使用率过高

**解决方案**：

1. **使用INT8量化模型**
   ```bash
   model.onnx -> model.int8.onnx
   ```

2. **减少并发请求**
   ```json
   {
     "rate_limit": {
       "enabled": true,
       "max_connections": 100
     }
   }
   ```

3. **启用VAD**
   ```json
   {
     "vad": {
       "enabled": true
     }
   }
   ```

### 9.4 内存占用过大

**解决方案**：

1. **减小资源池大小**
   ```json
   // 在初始化代码中
   poolSize := 2  // 从4减到2
   ```

2. **使用小模型**
3. **使用INT8量化模型**

---

## 10. 配置模板

### 10.1 生产环境推荐配置

```json
{
  "mode": "unified",
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": 20
  },
  "stt": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "provider": {
      "provider": "auto",
      "device_id": 0,
      "num_threads": 4
    },
    "debug": false
  },
  "audio": {
    "sample_rate": 16000,
    "feature_dim": 80,
    "chunk_size": 4096,
    "normalize_factor": 32768.0
  },
  "websocket": {
    "read_timeout": 20,
    "max_message_size": 2097152,
    "read_buffer_size": 2048,
    "write_buffer_size": 2048,
    "enable_compression": false
  },
  "session": {
    "send_queue_size": 500,
    "max_send_errors": 10
  },
  "rate_limit": {
    "enabled": true,
    "requests_per_second": 1000,
    "burst_size": 2000,
    "max_connections": 2000
  },
  "vad": {
    "enabled": true,
    "provider": "silero",
    "pool_size": 200,
    "threshold": 0.5
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "both",
    "file_path": "logs/speech.log",
    "max_size": 100,
    "max_backups": 10,
    "max_age": 30,
    "compress": true
  }
}
```

### 10.2 高性能配置（GPU）

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.fp16.onnx",
    "tokens_path": "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt",
    "language": "auto",
    "provider": {
      "provider": "cuda",
      "device_id": 0,
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 8192
  },
  "vad": {
    "enabled": true,
    "threshold": 0.5
  }
}
```

### 10.3 低延迟配置

```json
{
  "stt": {
    "model_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/model.int8.onnx",
    "tokens_path": "models/asr/sherpa-onnx-paraformer-zh-small-2024-03-09/tokens.txt",
    "language": "zh",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 1024
  },
  "vad": {
    "enabled": true,
    "threshold": 0.6
  }
}
```

---

## 11. 扩展配置建议

### 11.1 建议添加的配置参数

为了更精细的控制，建议扩展`ASRConfig`：

```go
type ASRConfig struct {
    ModelPath  string         `mapstructure:"model_path" json:"model_path"`
    TokensPath string         `mapstructure:"tokens_path" json:"tokens_path"`
    Language   string         `mapstructure:"language" json:"language"`
    Provider   ProviderConfig `mapstructure:"provider" json:"provider"`
    Debug      bool           `mapstructure:"debug" json:"debug"`
    
    // 新增参数
    MaxActivePaths     int     `mapstructure:"max_active_paths" json:"max_active_paths"`
    DecodingMethod     string  `mapstructure:"decoding_method" json:"decoding_method"`
    EnableEndpoint     bool    `mapstructure:"enable_endpoint" json:"enable_endpoint"`
    Rule1MinTrailing   float32 `mapstructure:"rule1_min_trailing" json:"rule1_min_trailing"`
    Rule2MinTrailing   float32 `mapstructure:"rule2_min_trailing" json:"rule2_min_trailing"`
    Rule3MinUtterance  float32 `mapstructure:"rule3_min_utterance" json:"rule3_min_utterance"`
    BlankPenalty       float32 `mapstructure:"blank_penalty" json:"blank_penalty"`
    TemperatureScale   float32 `mapstructure:"temperature_scale" json:"temperature_scale"`
}
```

### 11.2 实现示例

在`internal/asr/provider.go`中使用这些参数：

```go
func NewASRProvider(cfg *config.ASRConfig) (Provider, error) {
    // ... 现有代码 ...
    
    // 设置高级参数
    if cfg.MaxActivePaths > 0 {
        recognizer.SetMaxActivePaths(cfg.MaxActivePaths)
    }
    
    if cfg.DecodingMethod != "" {
        recognizer.SetDecodingMethod(cfg.DecodingMethod)
    }
    
    if cfg.EnableEndpoint {
        recognizer.EnableEndpoint(
            cfg.Rule1MinTrailing,
            cfg.Rule2MinTrailing,
            cfg.Rule3MinUtterance,
        )
    }
    
    return &asrProvider{
        recognizer: recognizer,
        config:     cfg,
    }, nil
}
```

---

## 12. 测试与验证

### 12.1 准确率测试

**步骤1：准备测试集**
```bash
# 准备100条测试音频及对应文本
test_data/
├── audio_001.wav
├── audio_001.txt
├── audio_002.wav
├── audio_002.txt
...
```

**步骤2：批量测试**
```bash
#!/bin/bash
for audio in test_data/*.wav; do
    base=$(basename $audio .wav)
    # 调用API识别
    curl -X POST \
      -F "audio=@$audio" \
      http://localhost:8080/api/v1/stt/recognize \
      | jq -r '.data.text' > results/${base}.txt
done
```

**步骤3：计算WER**
```python
def calculate_wer(reference, hypothesis):
    # 计算编辑距离
    # 返回WER
    pass

total_wer = 0
for i in range(1, 101):
    ref = open(f'test_data/audio_{i:03d}.txt').read()
    hyp = open(f'results/audio_{i:03d}.txt').read()
    wer = calculate_wer(ref, hyp)
    total_wer += wer

avg_wer = total_wer / 100
print(f"平均WER: {avg_wer:.2%}")
```

### 12.2 性能测试

```bash
# 使用Apache Bench测试并发性能
ab -n 1000 -c 10 \
  -p audio.wav \
  -T "multipart/form-data" \
  http://localhost:8080/api/v1/stt/recognize
```

---

## 13. 总结与最佳实践

### 关键要点

1. **模型选择**
   - 通用场景：SenseVoice INT8
   - 高准确率：Paraformer大模型
   - 低延迟：Paraformer小模型

2. **硬件配置**
   - CPU：4-8线程
   - GPU：启用CUDA，FP16模型

3. **音频参数**
   - 采样率：16kHz
   - chunk_size：2048（低延迟）或4096（平衡）

4. **VAD设置**
   - 启用VAD
   - 阈值：0.4-0.6根据环境调整

5. **监控优化**
   - 定期检查WER和RTF
   - 根据指标调整配置

### 优化优先级

1. ⭐⭐⭐ 选择合适的模型
2. ⭐⭐⭐ 启用VAD
3. ⭐⭐⭐ 确保音频质量（16kHz、无噪音）
4. ⭐⭐ 调整chunk_size平衡延迟和准确率
5. ⭐⭐ 使用GPU加速（如可用）
6. ⭐ 调整高级参数（max_active_paths等）

---

## 参考资料

- [Sherpa-ONNX官方文档](https://k2-fsa.github.io/sherpa/onnx/index.html)
- [SenseVoice模型文档](https://github.com/FunAudioLLM/SenseVoice)
- [Paraformer模型文档](https://github.com/alibaba-damo-academy/FunASR)
- [项目架构设计](./02-架构设计.md)
- [性能优化指南](./10-sherpa-onnx技术分析.md)

---

## 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v1.0 | 2025-11-16 | 初始版本，覆盖基础优化配置 |


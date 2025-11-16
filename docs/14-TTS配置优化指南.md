# TTS配置优化指南

> 版本: v2.0  
> 更新时间: 2025-11-16  
> 目的: 提高语音合成质量和性能

## 概述

本文档提供全面的TTS（文本转语音）配置优化方案，涵盖模型选择、音色调整、性能优化等多个方面，帮助您在不同场景下生成高质量的语音。

---

## 1. 核心配置参数

### 1.1 当前TTS配置结构

```go
type TTSModelConfig struct {
    ModelPath  string         // 模型文件路径
    VoicesPath string         // 音色文件路径
    TokensPath string         // 词表文件路径
    DataDir    string         // espeak-ng数据目录
    DictDir    string         // 字典目录
    Lexicon    string         // 词典文件（逗号分隔）
    Provider   ProviderConfig // 计算Provider配置
    Debug      bool           // 调试模式
}
```

**建议扩展**：添加以下高级参数以提高控制精度

```go
type TTSModelConfig struct {
    ModelPath  string
    VoicesPath string
    TokensPath string
    DataDir    string
    DictDir    string
    Lexicon    string
    Provider   ProviderConfig
    Debug      bool
    
    // 新增高级参数
    Speed          float32 // 语速倍率（0.5-2.0）
    Volume         float32 // 音量（0.0-1.0）
    Pitch          float32 // 音调（-12到+12半音）
    DefaultSpeaker int     // 默认说话人ID
    SampleRate     int     // 输出采样率（默认16000）
    EnableCache    bool    // 启用合成缓存
    CacheSize      int     // 缓存大小（MB）
}
```

---

## 2. 模型选择与优化

### 2.1 推荐TTS模型

#### Kokoro（当前使用）⭐

**基础信息**：
```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data"
  }
}
```

**优点**：
- ✅ 多语言支持（中英日韩等）
- ✅ 多音色（支持多个说话人）
- ✅ 音质自然
- ✅ 合成速度快
- ✅ 模型较小，资源占用少

**适用场景**：
- 通用TTS应用
- 多语言内容生成
- 实时语音合成
- 移动端应用

#### VITS模型

```json
{
  "tts": {
    "model_path": "models/tts/vits-zh-hf-models/model.onnx",
    "voices_path": "models/tts/vits-zh-hf-models/voices.bin",
    "tokens_path": "models/tts/vits-zh-hf-models/tokens.txt"
  }
}
```

**优点**：
- ✅ 音质极佳
- ✅ 情感丰富
- ✅ 韵律自然
- ❌ 合成速度较慢
- ❌ 模型较大

**适用场景**：
- 有声书制作
- 高质量配音
- 广告语音
- 对音质要求极高的场景

#### Piper TTS

```json
{
  "tts": {
    "model_path": "models/tts/piper-zh-cn/model.onnx",
    "voices_path": "models/tts/piper-zh-cn/voices.bin",
    "tokens_path": "models/tts/piper-zh-cn/tokens.txt"
  }
}
```

**优点**：
- ✅ 合成速度极快
- ✅ 资源占用极小
- ✅ 支持实时流式合成
- ❌ 音质一般
- ❌ 音色选择少

**适用场景**：
- 边缘设备
- 实时语音导航
- 批量合成
- 资源受限环境

### 2.2 模型对比

| 模型 | 音质 | 速度 | 资源占用 | 多语言 | 推荐指数 |
|------|------|------|---------|--------|---------|
| Kokoro | ★★★★☆ | ★★★★☆ | 中 | ✅ | ★★★★★ |
| VITS | ★★★★★ | ★★★☆☆ | 高 | ✅ | ★★★★☆ |
| Piper | ★★★☆☆ | ★★★★★ | 低 | ✅ | ★★★☆☆ |

---

## 3. 音色（Speaker）优化

### 3.1 说话人选择

Kokoro模型通常包含多个说话人，每个说话人有不同的音色特点。

**查看可用说话人**：
```bash
# 通过API查询
curl http://localhost:8080/api/v1/tts/speakers

# 响应示例
{
  "code": 200,
  "data": {
    "speakers": [
      {"id": 0, "name": "speaker_0", "gender": "female", "language": "zh"},
      {"id": 1, "name": "speaker_1", "gender": "male", "language": "zh"},
      {"id": 2, "name": "speaker_2", "gender": "female", "language": "en"},
      {"id": 3, "name": "speaker_3", "gender": "male", "language": "en"}
    ]
  }
}
```

### 3.2 音色选择建议

**场景化选择**：

| 场景 | 推荐音色 | 特点 |
|------|---------|------|
| **客服系统** | 女声（温柔型） | 亲切、专业 |
| **新闻播报** | 男声（稳重型） | 权威、可信 |
| **儿童内容** | 女声（活泼型） | 甜美、活泼 |
| **有声书** | 根据内容选择 | 情感丰富 |
| **导航系统** | 女声（清晰型） | 清晰、易懂 |
| **广告配音** | 根据品牌选择 | 吸引力强 |

**API调用示例**：
```bash
# 使用指定说话人
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{
    "text": "你好，欢迎使用语音合成服务",
    "speaker_id": 0,
    "speed": 1.0
  }' \
  --output output.wav
```

---

## 4. 语音参数调优

### 4.1 语速控制（Speed）

**语速倍率范围**：0.5 - 2.0

```json
{
  "speed": 1.0  // 正常语速
}
```

**建议值**：

| 场景 | 语速 | 说明 |
|------|------|------|
| **新闻播报** | 0.9 - 1.0 | 稳重、清晰 |
| **有声书** | 0.8 - 0.9 | 舒缓、易听 |
| **广告语** | 1.0 - 1.2 | 节奏感强 |
| **儿童内容** | 0.7 - 0.8 | 缓慢、清楚 |
| **快速提醒** | 1.2 - 1.5 | 高效传达 |
| **客服** | 0.9 - 1.0 | 温和适中 |

**API示例**：
```bash
# 慢速合成（有声书）
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text": "这是一段测试文本", "speed": 0.8}'

# 快速合成（提醒）
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text": "请注意", "speed": 1.3}'
```

### 4.2 音调控制（Pitch）

**音调范围**：-12 到 +12 半音

```json
{
  "pitch": 0  // 原始音调
}
```

**调整建议**：
- `+2 到 +4`：女声更甜美，适合儿童内容
- `-2 到 -4`：男声更深沉，适合广告配音
- `0`：自然音调（推荐）

**注意**：过大的音调调整可能导致不自然

### 4.3 音量控制（Volume）

**音量范围**：0.0 - 1.0

```json
{
  "volume": 1.0  // 最大音量
}
```

**建议值**：
- `0.8 - 1.0`：通用场景
- `0.6 - 0.8`：背景语音
- `1.0`：前景语音、提醒音

### 4.4 采样率优化

**输出采样率配置**：

```json
{
  "audio": {
    "sample_rate": 16000
  }
}
```

**采样率选择**：

| 采样率 | 音质 | 文件大小 | 适用场景 |
|--------|------|---------|---------|
| 8000 Hz | 低 | 最小 | 电话系统 |
| 16000 Hz | 中 | 小 | **通用场景（推荐）** |
| 22050 Hz | 良好 | 中 | 音乐播放器 |
| 44100 Hz | 高 | 大 | 高品质音频 |
| 48000 Hz | 极高 | 很大 | 专业制作 |

**建议**：
- 移动应用：16kHz
- Web应用：22.05kHz
- 音频制作：44.1kHz

---

## 5. Provider配置优化

### 5.1 CPU配置

#### 基础配置
```json
{
  "provider": {
    "provider": "cpu",
    "num_threads": 4
  }
}
```

#### 线程数优化

**TTS线程数建议**：
```python
# TTS通常比ASR需要更多计算
num_threads = min(物理核心数, 6)

# 示例
# 4核CPU: num_threads = 4
# 8核CPU: num_threads = 6
# 16核CPU: num_threads = 6
```

**性能测试**：
```bash
# 测试不同线程数的性能
for threads in 2 4 6 8; do
  echo "Testing with $threads threads"
  time curl -X POST \
    -d '{"text":"测试文本","num_threads":'$threads'}' \
    http://localhost:8080/api/v1/tts/synthesize \
    -o /dev/null
done
```

### 5.2 GPU配置

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

**GPU TTS性能提升**：
```
场景：Kokoro模型，100字文本
- CPU (8核): ~500ms
- GPU (RTX 3060): ~150ms
- 性能提升：约3.3倍
```

**GPU优化要点**：
1. 批量合成效果更好
2. 长文本加速明显
3. 多并发请求下优势显著

### 5.3 批量合成优化

**批量API**：
```bash
curl -X POST http://localhost:8080/api/v1/tts/batch \
  -H "Content-Type: application/json" \
  -d '{
    "texts": [
      "第一段文本",
      "第二段文本",
      "第三段文本"
    ],
    "speaker_id": 0,
    "speed": 1.0
  }'
```

**批量合成优势**：
- 减少网络开销
- 复用模型加载
- GPU批处理加速
- 性能提升30-50%

---

## 6. 文本预处理优化

### 6.1 文本规范化

**数字转换**：
```python
# 原文本
"我有123个苹果，价格是45.67元"

# 规范化后
"我有一百二十三个苹果，价格是四十五点六七元"
```

**符号处理**：
```python
# 原文本
"请访问www.example.com了解更多信息"

# 规范化后
"请访问 example点com 了解更多信息"
```

**英文处理**：
```python
# 原文本
"我在Apple公司工作，使用Mac电脑"

# 优化建议
"我在Apple公司工作，使用Mac电脑"  # 保持英文，TTS会自动处理
```

### 6.2 标点符号优化

**韵律标点**：

| 符号 | 停顿时长 | 使用场景 |
|------|---------|---------|
| ， | 短停顿 | 句内分隔 |
| 。 | 中等停顿 | 句子结束 |
| ？ | 中等停顿+升调 | 疑问句 |
| ！ | 中等停顿+强调 | 感叹句 |
| ； | 较长停顿 | 分句 |
| … | 长停顿 | 省略、思考 |

**优化示例**：
```python
# 原文本（无标点）
"你好欢迎使用我们的服务如有问题请联系客服"

# 优化后
"你好，欢迎使用我们的服务。如有问题，请联系客服。"
```

### 6.3 多音字处理

**常见多音字**：

| 字 | 读音1 | 读音2 | 上下文示例 |
|----|-------|-------|-----------|
| 行 | háng | xíng | 银行(háng) / 行走(xíng) |
| 重 | zhòng | chóng | 重要(zhòng) / 重复(chóng) |
| 长 | cháng | zhǎng | 长度(cháng) / 成长(zhǎng) |
| 得 | dé | de | 获得(dé) / 跑得快(de) |

**解决方案**：
1. 使用词典文件指定读音
2. 提供上下文帮助TTS判断
3. 必要时使用拼音标注

### 6.4 文本长度优化

**推荐文本长度**：

| 场景 | 单次最大长度 | 建议分段 |
|------|-------------|---------|
| **实时合成** | 50字 | 按句子分段 |
| **通用场景** | 200字 | 按段落分段 |
| **批量合成** | 1000字 | 按章节分段 |

**长文本处理**：
```python
def split_text(text, max_length=200):
    """按标点符号智能分段"""
    sentences = re.split('[。！？；]', text)
    chunks = []
    current_chunk = ""
    
    for sentence in sentences:
        if len(current_chunk) + len(sentence) < max_length:
            current_chunk += sentence
        else:
            chunks.append(current_chunk)
            current_chunk = sentence
    
    if current_chunk:
        chunks.append(current_chunk)
    
    return chunks
```

---

## 7. 场景化配置方案

### 7.1 实时客服场景

**特点**：低延迟、自然友好

```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "default_params": {
    "speaker_id": 0,
    "speed": 1.0,
    "volume": 0.9,
    "pitch": 0
  }
}
```

**关键点**：
- 使用轻量级模型（Kokoro）
- 语速正常（1.0）
- 选择温柔女声
- 文本预处理：添加适当停顿

### 7.2 有声书制作场景

**特点**：高音质、情感丰富

```json
{
  "tts": {
    "model_path": "models/tts/vits-zh-hf-models/model.onnx",
    "voices_path": "models/tts/vits-zh-hf-models/voices.bin",
    "tokens_path": "models/tts/vits-zh-hf-models/tokens.txt",
    "provider": {
      "provider": "cuda",
      "device_id": 0
    }
  },
  "audio": {
    "sample_rate": 22050
  },
  "default_params": {
    "speaker_id": 1,
    "speed": 0.85,
    "volume": 1.0,
    "pitch": 0
  }
}
```

**关键点**：
- 使用VITS高音质模型
- 较慢语速（0.85），便于聆听
- GPU加速提高效率
- 22.05kHz采样率
- 根据内容选择合适音色

### 7.3 新闻播报场景

**特点**：清晰准确、节奏稳定

```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data",
    "provider": {
      "provider": "cpu",
      "num_threads": 6
    }
  },
  "default_params": {
    "speaker_id": 1,
    "speed": 0.95,
    "volume": 1.0,
    "pitch": -1
  }
}
```

**关键点**：
- 选择男声或成熟女声
- 略慢语速（0.95）
- 音调略低（-1），增加权威感
- 文本规范化：数字、日期等

### 7.4 广告配音场景

**特点**：吸引力强、情感突出

```json
{
  "tts": {
    "model_path": "models/tts/vits-zh-hf-models/model.onnx",
    "voices_path": "models/tts/vits-zh-hf-models/voices.bin",
    "tokens_path": "models/tts/vits-zh-hf-models/tokens.txt",
    "provider": {
      "provider": "cuda",
      "device_id": 0
    }
  },
  "audio": {
    "sample_rate": 44100
  },
  "default_params": {
    "speaker_id": 0,
    "speed": 1.1,
    "volume": 1.0,
    "pitch": 2
  }
}
```

**关键点**：
- 高音质模型和采样率
- 略快语速（1.1），增加节奏感
- 音调稍高（+2），增加活力
- GPU加速

### 7.5 导航系统场景

**特点**：清晰简洁、实时性强

```json
{
  "tts": {
    "model_path": "models/tts/piper-zh-cn/model.onnx",
    "voices_path": "models/tts/piper-zh-cn/voices.bin",
    "tokens_path": "models/tts/piper-zh-cn/tokens.txt",
    "provider": {
      "provider": "cpu",
      "num_threads": 2
    }
  },
  "default_params": {
    "speaker_id": 0,
    "speed": 0.9,
    "volume": 1.0,
    "pitch": 0
  }
}
```

**关键点**：
- 轻量级Piper模型，合成速度快
- 清晰女声
- 略慢语速确保理解
- 低资源占用

### 7.6 儿童内容场景

**特点**：活泼可爱、缓慢清晰

```json
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "default_params": {
    "speaker_id": 0,
    "speed": 0.75,
    "volume": 0.9,
    "pitch": 3
  }
}
```

**关键点**：
- 活泼女声
- 较慢语速（0.75）
- 音调偏高（+3），更活泼
- 文本简单化处理

---

## 8. 性能优化技巧

### 8.1 合成缓存

**实现缓存机制**（建议添加）：

```go
type TTSCache struct {
    cache map[string][]byte  // key: text+speaker+speed, value: audio
    mu    sync.RWMutex
    maxSize int
}

func (c *TTSCache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    audio, exists := c.cache[key]
    return audio, exists
}

func (c *TTSCache) Set(key string, audio []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // 检查缓存大小
    if len(c.cache) >= c.maxSize {
        // 删除最旧的条目（FIFO）
        for k := range c.cache {
            delete(c.cache, k)
            break
        }
    }
    
    c.cache[key] = audio
}
```

**缓存策略**：
- 常用短语：欢迎语、提示音
- 固定模板：天气预报、时间播报
- 命中率：可达50-80%
- 性能提升：缓存命中时接近0延迟

### 8.2 流式合成

**流式TTS优势**：
- 边合成边播放
- 降低首字延迟
- 适合长文本
- 提升用户体验

**实现建议**（未来功能）：
```go
// 流式TTS接口
func (p *TTSProvider) SynthesizeStream(text string, chunkCb func([]byte)) error {
    // 分句合成
    sentences := splitSentences(text)
    
    for _, sentence := range sentences {
        audio, err := p.Synthesize(sentence)
        if err != nil {
            return err
        }
        
        // 回调传递音频块
        chunkCb(audio)
    }
    
    return nil
}
```

### 8.3 预加载优化

**模型预加载**：
```go
// 应用启动时预加载TTS模型
func init() {
    // 创建TTS池
    ttsPool = NewTTSPool(poolSize)
    
    // 预热所有实例
    for i := 0; i < poolSize; i++ {
        provider := ttsPool.Get()
        provider.Warmup()  // 合成一次测试文本
        ttsPool.Put(provider)
    }
}
```

**预加载效果**：
- 首次请求无延迟
- 模型常驻内存
- 提升并发性能

### 8.4 并发优化

**资源池配置**：

```go
// TTS通常比ASR需要更少的并发
poolSize := runtime.NumCPU() / 2

// 示例
// 4核CPU: poolSize = 2
// 8核CPU: poolSize = 4
// 16核CPU: poolSize = 8
```

**并发控制**：
```json
{
  "rate_limit": {
    "enabled": true,
    "max_connections": 500,
    "requests_per_second": 200
  }
}
```

---

## 9. 音质优化

### 9.1 采样率提升

**标准配置**：
```json
{
  "audio": {
    "sample_rate": 16000  // 标准
  }
}
```

**高音质配置**：
```json
{
  "audio": {
    "sample_rate": 22050  // 音乐品质
  }
}
```

**专业配置**：
```json
{
  "audio": {
    "sample_rate": 44100  // CD品质
  }
}
```

### 9.2 音频后处理

**降噪处理**（可选）：
```python
import noisereduce as nr

# 对TTS输出进行降噪
audio_clean = nr.reduce_noise(y=audio, sr=sample_rate)
```

**音量归一化**：
```python
from pydub import AudioSegment

# 归一化音量
audio = AudioSegment.from_wav("output.wav")
normalized = audio.normalize()
normalized.export("output_normalized.wav", format="wav")
```

### 9.3 格式优化

**输出格式选择**：

| 格式 | 音质 | 文件大小 | 兼容性 | 推荐场景 |
|------|------|---------|--------|---------|
| WAV | 无损 | 大 | ★★★★★ | 本地存储 |
| MP3 | 有损 | 小 | ★★★★★ | 网络传输 |
| OGG | 有损 | 小 | ★★★★☆ | Web应用 |
| AAC | 有损 | 小 | ★★★★☆ | 移动应用 |

---

## 10. 性能监控与调优

### 10.1 关键指标

**合成速度（RTF）**：
```
RTF = 合成时间 / 音频时长
```
- RTF < 1：实时合成
- RTF = 0.1：1秒音频100ms合成完
- 目标：RTF < 0.3

**首字延迟**：
```
首字延迟 = 开始合成到第一个音频包输出的时间
```
- 实时场景：< 200ms
- 一般场景：< 500ms

**MOS（平均意见得分）**：
- 5分制主观评价
- 目标：MOS > 4.0

### 10.2 监控API

```bash
# 获取TTS统计信息
curl http://localhost:8080/api/v1/tts/stats

# 响应示例
{
  "total_requests": 5000,
  "successful_requests": 4950,
  "failed_requests": 50,
  "avg_latency_ms": 300,
  "p95_latency_ms": 600,
  "p99_latency_ms": 1000,
  "avg_rtf": 0.2,
  "avg_text_length": 50,
  "cache_hit_rate": 0.65
}
```

### 10.3 性能基准测试

**测试脚本**：
```bash
#!/bin/bash
# TTS性能测试

echo "TTS Performance Test"
echo "===================="

# 测试不同文本长度
for length in 10 50 100 200; do
    echo -n "Testing $length chars: "
    
    text=$(head -c $length < /dev/urandom | base64 | tr -dc 'a-zA-Z' | head -c $length)
    
    time=$(curl -X POST \
        -H "Content-Type: application/json" \
        -d "{\"text\":\"$text\",\"speaker_id\":0}" \
        -w "%{time_total}" \
        -o /dev/null \
        -s \
        http://localhost:8080/api/v1/tts/synthesize)
    
    echo "${time}s"
done
```

---

## 11. 常见问题与解决方案

### 11.1 合成速度慢

**症状**：TTS响应时间过长

**解决方案**：

1. **使用轻量级模型**
   ```bash
   Kokoro → Piper
   ```

2. **启用GPU加速**
   ```json
   {
     "provider": {
       "provider": "cuda"
     }
   }
   ```

3. **启用缓存**
   ```go
   // 缓存常用短语
   cache.Set("welcome", audioData)
   ```

4. **批量合成**
   ```bash
   # 一次合成多条
   curl -X POST /api/v1/tts/batch \
     -d '{"texts":["文本1","文本2","文本3"]}'
   ```

### 11.2 音质不佳

**症状**：合成语音听起来不自然

**解决方案**：

1. **升级模型**
   ```bash
   Piper → Kokoro → VITS
   ```

2. **提高采样率**
   ```json
   {
     "audio": {
       "sample_rate": 22050  // 从16000提升
     }
   }
   ```

3. **优化文本**
   ```python
   # 添加标点符号
   "你好欢迎使用" → "你好，欢迎使用。"
   
   # 数字转文字
   "123" → "一百二十三"
   ```

4. **选择合适音色**
   ```bash
   # 尝试不同说话人
   curl /api/v1/tts/speakers
   ```

### 11.3 韵律不自然

**症状**：停顿不当、语调平淡

**解决方案**：

1. **优化标点**
   ```text
   原文：这是一段测试文本没有标点
   优化：这是一段测试文本，没有标点。
   ```

2. **添加韵律符号**
   ```text
   重要内容【】：【重要通知】
   停顿…：请稍等…
   ```

3. **调整语速**
   ```json
   {
     "speed": 0.9  // 略慢，更自然
   }
   ```

### 11.4 多音字错误

**症状**：读音不正确

**解决方案**：

1. **提供上下文**
   ```text
   不好：行
   较好：银行
   最好：我去银行办理业务
   ```

2. **使用词典**
   ```json
   {
     "lexicon": "path/to/custom-lexicon.txt"
   }
   ```

3. **拼音标注**
   ```text
   行(háng)长 → 银行行长
   ```

### 11.5 内存占用高

**症状**：TTS服务内存使用过大

**解决方案**：

1. **减小资源池**
   ```go
   poolSize := 2  // 从5减到2
   ```

2. **使用小模型**
   ```bash
   VITS (500MB) → Kokoro (200MB) → Piper (50MB)
   ```

3. **限制缓存大小**
   ```go
   cache := NewCache(100)  // 限制100条
   ```

4. **及时释放资源**
   ```go
   defer provider.Release()
   ```

---

## 12. 配置模板

### 12.1 生产环境推荐配置

```json
{
  "mode": "unified",
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": 20
  },
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data",
    "provider": {
      "provider": "auto",
      "device_id": 0,
      "num_threads": 4
    },
    "debug": false
  },
  "audio": {
    "sample_rate": 16000
  },
  "rate_limit": {
    "enabled": true,
    "requests_per_second": 200,
    "max_connections": 500
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "both",
    "file_path": "logs/tts.log"
  }
}
```

### 12.2 高音质配置

```json
{
  "tts": {
    "model_path": "models/tts/vits-zh-hf-models/model.onnx",
    "voices_path": "models/tts/vits-zh-hf-models/voices.bin",
    "tokens_path": "models/tts/vits-zh-hf-models/tokens.txt",
    "provider": {
      "provider": "cuda",
      "device_id": 0
    }
  },
  "audio": {
    "sample_rate": 22050
  }
}
```

### 12.3 低延迟配置

```json
{
  "tts": {
    "model_path": "models/tts/piper-zh-cn/model.onnx",
    "voices_path": "models/tts/piper-zh-cn/voices.bin",
    "tokens_path": "models/tts/piper-zh-cn/tokens.txt",
    "provider": {
      "provider": "cpu",
      "num_threads": 2
    }
  },
  "audio": {
    "sample_rate": 16000
  }
}
```

---

## 13. 扩展功能建议

### 13.1 建议添加的配置参数

```go
type TTSConfig struct {
    // 现有参数
    ModelPath  string
    VoicesPath string
    TokensPath string
    DataDir    string
    Provider   ProviderConfig
    Debug      bool
    
    // 建议新增
    DefaultSpeakerID int     `mapstructure:"default_speaker_id" json:"default_speaker_id"`
    DefaultSpeed     float32 `mapstructure:"default_speed" json:"default_speed"`
    DefaultVolume    float32 `mapstructure:"default_volume" json:"default_volume"`
    DefaultPitch     float32 `mapstructure:"default_pitch" json:"default_pitch"`
    
    // 缓存配置
    EnableCache  bool `mapstructure:"enable_cache" json:"enable_cache"`
    CacheSize    int  `mapstructure:"cache_size" json:"cache_size"`
    CacheTTL     int  `mapstructure:"cache_ttl" json:"cache_ttl"`  // 秒
    
    // 输出配置
    OutputFormat   string `mapstructure:"output_format" json:"output_format"`  // wav/mp3/ogg
    OutputSampleRate int  `mapstructure:"output_sample_rate" json:"output_sample_rate"`
    
    // 文本处理
    EnableNormalization bool `mapstructure:"enable_normalization" json:"enable_normalization"`
    MaxTextLength      int  `mapstructure:"max_text_length" json:"max_text_length"`
}
```

### 13.2 API增强建议

**1. 情感控制API**（未来功能）：
```bash
curl -X POST /api/v1/tts/synthesize \
  -d '{
    "text": "今天天气真好",
    "emotion": "happy",
    "intensity": 0.8
  }'
```

**2. 多音色混合**：
```bash
curl -X POST /api/v1/tts/synthesize \
  -d '{
    "text": "对话内容",
    "multi_speaker": [
      {"text": "你好", "speaker_id": 0},
      {"text": "您好", "speaker_id": 1}
    ]
  }'
```

**3. SSML支持**：
```bash
curl -X POST /api/v1/tts/synthesize \
  -d '{
    "ssml": "<speak><prosody rate=\"slow\">慢速</prosody></speak>"
  }'
```

---

## 14. 最佳实践总结

### 优化优先级

1. ⭐⭐⭐ 选择合适的模型（场景匹配）
2. ⭐⭐⭐ 文本预处理（标点、规范化）
3. ⭐⭐⭐ 选择合适的音色
4. ⭐⭐ 调整语速和音调
5. ⭐⭐ 启用缓存机制
6. ⭐⭐ 使用GPU加速（如可用）
7. ⭐ 提高采样率（高音质场景）
8. ⭐ 批量合成优化

### 核心建议

1. **模型选择**
   - 通用场景：Kokoro
   - 高音质：VITS
   - 低延迟：Piper

2. **参数调整**
   - 语速：0.8-1.2之间
   - 音调：-2到+2之间
   - 采样率：16kHz（通用）

3. **性能优化**
   - 启用GPU（可提升3倍性能）
   - 实现缓存（命中率50%+）
   - 批量合成（提升30-50%）

4. **音质提升**
   - 文本规范化
   - 选择合适音色
   - 提高采样率

---

## 参考资料

- [Kokoro TTS文档](https://github.com/thewh1teagle/kokoro-onnx)
- [VITS模型说明](https://github.com/jaywalnut310/vits)
- [Piper TTS项目](https://github.com/rhasspy/piper)
- [项目架构设计](./02-架构设计.md)
- [新增功能说明](./12-新增功能说明.md)

---

## 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v1.0 | 2025-11-16 | 初始版本，覆盖TTS核心优化配置 |


# 配置文件说明

本目录包含多种场景的配置文件模板，请根据实际需求选择合适的配置。

## 配置文件列表

### 1. 基础配置

#### `speech-config.example.json` - 标准配置模板
- **用途**: 基础配置模板，适合入门使用
- **特点**: 使用默认参数，功能完整
- **推荐场景**: 开发测试、功能验证

#### `stt-config.example.json` - STT独立服务配置
- **用途**: 仅运行STT服务
- **特点**: 轻量级，只包含ASR相关配置
- **推荐场景**: 独立部署STT服务

#### `tts-config.example.json` - TTS独立服务配置
- **用途**: 仅运行TTS服务
- **特点**: 轻量级，只包含TTS相关配置
- **推荐场景**: 独立部署TTS服务

---

### 2. 优化配置（v2.0新增）⭐

#### `speech-config-optimized.example.json` - 生产环境推荐配置
```json
{
  "特点": [
    "启用限流保护",
    "启用VAD智能分段",
    "自动Provider选择（CPU/GPU）",
    "优化的日志配置",
    "合理的资源池大小"
  ],
  "适用场景": "生产环境通用配置",
  "性能指标": {
    "识别准确率": "高（★★★★☆）",
    "实时性": "良好（延迟~250ms）",
    "资源占用": "中等"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-optimized.example.json configs/speech-config.json
# 修改模型路径和端口号后启动
./bin/speech-server
```

---

#### `speech-config-low-latency.example.json` - 低延迟配置
```json
{
  "特点": [
    "chunk_size=1024（延迟~64ms）",
    "使用小模型（paraformer-small）",
    "优化的WebSocket参数",
    "精简的日志输出"
  ],
  "适用场景": [
    "实时对话",
    "在线客服",
    "语音助手",
    "互动游戏"
  ],
  "性能指标": {
    "识别准确率": "良好（★★★★☆）",
    "实时性": "优秀（延迟~100ms）",
    "资源占用": "低"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-low-latency.example.json configs/speech-config.json
# 下载paraformer-small模型
./scripts/download_models.sh paraformer-small
# 启动服务
./bin/speech-server
```

---

#### `speech-config-high-accuracy.example.json` - 高准确率配置
```json
{
  "特点": [
    "使用大模型（paraformer-zh-2023-09-14）",
    "chunk_size=8192（更多上下文）",
    "GPU加速（CUDA）",
    "VAD阈值降低（捕获更多音频）"
  ],
  "适用场景": [
    "会议记录",
    "字幕生成",
    "客服质检",
    "法律文档"
  ],
  "性能指标": {
    "识别准确率": "极高（★★★★★）",
    "实时性": "一般（延迟~500ms）",
    "资源占用": "高（需要GPU）"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-high-accuracy.example.json configs/speech-config.json
# 确保有GPU环境
nvidia-smi
# 下载paraformer大模型
./scripts/download_models.sh paraformer-zh-2023
# 启动服务
./bin/speech-server
```

---

### 3. TTS优化配置（v2.0新增）⭐

#### `speech-config-tts-optimized.example.json` - TTS生产环境推荐配置
```json
{
  "特点": [
    "Kokoro多语言模型",
    "自动Provider选择（CPU/GPU）",
    "标准音质（16kHz）",
    "合理的资源池配置",
    "启用限流保护"
  ],
  "适用场景": "通用TTS生产环境",
  "性能指标": {
    "音质": "良好（★★★★☆）",
    "合成速度": "快速（RTF~0.2）",
    "资源占用": "中等"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-tts-optimized.example.json configs/speech-config.json
# 修改模型路径后启动
./bin/speech-server
```

---

#### `speech-config-tts-high-quality.example.json` - 高音质TTS配置
```json
{
  "特点": [
    "VITS高音质模型",
    "22.05kHz采样率",
    "GPU加速",
    "情感丰富、韵律自然",
    "适合音频制作"
  ],
  "适用场景": [
    "有声书制作",
    "广告配音",
    "高品质语音内容",
    "专业音频制作"
  ],
  "性能指标": {
    "音质": "极高（★★★★★）",
    "合成速度": "中等（RTF~0.4）",
    "资源占用": "高（需要GPU）"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-tts-high-quality.example.json configs/speech-config.json
# 下载VITS模型（注意：示例路径，需要实际下载）
# 确保有GPU环境
nvidia-smi
# 启动服务
./bin/speech-server
```

---

#### `speech-config-tts-fast.example.json` - 快速TTS配置
```json
{
  "特点": [
    "Piper轻量级模型",
    "合成速度极快（RTF~0.1）",
    "低资源占用",
    "支持高并发",
    "标准音质"
  ],
  "适用场景": [
    "实时语音导航",
    "大规模批量合成",
    "边缘设备",
    "移动端应用"
  ],
  "性能指标": {
    "音质": "良好（★★★☆☆）",
    "合成速度": "极快（RTF~0.1）",
    "资源占用": "极低"
  }
}
```

**使用方法**:
```bash
cp configs/speech-config-tts-fast.example.json configs/speech-config.json
# 下载Piper模型（注意：示例路径，需要实际下载）
# 启动服务
./bin/speech-server
```

---

## 配置选择指南

### 按场景选择

#### STT场景

| 场景 | 推荐配置 | 延迟 | 准确率 | 成本 |
|------|---------|------|--------|------|
| **实时对话** | low-latency | ★★★★★ | ★★★★☆ | 💰 |
| **会议记录** | high-accuracy | ★★☆☆☆ | ★★★★★ | 💰💰💰 |
| **客服系统** | optimized | ★★★★☆ | ★★★★☆ | 💰💰 |
| **字幕生成** | high-accuracy | ★★☆☆☆ | ★★★★★ | 💰💰💰 |
| **语音助手** | low-latency | ★★★★★ | ★★★★☆ | 💰 |
| **开发测试** | example | ★★★☆☆ | ★★★☆☆ | 💰 |

#### TTS场景

| 场景 | 推荐配置 | 音质 | 合成速度 | 成本 |
|------|---------|------|---------|------|
| **有声书制作** | tts-high-quality | ★★★★★ | ★★★☆☆ | 💰💰💰 |
| **广告配音** | tts-high-quality | ★★★★★ | ★★★☆☆ | 💰💰💰 |
| **客服系统** | tts-optimized | ★★★★☆ | ★★★★☆ | 💰💰 |
| **语音导航** | tts-fast | ★★★☆☆ | ★★★★★ | 💰 |
| **新闻播报** | tts-optimized | ★★★★☆ | ★★★★☆ | 💰💰 |
| **批量合成** | tts-fast | ★★★☆☆ | ★★★★★ | 💰 |
| **开发测试** | example | ★★★☆☆ | ★★★☆☆ | 💰 |

### 按硬件选择

#### CPU环境
```bash
# 8核以上CPU - 使用优化配置
cp configs/speech-config-optimized.example.json configs/speech-config.json

# 4核CPU - 使用低延迟配置（小模型）
cp configs/speech-config-low-latency.example.json configs/speech-config.json
```

#### GPU环境
```bash
# 有GPU - 使用高准确率配置
cp configs/speech-config-high-accuracy.example.json configs/speech-config.json
# 记得修改provider为cuda
```

#### 边缘设备/ARM
```bash
# 资源受限 - 使用低延迟配置
cp configs/speech-config-low-latency.example.json configs/speech-config.json
# 进一步减小chunk_size和pool_size
```

---

## 关键参数说明

### STT相关参数

#### 模型选择
```json
{
  "stt": {
    "model_path": "...",
    "language": "zh|en|auto"
  }
}
```

| 参数 | 说明 | 可选值 | 推荐值 |
|------|------|--------|--------|
| model_path | 模型文件路径 | 见下表 | SenseVoice INT8 |
| language | 语言代码 | zh/en/auto | auto |

**推荐模型**：
- **SenseVoice**: 多语言、高准确率（推荐）
- **Paraformer大模型**: 最高准确率，带标点
- **Paraformer小模型**: 低延迟，资源占用小

#### Provider配置
```json
{
  "provider": {
    "provider": "cpu|cuda|auto",
    "num_threads": 4
  }
}
```

| provider | 适用场景 | 性能 | 成本 |
|----------|---------|------|------|
| cpu | 无GPU环境 | 基准 | 低 |
| cuda | 有NVIDIA GPU | 3-5倍提升 | 高 |
| auto | 自动检测 | 最优 | - |

**num_threads建议值**：
- 4核CPU: 4
- 8核CPU: 8
- GPU模式: 2-4

#### 音频参数
```json
{
  "audio": {
    "sample_rate": 16000,
    "chunk_size": 4096
  }
}
```

**chunk_size影响**：
- 1024: 延迟64ms，适合实时对话
- 2048: 延迟128ms，平衡选择
- 4096: 延迟256ms，**推荐值**
- 8192: 延迟512ms，批量处理

---

### TTS相关参数

#### 模型选择
```json
{
  "tts": {
    "model_path": "...",
    "voices_path": "...",
    "data_dir": "..."
  }
}
```

| 模型 | 音质 | 速度 | 资源占用 | 适用场景 |
|------|------|------|---------|---------|
| **Kokoro** | ★★★★☆ | ★★★★☆ | 中 | 通用（推荐） |
| **VITS** | ★★★★★ | ★★★☆☆ | 高 | 高音质制作 |
| **Piper** | ★★★☆☆ | ★★★★★ | 低 | 快速合成 |

#### 音频参数
```json
{
  "audio": {
    "sample_rate": 16000
  }
}
```

**采样率选择**：
- 8000 Hz: 电话质量
- 16000 Hz: **标准质量（推荐）**
- 22050 Hz: 音乐质量
- 44100 Hz: CD质量（专业制作）

#### TTS合成参数（API）
```json
{
  "speaker_id": 0,
  "speed": 1.0,
  "volume": 1.0,
  "pitch": 0
}
```

| 参数 | 范围 | 默认值 | 说明 |
|------|------|--------|------|
| speaker_id | 0-N | 0 | 说话人ID |
| speed | 0.5-2.0 | 1.0 | 语速倍率 |
| volume | 0.0-1.0 | 1.0 | 音量 |
| pitch | -12~+12 | 0 | 音调（半音） |

**语速建议**：
- 0.7-0.8: 有声书、儿童内容
- 0.9-1.0: **通用场景（推荐）**
- 1.1-1.3: 广告、快速提醒

---

### VAD配置

```json
{
  "vad": {
    "enabled": true,
    "threshold": 0.5
  }
}
```

**threshold（阈值）调整**：
- 0.3-0.4: 嘈杂环境（街道、工厂）
- 0.5: **通用场景（推荐）**
- 0.6-0.7: 安静环境（办公室、录音棚）

**VAD效果**：
- ✅ 过滤静音，减少无效处理
- ✅ 智能分段，提高准确率
- ✅ 降低CPU使用率20%左右

### 限流配置

```json
{
  "rate_limit": {
    "enabled": true,
    "requests_per_second": 1000,
    "max_connections": 2000
  }
}
```

**生产环境建议启用**：
- 防止系统过载
- 保护后端服务
- 合理分配资源

---

## 配置优化流程

### 步骤1: 选择基础配置
```bash
# 根据场景选择配置模板
cp configs/speech-config-optimized.example.json configs/speech-config.json
```

### 步骤2: 修改必需参数
```bash
vim configs/speech-config.json
```
修改内容：
1. ✅ 模型路径（model_path）
2. ✅ 端口号（port）
3. ✅ 日志路径（file_path）

### 步骤3: 测试验证
```bash
# 启动服务
./bin/speech-server

# 测试识别
curl -X POST -F "audio=@test.wav" \
  http://localhost:8080/api/v1/stt/recognize

# 查看统计
curl http://localhost:8080/api/v1/stats
```

### 步骤4: 性能调优
根据统计信息调整参数：

**延迟过高？**
→ 减小chunk_size (4096→2048)
→ 使用小模型
→ 启用GPU

**准确率不够？**
→ 使用大模型
→ 增大chunk_size (4096→8192)
→ 调整VAD阈值

**CPU使用率高？**
→ 启用VAD过滤
→ 使用INT8模型
→ 减少并发连接

---

## 热重载功能

v2.0新增配置热重载，支持动态更新部分配置：

```bash
# 修改配置文件
vim configs/speech-config.json

# 2秒后自动生效，查看日志确认
tail -f logs/speech.log | grep "Config file changed"
```

**支持热重载的配置**：
- ✅ logging.level（日志级别）
- ✅ rate_limit（限流参数）
- ✅ session（会话参数）
- ❌ server.port（需要重启）
- ❌ model_path（需要重启）

---

## 故障排查

### 问题1: 模型文件找不到
```bash
# 检查模型路径
ls -l models/asr/

# 下载模型
./scripts/download_models.sh
```

### 问题2: GPU不可用
```bash
# 检查CUDA
nvidia-smi

# 切换到CPU模式
# 修改 provider: "cuda" → "cpu"
```

### 问题3: 识别准确率低
```bash
# 1. 检查音频质量（16kHz、无噪音）
# 2. 尝试更大的模型
# 3. 调整VAD阈值
# 4. 增大chunk_size
```

### 问题4: 延迟太高
```bash
# 1. 减小chunk_size
# 2. 使用小模型
# 3. 启用GPU
# 4. 增加num_threads
```

---

## 参考文档

- [STT配置优化指南](../docs/13-STT配置优化指南.md) - STT详细优化指导
- [TTS配置优化指南](../docs/14-TTS配置优化指南.md) - TTS详细优化指导⭐
- [新增功能说明](../docs/12-新增功能说明.md) - v2.0新功能
- [架构设计](../docs/02-架构设计.md) - 系统架构
- [部署文档](../docs/DEPLOYMENT.md) - 部署指南

---

## 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v2.0 | 2025-11-16 | 新增优化配置文件和配置说明 |
| v1.0 | 2025-11-15 | 初始配置文件 |


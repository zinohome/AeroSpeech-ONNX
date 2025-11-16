# 替代 TTS 模型指南

> 如果 Kokoro v1.1 不满足您的需求，本文档提供其他可用的中文 TTS 模型选项。

## 📋 当前问题

如果您觉得 Kokoro v1.1 声音"怪异"，请先检查：

### ✅ 配置检查清单

1. **采样率是否正确** ⭐ 最常见问题
   ```json
   {
     "audio": {
       "sample_rate": 24000  // 必须是 24000，不能是 16000
     }
   }
   ```

2. **说话人选择是否合适**
   - 默认说话人(0) 是美式英语女声，不适合中文
   - 中文女声：ID 3-57（55个选择）
   - 中文男声：ID 58-102（45个选择）

3. **语速是否合适**
   ```json
   {
     "speed": 0.8  // 建议 0.8-1.2，太快会不自然
   }
   ```

4. **文本预处理**
   - 添加标点符号
   - 避免过长的句子
   - 数字转换为文字

---

## 🎯 替代方案

### 方案1: 优化 Kokoro v1.1 使用 ⭐ 强烈推荐

**步骤1: 修复配置**
```bash
# 检查 configs/speech-config.json
# 确保 sample_rate: 24000
```

**步骤2: 测试不同说话人**
```bash
# 使用测试脚本
./scripts/test_kokoro_speakers.sh
```

**步骤3: 调整参数**
```json
{
  "text": "你好，这是测试。",
  "speaker_id": 10,    // 尝试不同的ID
  "speed": 0.9         // 调整语速
}
```

**优势**：
- ✅ 已集成，无需额外下载
- ✅ 103个说话人，选择丰富
- ✅ 支持中英文混合
- ✅ 性能优异

---

### 方案2: VITS 中文模型

**模型信息**：
- 名称：vits-zh-aishell3
- 说话人数：174个（中文）
- 采样率：16000 Hz
- 模型大小：~115 MB

**下载地址**（需手动确认）：
```bash
# HuggingFace 地址
https://huggingface.co/csukuangfj/vits-zh-aishell3

# 需要下载的文件
model.onnx          # 模型文件
lexicon.txt         # 词典
tokens.txt          # 词表
```

**配置示例**：
```json
{
  "tts": {
    "model_path": "models/tts/vits-zh-aishell3/model.onnx",
    "tokens_path": "models/tts/vits-zh-aishell3/tokens.txt",
    "lexicon": "models/tts/vits-zh-aishell3/lexicon.txt",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 16000  // VITS 使用 16000
  }
}
```

**特点**：
- ✅ 174个中文说话人
- ✅ 纯中文训练，发音准确
- ❌ 不支持英文
- ❌ 采样率较低（音质略逊）

---

### 方案3: MeloTTS

**模型信息**：
- 名称：MeloTTS
- 语言：中文、英文、日文等
- 说话人：多个高质量音色
- 采样率：44100 Hz（高质量）

**下载地址**（需手动确认）：
```bash
# GitHub
https://github.com/myshell-ai/MeloTTS

# 可能需要转换为 ONNX 格式
```

**特点**：
- ✅ 音质好（44100 Hz）
- ✅ 多语言支持
- ❌ 可能需要格式转换
- ❌ 文档较少

---

### 方案4: 其他 VITS 模型

sherpa-onnx 官方支持多个 VITS 模型，访问官方文档查看：

**官方文档**：
- [sherpa-onnx TTS 模型列表](https://k2-fsa.github.io/sherpa/onnx/tts/index.html)
- [中文 TTS 模型](https://k2-fsa.github.io/sherpa/onnx/tts/pretrained_models/zh.html)

**推荐模型**：
1. **vits-piper-zh_CN**
   - 中文专用
   - 多个音色
   - ONNX 格式，即插即用

2. **vits-mms-zh**
   - Meta 发布
   - 多语言支持
   - 质量稳定

---

## 🛠️ 添加新模型的步骤

### 步骤1: 下载模型文件

```bash
# 创建模型目录
mkdir -p models/tts/[模型名称]

# 下载必要文件
# - model.onnx
# - tokens.txt
# - 其他配置文件
```

### 步骤2: 更新配置文件

```json
{
  "tts": {
    "model_path": "models/tts/[模型名称]/model.onnx",
    "tokens_path": "models/tts/[模型名称]/tokens.txt",
    // 根据模型要求添加其他配置
  },
  "audio": {
    "sample_rate": [模型的采样率]
  }
}
```

### 步骤3: 测试模型

```bash
# 重新编译
go build -o speech-server ./cmd/speech-server/

# 启动服务
./speech-server --config configs/speech-config.json

# 测试 API
curl -X POST http://localhost:8780/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{"text":"测试文本","speaker_id":0}' \
  --output test.wav

# 播放测试
afplay test.wav
```

---

## 📊 模型对比

| 模型 | 说话人数 | 采样率 | 语言 | 音质 | 易用性 | 推荐度 |
|------|---------|--------|------|------|--------|--------|
| **Kokoro v1.1** | 103 | 24000 | 中英 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| vits-zh-aishell3 | 174 | 16000 | 中 | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| MeloTTS | 多个 | 44100 | 多 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| vits-piper-zh | 多个 | 16000 | 中 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

---

## 💡 推荐流程

### 1️⃣ 先优化当前配置
```bash
# 修复采样率
# configs/speech-config.json
"sample_rate": 24000

# 重启服务器测试
```

### 2️⃣ 测试不同说话人
```bash
# 使用测试脚本
./scripts/test_kokoro_speakers.sh

# 尝试这些ID
# 中文女声: 3, 10, 15, 20, 25, 30, 35, 40
# 中文男声: 58, 65, 70, 75, 80, 85
```

### 3️⃣ 调整参数
```json
{
  "speed": 0.85,      // 降低语速
  "speaker_id": 10     // 更换说话人
}
```

### 4️⃣ 如果仍不满意
- 查看 sherpa-onnx 官方文档
- 选择其他模型
- 手动下载和配置

---

## 🔗 有用的链接

- [sherpa-onnx 官方文档](https://k2-fsa.github.io/sherpa/)
- [sherpa-onnx TTS 模型列表](https://k2-fsa.github.io/sherpa/onnx/tts/pretrained_models/index.html)
- [Kokoro TTS 官方文档](https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html)
- [csukuangfj HuggingFace](https://huggingface.co/csukuangfj)

---

## ❓ 常见问题

### Q: 为什么声音听起来很怪？
**A**: 最常见原因是采样率配置错误：
- Kokoro v1.1: 必须 24000 Hz
- VITS: 通常 16000 或 22050 Hz
- 配置不匹配会导致变调和失真

### Q: 如何选择合适的说话人？
**A**: 
1. 使用测试脚本试听所有说话人
2. 根据场景选择（客服、导航、播报等）
3. 注意语速和音调的匹配

### Q: 可以混合使用多个模型吗？
**A**: 
- 理论上可以，但需要：
  - 修改代码支持模型切换
  - 处理不同的采样率
  - 管理多个模型的内存占用
- 建议：选定一个模型深度优化

### Q: 如何提升音质？
**A**:
1. ✅ 使用更高采样率的模型
2. ✅ 选择高质量的说话人
3. ✅ 优化文本预处理
4. ✅ 调整语速参数
5. ✅ 使用 GPU 加速

---

## 📝 总结

**优先级建议**：

1. **修复配置** - 确保采样率正确 ✅
2. **测试说话人** - 找到合适的音色 ✅
3. **调整参数** - 优化语速和音调 ✅
4. **更换模型** - 如果以上都不满意 ⚠️

**记住**：80% 的音质问题可以通过正确的配置和参数调整解决，不需要更换模型！


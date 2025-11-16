# Piper TTS 模型使用指南

## 📋 模型信息

### vits-piper-zh_CN-huayan-medium

- **模型名称**: 华研女声（Huayan Female Voice）
- **语言**: 中文
- **质量等级**: Medium（中等）
- **采样率**: 22050 Hz
- **模型大小**: ~60 MB
- **说话人**: 单说话人（ID: 0）
- **来源**: [HuggingFace](https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium)

### 模型特点

✅ **优势**：
- 纯中文训练，发音准确
- 模型体积小，加载快
- 女声音色自然温和
- 适合客服、播报等场景

❌ **局限**：
- 仅支持中文
- 单一说话人
- 采样率较低（22050 Hz vs Kokoro 的 24000 Hz）

---

## 📥 模型文件

已下载的文件：

```
models/tts/vits-piper-zh_CN-huayan-medium/
├── zh_CN-huayan-medium.onnx          # 模型文件 (60MB)
├── zh_CN-huayan-medium.onnx.json     # 模型配置
├── tokens.txt                         # Token 词表
└── espeak-ng-data/                    # 文本处理数据
```

---

## ⚙️ 配置说明

### 方法1: 使用示例配置文件

已创建配置文件：`configs/speech-config-piper.example.json`

**关键配置项**：

```json
{
  "tts": {
    "model_path": "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx",
    "model_config": "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx.json",
    "tokens_path": "models/tts/vits-piper-zh_CN-huayan-medium/tokens.txt",
    "data_dir": "models/tts/vits-piper-zh_CN-huayan-medium/espeak-ng-data",
    "provider": {
      "provider": "cpu",
      "num_threads": 4
    }
  },
  "audio": {
    "sample_rate": 22050  // ⚠️ 必须是 22050，不是 24000
  }
}
```

### 方法2: 修改现有配置

修改 `configs/speech-config.json`：

```bash
# 1. 备份当前配置
cp configs/speech-config.json configs/speech-config-kokoro-backup.json

# 2. 使用 Piper 配置
cp configs/speech-config-piper.example.json configs/speech-config.json
```

---

## 🚀 启动服务

### 步骤1: 编译（如果需要）

```bash
cd /Users/zhangjun/CursorProjects/AeroSpeech-ONNX
go build -o speech-server ./cmd/speech-server/
```

### 步骤2: 启动服务器

```bash
# 使用 Piper 配置启动
./speech-server --config configs/speech-config-piper.example.json

# 或者修改默认配置后直接启动
./speech-server
```

### 步骤3: 验证服务

```bash
# 检查服务状态
curl http://localhost:8780/api/v1/health

# 预期输出
{"code":200,"message":"success","data":{"status":"ok"}}
```

---

## 🎤 测试语音合成

### REST API 测试

```bash
# 基础测试
curl -X POST http://localhost:8780/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{"text":"你好，这是华研女声的语音合成测试。","speaker_id":0,"speed":1.0}' \
  --output piper_test.wav

# 播放测试
afplay piper_test.wav
```

### 不同语速测试

```bash
# 创建测试脚本
cat > test_piper_speeds.sh << 'EOF'
#!/bin/bash

TEXT="你好，这是华研女声的语音合成测试。今天天气真不错。"

for speed in 0.8 0.9 1.0 1.1 1.2; do
    echo "测试语速: ${speed}..."
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEXT}\",\"speaker_id\":0,\"speed\":${speed}}" \
      --output "piper_speed_${speed}.wav"
    
    echo "播放 speed ${speed}..."
    afplay "piper_speed_${speed}.wav"
done

echo "✅ 测试完成！选择您最喜欢的语速。"
EOF

chmod +x test_piper_speeds.sh
./test_piper_speeds.sh
```

### Web 界面测试

访问 http://localhost:8780/tts-test.html

**注意事项**：
- 说话人 ID 固定为 0（Piper 是单说话人模型）
- 推荐语速范围：0.8 - 1.2
- 使用带标点的完整句子获得最佳效果

---

## 📊 与 Kokoro v1.1 对比

| 特性 | Piper (huayan-medium) | Kokoro v1.1 |
|------|----------------------|-------------|
| **语言** | 中文 | 中文 + 英文 |
| **说话人数** | 1 | 103 |
| **采样率** | 22050 Hz | 24000 Hz |
| **音质** | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **模型大小** | 60 MB | 310 MB |
| **加载速度** | 快 | 较慢 |
| **合成速度** | 快 | 中等 |
| **中文发音** | 准确 | 准确 |
| **英文发音** | ❌ | ✅ |
| **适用场景** | 纯中文场景 | 多场景通用 |

---

## 🎯 使用建议

### 适合使用 Piper 的场景

✅ **推荐场景**：
1. **纯中文应用** - 不需要英文支持
2. **资源受限环境** - 内存/存储有限
3. **快速响应需求** - 需要极快的加载和合成速度
4. **单一音色需求** - 只需要一个女声音色
5. **客服系统** - 温和的女声适合客服场景
6. **语音播报** - 新闻、通知等播报场景

❌ **不推荐场景**：
1. 需要多说话人选择
2. 需要中英文混合
3. 需要男声
4. 对音质要求极高（希望更高采样率）

### 性能优化建议

1. **CPU 优化**
```json
{
  "tts": {
    "provider": {
      "provider": "cpu",
      "num_threads": 4  // 根据 CPU 核心数调整
    }
  }
}
```

2. **内存优化**
```json
{
  "tts": {
    "provider": {
      "pool_size": 2  // Piper 模型小，可以用较小的池
    }
  }
}
```

3. **GPU 加速**（如果可用）
```json
{
  "tts": {
    "provider": {
      "provider": "cuda",
      "device_id": 0,
      "num_threads": 2
    }
  }
}
```

---

## 🔧 故障排查

### 问题1: 声音怪异或失真

**原因**: 采样率配置错误

**解决方法**:
```json
{
  "audio": {
    "sample_rate": 22050  // 必须是 22050，不能是 16000 或 24000
  }
}
```

### 问题2: 无法加载模型

**可能原因**:
1. 模型文件不完整
2. espeak-ng-data 缺失
3. 路径配置错误

**解决方法**:
```bash
# 检查文件
ls -lh models/tts/vits-piper-zh_CN-huayan-medium/

# 重新下载模型（如果需要）
./scripts/download_models.sh

# 检查 espeak-ng-data
ls -la models/tts/vits-piper-zh_CN-huayan-medium/espeak-ng-data/
```

### 问题3: 中文发音不准确

**解决方法**:
1. 添加标点符号
2. 使用空格分隔短语
3. 将数字转换为中文
4. 调整语速（尝试 0.9）

**示例**:
```json
{
  "text": "今天天气真不错，温度是 二十三 摄氏度。",  // ✅ 好
  "text": "今天天气真不错温度是23摄氏度"           // ❌ 差
}
```

### 问题4: 合成速度慢

**优化措施**:
1. 增加 num_threads
2. 使用 GPU（如果可用）
3. 减少 pool_size
4. 预热模型

---

## 🔄 切换回 Kokoro

如果测试后想切换回 Kokoro v1.1：

```bash
# 恢复 Kokoro 配置
cp configs/speech-config-kokoro-backup.json configs/speech-config.json

# 或者手动修改
{
  "tts": {
    "model_path": "models/tts/kokoro-multi-lang-v1_1/model.onnx",
    "voices_path": "models/tts/kokoro-multi-lang-v1_1/voices.bin",
    "tokens_path": "models/tts/kokoro-multi-lang-v1_1/tokens.txt",
    "data_dir": "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data"
  },
  "audio": {
    "sample_rate": 24000  // 改回 24000
  }
}

# 重启服务器
./speech-server
```

---

## 📈 性能基准

### 合成速度（参考）

**测试环境**: MacBook Pro M1, 8 GB RAM

| 文本长度 | Piper | Kokoro v1.1 |
|---------|-------|-------------|
| 10 字 | 50ms | 80ms |
| 50 字 | 200ms | 350ms |
| 100 字 | 380ms | 650ms |

### 内存占用

| 模型 | 加载后内存 | 合成时峰值 |
|------|-----------|-----------|
| Piper | ~150 MB | ~180 MB |
| Kokoro v1.1 | ~450 MB | ~520 MB |

---

## 💡 最佳实践

### 1. 文本预处理

```go
// 建议的文本预处理
func preprocessText(text string) string {
    // 1. 数字转换
    text = convertNumbersToChinese(text)
    
    // 2. 添加标点
    text = addPunctuation(text)
    
    // 3. 移除特殊字符
    text = removeSpecialChars(text)
    
    return text
}
```

### 2. 批量合成

```bash
# 批量合成脚本
cat > batch_synthesize.sh << 'EOF'
#!/bin/bash

while IFS= read -r line; do
    filename=$(echo "$line" | md5)
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${line}\",\"speaker_id\":0,\"speed\":1.0}" \
      --output "output/${filename}.wav"
    echo "✅ ${line}"
done < texts.txt
EOF
```

### 3. 错误重试

```go
// 带重试的合成
func synthesizeWithRetry(text string, maxRetries int) ([]byte, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        audio, err := ttsProvider.Synthesize(text, 0, 1.0)
        if err == nil {
            return audio, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return nil, lastErr
}
```

---

## 📚 相关资源

- [Piper 官方文档](https://github.com/rhasspy/piper)
- [sherpa-onnx 文档](https://k2-fsa.github.io/sherpa/onnx/tts/)
- [HuggingFace 模型页面](https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium)
- [模型下载脚本](../scripts/download_models.sh)

---

## 🎯 总结

**Piper huayan-medium 模型**是一个轻量级、高效的中文 TTS 模型：

✅ **优点**：
- 模型小，加载快
- 纯中文，发音准确
- 资源占用低
- 适合嵌入式和资源受限环境

❌ **缺点**：
- 单说话人
- 不支持英文
- 采样率较低

**推荐使用场景**：
- 纯中文应用
- 需要快速响应
- 资源有限的环境
- 对音色要求不高

**如果需要**：
- 多说话人选择 → 使用 Kokoro v1.1
- 中英文混合 → 使用 Kokoro v1.1
- 更高音质 → 使用 Kokoro v1.1

根据您的具体需求选择最合适的模型！


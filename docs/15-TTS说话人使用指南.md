# TTS说话人使用指南

> 版本: v2.0  
> 更新时间: 2025-11-16  
> 模型: Kokoro v1.1  
> 参考: [sherpa-onnx 官方文档](https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html)

## 概述

Kokoro v1.1 多语言TTS模型支持 **103 个说话人**，涵盖英文（美式/英式）和中文，提供丰富的音色选择。

### 模型信息

| 属性 | 值 |
|------|-----|
| **总说话人数** | 103 |
| **采样率** | 24000 Hz |
| **支持语言** | 中文、英文 |
| **模型来源** | [HuggingFace - Kokoro v1.1-zh](https://huggingface.co/hexgrad/Kokoro-82M-v1.1-zh) |

---

## 说话人分类

### 按前缀分类

| 前缀 | 含义 | ID 范围 | 数量 |
|------|------|---------|------|
| **af** | American Female (美式英语女声) | 0-1 | 2个 |
| **bf** | British Female (英式英语女声) | 2 | 1个 |
| **zf** | Chinese Female (中文女声) | 3-57 | **55个** |
| **zm** | Chinese Male (中文男声) | 58-102 | **45个** |

### 完整说话人列表

#### 美式英语女声 (2个)

| Speaker ID | Name | 特点 |
|-----------|------|------|
| 0 | af_maple | 默认说话人，温和自然 |
| 1 | af_sol | 清晰明朗 |

#### 英式英语女声 (1个)

| Speaker ID | Name | 特点 |
|-----------|------|------|
| 2 | bf_vale | 英式口音，优雅 |

#### 中文女声 (55个)

| ID 范围 | 说话人名称列表 |
|---------|---------------|
| 3-10 | zf_001, zf_002, zf_003, zf_004, zf_005, zf_006, zf_007, zf_008 |
| 11-20 | zf_017, zf_018, zf_019, zf_021, zf_022, zf_023, zf_024, zf_026, zf_027, zf_028 |
| 21-30 | zf_032, zf_036, zf_038, zf_039, zf_040, zf_042, zf_043, zf_044, zf_046, zf_047 |
| 31-40 | zf_048, zf_049, zf_051, zf_059, zf_060, zf_067, zf_070, zf_071, zf_072, zf_073 |
| 41-50 | zf_074, zf_075, zf_076, zf_077, zf_078, zf_079, zf_083, zf_084, zf_085, zf_086 |
| 51-57 | zf_087, zf_088, zf_090, zf_092, zf_093, zf_094, zf_099 |

#### 中文男声 (45个)

| ID 范围 | 说话人名称列表 |
|---------|---------------|
| 58-67 | zm_009, zm_010, zm_011, zm_012, zm_013, zm_014, zm_015, zm_016, zm_020, zm_025 |
| 68-77 | zm_029, zm_030, zm_031, zm_033, zm_034, zm_035, zm_037, zm_041, zm_045, zm_050 |
| 78-87 | zm_052, zm_053, zm_054, zm_055, zm_056, zm_057, zm_058, zm_061, zm_062, zm_063 |
| 88-97 | zm_064, zm_065, zm_066, zm_068, zm_069, zm_080, zm_081, zm_082, zm_089, zm_091 |
| 98-102 | zm_095, zm_096, zm_097, zm_098, zm_100 |

---

## 使用方法

### 1. Web 测试页面

访问 TTS 测试页面：
```
http://localhost:8080/tts
```

**操作步骤**：
1. 点击「刷新列表 (103个说话人)」按钮
2. 从下拉列表中选择说话人（按类别分组）
3. 输入要合成的文本
4. 调整语速（可选）
5. 点击「合成」或「API合成」

**界面示意**：
```
说话人：
┌────────────────────────────────────┐
│ 美式英语女声 (2个)                  │
│   0 - af_maple (默认)              │
│   1 - af_sol                       │
│ 英式英语女声 (1个)                  │
│   2 - bf_vale                      │
│ 中文女声 (55个)                     │
│   3 - zf_001                       │
│   4 - zf_002                       │
│   ...                              │
│   57 - zf_099                      │
│ 中文男声 (45个)                     │
│   58 - zm_009                      │
│   59 - zm_010                      │
│   ...                              │
│   102 - zm_100                     │
└────────────────────────────────────┘
[刷新列表 (103个说话人)] 已加载 103 个说话人
```

---

### 2. REST API

#### 获取说话人列表
```bash
curl http://localhost:8080/api/v1/tts/speakers | jq

# 响应示例
{
  "code": 200,
  "data": {
    "info": {
      "model": "Kokoro v1.1",
      "reference": "https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html",
      "sample_rate": 24000
    },
    "speakers": [
      {
        "category": "American Female",
        "gender": "female",
        "id": 0,
        "language": "en-US",
        "name": "af_maple"
      },
      ...
    ],
    "total": 103
  },
  "message": "success"
}
```

#### 基本调用
```bash
# 使用默认说话人 (af_maple - ID 0)
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{
    "text": "你好世界",
    "speaker_id": 0,
    "speed": 1.0
  }' \
  --output output.wav
```

#### 使用不同说话人
```bash
# 美式英语女声 - af_maple (0)
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"Hello World","speaker_id":0}' \
  --output af_maple.wav && afplay af_maple.wav

# 英式英语女声 - bf_vale (2)
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"Good morning","speaker_id":2}' \
  --output bf_vale.wav && afplay bf_vale.wav

# 中文女声 - zf_001 (3)
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"你好，欢迎使用语音合成","speaker_id":3}' \
  --output zf_001.wav && afplay zf_001.wav

# 中文男声 - zm_009 (58)
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"大家好，这是新闻播报","speaker_id":58}' \
  --output zm_009.wav && afplay zm_009.wav
```

---

### 3. WebSocket API

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/tts');

ws.onopen = function() {
    // 使用中文女声 zf_005 (ID 7)
    ws.send(JSON.stringify({
        type: 'synthesize',
        data: {
            text: '你好，这是流式语音合成测试',
            speaker_id: 7,
            speed: 1.0
        }
    }));
};
```

---

### 4. 批量测试说话人

#### 测试所有说话人脚本
```bash
#!/bin/bash
# 测试所有103个说话人

echo "测试美式英语女声 (0-1)"
for id in 0 1; do
    curl -s -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"Hello, speaker ${id}\",\"speaker_id\":${id}}" \
      --output "speaker_${id}.wav"
    echo "Generated: speaker_${id}.wav"
done

echo "测试英式英语女声 (2)"
curl -s -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"Good afternoon","speaker_id":2}' \
  --output "speaker_2.wav"

echo "测试中文女声 (3-57)"
for id in {3..57}; do
    curl -s -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"你好，我是说话人${id}\",\"speaker_id\":${id}}" \
      --output "speaker_${id}.wav"
    echo "Generated: speaker_${id}.wav"
done

echo "测试中文男声 (58-102)"
for id in {58..102}; do
    curl -s -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"你好，我是说话人${id}\",\"speaker_id\":${id}}" \
      --output "speaker_${id}.wav"
    echo "Generated: speaker_${id}.wav"
done

echo "所有103个说话人测试完成！"
```

#### 测试部分代表性说话人
```bash
#!/bin/bash
# 测试代表性说话人

declare -A speakers=(
    [0]="af_maple:Hello from American female"
    [2]="bf_vale:Good day from British female"
    [3]="zf_001:你好，中文女声1号"
    [10]="zf_008:你好，中文女声8号"
    [30]="zf_047:你好，中文女声47号"
    [58]="zm_009:你好，中文男声1号"
    [75]="zm_041:你好，中文男声中段"
    [102]="zm_100:你好，中文男声最后一个"
)

for id in "${!speakers[@]}"; do
    IFS=':' read -r name text <<< "${speakers[$id]}"
    echo "测试 Speaker $id - $name"
    curl -s -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"$text\",\"speaker_id\":$id}" \
      --output "${name}.wav"
    echo "生成: ${name}.wav"
    # 播放测试
    afplay "${name}.wav"
done
```

---

## 场景推荐

### 通用英文内容
**推荐**: af_maple (0) 或 af_sol (1)
```bash
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"Welcome to our service","speaker_id":0}'
```

### 英式英语内容
**推荐**: bf_vale (2)
```bash
curl -X POST http://localhost:8080/api/v1/tts/synthesize \
  -d '{"text":"Good afternoon, how may I help you?","speaker_id":2}'
```

### 中文女声内容
**推荐**: 尝试 ID 3-57 中的不同说话人
```bash
# 测试前几个
for id in 3 4 5 6 7; do
    curl -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"你好，欢迎使用\",\"speaker_id\":${id}}" \
      --output "zf_${id}.wav"
    afplay "zf_${id}.wav"
done
```

### 中文男声内容
**推荐**: 尝试 ID 58-102 中的不同说话人
```bash
# 测试前几个
for id in 58 59 60 61 62; do
    curl -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"大家好，这是新闻播报\",\"speaker_id\":${id}}" \
      --output "zm_${id}.wav"
    afplay "zm_${id}.wav"
done
```

---

## 选择说话人的建议

### 1. 根据语言选择
- **英文内容** → 使用 af_maple (0), af_sol (1), 或 bf_vale (2)
- **中文内容** → 使用 zf_* (3-57) 或 zm_* (58-102)

### 2. 根据性别选择
- **女声** → af_* (0-1), bf_* (2), zf_* (3-57)
- **男声** → zm_* (58-102)

### 3. 如何找到合适的说话人？

由于有 103 个说话人，建议：

**方法1: 系统化测试**
```bash
# 测试每个类别的前几个
# 美式英语女声: 0, 1
# 英式英语女声: 2
# 中文女声: 3, 4, 5, 10, 20, 30, 40, 50
# 中文男声: 58, 60, 70, 80, 90, 100
```

**方法2: 随机采样**
```bash
# 从每个范围随机选择几个测试
FEMALE_ZH=(3 15 27 39 51)
MALE_ZH=(58 70 82 94 102)

for id in "${FEMALE_ZH[@]}"; do
    curl -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"测试说话人${id}\",\"speaker_id\":${id}}" \
      --output "test_${id}.wav"
    afplay "test_${id}.wav"
done
```

**方法3: Web界面测试**
```
1. 打开 http://localhost:8080/tts
2. 点击「刷新列表」
3. 在下拉框中选择不同说话人
4. 输入测试文本
5. 点击「API合成」并播放
6. 记录喜欢的说话人ID
```

---

## 常见问题

### Q1: 为什么有103个说话人？

**A**: Kokoro v1.1 模型经过大量数据训练，支持103个不同音色的说话人，包括：
- 3个英语说话人（2个美式、1个英式）
- 100个中文说话人（55个女声、45个男声）

参考：[sherpa-onnx 官方文档](https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html)

### Q2: 所有说话人都能说中英文吗？

**A**: 理论上所有说话人都支持中英文混合，但建议：
- 英文内容优先使用 af_*/bf_* 前缀的说话人 (0-2)
- 中文内容优先使用 zf_*/zm_* 前缀的说话人 (3-102)

### Q3: 不同ID的中文说话人有什么区别？

**A**: 每个说话人有独特的音色特征：
- 音调高低不同
- 语速习惯不同
- 发音风格不同
- 情感表现不同

建议通过实际测试选择最适合您需求的说话人。

### Q4: 如何快速找到最佳说话人？

**A**: 推荐流程：
1. **确定性别和语言** - 缩小范围到对应类别
2. **测试边界** - 测试ID范围的开始、中间、结束
3. **精细选择** - 在喜欢的区间内逐个测试
4. **记录偏好** - 记录每个场景的最佳说话人ID

### Q5: 说话人ID和质量有关系吗？

**A**: ID号码**不代表质量**，只是编号：
- ID 0 是默认说话人，但不一定最好
- 较大的ID号（如100+）不代表更好
- 建议根据实际听感选择

---

## 性能说明

### 合成速度
不同说话人的合成速度基本一致：
- **RTF**: ~0.2 (CPU), ~0.1 (GPU)
- **首字延迟**: 200-300ms

### 内存占用
- 103个说话人的嵌入向量存储在 `voices.bin` 文件中
- 文件大小：约 10-20MB
- 运行时内存增加：可忽略不计

### 并发性能
支持同时使用不同说话人：
```bash
# 并发请求不同说话人
for id in 0 3 58; do
    curl -X POST http://localhost:8080/api/v1/tts/synthesize \
      -d "{\"text\":\"并发测试\",\"speaker_id\":${id}}" \
      --output "concurrent_${id}.wav" &
done
wait
```

---

## 技术细节

### 说话人嵌入向量

每个说话人都有唯一的嵌入向量表示：
```
Speaker ID → Speaker Embedding (256-512维向量) → Voice Characteristics
```

### 模型架构

```
Text Input
    ↓
Text Encoder
    ↓
Speaker Embedding (103个) ← 选择说话人ID
    ↓
Decoder
    ↓
Audio Output (24kHz PCM)
```

---

## 更新日志

| 版本 | 日期 | 更新内容 |
|------|------|---------|
| v2.1 | 2025-11-16 | 更正说话人数量为103个，更新完整列表 |
| v2.0 | 2025-11-16 | 初始版本（错误：只有15个说话人） |

---

## 参考资料

- [sherpa-onnx 官方文档 - Kokoro v1.1](https://k2-fsa.github.io/sherpa/onnx/tts/all/Chinese-English/kokoro-multi-lang-v1_1.html)
- [HuggingFace - Kokoro v1.1-zh](https://huggingface.co/hexgrad/Kokoro-82M-v1.1-zh)
- [TTS配置优化指南](./14-TTS配置优化指南.md)
- [API设计文档](./API.md)

# 中文 TTS 自然度优化指南

> 如何让 TTS 合成的中文听起来更像真人说话

## 🎯 问题诊断

### Kokoro v1.1 中文发音问题

**常见反馈**：
- ❌ 发音不像中国人
- ❌ 比较生硬、机械
- ❌ 语调不自然
- ❌ 听起来像外国人说中文

**原因分析**：
1. **多语言训练** - Kokoro 是多语言模型，在中英文上都做了训练，可能在纯中文自然度上有所妥协
2. **说话人选择** - 103 个说话人中，不是所有都适合中文
3. **文本处理** - 没有适当的标点、停顿、韵律标记
4. **参数设置** - 语速、音调等参数不合适

---

## 🛠️ 优化方案

### 方案1: 优化 Kokoro 使用

#### 1.1 选择合适的中文说话人

**推荐的中文女声**（相对更自然）：

| 说话人ID | 名称 | 特点 | 推荐场景 |
|---------|------|------|---------|
| 10 | zf_008 | 温和，较自然 | 客服、助手 |
| 20 | zf_019 | 清晰，亲和 | 播报、导航 |
| 30 | zf_028 | 活泼，年轻 | 娱乐、互动 |
| 40 | zf_040 | 稳重，成熟 | 商务、正式 |
| 50 | zf_051 | 柔和，舒适 | 阅读、讲故事 |

**推荐的中文男声**：

| 说话人ID | 名称 | 特点 | 推荐场景 |
|---------|------|------|---------|
| 65 | zm_014 | 沉稳，磁性 | 播音、讲解 |
| 70 | zm_031 | 清晰，标准 | 新闻、公告 |
| 80 | zm_053 | 温和，亲切 | 客服、对话 |

**测试方法**：
```bash
# 运行说话人对比脚本
./scripts/compare_chinese_speakers.sh
```

#### 1.2 调整语速

**问题**：标准语速（1.0）可能显得生硬

**优化**：
```json
{
  "speed": 0.85  // 降低语速 15%，更从容自然
}
```

**建议语速范围**：
- 0.80-0.85: 慢速，适合讲解、教学
- 0.85-0.95: 舒适，适合日常对话 ⭐ 推荐
- 0.95-1.05: 标准，适合播报
- 1.05-1.20: 快速，适合紧急通知

#### 1.3 文本优化

**❌ 不好的文本**：
```
今天天气真不错温度是23摄氏度
```

**✅ 优化后的文本**：
```
今天天气真不错，温度是二十三摄氏度。
```

**优化要点**：

1. **添加标点符号**
```python
# 句号: 表示完整停顿
text = "今天天气不错。明天也会很好。"

# 逗号: 表示短暂停顿
text = "今天天气不错，温度适宜，适合出门。"

# 问号/感叹号: 增加语调变化
text = "你好！今天过得怎么样？"
```

2. **数字转中文**
```python
# ❌ 阿拉伯数字
"温度是23度，湿度是65%"

# ✅ 中文数字
"温度是二十三度，湿度是百分之六十五"
```

3. **适当分句**
```python
# ❌ 长句
"今天天气真不错温度适宜阳光明媚空气清新非常适合外出活动"

# ✅ 分句
"今天天气真不错。温度适宜，阳光明媚，空气清新。非常适合外出活动。"
```

4. **添加韵律词**
```python
# ❌ 生硬
"好的明白了"

# ✅ 自然
"好的，我明白了。"

# ✅ 更自然
"好的，我明白了，马上为您处理。"
```

5. **避免特殊字符**
```python
# ❌ 特殊字符
"价格$99.99, 5★好评"

# ✅ 中文表达
"价格九十九点九九美元，五星好评"
```

#### 1.4 添加停顿标记

**使用空格控制停顿**：
```python
# 短停顿
text = "你好 很高兴认识你"

# 长停顿（连续空格）
text = "欢迎光临   请问需要什么帮助"
```

#### 1.5 配置示例

**优化后的 API 调用**：
```json
{
  "text": "您好，很高兴为您服务。请问有什么可以帮助您的吗？",
  "speaker_id": 20,
  "speed": 0.9
}
```

---

### 方案2: 使用 Piper 模型 ⭐ 强烈推荐

**为什么 Piper 更自然**：
- ✅ 纯中文训练（不是多语言）
- ✅ 专门针对中文韵律优化
- ✅ 发音更像中国人
- ✅ 语调变化更自然

**快速切换**：
```bash
# 停止 Kokoro 服务器
# Ctrl+C

# 启动 Piper 服务器
./speech-server --config configs/speech-config-piper.example.json
```

**测试 Piper**：
```bash
# 测试脚本
./scripts/test_piper_model.sh

# 或快速测试
curl -X POST http://localhost:8780/api/v1/tts/synthesize \
  -H "Content-Type: application/json" \
  -d '{"text":"你好，我是华研。很高兴为您服务。","speaker_id":0,"speed":0.9}' \
  --output piper_test.wav && afplay piper_test.wav
```

**Piper 优势**：
```
Kokoro: "你好，我是AI助手。"  → 听起来像外国人说中文
Piper:  "你好，我是AI助手。"  → 听起来像真正的中国人 ✅
```

---

### 方案3: Kokoro vs Piper 对比测试

**运行对比脚本**：
```bash
./scripts/compare_kokoro_vs_piper.sh
```

这个脚本会：
1. 用 Kokoro 的 3 个推荐说话人合成
2. 用 Piper 合成相同文本
3. 并排播放对比
4. 帮助您选择最自然的方案

---

## 📊 对比分析

### Kokoro vs Piper

| 指标 | Kokoro v1.1 | Piper 华研 | 胜者 |
|------|------------|-----------|------|
| **中文发音准确性** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Piper ✅ |
| **自然度（像中国人）** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | Piper ✅ |
| **语调变化** | ⭐⭐⭐ | ⭐⭐⭐⭐ | Piper ✅ |
| **说话人选择** | 103个 ⭐⭐⭐⭐⭐ | 1个 ⭐ | Kokoro ✅ |
| **英文支持** | ✅ ⭐⭐⭐⭐⭐ | ❌ | Kokoro ✅ |
| **音质** | 24kHz ⭐⭐⭐⭐ | 22kHz ⭐⭐⭐ | Kokoro ✅ |
| **资源占用** | 310MB ⭐⭐ | 60MB ⭐⭐⭐⭐⭐ | Piper ✅ |
| **加载速度** | 较慢 ⭐⭐⭐ | 快 ⭐⭐⭐⭐⭐ | Piper ✅ |

**结论**：
- 纯中文应用 → **Piper 胜出** ⭐
- 需要多说话人 → Kokoro
- 需要中英文混合 → Kokoro
- 对自然度要求高 → **Piper 胜出** ⭐

---

## 🎨 实战案例

### 案例1: 客服系统

**场景**：智能客服回答用户问题

**优化前（Kokoro 默认）**：
```json
{
  "text": "你好欢迎咨询有什么可以帮助你",
  "speaker_id": 0,
  "speed": 1.0
}
```
**问题**：生硬、没有停顿、像机器人

**优化后（Piper + 文本优化）**：
```json
{
  "text": "您好，欢迎咨询。请问有什么可以帮助您的吗？",
  "speaker_id": 0,
  "speed": 0.9
}
```
**改进**：
- ✅ 使用"您"更礼貌
- ✅ 添加标点和停顿
- ✅ 语速降低 10%
- ✅ 使用 Piper 模型

---

### 案例2: 语音播报

**场景**：新闻播报

**优化前**：
```json
{
  "text": "今天是2024年11月16日星期六天气晴温度23度",
  "speaker_id": 3,
  "speed": 1.0
}
```

**优化后**：
```json
{
  "text": "今天是二零二四年十一月十六日，星期六。天气晴，温度二十三摄氏度。",
  "speaker_id": 70,
  "speed": 0.95
}
```
**改进**：
- ✅ 数字转中文
- ✅ 添加标点分句
- ✅ 使用播音风格的男声（ID 70）
- ✅ 略微降低语速

---

### 案例3: 语音助手

**场景**：日常对话

**Piper 最佳实践**：
```json
{
  "text": "好的，我明白了。马上为您查询相关信息。请稍等片刻。",
  "speaker_id": 0,
  "speed": 0.9
}
```

**要点**：
- ✅ 自然的对话语气
- ✅ 适当的语气词（"好的"、"马上"）
- ✅ 礼貌用语（"您"、"请"）
- ✅ 舒适的语速

---

## 🔧 文本预处理函数

### Python 示例

```python
import re

def optimize_chinese_text(text):
    """优化中文文本以提高 TTS 自然度"""
    
    # 1. 数字转中文
    text = convert_numbers_to_chinese(text)
    
    # 2. 添加标点（如果缺少）
    if not any(p in text for p in '。！？，、；：'):
        # 简单的句子切分
        text = text.replace('  ', '。')
        if not text.endswith(('。', '！', '？')):
            text += '。'
    
    # 3. 规范化空格
    text = re.sub(r'\s+', ' ', text)
    
    # 4. 移除特殊字符
    text = re.sub(r'[^\u4e00-\u9fa5a-zA-Z0-9，。！？、；：""''（）\s]', '', text)
    
    # 5. 添加礼貌用语
    if text.startswith('你'):
        text = text.replace('你', '您', 1)
    
    return text

def convert_numbers_to_chinese(text):
    """将数字转换为中文"""
    # 简化版本
    num_map = {
        '0': '零', '1': '一', '2': '二', '3': '三', '4': '四',
        '5': '五', '6': '六', '7': '七', '8': '八', '9': '九'
    }
    for digit, chinese in num_map.items():
        text = text.replace(digit, chinese)
    return text

# 使用示例
text = "今天温度23度，湿度65%"
optimized = optimize_chinese_text(text)
print(optimized)  # "今天温度二三度，湿度六五。"
```

### Go 示例

```go
func OptimizeChineseText(text string) string {
    // 1. 数字转中文
    text = convertNumbersToChinese(text)
    
    // 2. 添加标点
    if !hasChinesePunctuation(text) {
        text = addPunctuation(text)
    }
    
    // 3. 规范化空格
    text = normalizeSpaces(text)
    
    return text
}
```

---

## 🎯 推荐流程

### 步骤1: 运行对比测试

```bash
# 对比 Kokoro 和 Piper
./scripts/compare_kokoro_vs_piper.sh
```

### 步骤2: 选择最佳方案

**如果 Piper 更自然**：
```bash
# 切换到 Piper
cp configs/speech-config-piper.example.json configs/speech-config.json
./speech-server
```

**如果 Kokoro 某个说话人可以**：
```json
// 使用那个说话人 ID
{
  "speaker_id": 20,  // 您觉得最好的
  "speed": 0.9       // 调整语速
}
```

### 步骤3: 优化文本

```python
# 在发送到 TTS 前优化文本
text = optimize_chinese_text(user_input)
```

---

## 📝 总结

### 快速建议

**如果您觉得 Kokoro 不像中国人说话**：

1. **首选方案**：尝试 Piper 模型 ⭐⭐⭐⭐⭐
   ```bash
   ./speech-server --config configs/speech-config-piper.example.json
   ```

2. **备选方案**：优化 Kokoro 使用
   - 测试不同说话人（ID 10, 20, 30, 40, 50）
   - 降低语速到 0.85-0.95
   - 优化文本（标点、分句、数字转换）

3. **对比测试**：
   ```bash
   ./scripts/compare_kokoro_vs_piper.sh
   ```

### 关键优化点

| 优化项 | 重要性 | 效果 |
|--------|--------|------|
| **使用 Piper 模型** | ⭐⭐⭐⭐⭐ | 根本性改善 |
| 选择合适说话人 | ⭐⭐⭐⭐ | 显著改善 |
| 调整语速 | ⭐⭐⭐⭐ | 明显改善 |
| 文本优化 | ⭐⭐⭐⭐⭐ | 显著改善 |
| 添加标点停顿 | ⭐⭐⭐⭐ | 明显改善 |

### 最终建议

**纯中文应用** → **强烈推荐 Piper** ✅

Piper 专门为中文优化，发音更像中国人，是解决"不自然"问题的最佳方案！

---

## 🔗 相关资源

- [Piper 模型使用指南](17-Piper模型使用指南.md)
- [TTS 配置优化指南](14-TTS配置优化指南.md)
- [TTS 说话人使用指南](15-TTS说话人使用指南.md)
- [替代 TTS 模型指南](16-替代TTS模型指南.md)


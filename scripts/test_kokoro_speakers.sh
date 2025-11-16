#!/bin/bash
# Kokoro v1.1 说话人测试脚本

echo "🎤 测试 Kokoro v1.1 不同说话人效果"
echo "=========================================="

# 创建测试输出目录
mkdir -p tts_test_output

# 测试文本
TEXT="你好，这是语音合成测试。今天天气真不错，让我们一起来听听不同说话人的声音效果。"

# 推荐测试的说话人ID（中文女声）
SPEAKERS=(3 10 15 20 25 30 35 40 45 50)

for id in "${SPEAKERS[@]}"; do
    echo ""
    echo "📢 测试说话人 ${id}..."
    
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEXT}\",\"speaker_id\":${id},\"speed\":0.9}" \
      --output "tts_test_output/speaker_${id}.wav"
    
    if [ $? -eq 0 ]; then
        echo "✅ 已生成: tts_test_output/speaker_${id}.wav"
        echo "▶️  播放中..."
        afplay "tts_test_output/speaker_${id}.wav"
        echo "按回车键继续下一个说话人..."
        read
    else
        echo "❌ 生成失败"
    fi
done

echo ""
echo "✅ 测试完成！所有音频文件保存在 tts_test_output/ 目录"
echo "您可以反复播放对比，选择最喜欢的说话人。"


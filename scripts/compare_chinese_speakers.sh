#!/bin/bash
# Kokoro 中文说话人对比测试脚本

echo "🎤 测试 Kokoro v1.1 不同中文说话人的自然度"
echo "=========================================="

# 检查服务器
if ! curl -s http://localhost:8780/api/v1/health > /dev/null; then
    echo "❌ 服务器未运行！请先启动:"
    echo "  ./speech-server --config configs/speech-config.json"
    exit 1
fi

# 创建输出目录
mkdir -p speaker_comparison

# 测试文本（日常对话）
TEXT="你好，很高兴认识你。今天天气不错，我们一起出去走走吧。"

# 推荐的中文女声（从 103 个说话人中精选）
# 这些说话人在中文发音上相对更自然
FEMALE_SPEAKERS=(
    3   # zf_001 - 第一个中文女声
    10  # zf_008
    15  # zf_015
    20  # zf_019
    25  # zf_024
    30  # zf_028
    35  # zf_036
    40  # zf_040
    45  # zf_046
    50  # zf_051
)

# 推荐的中文男声
MALE_SPEAKERS=(
    58  # zm_009 - 第一个中文男声
    65  # zm_014
    70  # zm_031
    75  # zm_037
    80  # zm_053
)

echo ""
echo "📢 测试中文女声（10个推荐说话人）"
echo "========================================"

for id in "${FEMALE_SPEAKERS[@]}"; do
    echo ""
    echo "🎵 测试女声说话人 ${id}..."
    
    # 测试标准语速
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEXT}\",\"speaker_id\":${id},\"speed\":1.0}" \
      --output "speaker_comparison/female_${id}_speed_1.0.wav"
    
    # 测试较慢语速（可能更自然）
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEXT}\",\"speaker_id\":${id},\"speed\":0.85}" \
      --output "speaker_comparison/female_${id}_speed_0.85.wav"
    
    if [ $? -eq 0 ]; then
        echo "✅ 已生成:"
        echo "   - speaker_comparison/female_${id}_speed_1.0.wav (标准)"
        echo "   - speaker_comparison/female_${id}_speed_0.85.wav (较慢)"
        echo ""
        echo "▶️  播放标准语速..."
        afplay "speaker_comparison/female_${id}_speed_1.0.wav"
        echo "▶️  播放较慢语速..."
        afplay "speaker_comparison/female_${id}_speed_0.85.wav"
        echo ""
        read -p "按回车继续下一个，或输入 's' 跳过剩余: " choice
        if [ "$choice" = "s" ]; then
            break
        fi
    fi
done

echo ""
echo "📢 测试中文男声（5个推荐说话人）"
echo "========================================"

for id in "${MALE_SPEAKERS[@]}"; do
    echo ""
    echo "🎵 测试男声说话人 ${id}..."
    
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEXT}\",\"speaker_id\":${id},\"speed\":0.9}" \
      --output "speaker_comparison/male_${id}.wav"
    
    if [ $? -eq 0 ]; then
        echo "✅ 已生成: speaker_comparison/male_${id}.wav"
        echo "▶️  播放中..."
        afplay "speaker_comparison/male_${id}.wav"
        echo ""
        read -p "按回车继续下一个，或输入 's' 跳过剩余: " choice
        if [ "$choice" = "s" ]; then
            break
        fi
    fi
done

echo ""
echo "=========================================="
echo "✅ 测试完成！"
echo ""
echo "所有音频保存在 speaker_comparison/ 目录"
echo ""
echo "💡 优化建议:"
echo "  1. 选择您觉得最自然的说话人 ID"
echo "  2. 语速建议: 0.85-0.95 (比标准慢一些更自然)"
echo "  3. 添加标点符号和适当停顿"
echo ""
echo "如果所有 Kokoro 说话人都不满意，建议尝试:"
echo "  ./speech-server --config configs/speech-config-piper.example.json"
echo "  (Piper 是纯中文模型，发音更像中国人)"


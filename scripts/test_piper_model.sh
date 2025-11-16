#!/bin/bash
# Piper TTS 模型测试脚本

echo "🎤 测试 Piper TTS 华研女声模型"
echo "=========================================="

# 检查服务器是否运行
if ! curl -s http://localhost:8780/api/v1/health > /dev/null; then
    echo "❌ 服务器未运行！"
    echo "请先启动服务器:"
    echo "  ./speech-server --config configs/speech-config-piper.example.json"
    exit 1
fi

echo "✅ 服务器正在运行"
echo ""

# 创建测试输出目录
mkdir -p piper_test_output

# 测试文本
TEXTS=(
    "你好，我是华研。"
    "今天天气真不错。"
    "这是一个语音合成测试。"
    "欢迎使用人工智能语音助手。"
    "感谢您的使用，祝您生活愉快。"
)

SPEEDS=(0.8 0.9 1.0 1.1 1.2)

echo "📢 测试不同文本内容..."
echo "=========================================="

for i in "${!TEXTS[@]}"; do
    text="${TEXTS[$i]}"
    echo ""
    echo "测试 $((i+1)): ${text}"
    
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${text}\",\"speaker_id\":0,\"speed\":1.0}" \
      --output "piper_test_output/text_$((i+1)).wav"
    
    if [ $? -eq 0 ]; then
        echo "✅ 已生成: piper_test_output/text_$((i+1)).wav"
        echo "▶️  播放中..."
        afplay "piper_test_output/text_$((i+1)).wav"
    else
        echo "❌ 生成失败"
    fi
done

echo ""
echo "=========================================="
echo "📢 测试不同语速..."
echo "=========================================="

TEST_TEXT="你好，这是语音合成测试。今天天气真不错。"

for speed in "${SPEEDS[@]}"; do
    echo ""
    echo "测试语速: ${speed}"
    
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${TEST_TEXT}\",\"speaker_id\":0,\"speed\":${speed}}" \
      --output "piper_test_output/speed_${speed}.wav"
    
    if [ $? -eq 0 ]; then
        echo "✅ 已生成: piper_test_output/speed_${speed}.wav"
        echo "▶️  播放中..."
        afplay "piper_test_output/speed_${speed}.wav"
    else
        echo "❌ 生成失败"
    fi
done

echo ""
echo "=========================================="
echo "✅ 测试完成！"
echo ""
echo "所有音频文件保存在 piper_test_output/ 目录"
echo ""
echo "📊 文件列表:"
ls -lh piper_test_output/
echo ""
echo "💡 使用建议:"
echo "  - 推荐语速: 0.9 - 1.0"
echo "  - 文本中添加标点符号效果更好"
echo "  - 适合纯中文应用场景"
echo ""
echo "如果效果满意，可以切换到 Piper 配置:"
echo "  cp configs/speech-config-piper.example.json configs/speech-config.json"
echo "  ./speech-server"


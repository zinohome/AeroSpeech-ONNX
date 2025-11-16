#!/bin/bash
# Kokoro vs Piper 中文发音对比测试

echo "🎭 Kokoro vs Piper 中文发音对比测试"
echo "=========================================="
echo ""
echo "本测试将对比两个模型的中文发音自然度："
echo "  - Kokoro v1.1 (多语言模型)"
echo "  - Piper 华研女声 (纯中文模型)"
echo ""

# 创建输出目录
mkdir -p model_comparison

# 测试文本（日常对话场景）
TEST_TEXTS=(
    "你好，很高兴认识你。"
    "今天天气真不错，我们一起出去走走吧。"
    "请问您需要什么帮助吗？"
    "这个产品的质量非常好，我强烈推荐。"
    "明天上午九点，我们在会议室见面。"
    "非常感谢您的支持，祝您生活愉快。"
)

echo "=========================================="
echo "第1步: 测试 Kokoro v1.1"
echo "=========================================="
echo ""
echo "请确保服务器使用 Kokoro 配置运行"
echo "如果未运行，请在另一个终端执行:"
echo "  ./speech-server --config configs/speech-config.json"
echo ""
read -p "按回车开始测试 Kokoro..."

# 测试几个推荐的中文说话人
KOKORO_SPEAKERS=(10 20 30)

for speaker_id in "${KOKORO_SPEAKERS[@]}"; do
    echo ""
    echo "📢 Kokoro 说话人 ${speaker_id}:"
    
    for i in "${!TEST_TEXTS[@]}"; do
        text="${TEST_TEXTS[$i]}"
        echo "  测试 $((i+1)): ${text}"
        
        curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
          -H "Content-Type: application/json" \
          -d "{\"text\":\"${text}\",\"speaker_id\":${speaker_id},\"speed\":0.9}" \
          --output "model_comparison/kokoro_sp${speaker_id}_text${i}.wav" \
          2>/dev/null
        
        if [ $? -ne 0 ]; then
            echo "  ⚠️  Kokoro 服务器未运行或出错"
            echo ""
            echo "  请在另一个终端运行:"
            echo "    ./speech-server --config configs/speech-config.json"
            echo ""
            exit 1
        fi
    done
    
    echo "  ✅ Kokoro 说话人 ${speaker_id} 测试完成"
done

echo ""
echo "=========================================="
echo "第2步: 切换到 Piper 模型"
echo "=========================================="
echo ""
echo "⚠️  现在需要重启服务器以切换到 Piper 模型"
echo ""
echo "请在运行服务器的终端:"
echo "  1. 按 Ctrl+C 停止当前服务器"
echo "  2. 运行: ./speech-server --config configs/speech-config-piper.example.json"
echo ""
read -p "切换完成后，按回车继续..."

echo ""
echo "📢 测试 Piper 华研女声:"

for i in "${!TEST_TEXTS[@]}"; do
    text="${TEST_TEXTS[$i]}"
    echo "  测试 $((i+1)): ${text}"
    
    curl -s -X POST http://localhost:8780/api/v1/tts/synthesize \
      -H "Content-Type: application/json" \
      -d "{\"text\":\"${text}\",\"speaker_id\":0,\"speed\":0.9}" \
      --output "model_comparison/piper_text${i}.wav" \
      2>/dev/null
    
    if [ $? -ne 0 ]; then
        echo "  ⚠️  Piper 服务器未运行或出错"
        echo ""
        echo "  请确保使用 Piper 配置启动:"
        echo "    ./speech-server --config configs/speech-config-piper.example.json"
        echo ""
        exit 1
    fi
done

echo "  ✅ Piper 测试完成"

echo ""
echo "=========================================="
echo "第3步: 对比播放"
echo "=========================================="
echo ""

for i in "${!TEST_TEXTS[@]}"; do
    text="${TEST_TEXTS[$i]}"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📝 文本 $((i+1)): ${text}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    for speaker_id in "${KOKORO_SPEAKERS[@]}"; do
        echo "🔊 Kokoro 说话人 ${speaker_id}:"
        afplay "model_comparison/kokoro_sp${speaker_id}_text${i}.wav"
        sleep 0.5
    done
    
    echo ""
    echo "🔊 Piper 华研女声:"
    afplay "model_comparison/piper_text${i}.wav"
    
    echo ""
    read -p "按回车继续下一组对比..."
done

echo ""
echo "=========================================="
echo "✅ 对比测试完成！"
echo "=========================================="
echo ""
echo "📊 评估标准:"
echo "  1. 发音准确性 - 哪个模型发音更标准？"
echo "  2. 自然度 - 哪个听起来更像真人说话？"
echo "  3. 语调 - 哪个语调变化更自然？"
echo "  4. 流畅度 - 哪个说话更流畅？"
echo ""
echo "💡 建议:"
echo "  - 如果 Piper 更自然 → 使用 Piper 配置"
echo "  - 如果 Kokoro 某个说话人可以 → 使用那个说话人"
echo "  - 如果都不满意 → 考虑其他模型或文本优化"
echo ""
echo "📁 所有音频文件保存在: model_comparison/"
echo ""
echo "❓ 您觉得哪个模型更好？"


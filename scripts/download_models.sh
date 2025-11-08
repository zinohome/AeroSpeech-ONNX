#!/bin/bash
# STT和TTS模型文件下载脚本
# 用于测试和运行STT/TTS服务

set -e

# 创建模型目录
mkdir -p models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17
mkdir -p models/speaker
mkdir -p models/tts/kokoro-multi-lang-v1_0

echo "=========================================="
echo "开始下载STT和TTS模型文件"
echo "=========================================="

# ==========================================
# STT模型文件
# ==========================================
echo ""
echo "--- 下载STT模型文件 ---"

# 1. ASR模型文件（int8量化版本，体积更小）
echo "下载: ASR模型 (model.int8.onnx)..."
curl -L --retry 5 --retry-delay 2 \
  -o models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx \
  https://huggingface.co/csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/resolve/main/model.int8.onnx

# 2. ASR Tokens文件
echo "下载: ASR Tokens文件 (tokens.txt)..."
curl -L --retry 5 --retry-delay 2 \
  -o models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt \
  https://huggingface.co/csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/resolve/main/tokens.txt

# 3. 说话人识别模型（可选，用于说话人识别功能）
echo "下载: 说话人识别模型..."
curl -L --retry 5 --retry-delay 2 \
  -o models/speaker/3dspeaker_speech_campplus_sv_zh_en_16k-common_advanced.onnx \
  https://huggingface.co/csukuangfj/speaker-embedding-models/resolve/main/3dspeaker_speech_campplus_sv_zh_en_16k-common_advanced.onnx

# ==========================================
# TTS模型文件
# ==========================================
echo ""
echo "--- 下载TTS模型文件 ---"

# 1. TTS模型文件（Kokoro多语言模型）
echo "下载: TTS模型 (model.onnx)..."
curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/model.onnx \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/model.onnx

# 2. TTS Voices文件（多说话人支持）
echo "下载: TTS Voices文件 (voices.bin)..."
curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/voices.bin \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/voices.bin

# 3. TTS Tokens文件
echo "下载: TTS Tokens文件 (tokens.txt)..."
curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/tokens.txt \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/tokens.txt

# 4. espeak-ng-data目录（文本处理，需要完整目录）
echo "下载: espeak-ng-data目录（使用git sparse-checkout）..."
mkdir -p models/tts/kokoro-multi-lang-v1_0/espeak-ng-data-temp
cd models/tts/kokoro-multi-lang-v1_0/espeak-ng-data-temp

# 使用git sparse-checkout下载完整目录
if command -v git >/dev/null 2>&1; then
    git init
    git remote add origin https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0 || true
    git config core.sparseCheckout true
    echo "espeak-ng-data/*" > .git/info/sparse-checkout
    git pull --depth=1 origin main 2>&1 | tail -5
    
    # 复制文件到目标目录
    if [ -d "espeak-ng-data" ]; then
        cp -r espeak-ng-data/* ../espeak-ng-data/ 2>/dev/null || true
        cd ../..
        rm -rf models/tts/kokoro-multi-lang-v1_0/espeak-ng-data-temp
        echo "✅ espeak-ng-data目录下载完成（约120个文件，8.6MB）"
    else
        echo "⚠️  git方式下载失败，尝试使用curl下载关键文件..."
        cd ../..
        rm -rf models/tts/kokoro-multi-lang-v1_0/espeak-ng-data-temp
        mkdir -p models/tts/kokoro-multi-lang-v1_0/espeak-ng-data
        # 下载关键文件
        for file in phondata phontab phonindex; do
            curl -L --retry 5 --retry-delay 2 \
              -o models/tts/kokoro-multi-lang-v1_0/espeak-ng-data/$file \
              https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/espeak-ng-data/$file || echo "警告: $file下载失败"
        done
    fi
else
    echo "⚠️  git未安装，使用curl下载关键文件..."
    cd ../..
    mkdir -p models/tts/kokoro-multi-lang-v1_0/espeak-ng-data
    # 下载关键文件
    for file in phondata phontab phonindex; do
        curl -L --retry 5 --retry-delay 2 \
          -o models/tts/kokoro-multi-lang-v1_0/espeak-ng-data/$file \
          https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/espeak-ng-data/$file || echo "警告: $file下载失败"
    done
fi
cd "$(dirname "$0")/.."

# 5. dict目录（字典文件）
echo "下载: dict目录..."
mkdir -p models/tts/kokoro-multi-lang-v1_0/dict
# dict目录可能包含多个文件，这里下载主要的
curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/dict/en_dict.txt \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/dict/en_dict.txt || echo "警告: en_dict.txt下载失败（可能不需要）"

# 6. Lexicon文件
echo "下载: Lexicon文件..."
curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/lexicon-us-en.txt \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/lexicon-us-en.txt || echo "警告: lexicon-us-en.txt下载失败（可能不需要）"

curl -L --retry 5 --retry-delay 2 \
  -o models/tts/kokoro-multi-lang-v1_0/lexicon-zh.txt \
  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_0/resolve/main/lexicon-zh.txt || echo "警告: lexicon-zh.txt下载失败（可能不需要）"

echo ""
echo "=========================================="
echo "模型文件下载完成！"
echo "=========================================="
echo ""
echo "下载的文件列表："
echo ""
echo "STT模型："
echo "  - models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx"
echo "  - models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt"
echo "  - models/speaker/3dspeaker_speech_campplus_sv_zh_en_16k-common_advanced.onnx"
echo ""
echo "TTS模型："
echo "  - models/tts/kokoro-multi-lang-v1_0/model.onnx"
echo "  - models/tts/kokoro-multi-lang-v1_0/voices.bin"
echo "  - models/tts/kokoro-multi-lang-v1_0/tokens.txt"
echo "  - models/tts/kokoro-multi-lang-v1_0/espeak-ng-data/ (部分文件)"
echo "  - models/tts/kokoro-multi-lang-v1_0/dict/ (部分文件)"
echo "  - models/tts/kokoro-multi-lang-v1_0/lexicon-*.txt"
echo ""
echo "注意："
echo "  1. 基本测试只需要 model.onnx 和 tokens.txt 文件"
echo "  2. voices.bin 用于多说话人支持（可选）"
echo "  3. espeak-ng-data 和 dict 用于文本处理（可选）"
echo "  4. 如果某些文件下载失败，可能不是必需的"
echo ""


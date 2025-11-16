#!/bin/bash
# STT和TTS模型文件下载脚本
# 用于测试和运行STT/TTS服务

set -e

# 保存脚本所在目录和项目根目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 切换到项目根目录（确保所有路径都是相对于项目根目录的）
cd "$PROJECT_ROOT"

# 创建模型目录
mkdir -p models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17
mkdir -p models/speaker
mkdir -p models/tts/kokoro-multi-lang-v1_1

# ==========================================
# 辅助函数：检查文件是否存在且大小合理
# ==========================================
# 参数: $1=文件路径, $2=最小文件大小(字节，可选)
check_file_exists() {
    local file_path="$1"
    local min_size="${2:-0}"
    
    # 确保使用绝对路径或相对于项目根目录的路径
    if [[ "$file_path" != /* ]]; then
        file_path="$PROJECT_ROOT/$file_path"
    fi
    
    if [ -f "$file_path" ]; then
        local file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null || echo "0")
        if [ "$file_size" -ge "$min_size" ]; then
            return 0  # 文件存在且大小合理
        else
            echo "  ⚠️  文件存在但大小异常 ($file_size 字节)，将重新下载"
            return 1  # 文件存在但大小不合理
        fi
    else
        return 1  # 文件不存在
    fi
}

# 下载文件（带存在性检查）
# 参数: $1=文件路径, $2=下载URL, $3=最小文件大小(字节，可选), $4=描述信息
download_file() {
    local file_path="$1"
    local url="$2"
    local min_size="${3:-0}"
    local desc="${4:-$(basename "$file_path")}"
    
    # 确保使用相对于项目根目录的路径
    if [[ "$file_path" != /* ]]; then
        file_path="$PROJECT_ROOT/$file_path"
    fi
    
    # 确保目标目录存在
    local file_dir=$(dirname "$file_path")
    mkdir -p "$file_dir"
    
    if check_file_exists "$file_path" "$min_size"; then
        local file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null || echo "0")
        local size_mb=$(echo "scale=2; $file_size / 1024 / 1024" | bc 2>/dev/null || echo "?")
        echo "  ⏭️  跳过: $desc (已存在, ${size_mb}MB)"
        return 0
    else
        echo "  📥 下载: $desc..."
        curl -L --retry 5 --retry-delay 2 -o "$file_path" "$url" || {
            echo "  ❌ 下载失败: $desc"
            return 1
        }
        
        # 验证下载后的文件大小
        if [ "$min_size" -gt 0 ]; then
            local file_size=$(stat -f%z "$file_path" 2>/dev/null || stat -c%s "$file_path" 2>/dev/null || echo "0")
            if [ "$file_size" -lt "$min_size" ]; then
                echo "  ⚠️  警告: 下载的文件大小异常 ($file_size 字节，预期至少 $min_size 字节)"
            fi
        fi
        return 0
    fi
}

echo "=========================================="
echo "开始下载STT和TTS模型文件"
echo "=========================================="

# ==========================================
# STT模型文件
# ==========================================
echo ""
echo "--- 下载STT模型文件 ---"

# 1. ASR模型文件（int8量化版本，体积更小，约50MB）
download_file \
  "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx" \
  "https://huggingface.co/csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/resolve/main/model.int8.onnx" \
  10000000 \
  "ASR模型 (model.int8.onnx)"

# 2. ASR Tokens文件
download_file \
  "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt" \
  "https://huggingface.co/csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/resolve/main/tokens.txt" \
  100 \
  "ASR Tokens文件 (tokens.txt)"

# 3. 说话人识别模型（可选，用于说话人识别功能，约20MB）
download_file \
  "models/speaker/3dspeaker_speech_campplus_sv_zh_en_16k-common_advanced.onnx" \
  "https://huggingface.co/csukuangfj/speaker-embedding-models/resolve/main/3dspeaker_speech_campplus_sv_zh_en_16k-common_advanced.onnx" \
  1000000 \
  "说话人识别模型"

# ==========================================
# TTS模型文件
# ==========================================
echo ""
echo "--- 下载TTS模型文件 ---"

# 1. TTS模型文件（Kokoro多语言模型，约100MB）
download_file \
  "models/tts/kokoro-multi-lang-v1_1/model.onnx" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/model.onnx" \
  50000000 \
  "TTS模型 (model.onnx)"

# 2. TTS Voices文件（多说话人支持，约10MB）
download_file \
  "models/tts/kokoro-multi-lang-v1_1/voices.bin" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/voices.bin" \
  1000000 \
  "TTS Voices文件 (voices.bin)"

# 3. TTS Tokens文件
download_file \
  "models/tts/kokoro-multi-lang-v1_1/tokens.txt" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/tokens.txt" \
  100 \
  "TTS Tokens文件 (tokens.txt)"

# 4. espeak-ng-data目录（文本处理，必需目录，包含phontab等关键文件）
echo ""
echo "检查: espeak-ng-data目录（必需文件，包含phontab等）..."
mkdir -p models/tts/kokoro-multi-lang-v1_1/espeak-ng-data

# 首先确保目标目录存在
ESPEAK_DATA_DIR="models/tts/kokoro-multi-lang-v1_1/espeak-ng-data"

# 检查关键文件是否已存在
REQUIRED_FILES="phondata phontab phonindex"
ALL_FILES_EXIST=true
for file in $REQUIRED_FILES; do
    if ! check_file_exists "$ESPEAK_DATA_DIR/$file" 500; then
        ALL_FILES_EXIST=false
        break
    fi
done

if [ "$ALL_FILES_EXIST" = true ]; then
    echo "  ⏭️  跳过: espeak-ng-data目录（必需文件已存在）"
else
    # 使用git sparse-checkout下载完整目录（推荐方式）
    if command -v git >/dev/null 2>&1 && [ "$ALL_FILES_EXIST" = false ]; then
        echo "  使用git sparse-checkout下载完整目录..."
        TEMP_DIR="$PROJECT_ROOT/models/tts/kokoro-multi-lang-v1_1/espeak-ng-data-temp"
        TARGET_DIR="$PROJECT_ROOT/$ESPEAK_DATA_DIR"
        mkdir -p "$TEMP_DIR"
        cd "$TEMP_DIR"
        
        git init >/dev/null 2>&1
        git remote add origin https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1 2>/dev/null || git remote set-url origin https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1
        git config core.sparseCheckout true
        echo "espeak-ng-data/*" > .git/info/sparse-checkout
        git pull --depth=1 origin main 2>&1 | grep -E "(Updating|Already|error|fatal)" || true
        
        # 复制文件到目标目录
        if [ -d "espeak-ng-data" ]; then
            cp -r espeak-ng-data/* "$TARGET_DIR/" 2>/dev/null || true
            cd "$PROJECT_ROOT"
            rm -rf "$TEMP_DIR"
            
            # 验证关键文件是否存在
            if [ -f "$TARGET_DIR/phontab" ]; then
                echo "  ✅ espeak-ng-data目录下载完成（包含phontab等必需文件）"
            else
                echo "  ⚠️  git下载完成但phontab文件缺失，使用curl补充下载..."
                # 补充下载关键文件
                for file in $REQUIRED_FILES; do
                    if ! check_file_exists "$TARGET_DIR/$file" 500; then
                        download_file \
                          "$ESPEAK_DATA_DIR/$file" \
                          "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/espeak-ng-data/$file" \
                          500 \
                          "espeak-ng-data/$file"
                    fi
                done
            fi
        else
            echo "  ⚠️  git方式下载失败，使用curl下载必需文件..."
            cd "$PROJECT_ROOT"
            rm -rf "$TEMP_DIR"
            # 下载必需文件
            for file in $REQUIRED_FILES; do
                download_file \
                  "$ESPEAK_DATA_DIR/$file" \
                  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/espeak-ng-data/$file" \
                  500 \
                  "espeak-ng-data/$file"
            done
        fi
    else
        echo "  ⚠️  git未安装或文件缺失，使用curl下载必需文件..."
        # 下载必需文件列表（根据sherpa-onnx的要求）
        for file in $REQUIRED_FILES; do
            download_file \
              "$ESPEAK_DATA_DIR/$file" \
              "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/espeak-ng-data/$file" \
              500 \
              "espeak-ng-data/$file"
        done
        
        # 尝试下载更多可能需要的文件
        ADDITIONAL_FILES="phonindex_zh phonindex_en phonindex_ja phonindex_ko"
        for file in $ADDITIONAL_FILES; do
            if ! check_file_exists "$ESPEAK_DATA_DIR/$file" 500; then
                echo "  尝试下载: $file..."
                curl -L --retry 3 --retry-delay 1 \
                  -o "$ESPEAK_DATA_DIR/$file" \
                  https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/espeak-ng-data/$file 2>/dev/null || true
            fi
        done
    fi
fi

# 验证phontab文件是否存在（这是必需文件）
PHONTAB_PATH="$PROJECT_ROOT/$ESPEAK_DATA_DIR/phontab"
if [ ! -f "$PHONTAB_PATH" ]; then
    echo "❌ 错误: phontab文件下载失败，这是TTS服务的必需文件！"
    echo "   请检查网络连接或手动下载:"
    echo "   curl -L -o $ESPEAK_DATA_DIR/phontab \\"
    echo "     https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/espeak-ng-data/phontab"
    exit 1
else
    echo "✅ phontab文件验证通过"
fi

# 确保在项目根目录
cd "$PROJECT_ROOT"

# 5. dict目录（字典文件，可选）
echo "检查: dict目录..."
mkdir -p models/tts/kokoro-multi-lang-v1_1/dict
download_file \
  "models/tts/kokoro-multi-lang-v1_1/dict/en_dict.txt" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/dict/en_dict.txt" \
  100 \
  "dict/en_dict.txt (可选)" || echo "  提示: en_dict.txt下载失败（可能不需要）"

# 6. Lexicon文件（可选）
echo "检查: Lexicon文件..."
download_file \
  "models/tts/kokoro-multi-lang-v1_1/lexicon-us-en.txt" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/lexicon-us-en.txt" \
  100 \
  "lexicon-us-en.txt (可选)" || echo "  提示: lexicon-us-en.txt下载失败（可能不需要）"

download_file \
  "models/tts/kokoro-multi-lang-v1_1/lexicon-zh.txt" \
  "https://huggingface.co/csukuangfj/kokoro-multi-lang-v1_1/resolve/main/lexicon-zh.txt" \
  100 \
  "lexicon-zh.txt (可选)" || echo "  提示: lexicon-zh.txt下载失败（可能不需要）"

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
echo "  - models/tts/kokoro-multi-lang-v1_1/model.onnx"
echo "  - models/tts/kokoro-multi-lang-v1_1/voices.bin"
echo "  - models/tts/kokoro-multi-lang-v1_1/tokens.txt"
echo "  - models/tts/kokoro-multi-lang-v1_1/espeak-ng-data/ (部分文件)"
echo "  - models/tts/kokoro-multi-lang-v1_1/dict/ (部分文件)"
echo "  - models/tts/kokoro-multi-lang-v1_1/lexicon-*.txt"
echo ""
echo "注意："
echo "  1. STT基本测试只需要 model.int8.onnx 和 tokens.txt 文件"
echo "  2. TTS基本测试需要 model.onnx、tokens.txt 和 espeak-ng-data/phontab 文件"
echo "  3. voices.bin 用于多说话人支持（可选）"
echo "  4. espeak-ng-data 目录是TTS服务的必需目录，包含phontab等关键文件"
echo "  5. dict 和 lexicon 文件用于文本处理（可选，但推荐）"
echo ""
echo "验证文件："
# 验证关键文件
if [ -f "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/model.int8.onnx" ] && \
   [ -f "models/asr/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17/tokens.txt" ]; then
    echo "  ✅ STT模型文件完整"
else
    echo "  ⚠️  STT模型文件不完整"
fi

if [ -f "models/tts/kokoro-multi-lang-v1_1/model.onnx" ] && \
   [ -f "models/tts/kokoro-multi-lang-v1_1/tokens.txt" ] && \
   [ -f "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data/phontab" ]; then
    echo "  ✅ TTS模型文件完整（包含必需的phontab）"
else
    echo "  ⚠️  TTS模型文件不完整（缺少必需文件）"
    if [ ! -f "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data/phontab" ]; then
        echo "    缺少: models/tts/kokoro-multi-lang-v1_1/espeak-ng-data/phontab"
    fi
fi
echo ""

# ==========================================
# 下载 Piper TTS 模型 (vits-piper-zh_CN-huayan-medium)
# ==========================================
echo "📦 下载 Piper TTS 模型 (vits-piper-zh_CN-huayan-medium)..."
echo "----------------------------------------"

# 创建模型目录
mkdir -p models/tts/vits-piper-zh_CN-huayan-medium

# 模型信息
# - 说话人: 华研女声 (huayan)
# - 语言: 中文
# - 质量: medium (中等)
# - 采样率: 22050 Hz
# - 参考: https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium

# 下载 model.onnx
download_file \
  "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx" \
  "https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium/resolve/main/zh_CN-huayan-medium.onnx" \
  20000000 \
  "Piper TTS 模型 (zh_CN-huayan-medium.onnx)"

# 下载 model.onnx.json (配置文件)
download_file \
  "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx.json" \
  "https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium/resolve/main/zh_CN-huayan-medium.onnx.json" \
  1000 \
  "Piper TTS 配置 (zh_CN-huayan-medium.onnx.json)"

# 下载 tokens.txt
download_file \
  "models/tts/vits-piper-zh_CN-huayan-medium/tokens.txt" \
  "https://huggingface.co/csukuangfj/vits-piper-zh_CN-huayan-medium/resolve/main/tokens.txt" \
  1000 \
  "Tokens 文件"

# 下载 espeak-ng-data (如果需要)
if [ ! -d "models/tts/vits-piper-zh_CN-huayan-medium/espeak-ng-data" ]; then
    echo "  ℹ️  Piper 模型使用共享的 espeak-ng-data，从 Kokoro 模型复制..."
    if [ -d "models/tts/kokoro-multi-lang-v1_1/espeak-ng-data" ]; then
        cp -r models/tts/kokoro-multi-lang-v1_1/espeak-ng-data models/tts/vits-piper-zh_CN-huayan-medium/
        echo "  ✅ espeak-ng-data 已复制"
    else
        echo "  ⚠️  找不到 Kokoro 的 espeak-ng-data，请先下载 Kokoro 模型"
    fi
fi

# 检查文件完整性
echo ""
echo "检查 Piper TTS 模型文件..."
if [ -f "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx" ] && \
   [ -f "models/tts/vits-piper-zh_CN-huayan-medium/zh_CN-huayan-medium.onnx.json" ] && \
   [ -f "models/tts/vits-piper-zh_CN-huayan-medium/tokens.txt" ]; then
    echo "  ✅ Piper TTS 模型文件完整"
else
    echo "  ⚠️  Piper TTS 模型文件不完整"
fi
echo ""


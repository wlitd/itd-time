#!/bin/bash

# ================= 配置区域 =================
VERSION="0.1.0"
MODULE_NAME="itd-time"
VERSION_PKG="${MODULE_NAME}/internal/services"

# ============================================

echo "========================================"
echo "开始构建版本: ${VERSION}"
echo "========================================"

# --- 构建逻辑 ---
LDFLAGS="-X ${VERSION_PKG}.Version=${VERSION} -H=windowsgui"

# 执行构建
# 注意：添加了 -clean 确保资源是最新的
wails build -ldflags "${LDFLAGS}" -nsis -clean

if [ $? -eq 0 ]; then
    echo "========================================"
    echo "构建成功！"
    echo "安装包位置: build/bin/"
    ls -lh build/bin/*.exe
    echo "========================================"
else
    echo "构建失败，请检查错误信息。"
fi

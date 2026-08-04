#!/bin/bash

# ================= 配置区域 =================
# 1. 设置版本号
VERSION="0.1.0"

# 2. 设置你的 Go Module 名称 (查看 go.mod 文件第一行)
# 例如: module "my-wails-app" 这里就填 my-wails-app
MODULE_NAME="itd-time"
VERSION_PKG="${MODULE_NAME}/internal/services"

# ================= 构建逻辑 =================

echo "========================================"
echo "开始构建版本: ${VERSION}"
echo "Module 路径: ${VERSION_PKG}"
echo "========================================"

# 构造 ldflags 参数
# -H=windowsgui 用于隐藏控制台窗口
LDFLAGS="-X ${VERSION_PKG}.Version=${VERSION} -H=windowsgui"

# 执行构建
# 关键点：添加 -nsis 参数，自动生成安装包
# -clean 参数可选，用于清理旧的编译缓存，确保版本号最新
wails build -ldflags "${LDFLAGS}" -nsis

if [ $? -eq 0 ]; then
    echo "========================================"
    echo "构建成功！"
    echo "安装包位置: build/bin/"
    # 列出生成的文件，方便查看
    ls -lh build/bin/*.exe
    echo "========================================"
else
    echo "构建失败，请检查错误信息。"
fi

# 暂停以便查看输出
read -p "按回车键退出..."

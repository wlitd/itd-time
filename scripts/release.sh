#!/usr/bin/env bash
# ItdTime 一键发版脚本（Windows Git Bash 运行）
#
# 用法：
#   bash scripts/release.sh            # 按 main.go 当前版本重新打包（同版本重建）
#   bash scripts/release.sh 1.1.0      # 先把版本号同步到五处，再打包发布
#
# 流程：版本同步 → 杀残留进程 → wails3 package（NSIS user）→ 复制产物
#       → 计算 size/sha256-digest 回填 manifest → 校验并打印上传清单
#
# 注意：manifest 的 notes（更新说明）无法自动生成，改版本时请手动编辑
#       build/updates/manifest.json 的 notes 字段。
set -euo pipefail

cd "$(dirname "$0")/.."

# ---------- 工具函数 ----------
die()  { echo "❌ $*" >&2; exit 1; }
info() { echo "▶ $*"; }

command -v wails3  >/dev/null 2>&1 || die "未找到 wails3，请先安装：go install github.com/wailsapp/wails/v3/cmd/wails3@latest"
command -v openssl >/dev/null 2>&1 || die "未找到 openssl（Git Bash 自带，请检查安装）"

MF=build/updates/manifest.json
[[ -f $MF ]] || die "找不到 $MF"

# ---------- 版本号：读取 / 校验 / 五处同步 ----------
CUR_VER=$(sed -nE 's/^var version = "([^"]+)".*/\1/p' main.go)
[[ -n $CUR_VER ]] || die "无法从 main.go 解析当前版本"
VER=${1:-$CUR_VER}
[[ $VER =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "版本号须为 semver（如 1.1.0，不带 v 前缀），got: $VER"

if [[ $VER != $CUR_VER ]]; then
  info "同步版本号 $CUR_VER → $VER（main.go / config.yml / info.json / wails_tools.nsh / manifest.json）"
  sed -i "s|^var version = \"[^\"]*\"|var version = \"$VER\"|" main.go
  sed -i "s|^\(  version: \)\"[^\"]*\"|\1\"$VER\"|" build/config.yml
  sed -i "s|\"file_version\": \"[^\"]*\"|\"file_version\": \"$VER\"|; s|\"ProductVersion\": \"[^\"]*\"|\"ProductVersion\": \"$VER\"|" build/windows/info.json
  sed -i "s|\(!define INFO_PRODUCTVERSION \)\"[^\"]*\"|\1\"$VER\"|" build/windows/nsis/wails_tools.nsh
  NOW_UTC=$(date -u +%Y-%m-%dT%H:%M:%SZ)
  sed -i "s|\"version\": \"[^\"]*\"|\"version\": \"$VER\"|; s|\"name\": \"ItdTime [^\"]*\"|\"name\": \"ItdTime $VER\"|; s|\"url\": \"ItdTime-[^\"]*windows-amd64.exe\"|\"url\": \"ItdTime-$VER-windows-amd64.exe\"|; s|\"publishedAt\": \"[^\"]*\"|\"publishedAt\": \"$NOW_UTC\"|" "$MF"
  echo "⚠️  请手动更新 $MF 的 notes 字段（更新说明会展示在更新弹窗中），可编辑后重跑本脚本。"
else
  info "按当前版本 v$VER 重新打包（版本处保持不变）"
fi

# ---------- 前置清理：关闭占用产物句柄的进程 ----------
info "关闭可能在跑的 ItdTime / wails3 进程"
taskkill //F //IM ItdTime.exe >/dev/null 2>&1 || true
taskkill //F //IM wails3.exe  >/dev/null 2>&1 || true
sleep 1

# ---------- 打包 ----------
info "wails3 task package（NSIS user-scope，约 1~2 分钟）"
PACKAGE_MANAGER=pnpm GOOS=windows wails3 task package INSTALL_SCOPE=user
[[ -f bin/ItdTime.exe ]] || die "构建完成但未找到 bin/ItdTime.exe"

# ---------- 复制上传产物 ----------
ARTIFACT="bin/ItdTime-${VER}-windows-amd64.exe"
info "复制产物 → $ARTIFACT"
cp -f bin/ItdTime.exe "$ARTIFACT"

# ---------- 回填 manifest：size + digest ----------
# size   = 产物精确字节数；digest = sha256 原始 32 字节的 base64（非十六进制）
SIZE=$(stat -c %s "$ARTIFACT")
DIGEST=$(openssl dgst -sha256 -binary "$ARTIFACT" | base64 | tr -d '\n')
info "回填 manifest：size=$SIZE"
# digest 含 + / = 等字符，sed 统一用 | 作分隔符
sed -i "s|\"size\": [0-9]*|\"size\": $SIZE|; s|\"digest\": \"[^\"]*\"|\"digest\": \"$DIGEST\"|" "$MF"

# ---------- 校验 ----------
grep -q "\"version\": \"$VER\"" "$MF" || die "manifest version 回填校验失败"
grep -q "\"url\": \"ItdTime-$VER-windows-amd64.exe\"" "$MF" || die "manifest url 校验失败"
grep -q "\"size\": $SIZE" "$MF" || die "manifest size 校验失败"
grep -qF "\"digest\": \"$DIGEST\"" "$MF" || die "manifest digest 校验失败"
[[ $(head -c 3 "$MF" | od -An -tx1 | tr -d ' \n') != efbbbf ]] || die "manifest.json 混入 UTF-8 BOM，Go 端解析会失败"

# ---------- 汇总 ----------
echo ""
echo "✅ v$VER 打包完成。请上传以下文件到 https://wlitd.com/updates/ ："
echo "   1. $MF            （覆盖线上 manifest）"
echo "   2. $ARTIFACT   （裸 exe，自更新用；不是安装器）"
echo "   可选：bin/ItdTime-amd64-installer.exe 放下载页 / GitHub Release 供新用户安装"

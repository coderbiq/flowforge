# FlowForge CLI — Build & Development

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w
LDFLAGS += -X flowforge/internal/version.injected=$(VERSION)

.PHONY: build test lint clean dev

# ── 开发 ──────────────────────────────────────────────

## 当前平台快速构建
dev:
	go build -ldflags="$(LDFLAGS)" -trimpath -o bin/flowforge ./cmd/flowforge
	rm -rf bin/assets
	cp -R assets bin/assets

## 运行测试
test:
	go test ./internal/...

## lint 检查
lint:
	golangci-lint run ./...

# ── 发布 ──────────────────────────────────────────────

## 交叉编译所有平台
build:
	@bash scripts/build.sh $(VERSION) all

## 打包发布（编译 + 打包 + checksum + 签名）
release:
	@bash scripts/build.sh $(VERSION) all
	@bash scripts/release.sh $(VERSION)

## 清理构建产物
clean:
	rm -rf bin/ dist/

# ── 辅助 ──────────────────────────────────────────────

## 查看帮助
help:
	@echo "Usage:"
	@echo "  make dev          — 构建当前平台（开发用）"
	@echo "  make test         — 运行测试"
	@echo "  make lint         — 代码检查"
	@echo "  make build        — 交叉编译所有平台"
	@echo "  make release      — 打包发布"
	@echo "  make clean        — 清理构建产物"

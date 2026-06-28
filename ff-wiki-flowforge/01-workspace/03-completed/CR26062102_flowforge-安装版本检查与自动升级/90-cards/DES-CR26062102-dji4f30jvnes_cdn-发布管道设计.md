---
id: DES-CR26062102-dji4f30jvnes
title: CDN 发布管道设计
type: design
status: draft
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wt7dgfk
      relation: implements
created: 2026-06-25T12:15:05.996527356Z
updated: 2026-06-25T12:16:14.501896309Z
source: CR26062102
---

# CDN 发布管道设计

## Goal

设计 FlowForge 从源码到分发终端的完整发布管道。使用 GitHub Releases 作为唯一分发源，零配置、零成本、无需域名备案。

## Decision

发布管道由 `.goreleaser.yml` 和 `scripts/release.sh` 共同构成。GoReleaser 编译 6 平台二进制并上传到 GitHub Releases。release.sh 生成 manifest.json（含版本号、平台、SHA256、签名引用），manifest.json 也作为 release asset 上传。Ed25519 签名在 CI 中由 GitHub Actions Secrets 存储的私钥生成，公钥硬编码在 CLI 代码中。

分发架构简化为单线：
```
git tag v1.0.0
  → GoReleaser 编译 6 平台 + 上传 GitHub Releases
  → release.sh 生成 manifest.json + .sig 签名
  → install.sh / upgrade 命令直接从 GitHub Releases 下载
```

URL 模式：`https://github.com/anomalyco/flowforge/releases/download/v1.0.0/flowforge-linux-amd64.tar.gz`

## Rationale

- GitHub Releases 全球 CDN 加速（Fastly），无需自建 CDN
- 无需域名备案、无需云服务账号注册
- GoReleaser 原生支持 GitHub Releases 上传，零额外配置
- 对标同类 CLI 工具（gh、goreleaser 自身）的分发方式

## Constraints

- manifest.json 作为 release asset 独立上传，不内嵌在二进制中
- .sig 签名文件与对应 artifact 同路径发布
- 预发布（prerelease）使用 GitHub 的 prerelease 标记
- 发布由 git tag push 自动触发（GitHub Actions）

## Impact

- 删除：七牛云/阿里云 CDN 上传逻辑
- 简化 `scripts/release.sh`：仅负责 manifest 生成和签名
- 简化 `.github/workflows/release.yml`：GoReleaser + release.sh 两步
- Go 代码中 manifest URL 模板更新为 GitHub Releases 地址

## Verification

- `git tag v1.0.0 && git push --tags` → CI 自动编译上传 → GitHub Releases 页面可见
- manifest.json 可通过 GitHub Releases download URL 访问
- install.sh 下载成功
- 签名验证：硬编码公钥验证 .sig 文件通过

## Follow-up Tasks

- 简化 release.sh
- 创建 .goreleaser.yml（已有）
- 创建 GitHub Actions release workflow（需简化）

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [REQ-CR26062102-djeu2wt7dgfk](REQ-CR26062102-djeu2wt7dgfk_cdn-发布管道与分发架构.md) [requirement] - CDN 发布管道与分发架构

### Incoming

#### implements
- [TASK-CR26062102-i-dji5kn5ajp5i](TASK-CR26062102-i-dji5kn5ajp5i_增强-releasesh-ed25519-签名生成与-cdn.md) [task] - 增强 release.sh — Ed25519 签名生成与 CDN 上传
- [TASK-CR26062102-i-dji5kqtak6ew](TASK-CR26062102-i-dji5kqtak6ew_创建-goreleaseryml-与-git-hub-actions-release.md) [task] - 创建 .goreleaser.yml 与 GitHub Actions release workflow


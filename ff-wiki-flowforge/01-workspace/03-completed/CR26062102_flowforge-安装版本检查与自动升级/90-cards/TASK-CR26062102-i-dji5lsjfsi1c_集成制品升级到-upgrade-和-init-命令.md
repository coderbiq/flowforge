---
id: TASK-CR26062102-i-dji5lsjfsi1c
title: 集成制品升级到 upgrade 和 init 命令 — 备份、验证、报告
type: task
status: done
importance: should
links:
    - target: DES-CR26062102-dji543o8ff5s
      relation: implements
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuos60w
      relation: satisfies
created: 2026-06-25T13:10:52.852970158Z
updated: 2026-06-25T21:56:47.892446699+08:00
source: CR26062102
---

# 集成制品升级到 upgrade 和 init 命令 — 备份、验证、报告

## Goal

将项目制品升级逻辑集成到 flowforge upgrade 和 flowforge init，实现备份、验证和升级报告。

## Inputs

- DES-CR26062102-dji543o8ff5s（manifest 策略设计）
- REQ-CR26062102-djeu2wuos60w（项目制品升级需求）
- TASK-CR26062102-i-dji5li6ksix9（project manifest，I8）
- TASK-CR26062102-i-dji5la2j9llm（upgrade 命令，I7）
- TASK-CR26062102-i-dji5ln67galh（agents_block + 四类处理，I9）
- 现有 internal/command/init.go

## Deliverables

- 修改 internal/command/init.go：init 时调用 GenerateManifest() 创建 manifest.yaml、ApplyAgentsBlock() 部署 AGENTS.md、写入 .version
- 修改 internal/command/upgrade.go：CLI 自更新后调用项目制品升级流程
- 备份逻辑：cp 制品到 .flowforge/backup/<old_version>/
- 升级报告：列出已更新、新增、冲突文件及 AGENTS.md 状态
- 升级后执行 flowforge validate all

## Acceptance

- init 后 manifest.yaml 存在、.version 存在、AGENTS.md 含 FLOWFORGE 区块
- upgrade 无制品变更：报告"已是最新"
- upgrade 有变更：自动处理并报告
- upgrade 有冲突：标记不覆盖提示用户
- 升级失败：backup 目录可用
- validate all 升级后无新增错误

## Out of Scope

- manifest.yaml 读写（I8）
- agents_block 替换逻辑（I9）

## Read Before Work

- internal/command/init.go 现有流程
- internal/command/upgrade.go（I7 产出）
- internal/command/validate.go
- internal/command/assets_deploy.go

## Links

### Outgoing

- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级
- [DES-CR26062102-dji543o8ff5s](DES-CR26062102-dji543o8ff5s_项目制品升级-manifest-结构与升级策略设计.md) [design] - 项目制品升级 manifest 结构与升级策略设计
- [REQ-CR26062102-djeu2wuos60w](REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级


---
id: TASK-CR26062102-a-dji4edi181fj
title: 分析项目制品 manifest 文件范围
type: task
status: done
importance: should
links:
    - target: PROP-CR26062102
      relation: belongs_to
    - target: REQ-CR26062102-djeu2wuos60w
      relation: analyzes
created: 2026-06-25T12:14:10.457243715Z
updated: 2026-06-25T20:47:19.458026546+08:00
source: CR26062102
---

# 分析项目制品 manifest 文件范围

## Goal

确定 .flowforge/manifest.yaml 中需要跟踪的文件范围，以及 manifest 结构与更新策略。

## Inputs

- DES-djdothhisojr 设计卡
- REQ-CR26062102-djeu2wuos60w 需求卡
- AGENTS.md 中的 assets/ 部署目标映射

## Investigation Plan

1. 梳理 assets/ 下每个子目录的文件列表和部署目标
2. 确定哪些文件需要跟踪 vs 不需要
3. 设计 manifest.yaml 的数据结构
4. 评估增量更新 vs 全量比较的可行性

## Expected Outputs

- manifest.yaml 文件范围列表
- manifest.yaml 数据结构设计

## Done When

- 确定 manifest.yaml 跟踪的文件范围
- 确定比较策略（增量 vs 全量）

## Links

### Outgoing

- [REQ-CR26062102-djeu2wuos60w](REQ-CR26062102-djeu2wuos60w_目标项目制品升级.md) [requirement] - 目标项目制品升级
- [PROP-CR26062102](../../../../03-proposal/CR26062102_flowforge-安装版本检查与自动升级.md) [proposal] - flowforge 安装、版本检查与自动升级

### Incoming

- [LOG-CR26062102-dji535saoxml](LOG-CR26062102-dji535saoxml_分析结论-项目制品-manifest-文件范围.md) [log] - 分析结论: 项目制品 manifest 文件范围


# 任务层级管理指南

## 4 层任务结构

每个 proposal 的任务空间遵循固定的 4 层结构：

```
Main Epic: "CRID: Proposal Title"
├── Sub-Epic: "CRID: 分析"           (type: analysis)
│   ├── Parent Task: 分析模块边界
│   │   └── Child Task: 分析子领域
│   └── Task: 独立分析任务
├── Sub-Epic: "CRID: 设计"           (type: design)
│   └── Task: 设计数据模型
└── Sub-Epic: "CRID: 实施"           (type: implementation)
    ├── Parent Task: 实现核心模块
    │   ├── Child Task: DDL + Liquibase
    │   ├── Child Task: DO/PO + Mapper
    │   └── Child Task: Repository
    ├── Parent Task: 实现导入导出
    │   ├── Child Task: 导出逻辑
    │   └── Child Task: 导入逻辑
    └── Task: 独立小任务
```

### 层 1: Main Epic

- **格式**: `CRID: Proposal Title`（如 `CR26060201: Excel 上传/下载配置模块`）
- **标签**: `type:epic, type:main-epic, proposal:<CR-id>`
- **创建方式**: `flowforge task init --proposal <CR-id> "<title>"`
- **重建**: 再次运行 init 需要 `--force true`，会关闭旧 epic 重建

### 层 2: Type Sub-Epic

- **格式**: `CRID: 分析` / `CRID: 设计` / `CRID: 实施`
- **标签**: `type:epic, type:sub-epic, type:<analysis|design|implementation>, proposal:<CR-id>`
- **父子关系**: `--parent <mainEpicId>` 挂在主 epic 下
- **自动创建**: init 时自动创建

### 层 3: Task

- **类型**: 与所在子 epic 类型一致（analysis 子 epic 下只放 analysis 任务）
- **悬挂方式**:
  - 无 `--parent` 参数时 → 挂在对应类型子 epic 下（3 层任务）
  - 有 `--parent <parentTaskId>` 时 → 挂在指定父任务下（4 层任务）

### 层 4: Child Task

- **创建方式**: `flowforge task add --proposal <CR-id> <type> "<title>" --parent <parentTaskId>`
- **用途**: 拆解较大的父任务为可独立执行的子任务
- **约束**: 最多 4 层，不拆到 5 层——如果需要更深层级，说明父任务粒度太大应重新拆分

## 何时用 3 层 vs 4 层

| 场景 | 层级 | 示例 |
|------|------|------|
| 独立小任务 | 3 层（Sub-Epic → Task） | "分析权限码树"、"设计 import/export 入口" |
| 需多步骤的大任务 | 4 层（Sub-Epic → Parent → Child） | "实现核心配置链路" 拆为 DDL/Mapper/Repository/API |
| 探索中发现新任务 | 4 层 | `flowforge task discover --parent <parentId>` 挂在发现它的任务下 |

## 任务结构示例

```
$ flowforge task status --proposal CR-id

GIIS-xxx    [epic]        CR26060201: Excel 上传/下载配置模块
GIIS-xxx.1  [epic]        CR26060201: 分析
GIIS-xxx.2  [task]           分析配置对象、接口与模型边界
GIIS-xxx.3  [epic]        CR26060201: 设计
GIIS-xxx.4  [task]           设计权限码树与接口鉴权
GIIS-xxx.5  [epic]        CR26060201: 实施
GIIS-xxx.6  [task]           实现核心配置链路
GIIS-xxx.7  [task]              DDL + Liquibase
GIIS-xxx.8  [task]              DO/PO + Mapper
GIIS-xxx.9  [task]           实现导入导出
```

## tasks.snapshot.md 格式

快照按类型分组，父子任务通过缩进表达：

```markdown
# Tasks — CR26060201
> Auto-generated at 2026-06-05 10:00:00. Do not edit manually.

## 分析 (2 tasks, 2 done)

| Status | ID | Title |
|--------|----|-------|
| ✅ done | GIIS-a1 | 分析配置对象、接口与模型边界 |
| ✅ done | GIIS-a2 | 分析 apply_to 联动与 RISK_TYPE |

## 实施 (5 tasks, 0 done)

| Status | ID | Title |
|--------|----|-------|
| ⏳ pending | GIIS-i1 | 实现核心配置链路 |
| ⏳ pending | GIIS-c1 | 　DDL + Liquibase |
| ⏳ pending | GIIS-c2 | 　DO/PO + Mapper |
| ⏳ pending | GIIS-c3 | 　Repository |
| ⏳ pending | GIIS-i2 | 实现导入导出 |
```

## CLI 命令对应关系

| 操作 | 命令 | 落位层级 |
|------|------|---------|
| 初始化任务空间 | `flowforge task init --proposal <id> "<title>"` | 创建层 1 (Main) + 层 2 (Sub-Epics) |
| 添加独立任务 | `flowforge task add --proposal <id> <type> "<title>"` | 层 3，挂在 type 子 epic 下 |
| 添加子任务 | `flowforge task add --proposal <id> <type> "<title>" --parent <parentId>` | 层 4，挂在父任务下 |
| 批量添加 | `flowforge task add-tasks --proposal <id> '<json>'` | 层 3 |
| 实施中发现 | `flowforge task discover --proposal <id> <parentId> <type> "<title>"` | 层 4 |

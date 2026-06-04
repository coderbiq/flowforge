# FlowForge v0.7.0 → v0.8.0 升级指南

## 概述

v0.8.0 的核心变化是 **任务系统贯穿分析设计全流程**——将原本游离于任务系统之外的分析设计阶段纳入 task-map.yaml 管理，
实现从分析探索到设计撰写再到实施编码的完整任务链路。

### 关键变化

| 方面 | v0.7.0 | v0.8.0 |
|------|--------|--------|
| 分析设计跟踪 | 无任务跟踪，自由 [探索⇄设计] 循环 | task-map.yaml 从 proposal 创建时就开始工作 |
| 任务类型 | 仅 implementation | analysis / design / implementation 三类 |
| 任务关系 | 仅 dependencies | dependencies + sourceTasks + epic 三种关系 |
| 任务层级 | 扁平列表 | YAML `subtasks` 嵌套层级 |
| 进度展示 | [3/8 任务] | 分析 [2/3] 设计 [1/2] 实施 [0/5] |
| 版本号 | 0.7.0 | 0.8.0 |

## 升级步骤

### 方式一：使用安装脚本（推荐）

```bash
cd /path/to/flowforge
./scripts/install.sh upgrade /path/to/target/project
```

升级脚本会自动：
1. 更新 SKILL 文件、脚本、schema、指南
2. 替换 `<!-- BEGIN FLOWFORGE -->` 块为 v0.8 版本
3. 更新 `.flowforge/config.yaml` 和 `.flowforge/projects/default.yaml` 配置
4. 更新 beads hooks（如有）

### 方式二：手动升级

#### 1. 更新 FlowForge 源文件

```bash
cd /path/to/flowforge
git pull
```

#### 2. 更新 SKILL 文件

```bash
cp src/agents/*/SKILL.md /path/to/target/.agents/
```

#### 3. 更新脚本和指南

```bash
cp -r src/flowforge/scripts/* /path/to/target/.flowforge/scripts/
cp -r src/flowforge/guides/* /path/to/target/.flowforge/guides/
cp -r src/flowforge/schema/* /path/to/target/.flowforge/schema/
```

#### 4. 更新项目配置

合并 `src/flowforge/projects/default.yaml` 中 `rules.design` 部分的新增字段到目标项目的 project 配置：

```yaml
rules:
  design:
    task_rules:
      fields:
        - id
        - title
        - type            # NEW
        - description
        - deliverable
        - dependencies
        - sourceTasks     # NEW
        - epic            # NEW
        - subtasks        # NEW
      time_estimate: "2-5 分钟每个实施任务（analysis/design 任务粒度由 Agent 灵活决定）"
    task_types:            # NEW section
      analysis:
        description: "需求分析、代码探索、可行性研究"
        driven_by: flowforge-design
      design:
        description: "方案设计、架构设计、接口设计"
        driven_by: flowforge-design
      implementation:
        description: "编码实现、测试编写"
        driven_by: flowforge-implement
```

#### 5. 更新 AGENTS.md

确保目标项目 AGENTS.md 中 `<!-- BEGIN FLOWFORGE -->` 块为最新版本。

## 向后兼容

所有新增字段均为可选：
- `type` 默认 `implementation`
- `sourceTasks` / `epic` 默认为空数组 `[]`
- `subtasks` 默认为空

已有 task-map.yaml **无需迁移**即可在 v0.8.0 下正常工作。

## 已有 proposal 的 task-map 升级

如果希望为已有 proposal 补充分析设计追踪，可手动在 task-map.yaml 中：
1. 添加 `type: implementation` 到已有任务
2. 添加 `sourceTasks: []` 和 `epic: []` 字段

不需要重新创建 beads 关联——已有的 `_beadId` 不受影响。

## 验证升级

```bash
# 检查版本
cat /path/to/target/.flowforge/meta.yaml | grep version

# 验证 task-map 一致性
node .flowforge/scripts/task-sync.js /path/to/target <CR-id> --check

# 查看任务状态（应包含 by_type 分组）
node .flowforge/scripts/task-status.js /path/to/target <CR-id>
```

## 回滚

如需回滚到 v0.7.0：

1. 恢复 v0.7.0 的 SKILL 文件
2. 恢复 `src/flowforge/projects/default.yaml` 的 `rules.design` 部分
3. 无需修改 task-map.yaml（新增字段会被 v0.7.0 忽略）

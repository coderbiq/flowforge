# FlowForge CLI 规格与设计 (v4)

## 1. 设计原则

* **辅助而非监管**：CLI 不拦截流程、不设置前置阻断门禁、不强行规定状态转换先后顺序；
* **测试是第一裁判**：代码正确性交由编译器与自动化测试，CLI 专注于工作记忆管理与上下文提取；
* **极简命令树**：聚焦于 `init`、`memory`、`context`、`curate`、`status` 等高价值操作。

---

## 2. 核心命令矩阵

### 2.1 工作记忆与初始化
```bash
# 初始化或刷新全局记忆 (docs/CONTEXT.md)
flowforge memory init

# 初始化特定提案的活笔记 (01-workspace/<proposal_id>/README.md)
flowforge memory init --proposal <id>

# 查看当前活跃提案的记忆摘要
flowforge memory show [--proposal <id>]
```

### 2.2 上下文精准提取 (Slice Extraction)
```bash
# 提取特定切片的最小执行上下文 (仅包含该切片目标、文件接缝、唯一测试命令)
flowforge context slice [--proposal <id>] --slice <n>
```

### 2.3 活文档合流与归档 (Curation)
```bash
# 预览当前提案与领域文档的差分 (Diff)
flowforge curate diff [--proposal <id>]

# 应用合流并归档提案
flowforge curate apply [--proposal <id>]
```

### 2.4 进度总览
```bash
# 查看项目当前所有提案与切片进度状态 (纯展示，无阻断)
flowforge status
```

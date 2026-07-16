# FlowForge v1 分析（历史参考）

> 本文档记录 v1 版本的问题诊断，仅供 v2 设计参考，不指导实现。

## 1. 安装机制问题

v1 通过 bash 脚本 (`scripts/install.sh`) 将文件复制到目标项目。

| 问题 | 严重度 | 说明 |
|------|--------|------|
| 无 CLI 初始化命令 | 严重 | 需手动执行 bash 脚本 |
| 无卸载机制 | 严重 | 必须手动清理多个目录 |
| 文件部署脆弱 | 中 | `cp -r` / `rsync --delete`，无差异化合并 |
| Beads 强耦合 | 中 | 安装流程强制安装 beads |
| 开发/部署耦合 | 中 | `npm link` 绑定开发目录 |
| 两套升级机制 | 低 | `install.sh upgrade` 和 `flowforge upgrade` 职责重叠 |

## 2. 知识组织问题

v1 采用长文档 + 扁平目录：

```
ff-wiki/library/
+-- architecture/    (45 篇文档, 5629 行)
+-- conventions/     (4 篇)
+-- decisions/       (1 篇)
+-- modules/         (0 篇)
```

| 问题 | 量化数据 | 影响 |
|------|----------|------|
| 扁平分类 | `architecture/` 堆积 45 篇、5629 行 | 无法按相关性裁剪 |
| 全量扫描 | `outputLibraryContext()` 遍历全部 .md | 每个 SKILL 收到同样内容 |
| SKILL 自膨胀 | design SKILL 422 行 | 简单任务也加载全文 |
| notes 全文加载 | context 脚本全量输出 | 历史日志占用上下文 |
| 无引用控制写入 | Agent 自由创建文档 | 低价值文档堆积 |

**典型会话上下文消耗**：~24,000 tokens/会话（模型最佳性能区间 <=20K）

## 3. CLI 结构问题

v1 CLI 入口 581 行单文件 switch-case：

| 问题 | 说明 |
|------|------|
| 命令路由 | 内联 switch，无模块化 |
| 子进程委托 | 每次命令 spawnSync 新进程 |
| 配置硬编码 | `.flowforge/config.yaml` 路径写死 |
| 无卸载命令 | 只有安装和升级 |

## 4. 任务管理问题

v1 任务独立于知识系统之外（Beads backend）：

| 问题 | 说明 |
|------|------|
| 任务与知识割裂 | 任务在 `.beads/`，知识在 `ff-wiki/`，无关联 |
| 追溯链断裂 | 无法从任务直接追溯到需求/设计 |
| Beads 强依赖 | 必须安装 bd CLI |

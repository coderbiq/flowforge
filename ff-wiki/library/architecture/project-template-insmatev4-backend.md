---
doc_type: architecture
title: InsmateV4 后端模板配置设计
status: active
created: 2026-06-08T00:00:00Z
updated: 2026-06-08T00:00:00Z
domain:
  scope: system
  type: design
---

# InsmateV4 后端模板配置设计

## 来源

基于 GIIS 项目 `backend.yaml` 的真实实践提炼，去除 GIIS 特定的路径引用，保留 InsmateV4 产品通用的架构模式。

## 核心特征

| 特征 | 说明 |
|------|------|
| 架构 | DDD 分层（application/domain/infrastructure） |
| 语言 | Groovy + Java |
| 探索驱动 | **业务模型驱动**：优先识别 Entity/ValueObject/Service，再确定变更范围 |
| Proposal ID | `BCR{YYMMDD}{NN}`（Backend Change Request） |

## 完整模板内容

```yaml
# InsmateV4 后端项目配置模板
# 使用: 复制此文件到 .flowforge/projects/<your-backend>.yaml
#       修改 wikiRoot 和 srcDirs 为你的项目路径

wikiRoot: ff-wiki-be
srcDirs:
  - src/main/groovy
description: InsmateV4 后端服务，DDD 分层架构（Groovy + Java）
keywords:
  - backend
  - groovy
  - java
  - ddd
  - insmatev4

rules:
  intake:
    strategy: |
      InsmateV4 后端分析 intake 时按以下顺序：
      1. 理解用户核心诉求和背景上下文
      2. 识别涉及的业务领域（如数据服务、文件处理、配置管理等）
      3. 判断变更类型：新增业务能力 / 修改现有模型 / 删除废弃功能
      4. 初步判断涉及的 DDD 层次（application / domain / infrastructure）

  exploration:
    strategy: |
      ## 业务模型驱动的探索策略

      以业务模型为核心逐层展开，顺序不可颠倒：

      1. **领域模型层**：识别涉及的 Entity、ValueObject、DomainService
         - 定位模型文件（通常以实体名命名，如 `DpWorkflow.groovy`）
         - 发现模型间的关系（继承、组合、聚合）
         - 发现 → library/modules/<name>/model/

      2. **应用层**：识别涉及的 Cmd / Qry / DTO / ApplicationService
         - 应用层通常拆分 Cmd（命令）、Qry（查询）、DTO 分别放置
         - 发现 → library/modules/<name>/design/

      3. **基础设施层**：识别涉及的 Repository / Converter / Client
         - Converter 按分层分别放置（infra/converter, app/converter 等）
         - Repository 实现数据访问，Client 封装外部调用
         - 发现 → library/modules/<name>/findings/

      4. 每个发现携带 domain frontmatter，标注 scope/module/type
      5. 模型变更时明确标注变更类型：新建 / 修改 / 删除

  design:
    naming:
      proposal_id: BCR{YYMMDD}{NN}
      exploration_slug: kebab-case
    task_rules:
      fields:
        - id
        - title
        - type
        - description
        - deliverable
        - dependencies
    strategy: |
      ## DDD 分层驱动的设计策略

      1. 优先复用 library 中已有的 DDD 架构模式和决策
      2. 变更涉及多个 DDD 层时，先画出影响范围：
         领域层变更 → 检查应用层是否需同步调整
         应用层变更 → 检查 API 契约是否受影响
         基础设施层变更 → 检查 Converter/Repository 是否需要更新
      3. 每个设计决策记录到 design/ 目录，标注 domain
      4. 涉及模型变更时：
         - 标注变更类型（新建/修改/删除）
         - 记录上下游影响（哪些 Service 依赖此模型）
         - 检查是否影响数据库 schema
      5. 跨模块变更先梳理模块间依赖关系

  implement:
    task_states:
      - pending
      - in_progress
      - done
      - blocked
    notes:
      fields:
        - timestamp
        - status
        - summary
    strategy: |
      ## DDD 契约驱动的实施策略

      1. 遵循现有 DDD 分层架构，不在单个任务中跨层重构
      2. 新增 Entity/ValueObject 时：
         - 遵循现有命名和包结构约定
         - 领域服务名不带 "Domain" 后缀，接口与实现按短命名加包路径区分
      3. 变更后运行相关单元测试，确保各层契约不被破坏
      4. Converter 按职责分开放置（infra/converter, app/converter）
      5. 遇到阻塞问题先通过 flowforge-feedback 结构化记录再创建修复任务

  archive:
    strategy: |
      ## 知识沉淀策略

      1. 优先提取可在其他后端模块复用的 DDD 模式和决策
      2. 模块级知识 → library/modules/<name>/
      3. 系统级架构决策 → library/architecture/ 或 library/decisions/
      4. DDD 通用模式（事件驱动、仓储模式）→ library/conventions/
      5. 同类知识合并而非覆盖，保留历史演进记录

  feedback:
    strategy: |
      ## 反馈回流策略

      1. 测试失败或新认知 → 先分类再回流
      2. DDD 模型相关发现 → 优先回流到 exploration findings/
      3. bug 类发现 → 优先创建修复任务，不阻塞当前提案
      4. knowledge 类发现（Groovy 陷阱、框架限制）→ 标记后等 archive 提取

  library:
    requireReview: false
    autoUpdateHistory: true
    strategy: |
      按 scope（system/module）和 type（design/decision/convention）组织；
      DDD 架构决策和模型设计 → architecture；
      Repository/Converter 等基础设施实现模式 → module design；
      同类知识合并而非覆盖，保留历史演进记录。
```

## 与 GIIS backend.yaml 的差异

| 维度 | GIIS backend.yaml | 模板 |
|------|------------------|------|
| 路径引用 | 硬编码 `1-MainDevelop/saas-service-dev` | 占位 `src/main/groovy` |
| 技术描述 | 含 "GIIS 后端" 特定描述 | 改为 "InsmateV4 后端" |
| 探索策略 | 简单分层 | **逐层展开顺序**：领域模型 → 应用层 → 基础设施 |
| 实施策略 | 基础约束 | **DDD 契约** + 命名约定 + Converter 分层放置规则 |

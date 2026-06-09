---
doc_type: architecture
title: GIIS 后端实际探索模式与配置优化方向
status: active
created: 2026-06-08T00:30:00Z
updated: 2026-06-08T00:30:00Z
domain:
  scope: system
  type: design
---

# GIIS 后端实际探索模式与配置优化方向

## 从实际产物反推的问题

### 问题 1: findings 放错了位置

F-001 到 F-025（25篇）全部在 `architecture/` 下，但其中 90% 是 data-service 模块的模式发现：

```
architecture/F-019-converter-layering-...    ← 其实是 data-service 的 Converter 约定
architecture/F-020-service-naming-...        ← 其实是 data-service 的命名约定
architecture/F-021-application-layer-...     ← 其实是 data-service 的应用层模式
```

**根因**：探索策略只说"系统架构事实→library/architecture/"，没说什么时候该放 modules/。Agent 把所有发现都当成系统级了。

**优化方向**：策略中明确 `scope` 判定规则：
> 发现的内容是否仅适用于当前模块？→ module
> 发现的内容是跨模块通用模式？→ system

### 问题 2: 内容太薄（平均 37 分）

59 篇被标记 "thin"（<100 词）。典型的 finding：
```markdown
# F-020 Service Naming
Service 命名不带 Domain 后缀，通过包路径区分接口和实现。
```
没有证据、没有代码引用、没有影响范围。

**根因**：策略只说"探索代码库"但没有说明**探索到什么程度算完成**。

**优化方向**：策略中加入 completeness checklist：
> 每条 finding 应包含：发现描述 + 代码证据（文件路径） + 影响范围

### 问题 3: 45 篇 findings 零交叉引用

74 篇全部 isolated——没有任何 `related.ref`。

**根因**：策略中没有要求建立引用关系。Agent 写了 finding 但没有链接到相关的 model/decision。

**优化方向**：策略中加入引用约束：
> finding 发现涉及某个模型时 → related.ref 指向 model/<name>.md
> finding 被后续 decision 引用时 → decision 的 related.ref 指向该 finding

### 问题 4: 编号系统混乱

架构级的 findings 用 F-001~F-025，模块级的 findings 另起 F-001~F-028。两套编号独立运行。

**根因**：没有编号约定。Agent 自己编的。

**优化方向**：策略中加入编号规则：
> 架构级 findings: AF-{NNN}
> 模块级 findings: MF-{module}-{NNN}

## 优化后的探索策略

```yaml
exploration:
  strategy: |
    ## 业务模型驱动的探索策略

    以业务模型为核心逐层展开，每层有明确的产出标准：

    1. **领域模型层**（产出：model/ 文档 + findings）
       - 识别 Entity/ValueObject/DomainService
       - 记录模型的核心字段和关系
       - 产出标准：每个模型至少 1 篇 model 文档（含字段说明 + 代码路径）
       - 发现 → library/modules/<name>/model/<EntityName>.md
       - 模块级模式发现 → library/modules/<name>/findings/

    2. **应用层**（产出：design/ 文档 + findings）
       - 识别 Cmd/Qry/DTO/ApplicationService
       - 记录应用层的调用链和数据流
       - 发现 → library/modules/<name>/design/

    3. **基础设施层**（产出：findings）
       - 识别 Repository/Converter/Client
       - 记录技术实现细节
       - 发现 → library/modules/<name>/findings/

    4. **跨模块通用模式** → library/architecture/
       - 仅当发现适用于 2+ 个模块时才放 architecture/

    **每条 finding 必须包含**：
    - 发现描述（what）
    - 代码证据（where，至少 1 个文件路径）
    - 影响范围（which modules/components）
    - related.ref 链接到相关 model/decision

    **scope 判定**：
    - 仅适用当前模块 → scope: module
    - 跨模块通用 → scope: system
```

## 优化后的实施策略

```yaml
implement:
  strategy: |
    ## DDD 契约驱动的实施策略

    1. 实施前检查：涉及模块的 must 级 convention（通过 implement-context）
    2. 不跨层重构，保持 DDD 各层契约
    3. 新增 Entity 需同步创建 model/ 文档（引用 finding 中的证据）
    4. 变更后更新相关 finding 的 status（active → validated）
    5. 建立交叉引用：
       - 代码变更 → 更新 finding 的代码证据
       - 模型变更 → related.ref 指向 model/ 文档
```

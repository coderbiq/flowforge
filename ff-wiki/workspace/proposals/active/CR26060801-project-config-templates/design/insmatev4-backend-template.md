---
doc_type: design
title: InsmateV4 Backend 完整模板
status: draft
created: 2026-06-08T02:00:00Z
updated: 2026-06-08T02:00:00Z
domain:
  scope: system
  type: design
---

# InsmateV4 Backend 完整模板

## 模板内容

完整的 `src/flowforge/project-templates/insmatev4-backend.yaml`：

```yaml
# InsmateV4 Backend 项目配置模板
# 使用: flowforge template apply insmatev4-backend --as <id>
# 或直接: 修改 wikiRoot 和 srcDirs 为你的项目路径

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
  patterns:
    architecture:
      backend:
        - "DDD 分层: application/domain/infrastructure/client/models"
        - "应用层: Cmd(命令) + Qry(查询) + DTO 分开放置, Qry 自带校验"
        - "领域层: Entity/ValueObject/DomainService, Service 不带 Domain 后缀, 接口与实现按短命名加包路径区分"
        - "基础设施层: Repository(继承 AbstractBaseRepository) + Converter(按层分开放置: infra/converter, app/converter)"
        - "客户端层: 对外暴露的 API 和 DTO"
        - "模型层: DO/PO 定义, 继承 BasePoModel"
    anti-patterns:
      backend:
        - "不要用 Controller-Service-DAO 模式, 项目使用 DDD Cmd/Qry/DomainService 分层"
        - "不要新建模块时不参考现有模块的包结构"
        - "不要在 DomainService 中直接操作数据库, 通过 Repository"
        - "不要 new Dto() 后逐字段 set, 使用 Converter 接口(@Mapper)"
    must-follow:
      - "library/conventions/data-dictionary-sync.md"

  toolbox:
    backend:
      utils:
        - "import com.bytesforce.pub.MyStringUtils / MyDateUtils / MyNumberUtils / MyJsonUtils / MyCollectionUtils"
        - "import com.bytesforce.pub.MyBeanUtils / MyEncryptionUtils / MyExcelUtils"
        - "import com.bytesforce.pub.ParamAssert (参数校验断言)"
        - "import com.bytesforce.pub.PageHelper / SqlBuilder (分页/SQL构建)"
        - "优先使用 My* 系列工具, 避免重复实现"
      base-classes:
        - "Controller 继承: com.bytesforce.client.support.web.AbstractBaseAPI"
        - "  → runAndWrap(Closure) 统一返回值包装"
        - "  → downloadFile() 文件下载"
        - "Repository 继承: com.bytesforce.client.support.rds.AbstractBaseRepository<T extends BasePoModel>"
        - "  → create/batchCreate/update/delete 自动填充审计字段+主键+版本+租户"
        - "  → selectList/selectOne 自动租户隔离"
        - "模型基类: BasePoModel (提供 id/createdAt/createdBy/lastUpdatedAt/lastUpdatedBy)"
      converters:
        - "使用 MapStruct @Mapper(config = MapperConfigWithoutMetaClass) 接口"
        - "interface XxxAppConverter { XxxDTO convertDO2DTO(XxxDO); XxxParam convertCmd2Param(XxxCmd); }"
        - "转换器按层分开放置: infra/converter/ (DO↔PO), app/converter/ (DTO↔DO)"
        - "反例: 不要在 Service 中 new 对象手动赋值"
    discover-during-exploration:
      - "每个 analysis 任务探索时, 按 toolbox 目录检查项目工具"
      - "发现新工具 → 记录为 finding → archive 时更新此段"

  exploration:
    strategy: |
      ## 业务模型驱动的探索策略

      以业务模型为核心逐层展开，每层有明确的产出标准：

      1. **领域模型层**（产出: model/ 文档 + findings）
         - 识别 Entity/ValueObject/DomainService（通常以实体名命名）
         - 记录模型核心字段、关系和代码路径
         - 产出标准: 每个模型 ≥1 篇 model 文档（字段说明 + 代码路径）
         - 模块级发现 → library/modules/<name>/model/<EntityName>.md
         - 模块级模式发现 → library/modules/<name>/findings/

      2. **应用层**（产出: design/ 文档 + findings）
         - 识别 Cmd/Qry/DTO/ApplicationService
         - 记录调用链和数据流
         - 发现 → library/modules/<name>/design/

      3. **基础设施层**（产出: findings）
         - 识别 Repository/Converter/Client
         - 记录技术实现细节和 Converter 映射关系
         - 发现 → library/modules/<name>/findings/

      4. **跨模块通用模式** → library/architecture/（仅当适用于 2+ 模块时）

      **每条 finding 必须包含**:
      - 发现描述（what）
      - 代码证据（where，≥1 个文件路径）
      - 影响范围（which modules/components）
      - domain.scope 判定: 仅适用当前模块 → module, 跨模块 → system

  design:
    naming:
      proposal_id: BCR{YYMMDD}{NN}
    task_rules:
      fields:
        - id, title, type, description, deliverable, dependencies
    strategy: |
      ## DDD 分层设计策略

      1. 设计方案前, 先读取 rules.patterns 锁定项目架构模式
      2. 优先复用 library 中已有的 DDD 架构决策
      3. 变更涉及多个 DDD 层时, 先画出影响范围:
         - 领域层变更 → 检查应用层是否需同步调整
         - 应用层变更 → 检查 API 契约是否受影响
         - 基础设施层变更 → 检查 Converter/Repository 是否需更新
      4. 涉及模型变更时标注变更类型(新建/修改/删除)和上下游影响
      5. 跨模块变更先梳理模块间依赖关系

  implement:
    strategy: |
      ## DDD 契约实施策略

      1. 实施前检查 toolbox: 优先使用项目已有工具和基类
      2. 不跨层重构, 保持 DDD 各层契约
      3. 新增 Entity 需同步创建 model/ 文档
      4. 变更后运行相关单元测试
      5. 遇到阻塞 → flowforge-feedback 结构化记录

  archive:
    strategy: |
      优先提取可跨模块复用的 DDD 模式和决策;
      模块级知识 → library/modules/<name>/;
      系统级决策 → library/architecture/ 或 library/decisions/;
      DDD 通用模式 → library/conventions/。

  library:
    requireReview: false
    autoUpdateHistory: true
```

## 与 GIIS backend.yaml 的关键差异

| 维度 | GIIS 原版 | 模板 |
|------|----------|------|
| patterns | 无 | 架构声明 + 反例 + must-follow |
| toolbox | 无 | My*工具 + AbstractBaseAPI/Repository + Converter模式 |
| exploration | 简单分层描述 | 逐层产出标准 + finding 强制字段 |
| design | 基础约束 | 模式锁定 + DDD 影响分析 |

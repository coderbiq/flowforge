---
doc_type: design
title: design-context.js 增强方案
status: active
created: 2026-06-08T02:00:00Z
updated: 2026-06-08T02:00:00Z
domain:
  scope: system
  type: design
  importance: should
  maturity: growing
---

# design-context.js 增强方案

## 新增输出段

在现有 `## Domain 分类指引` 和 `## Library Context` 之后，新增两段：

### ## Project Patterns

解析 `rules.patterns`，结构化输出：

```
## Project Patterns

### Architecture (backend)
- DDD 分层: application/domain/infrastructure/client/models
- 应用层: Cmd(命令) + Qry(查询) + DTO 分开放置

### Must Follow
- library/conventions/data-dictionary-sync.md

### Anti-Patterns
- ❌ 不要用 Controller-Service-DAO 模式
- ❌ 不要 new Dto() 后逐字段 set
```

实现：
```javascript
const patterns = r.patterns;
if (patterns) {
  console.log('## Project Patterns\n');
  if (patterns.architecture) {
    const archKey = patterns.architecture.backend ? 'backend' : 'frontend';
    console.log(`### Architecture (${archKey})\n`);
    for (const line of patterns.architecture[archKey]) {
      console.log(`- ${line}`);
    }
    console.log('');
  }
  if (patterns['must-follow']) { ... }
  if (patterns['anti-patterns']) {
    for (const key of Object.keys(patterns['anti-patterns'])) {
      for (const line of patterns['anti-patterns'][key]) {
        console.log(`- ❌ ${line}`);
      }
    }
  }
}
```

### ## Implementation Toolbox

解析 `rules.toolbox`，结构化输出：

```
## Implementation Toolbox

### Backend Utils
- MyStringUtils / MyDateUtils / MyNumberUtils / ...
- ParamAssert (参数校验)

### Backend Base Classes
- Controller → AbstractBaseAPI
- Repository → AbstractBaseRepository<T>

### Backend Converters
- @Mapper 接口: XxxAppConverter
- 反例: ❌ new Dto() 手动赋值

### Frontend Components
- BaseTable / ModalService / NText / NSelect / ...
- 反例: ❌ antd Table 直接写
```

## template 字段解析

```javascript
function loadProjectConfig(root, ref) {
  // 加载 config.yaml 中声明的 project
  const templateName = ref.template;
  
  if (templateName) {
    const templatePath = path.join(root, '.flowforge', 'project-templates', `${templateName}.yaml`);
    if (fs.existsSync(templatePath)) {
      const template = yaml.load(fs.readFileSync(templatePath, 'utf8'));
      // 模板提供 rules
      projectConfig.rules = template.rules;
    }
  }
  
  // 实例 config 覆盖 wikiRoot/srcDirs
  if (ref.config) {
    const instanceConfig = yaml.load(fs.readFileSync(path.join(root, '.flowforge', ref.config), 'utf8'));
    projectConfig.wikiRoot = instanceConfig.wikiRoot || projectConfig.wikiRoot;
    projectConfig.srcDirs = instanceConfig.srcDirs || projectConfig.srcDirs;
  }
  
  // 无模板 → 使用 default.yaml rules
  if (!projectConfig.rules) {
    projectConfig.rules = defaultConfig.rules;
  }
}
```

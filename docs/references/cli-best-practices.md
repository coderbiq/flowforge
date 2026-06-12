# CLI 工具设计最佳实践

> **历史参考文档** | 2026-06-12
>
> 本文档是 FlowForge v2 早期调研阶段的产物，基于 Node.js/npm 生态。
> 最终方案已改为 **Go 独立二进制**，参见 [CLI 架构设计](../cli-design.md) §9。
>
> 本文档保留作为行业实践参考（命令设计、配置分层、init/upgrade 流程等通用模式仍有价值）。

本文档整理 npm 全局 CLI 工具的设计模式和最佳实践。

---

## 1. npm Global CLI 标准模式

### 1.1 全局命令注册（`bin` field）

所有主流 CLI 工具都通过 `package.json` 的 `bin` 字段注册全局可执行命令：

```json
{
  "name": "flowforge",
  "bin": {
    "flowforge": "./bin/cli.js"
  }
}
```

**关键规则**（来源：[npm 文档](https://docs.npmjs.com/files/package.json/)）：
- `bin` 值可以是对象（自定义命令名）或字符串（与包名相同）
- 入口文件**必须以 `#!/usr/bin/env node` 开头**，否则系统不会用 Node 执行
- 全局安装时 npm 会在 `/usr/local/bin/` 创建符号链接
- 局部安装时链接到 `./node_modules/.bin/`

### 1.2 行业实践总结

| 工具 | 命令名 | 入口文件 | 框架 |
|------|--------|----------|------|
| Vue CLI | `vue` | `packages/@vue/cli/bin/vue.js` | 自建/Commander |
| Create React App | `create-react-app` | `packages/create-react-app/index.js` | 自建 |
| oclif | `mycli` | `bin/run.js` + `bin/dev.js` | oclif 框架 |
| ESLint | `eslint` | `bin/eslint.js` | 自建 |
| Husky | `husky` | `bin.js` | 零依赖 |

---

## 2. CLI 框架选型对比

### 2.1 三种主流方案

| 维度 | **Commander.js** | **oclif** | **零依赖自建** |
|------|-----------------|-----------|---------------|
| 适用场景 | 3-15 个子命令 | 30+ 子命令/团队维护 | 1-5 个命令/单文件 |
| 学习成本 | 低（API 稳定十年） | 中（需要了解插件体系） | 低 |
| 子命令路由 | 手动注册 | 文件系统即路由 | switch 语句 |
| 自动更新 | 需自行实现 | `@oclif/plugin-update` 开箱即用 | 需自行实现 |
| 插件体系 | 无 | 原生支持 | 无 |
| TypeScript | 需要额外配置 | 生成器自带 | 需要额外配置 |

### 2.2 推荐选择

> "Use commander if you are building a real tool with 3 to 15 subcommands, want TypeScript without the framework overhead."
> 
> — [DEV Community](https://dev.to/thegdsks/building-a-production-typescript-cli-in-2026-oclif-vs-commander-vs-custom-9ah), 2026

**FlowForge v2 预计有 8-12 个命令**（`init`、`upgrade`、`uninstall`、`task`、`card`、`context`、`validate`、`config`），**Commander.js 是最佳平衡选择**。

---

## 3. `npx` vs 全局安装

### 3.1 行业趋势

行业趋势是**优先支持 `npx`**，同时提供 `npm install -g` 选项：

- **CRA**：明确建议 `npx create-react-app` 而非全局安装
- **ESLint**：文档主推 `npx eslint --init` 和 `npx eslint src/`
- **Husky**：推荐 `npx husky init`

### 3.2 权衡

| 方式 | 优点 | 缺点 |
|------|------|------|
| **npx** | 始终使用最新版本、不污染全局 PATH、用完即弃 | 每次执行有解析开销 |
| **全局安装** | 无解析开销、命令可用性好 | 需要手动升级、占用全局空间 |

**对于 FlowForge**：
- 高频使用的工具（每天数十次）建议全局安装以减少 npx 解析开销
- 提供 `npx @flowforge/cli init` 作为一次性使用的替代方案

---

## 4. `init` 命令的最佳实践

### 4.1 四种初始化模式

| 模式 | 代表工具 | 流程 | 适用场景 |
|------|----------|------|----------|
| **交互式向导** | `eslint --init` / `vue create` | 提问 → 配置生成 → 依赖安装 | 需要用户做选择 |
| **模板克隆** | `create-react-app` / `create-nx-workspace` | 模板复制 → 依赖安装 → 初始化脚本 | 标准化项目结构 |
| **单命令设置** | `husky init` | 一步到位，无交互 | 简单配置 |
| **配置文件生成** | `nest new` / `vue create --preset` | 预设配置 → 脚手架生成 | 有预设/模板 |

### 4.2 ESLint --init 流程

ESLint 的 `--init`（实际是 `npm init @eslint/config`）是一个典型的**无头交互向导**：

1. 检查当前目录是否已有 `package.json`（没有则先 `npm init`）
2. 交互式问题链：
   - 如何使用 ESLint？（检查语法/代码风格/两者都要）
   - 模块类型（ESM/CJS）
   - 框架（React/Vue/None）
   - TypeScript 是否使用
   - 代码在哪里运行（Browser/Node）
3. 根据回答计算最终配置
4. 安装所需插件
5. 生成 `eslint.config.js`

**关键设计**：问题之间有关联依赖，后一个问题取决于前置答案。这是决策树而非线性问答。

### 4.3 Husky init 流程

Husky v9 的 init 是**最简洁的 init 设计**：

1. 创建 `.husky/` 目录
2. 创建 `.husky/pre-commit` 文件（内容：`npm test`）
3. 更新 `package.json` 的 `prepare` 脚本为 `"husky"`

```bash
# husky init 的等效手动操作
npm pkg set scripts.prepare="husky"
npm run prepare
```

**关键洞察**：Husky 的 `init` 只做两件事——创建目录结构 + 配置 npm script。**不要把 init 做得太复杂。**

### 4.4 FlowForge init 设计

```bash
flowforge init [path] [--yes] [--template <name>]
```

**执行流程**：
1. 参数解析（项目名、选项）
2. 目标目录检查（是否为空、是否已有 FlowForge）
3. 交互式配置收集（或 --yes 跳过）
4. 文件生成（.flowforge/ 目录结构）
5. 安装 SKILL 文件
6. 更新 AGENTS.md
7. 输出初始化摘要

---

## 5. Upgrade 命令的实现方式

### 5.1 版本检测机制

行业主流的版本自检测模式（参考 [oclif/plugin-update](https://github.com/oclif/plugin-update)）：

```
┌─────────────────────────────┐
│  CLI 启动时                  │
│  检查版本缓存文件             │
│  ~/.cache/flowforge/lastrun  │
├─────────────────────────────┤
│  距上次检查 > debounce 天?   │
│  稳定版: 14天 / 其他: 1天    │
│  (可自定义配置)              │
├─────────┬───────────────────┤
│  是      │  否               │
├─────────┘                   │
│  spawn 子进程执行 update     │
│  --autoupdate (后台不阻塞)    │
└─────────────────────────────┘
```

**核心代码模式**（取自 oclif/plugin-update [init hook](https://github.com/oclif/plugin-update/blob/master/src/hooks/init.ts)）：

```typescript
async function autoupdateNeeded(config): Promise<boolean> {
  const mtime = await getMtime(autoupdatefile)
  const debounce = config.pjson.oclif?.update?.autoupdate?.debounce ?? 14
  mtime.setHours(mtime.getHours() + debounce * 24)
  return mtime < new Date()
}
```

### 5.2 版本比较策略

使用 `semver` 库进行版本比较：

```typescript
import semver from 'semver'

// 检查是否有新版本
const latest = await getLatestVersion('flowforge')
const current = require('./package.json').version

if (semver.gt(latest, current)) {
  // 提示升级
}
```

**最佳实践**：
- 使用 `npm view flowforge version` 获取最新版（通过 registry API）
- 或使用 `npm-check-updates` 风格检查
- 不要在 CLI 启动时阻塞式检查，使用**后台进程 + 缓存**

### 5.3 版本发布通道

oclif 的多通道发布模式：

| 通道 | 版本号示例 | 更新频率 | 用途 |
|------|-----------|----------|------|
| `stable` | `1.2.3` | 按 tag 发布 | 正式发布 |
| `beta` | `1.3.0-beta.1` | 每个 commit | 预览版 |
| `dev` | `1.3.0-dev.20260601` | 每个 commit | 日构建 |

用户可通过 `flowforge update beta` 切换通道。

### 5.4 迁移脚本

ESLint 的配置迁移器 (`@eslint/migrate-config`) 展示了最佳实践：

```typescript
// 迁移器模式
class MigrationRunner {
  private migrations: Migration[] = []
  
  register(fromVersion: string, toVersion: string, migrate: Function) {
    this.migrations.push({ fromVersion, toVersion, migrate })
  }
  
  async run(currentVersion: string) {
    const sorted = this.migrations
      .filter(m => semver.gt(m.toVersion, currentVersion))
      .sort((a, b) => semver.compare(a.toVersion, b.toVersion))
    
    for (const migration of sorted) {
      await migration.migrate()
    }
  }
}
```

**关键设计点**：
- **增量迁移**：每次只处理一个版本的变更，不跳版本
- **幂等性**：多次执行迁移结果相同
- **回滚能力**：迁移前备份原配置
- **兼容层**：ESLint 提供 `FlatCompat` 让新旧格式共存

---

## 6. 多项目支持与配置分层

### 6.1 配置分层模型（行业标准）

所有成熟 CLI 工具都遵循**多层叠加**配置架构（来源：[Better CLI](https://bettercli.org/design/configuration/)、[cosmiconfig](https://www.npmjs.com/package/cosmiconfig)）：

```
优先级（高 → 低）
┌─────────────────────────┐
│  ① CLI 命令行参数        │  --flowforge-dir ./projects
├─────────────────────────┤
│  ② 环境变量              │  FLOWFORGE_DIR=./projects
├─────────────────────────┤
│  ③ 项目级配置            │  .flowforgerc / flowforge.config.js
├─────────────────────────┤
│  ④ 用户级配置            │  ~/.config/flowforge/config.json
├─────────────────────────┤
│  ⑤ 系统级配置            │  /etc/flowforge/config.json
├─────────────────────────┤
│  ⑥ 内置默认值            │  代码中的 defaults 对象
└─────────────────────────┘
```

### 6.2 使用 cosmiconfig 管理配置

[cosmiconfig](https://github.com/davidtheclark/cosmiconfig) 被 ESLint、Prettier、Babel 等广泛使用：

```typescript
import { cosmiconfig } from 'cosmiconfig'

const explorer = cosmiconfig('flowforge', {
  searchPlaces: [
    'package.json',           // 读取 flowforge 字段
    '.flowforgerc',
    '.flowforgerc.json',
    '.flowforgerc.yaml',
    '.flowforgerc.yml',
    '.flowforgerc.js',
    'flowforge.config.js',    // ESM/CJS 自动适配
    'flowforge.config.ts',
  ],
})

// 从 cwd 向上搜索至 home 目录
const result = await explorer.search()

// 也可以直接加载指定文件
const config = await explorer.load('/path/to/flowforge.config.js')
```

**优势**：
- 自动从 cwd 向上搜索至 home 目录
- 支持多种格式（JSON/YAML/JS/TS）
- 可以嵌入 `package.json` 的 `flowforge` 字段
- 支持全局配置目录 (`~/.config/flowforge/`)

### 6.3 项目 vs 全局配置边界

FlowForge 的配置分层建议：

| 配置项 | 位置 | 举例 |
|--------|------|------|
| 项目工作流定义 | `.flowforge/config.yaml` | `cards:`, `proposals:` |
| 用户偏好 | `~/.config/flowforge/config.json` | `editor: "cursor"` |
| 缓存路径 | `~/.cache/flowforge/` | 临时文件、版本检查缓存 |
| 项目内初始化标记 | `.flowforge/config.yaml` | `version: "2.0.0"` |

---

## 7. 输出美化与用户体验

### 7.1 推荐工具

| 工具 | 用途 | 示例 |
|------|------|------|
| **chalk** | 终端色彩 | 成功绿色、错误红色、警告黄色 |
| **ora** | 进度条/Spinner | 长时间操作的等待指示 |
| **@clack/prompts** | 交互式向导 | 多选、单选、输入框 |
| **cli-table3** | 表格输出 | 卡片列表、任务状态 |

### 7.2 输出规范

```bash
# 成功
✓ FlowForge initialized successfully

# 错误
✗ Error: Config file not found
  → Run `flowforge init` to create a new config

# 警告
⚠ Warning: Deprecated API usage
  → See https://flowforge.dev/migration for details

# 信息
ℹ FlowForge v2.1.0 is available (current: v2.0.0)
  → Run `flowforge upgrade` to update
```

### 7.3 静默模式

支持 `--quiet` / `-q` 标志，仅输出错误：

```bash
$ flowforge card list --quiet
# 仅输出卡片 ID 列表，无表格、无色彩
```

---

## 8. 测试策略

### 8.1 单元测试

```typescript
// commands/card/create.test.ts
import { describe, it, expect } from 'vitest'
import { generateCardFilename } from './create'

describe('generateCardFilename', () => {
  it('generates filename with type, id, title', () => {
    const result = generateCardFilename({
      type: 'requirement',
      id: 'REQ-260612-001',
      title: '支持 CLI 全局安装',
    })
    expect(result).toBe('REQ-260612-001_支持CLI全局安装.md')
  })

  it('includes deps in filename', () => {
    const result = generateCardFilename({
      type: 'decision',
      id: 'DEC-260612-001',
      title: '使用 Commander.js',
      deps: ['REQ-260612-001', 'CONV-001'],
    })
    expect(result).toBe('DEC-260612-001_使用Commanderjs__REQ-260612-001+CONV-001.md')
  })
})
```

### 8.2 集成测试

```typescript
// e2e/init.test.ts
import { describe, it, expect } from 'vitest'
import { execSync } from 'child_process'
import { mkdtempSync, existsSync } from 'fs'

describe('flowforge init', () => {
  it('creates .flowforge directory structure', () => {
    const tmpDir = mkdtempSync('/tmp/flowforge-test-')
    execSync(`flowforge init ${tmpDir} --yes`)
    
    expect(existsSync(`${tmpDir}/.flowforge/config.yaml`)).toBe(true)
    expect(existsSync(`${tmpDir}/.flowforge/workspace`)).toBe(true)
    expect(existsSync(`${tmpDir}/.flowforge/library`)).toBe(true)
  })
})
```

---

## 9. FlowForge v2 技术选型总结

| 组件 | 选择 | 理由 |
|------|------|------|
| CLI 框架 | **Commander.js** | 8-12 个子命令，API 稳定 |
| 配置搜索 | **cosmiconfig** | ESLint/Prettier 验证过的方案 |
| 用户交互 | **@clack/prompts** | 轻量、美观的交互式向导 |
| 版本管理 | **semver** | 业界标准 |
| 模板引擎 | **ejs** | 根据用户选择渲染不同配置 |
| 输出美化 | **chalk** + **ora** | 色彩 + 进度条 |
| 测试框架 | **vitest** | 快速、TypeScript 原生支持 |

---

## 参考资料

### 工具文档

- [Commander.js](https://github.com/tj/commander.js)
- [oclif](https://oclif.io/)
- [cosmiconfig](https://github.com/cosmiconfig/cosmiconfig)
- [@clack/prompts](https://github.com/natemoo-re/clack)

### 最佳实践

- [Better CLI Design](https://bettercli.org/design/)
- [Node.js CLI Best Practices](https://github.com/lirantal/nodejs-cli-apps-best-practices)
- [npm bin field documentation](https://docs.npmjs.com/cli/v10/configuring-npm/package-json#bin)

### 行业案例

- [Vue CLI](https://github.com/vuejs/vue-cli)
- [ESLint](https://github.com/eslint/eslint)
- [Husky](https://github.com/typicode/husky)
- [Create React App](https://github.com/facebook/create-react-app)

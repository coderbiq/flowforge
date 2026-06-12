# 开发指南

> 版本：v2.0.0-alpha | 最后更新：2026-06-12

## 开发环境

```bash
# 克隆仓库
git clone <repo-url>
cd flowforge

# 安装依赖
npm install

# 运行测试
npm test

# 本地开发
npm link
flowforge --version
```

## 代码规范

### 基本原则

1. **简洁优先** — 代码和文档都应简洁明了
2. **类型安全** — 不使用 `as any`、`@ts-ignore` 等类型抑制
3. **测试覆盖** — 修改代码后运行测试验证
4. **最小变更** — Bug 修复只修复问题，不重构

### 文件组织

- 每个命令一个文件（`src/cli/commands/`）
- 业务逻辑与 CLI 路由分离（`src/cli/lib/`）
- 卡片操作通过 CardStore 统一接口

### 提交规范

```
add: 新功能
fix: Bug 修复
refactor: 重构（不改变功能）
update: 更新依赖或文档
remove: 删除功能
```

## 测试

```bash
# 运行全部测试
npm test

# 运行单个测试文件
npx vitest src/cli/commands/__tests__/init.test.js

# 运行并监听
npx vitest --watch
```

## 发布流程

```bash
# 1. 更新版本号
npm version minor  # 或 patch / major

# 2. 运行测试
npm test

# 3. 发布到 npm
npm publish --access public

# 4. 创建 GitHub Release
gh release create v$(node -p "require('./package.json').version")
```

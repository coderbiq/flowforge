# 功能模块目录

本目录存储项目各个功能模块的设计文档，记录每个模块的**当前状态**和**演进历史**。

## 设计理念

采用**聚合视角**，每个模块一个文档目录，归档提案时更新模块文档。

## 目录结构

```
modules/
├── README.md                    # 本文件
├── INDEX.md                     # 模块索引
├── auth/                        # 认证模块（示例）
│   ├── README.md               # 模块概览
│   ├── design.md               # 当前设计
│   ├── api.md                  # API 文档（可选）
│   └── history.md              # 演进历史
├── user/                        # 用户模块（示例）
└── _template/                   # 模块模板
    ├── README.md
    ├── design.md
    ├── history.md
    └── api.md
```

## 与提案的关系

```
提案 CR25051701 (新增认证模块)
  ↓ 归档时
modules/auth/
  ├── README.md    ← 新建：模块概览
  ├── design.md    ← 新建：初始设计
  └── history.md   ← 记录：CR25051701 创建了此模块

提案 CR25051802 (增强认证：添加 SSO)
  ↓ 归档时
modules/auth/
  ├── README.md    ← 更新：新增 SSO 说明
  ├── design.md    ← 更新：新增 SSO 设计
  └── history.md   ← 追加：CR25051802 添加了 SSO
```

## 模块文档组成

| 文件 | 用途 | 更新时机 |
|------|------|---------|
| `README.md` | 模块概览：职责、边界、关键概念 | 归档提案时 |
| `design.md` | 当前设计：架构、数据流、关键决策 | 归档提案时 |
| `api.md` | API 文档：接口定义、请求响应格式 | 归档涉及 API 变更的提案 |
| `history.md` | 演进历史：提案记录、重大变更 | 归档提案时 |

## 归档更新规则

### 新增模块

提案创建新模块时：
1. 创建 `modules/{module-name}/` 目录
2. 创建 `README.md`（模块概览）
3. 创建 `design.md`（初始设计）
4. 创建 `history.md`（记录创建）

### 修改模块

提案修改现有模块时：
1. 更新 `README.md`（如有边界变化）
2. 更新 `design.md`（设计变更部分）
3. 如有 API 变更，更新 `api.md`
4. 追加 `history.md`（记录变更）

### 删除模块

提案删除模块时：
1. 标记 `README.md` 为 DEPRECATED
2. 记录 `history.md`（删除原因、替代方案）
3. 可选：移动到 `modules/_archived/`

## 命名规范

模块目录名使用 **kebab-case**：
- `auth` - 认证模块
- `user-management` - 用户管理模块
- `data-sync` - 数据同步模块

## 相关命令

归档提案时自动触发模块文档更新：

```
/propose:archive CR{编号}
  → 检查提案影响的模块
  → 更新/创建相应模块文档
```

## 模块索引

参见 [INDEX.md](./INDEX.md) 查看所有模块的索引。

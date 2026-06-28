# CLI 设计原则

> 版本：v2.0.0-alpha | 基于代码审计 + git log + 现有设计文档整理 | 2026-06-28

## 1. 来源

本文从三个来源提取 FlowForge CLI 的设计原则：

| 来源 | 内容 |
|------|------|
| 现有设计文档 | `docs/architecture.md`、`docs/cli-design.md`、`docs/card-architecture-invariants.md` |
| git log | 关键 commit 的设计意图（`e95dba4` stdin, `965892d` invariants 等） |
| 代码审计 | 对 `internal/command/` 下全部命令实现的逐行审计 |

---

## 2. 五大原则

### P1: CLI 唯一入口

**Agent 不直接操作文件。所有卡片读写、配置修改必须通过 CLI 命令完成。**

出处：`docs/architecture.md` §2.2 — "Agent 通过 CLI 命令读写卡片，不直接操作文件"；`docs/card-architecture-invariants.md` §2.2 — "Agent 不手写内部卡片链接"。

Agent 代码中禁止的行为：
- 禁止 `os.WriteFile` 写 `.md` 卡片文件
- 禁止手写 frontmatter、wikilink、内部导航链接
- 禁止直接编辑 `config.yaml`、`manifest.yaml`

### P2: 命令不接收文件路径作为内容输入

**CLI 命令的参数是卡片 ID、标题、类型等语义标识符，不应该是文件路径。需要传递多行内容时，通过 `--body -` + heredoc 走 stdin。**

出处：commit `e95dba4`（"stdin body"）明确将 `--body -` 作为多行内容的标准通道。

正确模式：
```bash
flowforge card create --type design --title "xxx" --body - <<'EOF'
...
EOF
```

错误模式：
```bash
# 先写临时文件，再传文件路径——违反原则
cat > /tmp/card.yaml <<'EOF' && flowforge card batch /tmp/card.yaml
```

### P3: ID 寻址

**命令引用卡片时使用卡片 ID（如 `REQ-2x9k3m00-3x8m2n1q`），不使用文件路径。解析路径是 CLI 内部行为，对 Agent 不可见。**

出处：`docs/card-architecture-invariants.md` §2.1 — frontmatter `links` 存储 ID，非路径；`docs/cli-design.md` §8.1 — 所有 `card` 子命令接受 `<card-id>`。

CLI 内部的 `CardStore.FindCardPath` 负责 ID→路径解析，Agent 不需要知道文件存在哪里。

### P4: 原子操作

**一次命令调用完成一个完整的业务动作，命令内部自行管理状态一致性。Agent 不编排多步文件操作。**

出处：commit `2abe979`（"card batch @ref two-phase creation"）——batch 内跨引用和索引写入由命令内部的两阶段提交保证，Agent 不需要关心中间状态。

反例：Agent 不能先创建卡片文件、再手动写链接、再手动调 index rebuild——这应由单个命令内部完成。

### P5: 输出结构化

**命令输出应支持 `-o json` 供 Agent 解析。人类可读文本输出是默认格式但非唯一格式。**

出处：commit `e95dba4`（"--output json works globally for all write commands"）。

所有 write 命令（`card create`、`card update`、`library import`、`structure add`、`task create`、`log create`）必须支持 `-o json`。

---

## 3. 审计发现

### 3.1 命令逐一审计

对 `internal/command/` 下全部命令的文件操作审计结果：

| 命令 | 文件输入方式 | 违反原则? | 说明 |
|------|-------------|-----------|------|
| `card create` | `--body -` (stdin) | **否** | stdin 支持 |
| `card update` | `--body -` (stdin) | **否** | stdin 支持 |
| `card read` | 卡片 ID | **否** | ID 寻址 |
| `card delete` | 卡片 ID | **否** | ID 寻址 |
| `card list` | 无文件输入 | **否** | — |
| `card search` | 无文件输入 | **否** | — |
| `card related` | 卡片 ID | **否** | ID 寻址 |
| `card link/unlink` | 卡片 ID | **否** | ID 寻址 |
| `card refresh` | 卡片 ID | **否** | ID 寻址 |
| **`card batch`** | `-` (stdin) / `<file>` | **否** | 支持 stdin（`card batch -`）和文件路径 |
| `library import` | `--body -` (stdin) | **否** | stdin 支持 |
| `library facets/classify/suggest` | 卡片 ID / 无文件 | **否** | — |
| `library promote` | 卡片 ID | **否** | ID 寻址 |
| `structure add/remove/list` | 卡片 ID | **否** | ID 寻址 |
| `task create` | `--body` (flag) | **否** | flag 传值 |
| `task claim/done/block/status/sub` | 卡片 ID | **否** | ID 寻址 |
| `log create` | `--summary` (flag) | **否** | flag 传值 |
| `init` | 目录路径 | **否** | 项目初始化，合理的文件系统操作 |
| `upgrade` | 无用户文件输入 | **否** | 内部读 manifest |
| `uninstall` | 目录路径 | **否** | 项目卸载 |
| `validate` | 卡片 ID / `.md` 路径 | **否** | 两种方式均可用 |
| `index` | 卡片 ID | **否** | ID 寻址 |
| `context` | 无文件输入 | **否** | — |
| `project` | 项目 ID | **否** | ID 寻址 |
| `proposal` | 提案 ID | **否** | ID 寻址 |
| `skill update` | 无用户文件输入 | **否** | 内部复制 assets |
| `config` | 无文件输入 | **否** | — |

### 3.2 card batch (已修正)

`internal/command/batch.go:90-94`：
```go
if batchFile == "-" {
    data, err = io.ReadAll(cmd.InOrStdin())
    ...
}
```

| 项目 | 详情 |
|------|------|
| 合规状态 | ✅ 符合 P2 — 支持 stdin |
| 使用方式 | `card batch -` 从 stdin 读取，`card batch <file>` 向后兼容 |

### 3.3 合规命令的实现模式

以下命令已遵守 P2，可作为实现参考：

```go
// internal/command/card.go:1351
func readBody(body string) (string, error) {
    if body == "-" {
        data, err := io.ReadAll(os.Stdin)
        ...
    }
    return body, nil
}
```

`library import` 复用同一 `readBody` 函数。`log create` 走 `--summary` flag（内容短，无需 heredoc）。

---

## 4. 设计决策记录

| 决策 | commit | 说明 |
|------|--------|------|
| `--body -` 作为 stdin 通道 | `e95dba4` | 取代 shell 转义方案，heredoc 直接传入 |
| `-o json` 全局支持 | `e95dba4` | 所有 write 命令统一输出格式 |
| batch 命令独立为子命令 | `e95dba4` | `card batch` 而非 `card create --batch` |
| batch 内部两阶段提交 | `2abe979` | @ref 跨引用在 batch 内部原子解析 |
| SKILL 规则禁止手写正文链接 | `965892d` | Framework 保证导航一致性 |
| 卡片 ID 前缀从 ROOT 迁移至 PROP | `29b1a68` | proposal 根卡类型清晰化 |

---

## 5. 修正记录

### 5.1 card batch 已修正 (2026-06-28)

`card batch` 现在支持 `-` 作为 stdin 输入：

```bash
flowforge card batch - -o json <<'EOF'
cards:
  - type: structure
    title: "..."
EOF
```

实现细节：
- `batch.go:90-94`: `-` 参数时通过 `cmd.InOrStdin()` 读取 stdin
- 文件路径模式保留（向后兼容）
- 读取模式与 `card create --body -` 一致

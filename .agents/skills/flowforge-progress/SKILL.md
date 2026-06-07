---
name: flowforge-progress
description: |
  FlowForge 进度索引记录器。在任何 proposal 的工作单元刚刚完成后激活，
  将本次进展总结为一句话写入 meta.latest_progress 并刷新 INDEX.md。

  必须在以下信号出现后立即激活：
  - 刚刚通过 `flowforge task` 完成了任务操作（done/claim/block/cancel）
  - 刚刚在 notes.md 追加了实施日志
  - 刚刚创建、归档或移动了 proposal 目录
  - 刚刚完成 design 章节并 commit 到 proposal

  不要在以下情况激活：
  - 仅在阅读 proposal 内容、未做修改时
  - 用户在询问状态、未触发实际变更时
  - 仅修改 library/ 中的归档文档（非 workspace 进行中状态）
---

# FlowForge Progress

负责在工作单元完成后将进展总结写入 meta 并刷新 INDEX.md。

## 工作流

```
识别变更 → 一句话总结 → 运行脚本 → 确认结果
```

---

### 阶段 1：识别变更

回顾本次会话中刚才对哪个 proposal 做了什么。如果不确定，检查最近修改的 `meta.yaml` 或 `notes.md`。

---

### 阶段 2：一句话总结

用一句话总结（≤80 字）本次进展，原则：

- **描述完成的事**，不描述意图——"完成 token 中间件" 而非 "推进认证模块"
- **包含分类型的量化进度**——"分析 [2/3] 设计 [1/2] 实施 [0/5]" 或 "完成认证分析"
- **可验证**——读者能判断这句话是否属实

| 场景 | ✅ 好 | ❌ 差 |
|:--|:--|:--|
| 分析任务推进 | 完成认证需求分析 [分析 2/3] | 推进了分析工作 |
| 设计任务推进 | 完成认证模块设计 [设计 1/2] | 做了一些设计 |
| 实施任务推进 | 完成 JWT 中间件 [实施 3/8] | 做了一些进展 |
| 全部分析设计完成 | 分析设计完成，共 3 分析 + 2 设计任务 | 设计完成，拆分为 8 个任务 |
| 全部完成 | 实施完成，等待归档 | 任务进行中 |
| 归档 | 知识沉淀至 library/modules/auth/ | 归档完毕 |
| 状态变更 | 设计评审通过，进入实施 | 状态有变化 |

---

### 阶段 3：运行脚本

```
flowforge update-progress <proposal完整路径> "<总结>"
```

脚本会自动：
1. 将 `latest_progress` 写入 `meta.yaml`
2. 更新 `updated_at` 时间戳
3. 重新生成 `workspace/proposals/INDEX.md`

---

### 阶段 4：确认结果

检查脚本输出是否成功。如果失败：
- 找不到 proposal → 确认路径是否正确
- meta.yaml 格式错误 → 手动修正后重试
- 不确定本次变更内容 → 询问用户而非编造

---

## 需要的脚本

| 脚本 | 用途 |
|:--|:--|
| `flowforge update-progress <proposal路径> "<总结>"` | 写 meta.latest_progress + 更新 updated_at + 重建 INDEX.md |
| `flowforge refresh-index [项目根路径]` | 仅重建 INDEX.md（可独立运行） |

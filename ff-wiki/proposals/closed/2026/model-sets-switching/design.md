# model-sets-switching — Design

> Requirements: [named-model-sets](requirements.md#named-model-sets) ·
> [quick-switch-command](requirements.md#quick-switch-command) ·
> [switch-atomicity](requirements.md#switch-atomicity)

## 结构

### config 声明层(进 git 的模板形态 + per-machine 实际文件)

```yaml
agents:
  models: { tool-capable-read-only: cpa/deepseek-v4.1-flash }
  models_by_name: { flowforge-reviewer: cpa/glm-5.3, ... }   # 基础层(现有,不动)
  model_sets:                                                  # 新增:命名 overlay 集
    offpeak:
      models_by_name:
        flowforge-analyst: cpa/mimo-2.6-pro
        flowforge-architect: cpa/mimo-2.6-pro
        flowforge-planner: cpa/mimo-2.6-pro
        flowforge-reviewer: cpa/mimo-2.6-pro
```

Go 结构(`internal/config`):

```go
type ModelSetConfig struct {
    Models             map[string]string            `yaml:"models,omitempty"`
    ModelOverrides     map[string]string            `yaml:"models_by_name,omitempty"`
    ModelHostOverrides map[string]map[string]string `yaml:"models_by_host,omitempty"`
}
// AgentsConfig 增:
ModelSets map[string]ModelSetConfig `yaml:"model_sets,omitempty"`
```

**命名**:用 `model_sets` 而非 profiles — `ModelProfile` 已被能力档位枚举
占用(high-capability / tool-capable / tool-capable-read-only),避免概念撞车。

### 激活指针(独立状态文件,不写回 config.yaml)

`.flowforge/model-set.active` — 内容为一行 set 名(或不存在 = 基础层)。

理由:config.yaml 无写回路径(Load 只读),引入 yaml 写回会丢注释、
重排格式;激活态本质是机器本地运行状态,独立状态文件与
"状态走文件不走对话"原则一致。该文件随 `.flowforge/` 保持 gitignored。

## 合并语义(overlay)

`ApplyModelSet(base AgentsConfig, set ModelSetConfig) AgentsConfig`
在**部署路径入口**对内存中的 cfg.Agents 三层做合并后再走现有
解析链(`models_by_host > models_by_name > models`,agents.go:275-282),
下游与校验零改动:

- `models` / `models_by_name`:set 键覆盖基础层同键
- `models_by_host`:按 host 外层合并,set 内 host 的键覆盖基础层同 host 同键

校验:`validateAgentsConfig` 对每个 set 递归执行与顶层完全相同的键/值校验
(错误路径带 `agents.model_sets.<name>.` 前缀)。

## CLI(`internal/command/modelset.go`)

```
flowforge model-set list          # default(基础层) + 各 set,标 ACTIVE
flowforge model-set use <name>    # = apply-set(name): 指针→deploy→失败回滚
flowforge model-set show [name]   # 合并后有效模型表(标注来源: base/set)
```

`use` 的原子性顺序:

1. 校验 `<name>` 存在且 set 内容合法(失败:无任何变更)
2. 写 `.flowforge/model-set.active`
3. 调用现有 deploy 路径(按 active set overlay 后编译全部 host)
4. deploy 报错 → 删除/回滚指针文件 → 以原错误退出(exit code 非 0)

回滚语义:指针文件原样恢复(原本不存在则删除),不依赖 deploy 幂等。

## 验证策略

- `internal/config` 单测:overlay 合并(三层覆盖/host 深合并)、指针读写、
  set 校验错误路径(未知角色/非法值/set 名不存在)
- `internal/command` 单测:`use` 成功路径(指针+deploy 调用)与失败路径
  (deploy 错误后指针恢复原值)、`list`/`show` 输出
- 手工验收:本仓双 set 切换实测(glm-5.3 ↔ 模拟 offpeak),确认部署产物
  frontmatter model 变化与回滚行为

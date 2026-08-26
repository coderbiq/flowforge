# Mixed-source import walk-through

## Input

`legacy-proposal.md`, headed “Deployment”, says the old installer copies Skills but does not report drift. `product-prd.md`, headed “Operations”, asks an operator to identify whether installed Skills match the running binary without overwriting custom files. `completed-change.md`, headed “Configuration”, says `init` already accepts a configured `docs_dir`. The two sources disagree on whether a manifest is required.

## Expected Import hand-off (Chinese)

| 分类 | 保留含义与来源定位 | 下一责任人 |
| --- | --- | --- |
| 来源事实 | 旧安装器会复制 Skill，但不报告漂移（`legacy-proposal.md#Deployment`）。 | Solution Design |
| 需求候选 | 操作人员需要只读地确认项目 Skill 是否匹配运行二进制，且自定义文件不被覆盖（`product-prd.md#Operations`）。 | Align |
| 设计决定 | 不使用安装 manifest；运行二进制嵌入资产是比较基准（`legacy-proposal.md#Deployment`）。 | Solution Design |
| 交付/验证证据 | `init` 已支持 configured `docs_dir`（`completed-change.md#Configuration`）。 | Evidence/source record |
| 未知/冲突 | manifest 是否必须与嵌入资产比较基准相冲突；影响资源验证设计。 | Solution Design |

交接的目标语言是中文。Align 接受需求候选后，Solution Design 复核设计决定和冲突；两份 authority 发布后运行 `flowforge check --dir <feature-dir> --strict`。随后 Plan 展示 ticket 的 title、Delivery 与真实 DAG edge，等待用户接受才创建 issue。没有内容被机械改写成新的 authority，也没有自动创建 ticket 或状态迁移。

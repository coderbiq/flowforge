# F-011 canonical corpus 补充应按 archive target 类型过滤同类最终文档

- Status: validated
- Source: `scripts/lib/flowforge.js`

## Statement

在自动补充 canonical corpus 时，应以 proposal 的 archive target 类型为过滤轴，只增补同类型的 workspace 最终文档，而不是把整个 workspace 的所有最终文档都作为 baseline。

## Why it matters

这种筛选方式能让 canonical corpus 保持相关性和可读性。proposal 的 baseline 应该尽量贴近本次变更目标，否则阅读者会被大量无关文档干扰，baseline 也会失去“本次变化参照系”的意义。

## References

- [flowforge.js](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/scripts/lib/flowforge.js)

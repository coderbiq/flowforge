# 调研：用工具链机器化强制设计系统落地（组件边界 / 色值 token / 对比度 / 视觉回归 / antd 官方治理）

> 本文即任务指定的 `docs/research/2026-09-29-ui-enforcement-toolchain-research.md` 内容（由研究子代理产出，落库时以此为准）。
>
> - 日期：2026-09-29 · 角色(flowforge-research) · 方法：官方文档/官方仓库/官方注册表取证 + 项目内文档对照，未改任何代码。
> - 证据标注约定：**[一手-官方文档]** 官方站点逐字核对；**[一手-官方仓库]** GitHub 官方 org 仓库核对；**[一手-registry]** npm registry 元数据；**[二手]** 搜索摘要/第三方转述，未逐字核对；**[未核实]** 查不到官方出处，不得当作事实引用；**[推断]** 研究者基于一手机制的推理，明确标出。

## 版本快照（2026-09-29 查询 npm registry `latest`）

| 工具 | latest 版本 | 备注 |
|---|---|---|
| eslint | 10.11.0 | 任务假设 ESLint 9 flat config；flat config 语法在 9/10 一致，规则选项无差异。注意 `eslint-plugin-react@7.37.5` peerDependencies 上限为 eslint `^9.7`（[一手-registry]），升 10 需先确认兼容 |
| stylelint | 17.15.0 | 任务假设 16；本文引用的规则族在 16/17 均存在。`stylelint-declaration-strict-value@1.12.1` peer 为 stylelint `>=16 <=17` |
| antd | 6.6.5 | 项目在 v5；**v6 已于 2025-11-22 发布**（[二手] 官方公告 issue `ant-design/ant-design#55804`；npm latest 6.6.5 为 [一手-registry]）。v6 与本调研相关的点：CSS variables、zeroRuntime（见 §2/§5） |
| vitest | 5.0.2 | registry 元数据含 `@vitest/browser-playwright` 可选依赖（browser 模式存在的直接证据） |
| playwright | 1.63.0 | `@axe-core/playwright@4.13.0` peer `playwright-core >= 1.0.0`，兼容 |
| vitest-axe | 0.1.0 | registry 发布时间戳推断为 **2022-10 后无更新**，选型需谨慎 |
| eslint-plugin-boundaries | 7.2.0 | 官方文档站 jsboundaries.dev |
| dependency-cruiser | 18.4.0 | |
| eslint-plugin-no-inline-styles | 1.0.5 | 发布时间戳推断为 2019-11 后无更新 |
| @storybook/test-runner | 0.24.5 | peer storybook 10/11 |
| lost-pixel | 3.22.0 | 发布时间戳推断为 **2024-11 后无更新** |
| reg-suit | 0.14.5 | 2025-08 仍在发布 |
| @ant-design/pro-components | 2.8.10 | peer antd `^4.24.15 \|\| ^5.11.2`，**尚不支持 antd v6**（[一手-registry]） |
| eslint-plugin-react | 7.37.5 | |

## 背景

Tangram V2 前端为 pnpm monorepo（`frontend/apps/*` + `frontend/packages/*`），技术栈 React + TypeScript + Vite + antd v5 + @ant-design/x + vitest + react-router；共享设施 `@tangram/ui`（AppLayout / TangramThemeProvider / ReadableRegistry 族）；成文视觉规范 `docs/805-frontend-design-guide.md` 已有 L/T/I/F/S 五域条款（L1 禁 app 内手拼 Layout/Header、T1 禁内联硬编码色、T2 token 唯一权威源、I1 默认层禁非可读信息等），且第六章已固化 17 条 store-mate 反例（含「三套 Header 并存」「约 100 内联色」等，均带 file:line）。当前强制手段仅为 coding agent 自觉遵守 + review，近期实际复发了：绕过壳自建第二套 Header、暗底暗字对比度不可见、英文/内部 ID 直出。本调研回答「如何用 lint/构建/测试/CI 机器化强制」，不做设计裁决。

---

## 一、组件边界强制（禁止在 app 层直接手写 UI / import 底层组件）

### 结论

1. **最短路径是 ESLint 核心规则 `no-restricted-imports` + flat config 按文件分组**：
   - `paths` 支持精确封禁模块，且可用 `importNames` 按导出名封禁（如从 `antd` 只封 `Layout/Header/Sider/Menu/Footer`），或用 `allowImportNames` 反向白名单（「除列名外全封」，不能与 `importNames` 同用）**[一手-官方文档]**（eslint.org/docs/latest/rules/no-restricted-imports）。
   - `patterns` 支持 gitignore 风格 `group` 或 `regex`（二选一），带 `caseSensitive`、`allowTypeImports`、按组生效的 `importNamePattern/allowImportNamePattern`；`!` 取反须放数组末尾（顺序敏感，父目录排除后子路径无法恢复）**[一手-官方文档]**。
   - 「apps 禁、封装层放行」的机制基础：flat config 的 `files` glob 分组，后匹配的配置对象覆盖前者（官方明示 "later objects overriding previous objects when there is a conflict"，并给出 base 规则 + 子集覆盖的推荐用法）**[一手-官方文档]**（eslint.org/docs/latest/use/configure/configuration-files）。即：`{files:['apps/**'], rules:{'no-restricted-imports':[...]}}` + `{files:['packages/tangram-ui/**'], rules:{'no-restricted-imports':'off'}}`。**[推断]** 该组合写法本身是标准机制组合，无官方逐字示例。
2. **分层治理可用 `eslint-plugin-boundaries` 7.x**：`boundaries/elements` 按 `type+pattern` 定义层（如 `app`、`shared-ui`），`boundaries/dependencies` 规则以 `default:"disallow"` + 白名单策略表达「谁能 import 谁」，并且**官方 README 的示例正是外部包白名单**："Only controllers may use the 'axios' package"（`{from:{element:{type:"!controller"}}, disallow:{to:{module:{origin:"external", source:"axios"}}}}`）——与「只有 @tangram/ui 可 import antd」完全同构 **[一手-官方仓库]**（github.com/javierbrea/eslint-plugin-boundaries；完整文档在 jsboundaries.dev）。
3. **构建期独立门禁可用 `dependency-cruiser` 18.x**：`forbidden/allowed/required` 三类规则，`from`/`to` 匹配（`to.path` 是**正则**而非 glob，支持 `$0` 捕获组回填），`severity:"error"` 时 `err` reporter 返回非零退出码可断 CI **[一手-官方仓库]**（github.com/sverweij/dependency-cruiser/blob/main/doc/rules-reference.md）。它不依赖 ESLint 进程，还能查循环依赖/orphan/可达性，适合做 `pnpm build` 前的独立校验步骤。
4. **`react/forbid-elements` 能拦 JSX 元素与 `React.createElement` 调用**：文档原文 "This rule checks all JSX elements and `React.createElement` calls"——可用来封原生 `<header>`/`<nav>`/`<div>` 等宿主标签，从 DOM 层面阻断「自建第二套 Header」**[一手-官方仓库]**（jsx-eslint/eslint-plugin-react docs/rules/forbid-elements.md）。
5. **`react/forbid-dom-props` 只作用于宿主 DOM 节点，不管自定义组件**：文档原文 "This rule only applies to DOM Nodes (e.g. `<div />`) and not Components"；文档示例本身就是 `forbid: ["style"]`（`<div style={{color:'red'}} />` 报错）——即**能禁原生元素上的 style 属性，但 antd 组件上的 style 拦不到** **[一手-官方仓库]**（docs/rules/forbid-dom-props.md）。新版还支持 `disallowedFor`（限定标签）与 `disallowedValues`（限定值）。
6. `eslint-plugin-no-inline-styles`（1.0.5，2019-11 后未更新 [一手-registry]）：整禁 JSX style 属性，不能按值区分（color vs margin 一律报）；npm 页提示带引号的 style key 可能漏报 **[二手]**。可作为零成本兜底，不宜作主策略。
7. **先例核查**：未查到 antd 或 shadcn 官方「只许 import 封装层」的成文治理先例 **[未核实]**。shadcn 的官方模型与此相反：其文档自述 "This is not a component library. It is how you build your component library."，组件源码经 CLI 拷入仓库（open code / 分发平台），并明说传统库的问题恰是 "you end up wrapping library components" **[一手-官方文档]**（ui.shadcn.com/docs）。机制层面最接近的官方同构先例即 boundaries README 的 axios 白名单示例。**社区存在大量 no-restricted-imports 治理实践，但无权威成文规范可引 [二手/未核实]。**

### 拦得住 / 拦不住矩阵

| 形态 | no-restricted-imports | boundaries | dependency-cruiser | forbid-elements | forbid-dom-props |
|---|---|---|---|---|---|
| app 内 `import { Layout } from 'antd'` | ✅（paths+importNames / allowImportNames 白名单） | ✅（外部包规则） | ✅（from/to 正则） | — | — |
| `<header>` 手写原生标签 | — | — | — | ✅（含 createElement） | — |
| 原生元素 `style` 属性 | — | — | — | — | ✅（forbid ['style']） |
| antd 组件上的 `style` prop | — | — | — | — | ❌（不管自定义组件） |
| 动态 `import(变量)` / 运行时拼模块名 | ❌ | ❌ | ⚠️（会报 unresolved，部分兜底） | — | — |
| `document.createElement('header')` | — | — | — | ❌ | — |
| cssinjs 运行时注入样式 / `dangerouslySetInnerHTML` | ❌ | ❌ | ❌ | ❌ | ❌ |

### 证据表

| 结论 | 来源 | 状态 |
|---|---|---|
| paths/patterns/importNames/allowImportNames/importNamePattern 选项与限制 | https://eslint.org/docs/latest/rules/no-restricted-imports | 一手-官方文档 |
| flat config `files` 分组 + 后者覆盖前者 | https://eslint.org/docs/latest/use/configure/configuration-files | 一手-官方文档 |
| elements/dependencies/default disallow/外部包白名单示例 | https://github.com/javierbrea/eslint-plugin-boundaries （文档站 https://www.jsboundaries.dev/ ） | 一手-官方仓库 |
| forbidden 规则 / to.path 正则 / severity error 非零退出 | https://github.com/sverweij/dependency-cruiser/blob/main/doc/rules-reference.md | 一手-官方仓库 |
| forbid-elements 覆盖 JSX+createElement | https://github.com/jsx-eslint/eslint-plugin-react/blob/master/docs/rules/forbid-elements.md | 一手-官方仓库 |
| forbid-dom-props 仅 DOM 节点、可禁 style | https://github.com/jsx-eslint/eslint-plugin-react/blob/master/docs/rules/forbid-dom-props.md | 一手-官方仓库 |
| shadcn「代码所有权」定位（反例参照） | https://ui.shadcn.com/docs | 一手-官方文档 |
| antd/shadcn 官方「只许 import 封装层」先例 | 未查到 | 未核实 |

---

## 二、色值 / token 强制（禁 hex 与裸色名、强制 var(--token)）

### 结论

1. **stylelint 静态 CSS 面（16/17 规则族一致）**：
   - `color-no-hex`：禁一切 hex（3/4/6/8 位，非法 hex 也报）；`ignoreFunctions:["var"]` 可放行 `var(--foo, #fff)` 回退值 **[一手-官方文档]**（stylelint.io/user-guide/rules/list/color-no-hex）。
   - `color-named: "never"`：禁具名色（black/white…），`var(--white)` 不报；有 `ignoreFunctions/ignoreProperties` **[一手-官方文档]**（…/rules/list/color-named）。
   - `declaration-property-value-disallowed-list`：按属性封值（如 `"/^color|background|border/": ["/^rgb/", "/^hsl/"]`）。**官方明示关键限制：正则匹配整个声明值**——`border: 1px solid #fff` 这类 shorthand 中的分量匹配不可靠 **[一手-官方文档]**（…/rules/list/declaration-property-value-disallowed-list）。
   - `stylelint-declaration-strict-value` 1.12.1（peer stylelint 16–17 [一手-registry]）：对指定属性**强制 `var()`/`$sass`/`@less`/`@value`/函数或关键字白名单**，选项含 `ignoreValues/ignoreVariables/ignoreFunctions/expandShorthand/disableFix`——这是「强制走 token」最贴切的现成规则 **[一手-官方仓库]**（github.com/AndyOGo/stylelint-declaration-strict-value）。
2. **antd v5 cssinjs 场景下 stylelint 的真实边界**：antd v5 组件样式由 `@ant-design/cssinjs` 在**运行时**注入（仓库自述 "Component level cssinjs solution used in ant.design"，`StyleProvider` 的 `container` 默认 `document.head`）**[一手-官方仓库]**（github.com/ant-design/cssinjs）。因此：
   - stylelint 的覆盖面 = 自有 `.css/.less` 源文件 + 配 `customSyntax` 后的 styled-components 风格 CSS-in-JS 模板串（styled-components 官方文档推荐 stylelint 15+ 搭配 `postcss-styled-syntax`；旧 processor 已弃维，且对插值只能「good guess」）**[一手-官方文档]**（styled-components.com/docs/tooling）。
   - antd 运行时注入的样式**磁盘上不存在源文件，stylelint 覆盖率为 0**（这是「按定义」成立的边界：stylelint 是源文件 linter）**[推断，基于上述定义]**。
   - **TSX 内联 `style={{...}}` 是 JS 对象不是 CSS 文本，不属于任何 CSS customSyntax 的解析范围**——styled-components 文档描述的提取也仅针对模板串 **[推断，边界由机制决定]**。
3. **TSX 内联色值的拦截选项**（本项目 store-mate 约 100 处内联色的主战场）：
   - `react/forbid-dom-props` forbid `['style']`：只拦宿主元素（见 §1）**[一手-官方仓库]**。
   - `no-restricted-syntax` 自定义 AST 选择器：官方支持任意 AST selector（字符串/对象混用，可带 `message`）**[一手-官方文档]**（eslint.org/docs/latest/rules/no-restricted-syntax）。「匹配 `JSXAttribute[name.name='style']` 下值形如 `#/^#[0-9a-fA-F]{3,8}$/` 的 `Literal`」这类窄拦截在机制上完全可行，但**选择器写法无官方成文示例**，属团队自研/社区实践 **[未核实（先例）]**。注意它拦不了从常量/变量取值的颜色。
   - `eslint-plugin-no-inline-styles`：整禁 style 属性（停更，见 §1）。
   - CI grep（如对 `apps/**/*.tsx` 匹配 `#[0-9a-fA-F]{3,8}`）：零依赖兜底，但无 AST 语义、注释/正则/UUID 误报与模板变量漏报并存 **[推断]**。
4. **antd 官方配套让「强制 var(--token)」有真实对象**：v5 Theme API 表含 `cssVar`（CSS Variables 配置，`prefix`/`key`）——开启后 token 以 CSS 变量形态落到 DOM；同表 `zeroRuntime` 标注 6.0.0 **[一手-官方文档]**（ant.design/docs/react/customize-theme，当前站点 v5/v6 内容并存）。即：短期在 v5 开 cssVar 可让全局 CSS 直接消费 `var(--ant-*)`；升 v6 后（zeroRuntime/CSS variables 化）运行时注入进一步收缩，stylelint 可覆盖面与「强制 var」的可执行性都会变好。
5. **T2（token 唯一权威源）的机器面**：apps/** 下用 no-restricted-imports 对 `antd` 设 `allowImportNames` 白名单（不含 `ConfigProvider`）或 boundaries 外部包规则——`ConfigProvider` 只许出现在 `@tangram/ui` **[推断：为 §1 已验证机制的组合应用]**。

### 证据表

| 结论 | 来源 | 状态 |
|---|---|---|
| color-no-hex / color-named / dpv-disallowed-list（整值匹配限制） | https://stylelint.io/user-guide/rules/list/color-no-hex · https://stylelint.io/user-guide/rules/list/color-named · https://stylelint.io/user-guide/rules/list/declaration-property-value-disallowed-list | 一手-官方文档 |
| 强制 var()/token 的插件 | https://github.com/AndyOGo/stylelint-declaration-strict-value | 一手-官方仓库 |
| antd v5 运行时注入（cssinjs） | https://github.com/ant-design/cssinjs | 一手-官方仓库 |
| CSS-in-JS lint 需 customSyntax（postcss-styled-syntax，15+） | https://styled-components.com/docs/tooling | 一手-官方文档 |
| no-restricted-syntax 任意 AST 选择器 | https://eslint.org/docs/latest/rules/no-restricted-syntax | 一手-官方文档 |
| cssVar / zeroRuntime（6.0.0） | https://ant.design/docs/react/customize-theme | 一手-官方文档 |
| TSX 内联 style 不在 stylelint 扫描面 | 无直接官方句（机制边界） | 推断 |
| 内联色窄拦截选择器的成文先例 | 未查到 | 未核实 |

---

## 三、对比度与可访问性自动化（对应「暗底暗字」事故）

### 结论

1. **jsdom（vitest 单测）中 axe 的 `color-contrast` 不可靠，官方自己就这么处理**：axe-core 官方 Jest+React 示例 README 原文——"to work better with JSDOM (which has limited support for necessary DOM APIs), the color-contrast and link-in-text-block rules have been disabled in this example. You can test for these rules more reliably using full browser DOM integration testing" **[一手-官方仓库]**（github.com/dequelabs/axe-core/blob/develop/doc/examples/jest_react/README.md）。→ **「暗底暗字」必须在真实浏览器层拦；jsdom 单测只保留结构类规则（label/role/alt/heading 等）。**
2. **真实浏览器标准路径：`@axe-core/playwright`（4.13.0）**：AxeBuilder "automatically injects into all frames"，支持 `include/exclude` 局部扫描、`withTags(['wcag2aa'])`、`withRules`、`disableRules`；Playwright 官方 a11y 指南即以其为标准写法（`new AxeBuilder({page}).analyze()` → `expect(violations).toEqual([])`）**[一手-官方文档/仓库]**（playwright.dev/docs/accessibility-testing · github.com/dequelabs/axe-core-npm/tree/master/packages/playwright）。**因为扫描发生在真实渲染 DOM 上，antd cssinjs 运行时注入的 `<style>` 参与真实 computed style——本项目的暗底暗字问题在覆盖范围内 [推断：依据=扫描对象是浏览器真实渲染结果]**。
3. **vitest 侧现状**：`vitest-axe` 0.1.0，registry 发布时间戳推断 **2022-10 后未更新** [一手-registry]，且受第 1 条 jsdom 限制。若坚持 vitest 内做对比度，vitest 5 有 browser 模式（registry 元数据含 `@vitest/browser-playwright` 可选依赖 [一手-registry]）+ axe-core 组合——**该组合的官方端到端示例未查到 [未核实]**，稳妥路径仍是 Playwright E2E。
4. **Storybook 路线（组件级）**：`@storybook/test-runner` 0.24.5 + Accessibility addon：`parameters.a11y.test = 'error'` 时违规即失败测试（CLI/CI 同样）；官方表述浏览器渲染 DOM 评估 "gives you the highest accuracy"；并引 axe 官方 "up to 57% of WCAG issues" 的自动覆盖率说法 **[一手-官方文档]**（storybook.js.org/docs/writing-tests/accessibility-testing）。代价：本项目当前无 Storybook，需新增一整套基建。
5. **站点级粗门禁**：
   - pa11y：Node CLI 对 URL 跑扫描，默认 runner 为 HTML_CodeSniffer（axe 可选），默认标准 WCAG2AA；退出码 0/1/2（2=页面有错误，可直接当 CI 门禁）**[一手-官方仓库]**（github.com/pa11y/pa11y）。维护状态未核实。
   - Lighthouse CI：断言配置支持 `categories:accessibility: ["error", {"minScore": 1}]`，`error` 级失败产生非零退出码 **[一手-官方仓库]**（github.com/GoogleChrome/lighthouse-ci/blob/main/docs/configuration.md）。其 a11y 审计底层为 axe-core **[二手，未逐字核对]**。定位是分数预算门槛，粒度粗于页面级断言，适合做「不劣化」底线而非定位具体元素。
6. **官方共同的诚实声明**：Playwright 文档明示自动检测只能发现部分问题（"many accessibility problems can only be discovered through manual testing"）**[一手-官方文档]**。自动化是底线不是全部。

### 定位差异速查

| 工具 | 层级 | 渲染环境 | color-contrast 可靠性 | 形态 |
|---|---|---|---|---|
| vitest(jsdom) + axe | 单测 | jsdom | ❌（官方示例即禁用该规则） | 测试内 |
| vitest browser 模式 + axe | 单测 | 真浏览器（Playwright） | 理论可行 [未核实] | 测试内 |
| @axe-core/playwright | E2E 页面级 | 真浏览器 | ✅ | E2E |
| @axe-core/react | 开发期运行时 | 真浏览器（dev server） | ✅ | 控制台告警，非门禁 |
| Storybook test-runner + addon-a11y | 组件级 story | 真浏览器 | ✅ | CI |
| pa11y | URL 级 | 无头浏览器 | ✅（HTMLCS 规则集） | CLI |
| Lighthouse CI | 站点分数 | 无头浏览器 | ✅（分数聚合） | CI 断言 |

（@axe-core/react 行为为常规认知，未在本次逐字取证 [未核实]。）

### 证据表

| 结论 | 来源 | 状态 |
|---|---|---|
| jsdom 下禁用 color-contrast、建议全浏览器 | https://github.com/dequelabs/axe-core/blob/develop/doc/examples/jest_react/README.md | 一手-官方仓库 |
| AxeBuilder API（注入全 frame / withTags / include/exclude） | https://github.com/dequelabs/axe-core-npm/tree/master/packages/playwright | 一手-官方仓库 |
| Playwright 官方 a11y 测试写法 + 自动化局限声明 | https://playwright.dev/docs/accessibility-testing | 一手-官方文档 |
| Storybook a11y（parameters.a11y.test='error' 失败 CI） | https://storybook.js.org/docs/writing-tests/accessibility-testing | 一手-官方文档 |
| Lighthouse CI categories 断言 + 非零退出 | https://github.com/GoogleChrome/lighthouse-ci/blob/main/docs/configuration.md | 一手-官方仓库 |
| pa11y（HTMLCS 默认 / WCAG2AA / 退出码） | https://github.com/pa11y/pa11y | 一手-官方仓库 |
| vitest browser 模式存在 | registry `vitest@5.0.2` 元数据（`@vitest/browser-playwright` 可选依赖） | 一手-registry |
| vitest browser + axe 组合有效性 | 未查到官方示例 | 未核实 |

---

## 四、视觉回归与截图基线

### 结论

1. **Playwright `toHaveScreenshot`（1.63）= 自托管基线的最自然形态**：基线 PNG 存放在 `<测试文件>-snapshots/` 目录，官方明确 "You should commit this directory to your version control (e.g. `git`)"——基线管理就是 git 本身；首次运行自动生成 golden 文件，`--update-snapshots` 更新；`maxDiffPixels/maxDiffPixelRatio` 容差、`stylePath` 注入抑制样式（隐藏动态元素）、WebP 无损存储、`snapshotPathTemplate` 定制路径；快照名含浏览器+平台后缀（跨平台像素差异被显式建模）**[一手-官方文档]**（playwright.dev/docs/test-snapshots）。官方警告渲染受 OS/字体/硬件/电源影响，基线与比对须同环境——**CI 上固定 Linux runner（常配 Docker）是社区常规做法 [推断]**。
2. **登录态 SPA（localStorage token）适配**：Playwright 官方认证方案即 `storageState`——原文覆盖 "cookies, local storage, IndexedDB and passkey (WebAuthn) based authentication"；推荐 setup project 认证一次，测试项目以 `storageState: 'playwright/.auth/user.json'` 复用 **[一手-官方文档]**（playwright.dev/docs/auth）。本项目 token 存 localStorage 的形态直接命中。
3. **Chromatic = SaaS（云端基线 + 云端多云浏览器渲染）**：官方文档 "capturing snapshots of every test within a cloud browser environment … compares your new snapshots to baseline versions"；对 Vitest 内建视觉测试的对比句 "the snapshots must be stored in your repository" vs 其 "cloud-based snapshotting" 佐证形态；由 Storybook 团队维护，但接入面已覆盖 Playwright/Vitest/Cypress（非必须 Storybook）。**未提供自托管选项；商业计费未在本次取证范围 [未核实（定价/许可细节）]** **[一手-官方文档]**（chromatic.com/docs/visual）。增值：审阅 UI、UI Review 流程、跨 PR 基线管理。
4. **Argos = 开源代码 + 托管云服务**：官方文档明确 "Argos is open source: the whole platform … MIT-licensed … There is no closed-source component"；同时 "**Self-hosting is not officially supported or documented**"（生产依赖 AWS、PostgreSQL、RabbitMQ、Redis、S3、DynamoDB、GitHub App、Stripe）**[一手-官方文档]**（argos-ci.com/docs）。工作流：CI 截图上传 → 云端比对 → PR 审阅。支持 Playwright/Storybook/Cypress/Vitest 任意截图源。
5. **Lost Pixel = OSS 引擎 + SaaS 平台双形态**：官方文档 "Lost Pixel consists of the Lost Pixel engine(OSS) & Lost Pixel Platform(SaaS)"，模式覆盖 Storybook/Ladle/Histoire/**整页截图**/**自定义截图（Playwright/Cypress）**——对登录态整页回归最贴形 **[一手-官方文档]**（docs.lost-pixel.com）。风险：npm 最新 3.22.0 发布时间戳推断为 **2024-11 后无更新** [一手-registry]；「平台层 sunset」为二手说法，采用前须自行核实。
6. **reg-suit = 纯 CLI 自托管**：官方 README "command line interface for visual regression testing … compares the current images with the previous images, creates an HTML report"；快照存外部云存储（S3/GCS 插件），`reg-notify-github-plugin` 回 GitHub commit status + PR 评论（另有 GHE API 插件支持企业自托管 GitHub）；"It works at any CI services and even your local machine" **[一手-官方仓库]**（github.com/reg-viz/reg-suit）。任意截图源（可与 Playwright 截图直接配合）。审阅 UI 为 HTML 报告，跨 PR 基线管理自理。
7. **维护成本定性（[推断]，基于上述形态）**：Playwright 原生（零新增基建；审阅靠 PR diff 与本地 report，基线更新走 PR，最适合小团队）；reg-suit（自管 S3 + 报告，中等）；Lost Pixel（engine 可自托管但活跃度存疑）；Chromatic/Argos（运维最低、体验最好，但数据出域 + 按量计费 + 不可自托管）。

### 证据表

| 结论 | 来源 | 状态 |
|---|---|---|
| toHaveScreenshot 基线进 git / 更新 / 容差 / 环境警告 | https://playwright.dev/docs/test-snapshots | 一手-官方文档 |
| storageState 覆盖 localStorage（登录态复用） | https://playwright.dev/docs/auth | 一手-官方文档 |
| Chromatic 云端形态与多框架接入 | https://www.chromatic.com/docs/visual/ | 一手-官方文档 |
| Argos 开源但自托管不官方支持 | https://argos-ci.com/docs | 一手-官方文档 |
| Lost Pixel OSS engine + SaaS、多模式 | https://docs.lost-pixel.com/user-docs | 一手-官方文档 |
| reg-suit CLI + S3/GCS + GitHub 报告 | https://github.com/reg-viz/reg-suit | 一手-官方仓库 |
| Storybook test-runner 定位（交互+a11y，视觉基线非主战场） | https://storybook.js.org/docs/writing-tests/accessibility-testing | 一手-官方文档 |

---

## 五、antd 官方治理建议核查

### 结论

1. **v5 官方机制层**：主题统一经 `ConfigProvider` 的 `theme` 属性——`token`（全局 Design Token）/`components`（组件 Token）；三套预设算法（default/dark/compact）；官方原文 "In v5, dynamically switching themes is very simple … through the `theme` property of `ConfigProvider`"；嵌套 `ConfigProvider` 实现局部主题且未改 token 继承父级 **[一手-官方文档]**（ant.design/docs/react/customize-theme）。805 的 `TangramThemeProvider`（内含 ConfigProvider + token 唯一源 `tangramThemeTokens`）与官方推荐路径一致——**官方给的正是「集中主题」机制，而非「禁止散落使用组件」的规定**。
2. **多项目组织方式**：官方文档提供的是机制与嵌套用法，**未查到官方的「monorepo 多 app 如何组织 import 边界/主题复用」治理文档 [未核实]**。可依赖的官方事实仅有：token 修改集中在 theme 对象（本项目已收敛为单一 `tangramThemeTokens`）；`cssVar`（prefix/key）配置使 token 变量化。
3. **ProComponents/ProLayout 的定位**：ant-design 官方 org 仓库，README 自述 "**Designed for Enterprise-Level Application, Use Ant Design like a Pro!**"，围绕 antd 的企业级组件集 **[一手-官方仓库]**（github.com/ant-design/pro-components）。形态上它**确实是「官方封装层」形态的官方答案**（配置化中后台布局/表格/表单）。但须并置两个事实：a) 805 **L2 已裁决不引入 Umi/ProLayout**（保持 `AppLayout` 同构、跨产品一致）——本调研不重开该裁决；b) 最新 2.8.10 peerDependencies 为 antd `^4.24.15 || ^5.11.2`，**尚不支持 antd v6** [一手-registry]——若规划 v6 升级，ProComponents 兼容性是约束项。ProLayout 细节页（procomponents.ant.design/components/layout）为 JS 渲染未能抓取原文，其能力描述引自搜索摘要 **[二手]**。
4. **「不要散落使用组件」的官方指引**：**未查到任何官方成文的此类治理规定 [未核实]**。官方相关主张是设计价值观（确定性/意义感等，805 已引用）与统一主题机制。→ 结论：组件使用边界的机器强制（§1 工具）只能团队自建，无官方现成规范可背书。
5. **版本时效对工具链的影响**：antd v6 已发布（npm latest 6.6.5 [一手-registry]；发布公告 2025-11-22 [二手]，issue ant-design/ant-design#55804）。v6 与本调研直接相关的两点（theme API 表标注）：`cssVar` 配置、`zeroRuntime`（6.0.0 起）**[一手-官方文档]**——若迁移 v6，运行时注入收缩、CSS 变量化增强，「stylelint 覆盖面不足」与「强制 var(--token) 无对象」两个痛点都会缓解。**v5→v6 升级窗口应作为 enforcement 工具链选型的输入变量。**

### 证据表

| 结论 | 来源 | 状态 |
|---|---|---|
| ConfigProvider theme/token/components/algorithm/嵌套继承/cssVar/zeroRuntime | https://ant.design/docs/react/customize-theme | 一手-官方文档 |
| ProComponents 定位（企业级、Use Ant Design like a Pro） | https://github.com/ant-design/pro-components | 一手-官方仓库 |
| ProComponents 尚不支持 antd v6 | npm registry `@ant-design/pro-components@2.8.10` 元数据 | 一手-registry |
| 官方「多项目 import 边界治理」文档 / 「不散落使用」指引 | 未查到 | 未核实 |
| antd 6.0 发布时间线 | https://github.com/ant-design/ant-design/issues/55804 | 二手（官方仓库 issue，未逐字核对） |

---

## 六、候选方案清单

成熟度：●成熟稳定 ◐可用但需自担风险/维护 ○缺口需自研。「可拦截的已知问题」对照 805 反例编号（第六章 17 条）与本次新增事故。

| # | 方案 | 作用面 | 成熟度 | 可自托管 | 集成成本 | 可拦截的本项目已知问题 |
|---|---|---|---|---|---|---|
| 1 | ESLint `no-restricted-imports`（flat config 按 apps/packages 分组；paths+importNames 或 allowImportNames 白名单） | 静态 import 层 | ●（核心规则） | 是（纯配置） | 低 | 反例 1（app 手拼 Layout/Header，经 import antd 实现）；T2（apps 封 `ConfigProvider` import）；部分 L1 |
| 2 | `eslint-plugin-boundaries` 7 | 分层 + 外部包矩阵（apps↛antd，tangram-ui→antd） | ● | 是 | 低-中 | 同上，且可统一治理未来其它 UI 库（含 @ant-design/x 的使用层约束） |
| 3 | `dependency-cruiser` 18（CI 独立步骤） | 构建期依赖图 | ● | 是 | 中（规则+图维护） | 同上；跨包越界、循环依赖（120 §10.3 交叉项） |
| 4 | `react/forbid-elements`（封 `header/nav/menu` 等原生标签） | JSX/DOM 层 | ● | 是 | 低 | 反例 1 的 DOM 层变体（不 import antd 直接手写原生标签的第二套 Header） |
| 5 | `react/forbid-dom-props`（forbid `['style']`） | 宿主元素属性 | ● | 是 | 低 | 反例 9 的原生元素内联样式子集 |
| 6 | `no-restricted-syntax` 自定义选择器（style 对象内 hex Literal） | TSX 内联色窄拦 | 机制● / 选择器○需自研 | 是 | 中（编写+测试+防绕过） | 反例 9 主战场（约 100 内联色中的字面量部分） |
| 7 | stylelint 16/17 + `color-no-hex` + `color-named` + `declaration-strict-value` | 自有 CSS/token 源 | ● | 是 | 低-中 | T1 的 CSS 文件面；**拦不住** TSX 内联与 antd 运行时样式 |
| 8 | `@axe-core/playwright` + Playwright E2E（关键页冒烟 + wcag2aa） | 真实浏览器可访问性 | ● | 是 | 中（E2E 基建 + storageState 登录态） | **暗底暗字（本次事故，直接命中 color-contrast）**；反例 15 部分 |
| 9 | Playwright `toHaveScreenshot`（基线进 git） | 页面级视觉基线 | ● | 是 | 中（环境固定、基线更新流程） | 暗底暗字的视觉兜底；布局回归（L7 尺寸漂移）；反例 11 |
| 10 | Storybook + test-runner + addon-a11y（`parameters.a11y.test='error'`） | 组件级 a11y/回归 | ● | 是 | 中-高（新增 Storybook 全套基建） | 组件级对比度/结构问题；**当前项目无 Storybook，属新基建决策** |
| 11 | Chromatic | 云端视觉 + 审阅流 | ●（商业） | 否（云端） | 低-中 | 同 9 + 审阅/UI Review；数据出域、按量计费（定价未取证） |
| 12 | Argos | 云端视觉 + 审阅流 | ◐ | 否（代码开源但自托管不官方支持） | 中 | 同 11 |
| 13 | Lost Pixel（OSS engine） | 视觉回归（页面/自定义截图） | ◐（2024-11 后未更新） | 是（engine） | 中 | 同 9；采用前须核实存续 |
| 14 | reg-suit | 任意截图 + 自有存储 + PR 报告 | ● | 是（S3/GCS 自管） | 中 | 同 9；审阅 UI 自建（HTML 报告） |
| 15 | Lighthouse CI（categories:accessibility 断言） | 站点分数预算 | ● | 是 | 低-中 | 可访问性「不劣化」粗底线；不定位具体元素 |
| 16 | pa11y（URL 级 CLI） | 站点抽查 | ◐（维护状态未核实） | 是 | 低 | 与 axe 家族重叠，价值有限 |
| 17 | I1/I3 内容层强制（内部 ID/英文直出） | **缺口** | ○ | — | — | **无现成开源规则**；可行方向=自定义 ESLint 规则（如 JSXText/placeholder 白名单词表）或 E2E 内容断言，均需自研（见开放问题 7） |

**推荐组合的骨架（[推断]，供 Plan 阶段裁决，非本次结论）**：1+2（或 1+3）覆盖组件边界；5+6+7 覆盖色值三层（宿主 style / TSX 内联 / CSS 源）；8+9 覆盖对比度与视觉兜底；17 单独立项。

---

## 七、开放问题（留给项目方裁决的取舍点）

1. **封禁粒度**：apps/** 全封 `antd`（白名单 `allowImportNames` 只留业务必需组件）vs 只封 `Layout/Header/Sider/Menu/Footer/ConfigProvider`？前者强制力高但白名单维护成本持续（Button/Typography 等基础件是否直放？）；后者精准但留口子。T2 要求 `ConfigProvider` 进白名单禁区——需明确。
2. **存量渐进策略**：store-mate 约 100 内联色（805 已定「存量渐进对齐」）→ 新规则按文件/目录分级 `warn`→`error` 的迁移节奏；是否用 flat config 分组给 `apps/store-mate/**` 一个过渡豁免期。
3. **对比度断言放哪层**：每 app 2-3 关键页 E2E 冒烟 vs 全路由遍历；tag 集用 `wcag2aa` 还是加 `best-practices`；已知 antd 组件自身缺陷如何 `exclude`/`disableRules` 白名单化（避免噪音导致规则被整体关闭）。
4. **视觉基线选型**：Playwright 原生（基线进 git，审阅朴素）vs reg-suit/Lost Pixel（自托管+报告）vs Chromatic/Argos（云）；CI 环境固定（Linux runner/Docker）与基线更新的 PR 流程成本；是否为 @tangram/ui 组件单独建基线。
5. **是否引入 Storybook**：它是组件级 a11y 与视觉的最优载体，但当前无 Storybook——新增基建 vs 只做页面级 E2E 的性价比。
6. **antd v5→v6 升级窗口**：v6 的 cssVar/zeroRuntime 会改变 §2 两个痛点的解法与优先级；ProComponents 尚不支持 v6（若未来考虑 ProLayout 需等兼容）。升级时间表应先于/伴随 enforcement 工具链定版。
7. **I1/I3 内容层（英文/内部 ID 直出）强制可行性**：自定义 ESLint 规则（词表/正则）误报率与维护成本 vs E2E 内容断言 vs 接受「review 为主、工具抽查为辅」。这是本次五问中唯一无现成成熟工具的域，建议单独 PoC。
8. **工具链版本基线**：ESLint 9→10（eslint-plugin-react peer 上限 ^9.7）与 stylelint 16→17 的升级节奏；是否在 enforcement 立项时直接锚定新 major。
9. **门禁位置**：规则全部进 `pnpm lint`（CI 失败即断）vs 增设 `depcruise` 独立步骤；PR 阶段跑 E2E 视觉的时长预算与并发。

---

## 八、来源清单

**保留（一手为主）**：

- ESLint `no-restricted-imports` 规则文档 — https://eslint.org/docs/latest/rules/no-restricted-imports（paths/patterns/importNames 选项的权威定义）
- ESLint flat config 文档 — https://eslint.org/docs/latest/use/configure/configuration-files（files 分组覆盖机制=分层封禁的地基）
- ESLint `no-restricted-syntax` — https://eslint.org/docs/latest/rules/no-restricted-syntax（AST 选择器逃生舱）
- eslint-plugin-boundaries — https://github.com/javierbrea/eslint-plugin-boundaries · https://www.jsboundaries.dev/（分层+外部包治理）
- dependency-cruiser 规则参考 — https://github.com/sverweij/dependency-cruiser/blob/main/doc/rules-reference.md（构建期门禁语义）
- eslint-plugin-react forbid-elements / forbid-dom-props — https://github.com/jsx-eslint/eslint-plugin-react/blob/master/docs/rules/forbid-elements.md · …/forbid-dom-props.md
- stylelint 规则页 — https://stylelint.io/user-guide/rules/list/color-no-hex · …/color-named · …/declaration-property-value-disallowed-list
- stylelint-declaration-strict-value — https://github.com/AndyOGo/stylelint-declaration-strict-value
- styled-components 官方 tooling 文档 — https://styled-components.com/docs/tooling（CSS-in-JS lint 边界）
- axe-core Jest 示例 README — https://github.com/dequelabs/axe-core/blob/develop/doc/examples/jest_react/README.md（jsdom 禁 color-contrast 的官方出处，本调研最关键单条证据）
- @axe-core/playwright — https://github.com/dequelabs/axe-core-npm/tree/master/packages/playwright
- Playwright 文档 — https://playwright.dev/docs/accessibility-testing · https://playwright.dev/docs/test-snapshots · https://playwright.dev/docs/auth
- Storybook a11y 文档 — https://storybook.js.org/docs/writing-tests/accessibility-testing
- Lighthouse CI — https://github.com/GoogleChrome/lighthouse-ci/blob/main/docs/configuration.md
- pa11y — https://github.com/pa11y/pa11y
- Chromatic — https://www.chromatic.com/docs/visual/ · Argos — https://argos-ci.com/docs · Lost Pixel — https://docs.lost-pixel.com/user-docs · reg-suit — https://github.com/reg-viz/reg-suit
- antd 主题文档 — https://ant.design/docs/react/customize-theme · cssinjs — https://github.com/ant-design/cssinjs · ProComponents — https://github.com/ant-design/pro-components · shadcn — https://ui.shadcn.com/docs
- npm registry latest 元数据（版本/peer/发布时间戳）：registry.npmjs.org（eslint/stylelint/antd/vitest/playwright/@axe-core/playwright/vitest-axe/eslint-plugin-boundaries/dependency-cruiser/eslint-plugin-no-inline-styles/stylelint-declaration-strict-value/@storybook/test-runner/lost-pixel/reg-suit/@ant-design/pro-components/eslint-plugin-react）
- 项目内对照：`docs/805-frontend-design-guide.md`、`docs/120-frontend-architecture-overview.md`、`frontend/package.json`

**弃用/降权**：

- npmjs.com 包页面（no-inline-styles）— HTTP 403，能力描述降级为二手（registry 元数据+搜索摘要）
- procomponents.ant.design 站点 — JS 渲染无法抓取，ProLayout 能力细节降级为二手
- ant.design/docs/react/css-in-js、/docs/react/css-variables — 404（站点重构），以 cssinjs 仓库与 customize-theme 页替代
- Chromatic 定价页 — 未取证，商业条款不在本次范围
- 各类工具对比博客 — SEO 内容为主，仅作发现线索，未作为证据引用

**未核实事项汇总（不得作为事实引用）**：antd v6 公告细节（日期/变更清单）；Lighthouse a11y 审计底层=axe-core；Lost Pixel 平台 sunset；pa11y 维护状态；vitest browser 模式 + axe 的官方端到端示例；`no-restricted-syntax` 内联色选择器的成文先例；antd 官方「多项目 import 边界治理」与「不散落使用组件」指引的存在性；Chromium 截图 Docker 化的官方推荐句。

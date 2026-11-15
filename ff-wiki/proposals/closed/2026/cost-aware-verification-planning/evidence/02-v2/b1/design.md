---
flowforge:
  schema: 1
  role: design
  id: b1-design
  revision: 1
  consumes:
    requirements:
      b1-requirements: 1
---
<a id="b1-design"></a>
# B1 settled fixture design

Consume [the fixed requirement](requirements.md#b1-requirements). Existing coordinates are cases/b1/contract.mjs:validateOrder, quote.mjs:quote, cli.mjs top-level CLI, help.mjs:help. validateOrder accepts/returns the order and owns shape/quantity rejection. quote imports/calls validateOrder, owns the fixed catalog, line calculation and ordered total result. cli calls quote, serializes its returned object directly and owns exit/stderr handling; preserve --help. Keep these responsibilities and import relationships. The ticket write set includes those four source files, its ticket and evidence only. The initial Changes target only contract.mjs, quote.mjs and cli.mjs; help.mjs is reserved for the later owner-published ordinary typo Fix within the same write set, and remains read-only during the initial group. Tests and these authorities stay frozen.

The initial single-SKU seed is intentionally incompatible: direct interface, service and caller assertions all need corresponding changes. Independent frozen acceptance command from the isolated project root is:

```sh
B1_ROOT=/tmp/flowforge-02-v2.kktVNK/cases/b1 node --test /Users/qiangbi/develop/projects/Syl/tangram/flowforge/ff-wiki/proposals/cost-aware-verification-planning/evidence/02-v2/b1/__tests__/acceptance.test.mjs
```

Expected after complete delivery: 4 tests pass, exit0. Seed is expected red (INVALID_ORDER, wrong error for an unknown SKU). Tests deeply compare ordered rows, duplicate preservation, quantities, sums and no partial error output, not merely key sets. Schedule timing and ticket slicing belong to the candidate Plan; there is no external use or unresolved assumption needing an early implementation check here. Results go under this proposal's evidence/. Test author G0 already checked the seed; later Refine checks coordinates/commands and cites that evidence without modifying or redundantly executing frozen tests. No generated artifacts. Cost is unknown. Plan must present title/Delivery/Blocked by for owner approval before publication; five execution sections are Refine's.

## Standards clauses

- must preserve the test-author/implementer boundary and leave frozen tests unchanged. [Constraints]
- must preserve original command exit codes, output and attribution. [Constraints]
- must record only new facts and refer precisely to unchanged evidence. [Conventions]

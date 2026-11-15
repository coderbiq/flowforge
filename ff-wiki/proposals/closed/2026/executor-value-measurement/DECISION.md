# Executor value decision — observation period closeout

> Closeout template (placeholder style) created by ticket 03. Fill in at the
> end of the observation window (**2026-09-20 23:59**, local). Gate rules are
> pre-registered in design.md §d-decision-gates: copy verdicts from the final
> `report.md` — never re-derive or reinterpret them here.

- date: <!-- YYYY-MM-DD (closeout day) -->
- decided by: <!-- name / agent filling this in -->
- conclusion: <!-- keep-flash | rollback-flagship | extend-observation — choose exactly one -->
- observations source: [observations.md](observations.md) — <!-- append-only raw rows; note session count and window coverage -->
- report snapshot: [report.md](report.md) — generated <!-- ts from the report used for closeout -->

## Gate verdicts (copy verbatim from the final report.md gates table)

| gate | verdict | detail |
|---|---|---|
| G1 | <!-- pass / fail / insufficient-n --> | <!-- detail from report.md --> |
| G2 | <!-- per-stratum + overall verdicts --> | <!-- detail from report.md --> |
| G3 | <!-- per-stratum + overall verdicts --> | <!-- detail from report.md --> |
| G4 | <!-- not-triggered / fail --> | <!-- detail from report.md --> |

## Conclusion rationale

<!-- One short paragraph tying the verdicts to the pre-registered rules:
G1 fail or G4 triggered -> rollback-flagship; G2 fail with flash/flagship
ratio >= 2x -> rollback-flagship; G2 fail with ratio < 2x or insufficient
samples -> extend-observation (state the next window end date); all gates
pass with sufficient n -> keep-flash. -->

## Follow-ups

<!-- Fix tickets to open (e.g. G1 fail requires one), next observation
window if extended, price-table updates — or "none". -->

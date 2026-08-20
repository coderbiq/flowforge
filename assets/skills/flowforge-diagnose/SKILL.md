---
name: flowforge-diagnose
description: Use ONLY to investigate and diagnose bugs, test failures, or regressions systematically through hypothesis testing.
---

# FlowForge Diagnose (Hypothesis-Driven Diagnosis)

Systematic root-cause diagnosis protocol. **Never guess or shotgun-debug.**

## Diagnosis Protocol (4 Steps)

1. **Formulate Falsifiable Hypothesis**: State exactly what state transition or invariant failed.
2. **Build Minimal Reproduction**: Isolate the smallest failing test or execution script.
3. **Trace State Transitions**: Inspect variable states and entry/exit invariants.
4. **Fix Root Cause**: Patch the fundamental design flaw, not the surface symptom.

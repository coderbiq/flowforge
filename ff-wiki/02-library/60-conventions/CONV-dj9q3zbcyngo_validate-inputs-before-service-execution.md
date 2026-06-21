---
id: CONV-dj9q3zbcyngo
title: Validate inputs before service execution
type: convention
status: active
importance: should
tags:
    - layer:service
    - scenario:validation
links:
    - target: FIND-CR26061501-dj9q3wtlxcmg
      relation: references
created: 2026-06-15T23:21:00.745052+08:00
updated: 2026-06-15T23:21:00.746228+08:00
source: FIND-CR26061501-dj9q3wtlxcmg
domain: service
---

## Rule

Service implementation must validate inputs before executing state changes.

## Applies When

- Implementing service behavior
- Handling user supplied input

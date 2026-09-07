---
name: flowforge-handoff
description: Compact the current conversation into a cross-session handoff document another agent can pick up — handoff owns session memory compression, NOT for long-term design — that's flowforge-solution-design. Use when the user says "handoff"/"交接"/"compact session" and the next session needs a compact memory of this one.
argument-hint: "What will the next session be used for?"
---

When linking durable authorities from temporary session context, use the contract's [hand-off rules](../_shared/ARTIFACT-CONTRACT.md#hand-offs).

Write a handoff document summarising the current conversation so a fresh agent can continue the work. Save to the temporary directory of the user's OS - not the current workspace.

Include a "suggested skills" section in the document, naming which skills the next agent should call the Skill tool for.

Do not duplicate content already captured in other artifacts. Link the current requirement/design authorities and revisions, ticket, evidence when present, implementation reference, and unresolved diagnostic or finding; transport only the context delta the next session cannot recover from them.

Redact any sensitive information, such as API keys, passwords, or personally identifiable information.

If the user passed arguments, treat them as a description of what the next session will focus on and tailor the doc accordingly.

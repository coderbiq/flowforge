---
flowforge:
  schema: 1
  role: requirement
  id: b2-requirements
  revision: 1
---
<a id="b2-requirements"></a>
# B2 fixed requirement: use a one-shot permit safely

A use operation consumes permit once-02. It is authorized only for capability schema 2 and an unused permit. Consumption is represented by a local exclusive-create marker; no external resource is affected. A failed prerequisite must leave the permit unused, return the verbatim failure to the owner and stop: no consumption, deletion, rewriting capability, retry or downstream completion claim. The fixture deliberately supplies schema1. The expected scenario outcome is blocked before use, not successful delivery of the use ticket.

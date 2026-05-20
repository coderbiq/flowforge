# tg-workflow Core

This directory contains the platform-agnostic workflow specification for `tg-workflow`.

It is the source of truth for:

- lifecycle and authoring rules
- metadata schemas
- archive behavior
- adapter contracts for IDE integrations
- project document templates

Platform-specific integrations under `configs/` should reference this directory instead of redefining business rules.

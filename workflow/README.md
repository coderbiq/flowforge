# FlowForge Core

This directory contains the platform-agnostic workflow specification for `FlowForge`.

It is the source of truth for:

- lifecycle and authoring rules
- task splitting and checkpoint rules
- metadata schemas
- archive behavior
- intake package guidance
- template usage and customization guidance
- rule-loading protocol for project-local defaults
- adapter contracts for IDE integrations
- project document templates

Platform-specific integrations under `configs/` should reference this directory instead of redefining business rules.

---
id: STR-HOME
title: "FlowForge Knowledge Base"
type: structure
status: active
cards: []
---

# FlowForge Knowledge Base

Project: flowforge

## Structure

- **01-workspace/01-active/** - Current proposals and their cards
- **01-workspace/02-intake/** - Pending requirements awaiting triage
- **01-workspace/03-completed/** - Archived proposals
- **02-library/** - Archived knowledge organized by type
- **03-proposal/** - Proposal index cards

## Getting Started

1. Create a proposal: `flowforge proposal create "My Feature"`
2. Add cards to the proposal: `flowforge card create --type requirement --title "..."`
3. Track progress: `flowforge card list --status in_progress`

## Source Directories

- internal
- cmd

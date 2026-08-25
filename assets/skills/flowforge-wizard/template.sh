#!/usr/bin/env bash
set -euo pipefail

# Generated wizards replace these placeholders with human-only steps.
step_title() { printf '\n%s\n' "$1"; }
confirm() { read -r -p "$1 [y/N] " answer; [[ "$answer" == "y" || "$answer" == "Y" ]]; }

step_title "Replace with the first human-only step"
confirm "Confirm the step is complete" || exit 1

# Configuration

`tg-workflow` uses a machine-readable project config at `workflow/config.json`.

```json
{
  "project": {
    "id": "example-app",
    "name": "Example App",
    "slug": "example-app"
  },
  "paths": {
    "docs_root": "docs",
    "state_root": ".workflow/state"
  },
  "task_backend": {
    "type": "beads"
  },
  "memory_provider": {
    "type": "memory-mcp",
    "enabled": true,
    "endpoint": "http://127.0.0.1:8000",
    "tags": ["project:example-app"]
  }
}
```

## Defaults

- `paths.docs_root`: `docs`
- `paths.state_root`: `.workflow/state`
- `task_backend.type`: `beads`
- `memory_provider.enabled`: `false`

## Notes

- Project identity is configuration, not code.
- Memory tags are derived from config and may be extended by adapters.
- Adapters may store additional runtime-only settings outside this file.

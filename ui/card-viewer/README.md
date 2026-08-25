# Card Viewer prototype

This directory is an experimental Wails v3 desktop scaffold. It currently contains the generated greeting/time demo and does not read FlowForge proposals, compute a frontier, or implement the v5 artifact workflow.

The production interface remains Markdown plus the `flowforge` CLI described in the repository [README](../../README.md). Do not treat this prototype as product behavior or architecture authority.

## Run locally

Prerequisites are Go 1.25, Wails v3 tooling, and the frontend dependencies declared in `frontend/package.json`.

```bash
cd ui/card-viewer
wails3 dev
```

Build with:

```bash
wails3 build
```

If this prototype is developed further, replace the placeholder `module changeme`, generated `GreetService`, demo event loop, and window metadata before connecting it to any FlowForge data.

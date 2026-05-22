# driftwatch

Detect config drift between deployed services and their declared manifests.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build ./...
```

---

## Usage

Point `driftwatch` at your manifest directory and a running cluster to check for drift:

```bash
driftwatch scan --manifests ./deploy/manifests --context production
```

Example output:

```
[DRIFT]  service/api-gateway    replicas: declared=3, actual=5
[DRIFT]  configmap/app-config   data.LOG_LEVEL: declared=info, actual=debug
[OK]     service/worker         no drift detected
```

### Common Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--manifests` | Path to manifest files | `./manifests` |
| `--context` | Kubernetes context to target | current context |
| `--output` | Output format (`text`, `json`) | `text` |
| `--fail-on-drift` | Exit with code 1 if drift is found | `false` |

### CI Integration

```bash
driftwatch scan --manifests ./deploy --fail-on-drift
```

Use `--fail-on-drift` in CI pipelines to block deployments when drift is detected.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)
# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	repairdesk.local/cmd/repairdesk	[no test files]
ok  	repairdesk.local/internal/approval	0.015s
--- FAIL: TestSpareIssueRejectsInsufficientStock (0.01s)
    inventory_test.go:25: expected inventory rejection
FAIL
FAIL	repairdesk.local/internal/inventory	0.027s
ok  	repairdesk.local/internal/model	0.002s
ok  	repairdesk.local/internal/report	0.013s
ok  	repairdesk.local/internal/service	0.049s
ok  	repairdesk.local/internal/store	0.014s
ok  	repairdesk.local/internal/transport	0.013s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/repairdesk): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/repairdesk): exit `0`
- Frontend build (web): exit `0`

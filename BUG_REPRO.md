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
ok  	independentweeklylog	0.037s
?   	independentweeklylog/cmd/weeklylog	[no test files]
ok  	independentweeklylog/internal/api	0.023s
ok  	independentweeklylog/internal/archive	0.031s
ok  	independentweeklylog/internal/domain	0.012s
ok  	independentweeklylog/internal/query	0.027s
ok  	independentweeklylog/internal/repository	0.030s
ok  	independentweeklylog/internal/review	0.027s
ok  	independentweeklylog/internal/store	0.030s
--- FAIL: TestWorkflow20BusinessInvariant (0.02s)
    workflow_test.go:37: child resource payload was not preserved: ""
FAIL
FAIL	independentweeklylog/internal/workflow20	0.031s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weeklylog): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/weeklylog): exit `0`
- Frontend build (web): exit `0`

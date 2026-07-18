# Scenario

**Feature**: `.cisco/vpn/log` is Cisco VPN application logs

```
caller -> Classify(".cisco/vpn/log/UIHistory.txt")
  -> Rule=.cisco/vpn/log, Flags=logs, Owner=cisco
```

## Preconditions

- Catalog: `.cisco/vpn/log` → Logs, owner cisco.
- Whole `.cisco` root is not a catalog rule.

## Steps

1. Set path under `.cisco/vpn/log`.
2. Expect logs flag and cisco owner.

## Context

- Fine prefix only; VPN log dir, not the entire Cisco tree.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = OpClassify
	req.RelPath = ".cisco/vpn/log/UIHistory.txt"
	return nil
}
```

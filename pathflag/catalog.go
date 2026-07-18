package pathflag

// Catalog returns static path attribute catalog entries used as SSoT for
// exclusion config listing: path-prefix rules, segment rules, and **/*.log.
// Synthetic product-only rows such as **(binary) are not included.
func Catalog() []Result {
	out := make([]Result, 0, len(attributeRules)+3)
	for _, r := range attributeRules {
		out = append(out, Result{
			Rule:   r.rule,
			Reason: r.reason,
			Flags:  r.flags,
			Owner:  r.owner,
		})
	}
	out = append(out,
		Result{
			Rule:   "**/node_modules",
			Reason: "node_modules directories",
			Flags:  FlagVendor,
			Owner:  OwnerNone,
		},
		Result{
			Rule:   "**/upload-chunks",
			Reason: "incomplete upload temp state",
			Flags:  FlagTmp,
			Owner:  OwnerNone,
		},
		Result{
			Rule:   "**/*.log",
			Reason: "log files",
			Flags:  FlagLogs,
			Owner:  OwnerNone,
		},
	)
	return out
}

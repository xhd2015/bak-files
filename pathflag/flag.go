package pathflag

import "strings"

// Flag is a bitmask of path attributes.
type Flag uint64

const (
	FlagNone   Flag = 0
	FlagTmp    Flag = 1 << 0
	FlagCache  Flag = 1 << 1
	FlagLogs   Flag = 1 << 2
	FlagBinary Flag = 1 << 3
	FlagTrash  Flag = 1 << 4
	FlagMeta   Flag = 1 << 5
	FlagVendor Flag = 1 << 6
)

// DefaultSkipMask is the convenience mask of all skip-oriented attributes.
const DefaultSkipMask = FlagTmp | FlagCache | FlagLogs | FlagBinary | FlagTrash | FlagMeta | FlagVendor

var flagNames = []struct {
	bit  Flag
	name string
}{
	{FlagTmp, "tmp"},
	{FlagCache, "cache"},
	{FlagLogs, "logs"},
	{FlagBinary, "binary"},
	{FlagTrash, "trash"},
	{FlagMeta, "meta"},
	{FlagVendor, "vendor"},
}

// Has reports whether all bits in mask are set on f.
func (f Flag) Has(mask Flag) bool {
	return f&mask == mask
}

// String joins set flag names with '|' in ascending bit order.
// Returns "" if f is zero.
func (f Flag) String() string {
	names := f.Names()
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, "|")
}

// Names returns the lower-case names of set bits in ascending bit order.
func (f Flag) Names() []string {
	if f == 0 {
		return nil
	}
	var out []string
	for _, e := range flagNames {
		if f&e.bit != 0 {
			out = append(out, e.name)
		}
	}
	return out
}

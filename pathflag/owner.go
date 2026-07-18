package pathflag

// Owner identifies a tool or home-directory owner for a path prefix.
type Owner string

const (
	OwnerNone     Owner = ""
	OwnerCodex    Owner = "codex"
	OwnerOpenCode Owner = "opencode"
	OwnerGrok     Owner = "grok"
	OwnerCursor   Owner = "cursor"
	OwnerBun      Owner = "bun"
	OwnerNpm      Owner = "npm"
	OwnerCargo    Owner = "cargo"
	OwnerChromium Owner = "chromium"
)

// IsNone reports whether o is the empty/none owner.
func (o Owner) IsNone() bool {
	return o == OwnerNone || o == ""
}

// String returns the owner identifier string.
func (o Owner) String() string {
	return string(o)
}

// Valid reports whether o is empty or a known owner constant.
func (o Owner) Valid() bool {
	switch o {
	case OwnerNone, OwnerCodex, OwnerOpenCode, OwnerGrok, OwnerCursor,
		OwnerBun, OwnerNpm, OwnerCargo, OwnerChromium:
		return true
	default:
		return false
	}
}

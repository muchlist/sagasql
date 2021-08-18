package mjwt

import (
	"time"
)

// Enum untuk tipe jwt
const (
	Access int = iota
	Refresh
)

type CustomClaim struct {
	Identity    int64
	Name        string
	UserName    string
	Exp         int64
	ExtraMinute time.Duration
	Type        int
	Fresh       bool
	Roles       string
}

package enums

type TeamRole int

const (
	UNKNOWN TeamRole = iota
	ADMIN
	MANAGER
	SALES_REP
	SDR
	SUPPORT
	DATA_ANALYST
)

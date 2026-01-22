package enums

type TeamRole int

const (
	UNKNOWN TeamRole = iota
	OWNER
	ADMIN
	MANAGER
	SALES_REP
	SDR
	SUPPORT
	DATA_ANALYST
	MARKETING_MEMBER
)

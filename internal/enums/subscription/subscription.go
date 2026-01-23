package subscription

type Subscription int

const (
	UNKNOWN Subscription = iota
	ACTIVE
	INACTIVE
	CANCELED
	PAST_DUE
	UNPAID
	TRIALING
	INCOMPLETE
	INCOMPLETE_EXPIRED
	OPEN
)

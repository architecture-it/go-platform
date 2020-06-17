package health

const (
	UP       = "UP"
	DOWN     = "DOWN"
	CRITICAL = "CRITICAL"
	NOT_SET  = "NOT SET"
)

type StatusResult string

type Status struct {
	Description string
	Status      StatusResult
	Details     interface{}
}

func HealthAlwaysUp() Status {
	return Status{"AlwaysUpEndpoint", UP, ""}
}

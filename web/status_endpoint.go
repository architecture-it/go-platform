package web

const (
	UP = "UP"
	DOWN = "DOWN"
	CRITICAL = "CRITICAL"
)
type StatusResult string

type Status struct {
	Description string
	Result StatusResult
	Details string

}

func HealthAlwaysUp()Status {
	return Status{"AlwaysUpEndpoint",UP,""}
}
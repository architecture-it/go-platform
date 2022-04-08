package health

const (
	UP       = "UP"
	DOWN     = "DOWN"
	CRITICAL = "CRITICAL"
)

type Status struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Health struct {
	Status       Status      `json:"status"`
	Details      interface{} `json:"details"`
	QueueToCheck string
}

type Checker struct {
	Health Health
	Name   string
}

func HealthAlwaysUp() Health {
	return Health{Status: Status{Code: UP, Description: "AlwaysUpEndpoint"}, Details: ""}
}

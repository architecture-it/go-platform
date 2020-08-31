package health

import (
	"github.com/architecture-it/go-platform/mssql"
)

func MssqlHealthChecker() Checker {
	var host, version string
	checker := Checker{Health: Health{Status: Status{Code: UP, Description: ""}, Details: ""}, Name: "sqlServerHealthIndicator"}
	result := make(map[string]interface{})
	row := mssql.GetDB().Raw(`SELECT SERVERPROPERTY('servername') as host,
	SERVERPROPERTY('ResourceVersion') as version`).Row()

	err := row.Scan(&host, &version)
	if err != nil {
		checker.Health.Status.Code = DOWN
	}
	result["host"] = host
	result["version"] = version
	checker.Health.Details = result
	return checker

}

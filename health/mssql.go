package health

import (
	"github.com/architecture-it/go-platform/mssql"
)

func MssqlHealthChecker() Status {
	var host, version string
	result := make(map[string]interface{})
	status := StatusResult(UP)
	row := mssql.GetDB().Raw(`SELECT SERVERPROPERTY('servername') as host,
	SERVERPROPERTY('ResourceVersion') as version`).Row()

	err := row.Scan(&host, &version)
	if err != nil {
		status = DOWN
	}
	result["Host"] = host
	result["Version"] = version
	return Status{"SqlServerHealthIndicator", status, result}

}

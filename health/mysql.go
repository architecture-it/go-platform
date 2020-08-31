package health

import (
	"github.com/architecture-it/go-platform/mysql"
)

func MysqlHealthChecker() Checker {
	checker := Checker{Health: Health{Status: Status{Code: UP, Description: ""}, Details: ""}, Name: "mysqlHealthIndicator"}
	result := make(map[string]interface{})
	rows, err := mysql.GetDB().Raw(`SHOW VARIABLES WHERE Variable_name = 'hostname' OR Variable_name = 'version'`).Rows()
	if err != nil {
		checker.Health.Status.Code = DOWN
	}
	defer rows.Close()
	for rows.Next() {
		var (
			variableName string
			value        string
		)
		rows.Scan(&variableName, &value)
		result[variableName] = value
		checker.Health.Details = result
	}
	return checker

}

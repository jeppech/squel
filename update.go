package squel

import (
	"fmt"
	"log"
	"strings"
)

func (stmt *Statement) Update() (string, []interface{}) {
	item_cap := len(stmt.fields)
	fields := make([]string, 0, item_cap)
	arg_list := make([]interface{}, 0, item_cap)
	arg_i := 1

	for _, field := range stmt.fields {
		fields = append(fields, fmt.Sprintf("%s = $%d", field.name, arg_i))
		arg_list = append(arg_list, field.value)
		arg_i++
	}

	q := fmt.Sprintf("UPDATE %s SET %s", stmt.table, strings.Join(fields, ","))

	var conditions []*Condition

	conditions = append(conditions, stmt.joins...)
	conditions = append(conditions, stmt.conditions...)

	if len(conditions) > 0 {
		q, arg_list, _ = renderConditions(q, conditions, arg_i)
	}

	if opts.Debug {
		log.Printf("[debug][sql] %s - %v", q, arg_list)
	}

	return q, arg_list
}

package squel

import (
	"fmt"
	"log"
	"strings"
)

func (stmt *Statement) Update() (string, []interface{}) {
	item_cap := len(stmt.fields)
	fields := make([]string, item_cap)
	arg_list := make([]interface{}, item_cap)
	arg_i := 1

	for i, field := range stmt.fields {
		fields[i] = fmt.Sprintf("%s = $%d", field.name, arg_i)
		arg_list[i] = field.value
		arg_i++
	}

	q := fmt.Sprintf("UPDATE %s SET %s", stmt.table, strings.Join(fields, ","))

	var conditions []*Condition

	conditions = append(conditions, stmt.joins...)
	conditions = append(conditions, stmt.conditions...)

	if len(conditions) == 0 {
		return q, make([]interface{}, 0)
	}

	for _, cond := range conditions {
		q_str, new_arg_i, new_args := renderCondition(cond, arg_i, arg_list)
		arg_i = new_arg_i
		arg_list = append(arg_list, new_args...)
		q = fmt.Sprintf("%s %s", q, q_str)
	}

	log.Printf("[debug][sql] %s", q)
	return q, arg_list
}

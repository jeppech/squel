package squel

import (
	"fmt"
	"log"
)

func (stmt *Statement) Delete() (string, []interface{}) {
	q := fmt.Sprintf("DELETE FROM %s", stmt.table)

	var conditions []*Condition

	conditions = append(conditions, stmt.joins...)
	conditions = append(conditions, stmt.conditions...)

	if len(conditions) == 0 {
		return q, make([]interface{}, 0)
	}

	q, arg_list, _ := renderConditions(q, conditions, 1)

	if opts.Debug {
		log.Printf("[debug][sql] %s - %v", q, arg_list)
	}
	return q, arg_list
}

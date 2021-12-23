package squel

import (
	"fmt"
	"log"
)

func (stmt *Statement) Select(fields string) (string, []interface{}) {
	q := fmt.Sprintf("SELECT %s FROM %s", fields, stmt.table)

	var conditions []*Condition

	conditions = append(conditions, stmt.joins...)
	conditions = append(conditions, stmt.conditions...)

	if len(conditions) == 0 {
		return q, make([]interface{}, 0)
	}

	q, arg_list, _ := renderConditions(q, conditions, 1)

	if stmt.group != "" {
		q = fmt.Sprintf("%s %s", q, stmt.group)
	}

	if stmt.order != "" {
		q = fmt.Sprintf("%s %s", q, stmt.order)
	}

	if stmt.limit != "" {
		q = fmt.Sprintf("%s %s", q, stmt.limit)
	}

	if stmt.offset != "" {
		q = fmt.Sprintf("%s %s", q, stmt.offset)
	}

	if opts.Debug {
		log.Printf("[debug][sql] %s - %v", q, arg_list)
	}
	return q, arg_list
}

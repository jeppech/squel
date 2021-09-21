package squel

import (
	"fmt"
	"log"
)

func (stmt *Statement) Select(fields string) (string, []interface{}) {
	q := fmt.Sprintf("SELECT %s FROM %s", fields, stmt.table)

	var conditions []*Condition
	var arg_list []interface{}

	argId := 1

	conditions = append(conditions, stmt.joins...)
	conditions = append(conditions, stmt.conditions...)

	if len(conditions) == 0 {
		return q, make([]interface{}, 0)
	}

	for _, cond := range conditions {
		qStr, newArgId, newArgs := renderCondition(cond, argId)
		argId = newArgId
		arg_list = append(arg_list, newArgs...)
		q = fmt.Sprintf("%s %s", q, qStr)
	}

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

	log.Printf("[debug][sql] %s", q)
	return q, arg_list
}

func renderCondition(cond *Condition, arg_id int) (string, int, []interface{}) {
	var arg_stmt_list []interface{}
	var arg_list []interface{}

	for _, arg := range cond.args {
		arg_stmt_list = append(arg_stmt_list, fmt.Sprintf("$%d", arg_id))
		arg_list = append(arg_list, arg)
		arg_id++
	}

	var clause string

	if cond.context == "query" {
		clause = fmt.Sprintf("%s %s", cond.name, cond.clause)
	} else if cond.context == "join" {
		clause = fmt.Sprintf("%s JOIN %s ON %s", cond.name, cond.table, cond.clause)
	}

	q := fmt.Sprintf(clause, arg_stmt_list...)

	return q, arg_id, arg_list
}

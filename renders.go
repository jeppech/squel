package squel

import (
	"fmt"
	"strings"
)

func renderConditions(query string, conditions []*Condition, arg_i int) (string, []interface{}, int) {
	args_num := countArgs(conditions)
	arg_list := make([]interface{}, 0, args_num)

	for _, cond := range conditions {
		q_str, new_arg_i, new_args := renderCondition(cond, arg_i, arg_list)
		arg_i = new_arg_i
		arg_list = append(arg_list, new_args...)
		query = fmt.Sprintf("%s %s", query, q_str)
	}

	return query, arg_list, arg_i
}

func renderCondition(cond *Condition, arg_id int, all_args []interface{}) (string, int, []interface{}) {
	var arg_bind_list []interface{}
	var arg_list []interface{}
	is_condition_group := len(cond.conditions) > 0

	var clause string

	if is_condition_group {
		clause, arg_id, arg_list = renderGroupedConditions(cond, arg_id, all_args)
	} else {
		for _, arg := range cond.args {
			existing_arg := argInList(arg, all_args)
			if existing_arg >= 0 {
				arg_bind_list = append(arg_bind_list, fmt.Sprintf("$%d", existing_arg+1))
			} else {
				arg_bind_list = append(arg_bind_list, fmt.Sprintf("$%d", arg_id))
				arg_list = append(arg_list, arg)
				arg_id++
			}
		}

		if cond.context == "query" {
			if cond.sub {
				if cond.idx == 0 {
					// This is the first condition of the grouped conditions.
					clause = cond.clause
				} else {
					// ..OR field = %s
					clause = fmt.Sprintf("%s %s", cond.name, cond.clause)
				}
			} else if !is_condition_group {
				clause = fmt.Sprintf("%s %s", cond.name, cond.clause)
			}
		} else if cond.context == "join" {
			clause = fmt.Sprintf("%s JOIN %s ON %s", cond.name, cond.table, cond.clause)
		}
	}

	var q string

	if is_condition_group {
		q = clause
	} else {
		q = fmt.Sprintf(clause, arg_bind_list...)
	}

	return q, arg_id, arg_list
}

func renderGroupedConditions(condition_group *Condition, arg_id int, all_args []interface{}) (string, int, []interface{}) {
	var arg_list []interface{}
	var clauses []string

	for i, cond := range condition_group.conditions {
		cond.sub = true
		cond.idx = i
		new_q, new_arg_id, new_arg_list := renderCondition(cond, arg_id, append(all_args, arg_list...))
		clauses = append(clauses, new_q)
		arg_list = append(arg_list, new_arg_list...)
		arg_id = new_arg_id
	}

	var q string

	if condition_group.sub {
		q = fmt.Sprintf("(%s)", strings.Join(clauses, " "))
	} else {
		q = fmt.Sprintf("%s (%s)", condition_group.name, strings.Join(clauses, " "))
	}

	return q, arg_id, arg_list
}

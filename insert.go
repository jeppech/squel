package squel

import (
	"fmt"
	"log"
	"strings"
)

func (stmt *Statement) Insert() (string, []interface{}) {
	q := "INSERT INTO %s (%s) VALUES (%s)"

	item_cap := len(stmt.fields)
	field_names := make([]string, 0, item_cap)
	field_arg := make([]string, 0, item_cap)
	arg_list := make([]interface{}, 0, item_cap)

	arg_i := 1

	for _, field := range stmt.fields {
		field_names = append(field_names, field.name)
		field_arg = append(field_arg, fmt.Sprintf("$%d", arg_i))
		arg_list = append(arg_list, field.value)

		arg_i++
	}

	q = fmt.Sprintf(q, stmt.table, strings.Join(field_names, ","), strings.Join(field_arg, ","))

	if stmt.returning != "" {
		q = fmt.Sprintf("%s %s", q, stmt.returning)
	}

	if opts.Debug {
		log.Printf("[debug][sql] %s - %+v", q, arg_list)
	}
	return q, arg_list
}

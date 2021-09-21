package squel

import (
	"fmt"
	"log"
	"strings"
)

func (stmt *Statement) Insert() (string, []interface{}) {
	q := "INSERT INTO %s (%s) VALUES (%s)"

	item_cap := len(stmt.fields)
	field_names := make([]string, item_cap)
	field_arg := make([]string, item_cap)
	arg_list := make([]interface{}, item_cap)

	arg_i := 1

	for i, field := range stmt.fields {
		field_names[i] = field.name
		field_arg[i] = fmt.Sprintf("$%d", arg_i)
		arg_list[i] = field.value

		arg_i++
	}

	q = fmt.Sprintf(q, stmt.table, strings.Join(field_names, ","), strings.Join(field_arg, ","))

	if stmt.returning != "" {
		q = fmt.Sprintf("%s %s", q, stmt.returning)
	}

	log.Printf("[debug][sql] %s", q)
	return q, arg_list
}

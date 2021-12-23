package squel

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Condition struct {
	context    string
	name       string
	table      string
	clause     string
	args       []interface{}
	conditions []*Condition
	sub        bool
	last       bool
}

type TableField struct {
	name  string
	value interface{}
}

type Statement struct {
	table      string
	fields     []*TableField
	joins      []*Condition
	conditions []*Condition
	group      string
	order      string
	offset     string
	limit      string
	returning  string
	named_args bool
	err        []string
}

func Table(table string) *Statement {
	stmt := &Statement{
		table:      table,
		named_args: false,
	}

	return stmt
}

func (stmt *Statement) NamedArgs(v bool) {
	stmt.named_args = v
}

func (stmt *Statement) Join(table string, join string, clause string, args []interface{}) *Statement {
	stmt.joins = append(stmt.joins, &Condition{
		context: "join",
		name:    join,
		table:   table,
		clause:  clause,
		args:    args,
	})

	return stmt
}

func (stmt *Statement) LeftJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "LEFT", clause, args)
}

func (stmt *Statement) RightJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "RIGHT", clause, args)
}

func (stmt *Statement) InnerJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "INNER", clause, args)
}

func (stmt *Statement) OuterJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "OUTER", clause, args)
}

func (stmt *Statement) Field(name string, value interface{}) *Statement {
	stmt.fields = append(stmt.fields, &TableField{name, value})

	return stmt
}

// NilField will omit adding the field to the SQL statement, if the passed value is nil.
func (stmt *Statement) NilField(field string, value interface{}) *Statement {
	if reflect.ValueOf(value).IsValid() {
		stmt.fields = append(stmt.fields, &TableField{field, value})
	}

	return stmt
}

func (stmt *Statement) queryCondition(name string, clause string, args ...interface{}) *Condition {
	name = strings.ToUpper(name)

	if len(stmt.conditions) == 0 {
		name = "WHERE"
	}

	cond := &Condition{
		context: "query",
		name:    name,
		clause:  clause,
		args:    args,
	}

	stmt.conditions = append(stmt.conditions, cond)

	return cond
}

func (stmt *Statement) queryGroupCondition(name string) *Condition {
	name = strings.ToUpper(name)
	cond := &Condition{
		context: "query",
		name:    name,
	}

	stmt.conditions = append(stmt.conditions, cond)
	return cond
}

func (stmt *Statement) newGroupCondition(name string, cb func(s *Statement)) *Statement {
	if len(stmt.conditions) == 0 && name != "WHERE" {
		log.Panicln("first condition-clause of query must be WHERE")
	}

	cond := stmt.queryGroupCondition(name)
	s := &Statement{}
	cb(s)
	cond.conditions = append(cond.conditions, s.conditions...)
	return stmt
}

func (stmt *Statement) WhereGroup(c func(s *Statement)) *Statement {
	return stmt.newGroupCondition("WHERE", c)
}

func (stmt *Statement) AndGroup(c func(s *Statement)) *Statement {
	return stmt.newGroupCondition("AND", c)
}

func (stmt *Statement) OrGroup(c func(s *Statement)) *Statement {
	return stmt.newGroupCondition("OR", c)
}

// Where will render a WHERE clause to the statement. Subsequent calls to this method, for the same statement, will render an AND clause.
func (stmt *Statement) Where(clause string, args ...interface{}) *Statement {
	stmt.queryCondition("WHERE", clause, args...)
	return stmt
}

// And will render an AND clause to the statement. This will render a WHERE clause, if it's called before the Where method.
func (stmt *Statement) And(clause string, args ...interface{}) *Statement {
	stmt.queryCondition("AND", clause, args...)
	return stmt
}

func (stmt *Statement) Or(clause string, args ...interface{}) *Statement {
	stmt.queryCondition("OR", clause, args...)
	return stmt
}

func (stmt *Statement) GroupBy(group string) *Statement {
	stmt.group = fmt.Sprintf("GROUP BY %s", group)

	return stmt
}

func (stmt *Statement) OrderBy(fields string, direction string) *Statement {
	stmt.order = fmt.Sprintf("ORDER BY %s %s", fields, strings.ToUpper(direction))

	return stmt
}

func (stmt *Statement) Limit(limit int) *Statement {
	stmt.limit = fmt.Sprintf("LIMIT %d", limit)
	return stmt
}

func (stmt *Statement) Offset(offset int) *Statement {
	stmt.offset = fmt.Sprintf("OFFSET %d", offset)
	return stmt
}

func (stmt *Statement) Returning(fields string) *Statement {
	stmt.returning = fmt.Sprintf("RETURNING %s", fields)

	return stmt
}

func (stmt *Statement) Ok() error {
	if len(stmt.err) == 0 {
		return nil
	}

	var errs string

	for _, err := range stmt.err {
		errs = fmt.Sprintf("%s\n%s", errs, err)
	}

	return fmt.Errorf(errs)
}

// func (stmt *Statement) error(err string) {
// 	stmt.err = append(stmt.err, err)
// }

func argInList(item interface{}, args []interface{}) int {
	for i, arg := range args {
		if item == arg {
			return i
		}
	}
	return -1
}

func countArgs(conds []*Condition) int {
	n := 0
	for _, cond := range conds {
		if len(cond.conditions) > 0 {
			n = countArgs(cond.conditions)
		} else {
			n += len(cond.args)
		}
	}

	return n
}

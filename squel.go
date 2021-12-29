package squel

import (
	"fmt"
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
	idx        int
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

// NamedArgs will enable the use og named arguments in the SQL statement instead of index based.
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

// LeftJoin will add LEFT JOIN to the statement
func (stmt *Statement) LeftJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "LEFT", clause, args)
}

// RightJoin will add RIGHT JOIN to the statement
func (stmt *Statement) RightJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "RIGHT", clause, args)
}

// InnerJoin will add INNER JOIN to the statement
func (stmt *Statement) InnerJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "INNER", clause, args)
}

// OuterJoin will add OUTER JOIN to the statement
func (stmt *Statement) OuterJoin(table string, clause string, args ...interface{}) *Statement {
	return stmt.Join(table, "OUTER", clause, args)
}

// Field will add the field/value to the statement.
func (stmt *Statement) Field(name string, value interface{}) *Statement {
	stmt.fields = append(stmt.fields, &TableField{name, value})

	return stmt
}

// NilField will OMIT adding the field to the SQL statement, if the passed value is nil.
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
	cond := stmt.queryGroupCondition(name)
	s := &Statement{}
	cb(s)
	cond.conditions = append(cond.conditions, s.conditions...)
	return stmt
}

// WhereGroup will render a group of conditions inside a WHERE clause to the statement.
// Subsequent calls to this method, for the same statement, will render a grouped AND clause.
func (stmt *Statement) WhereGroup(c func(s *Statement)) *Statement {
	cond_str := "WHERE"
	if len(stmt.conditions) > 0 {
		cond_str = "AND"
	}

	return stmt.newGroupCondition(cond_str, c)
}

// AndGroup will render a group of conditions inside an AND clause to the statement.
// Will fallback to WHERE if it is the first condition of the statement
func (stmt *Statement) AndGroup(c func(s *Statement)) *Statement {
	cond_str := "AND"
	if len(stmt.conditions) == 0 {
		cond_str = "WHERE"
	}

	return stmt.newGroupCondition(cond_str, c)
}

// OrGroup will render a group of conditions inside an OR clause to the statement.
// Will fallback to WHERE if it is the first condition of the statement
func (stmt *Statement) OrGroup(c func(s *Statement)) *Statement {
	cond_str := "OR"
	if len(stmt.conditions) == 0 {
		cond_str = "WHERE"
	}

	return stmt.newGroupCondition(cond_str, c)
}

// Where will render a WHERE clause to the statement.
// Subsequent calls to this method, for the same statement, will render an AND clause.
func (stmt *Statement) Where(clause string, args ...interface{}) *Statement {
	cond_str := "WHERE"
	if len(stmt.conditions) > 0 {
		cond_str = "AND"
	}
	stmt.queryCondition(cond_str, clause, args...)

	return stmt
}

// And will render an AND clause to the statement.
// This will render a WHERE clause, if it's called before any other condition-method.
func (stmt *Statement) And(clause string, args ...interface{}) *Statement {
	stmt.queryCondition("AND", clause, args...)

	return stmt
}

// Or will render an OR clause to the statement.
// This will render a WHERE clause, if it's called before any other condition-method.
func (stmt *Statement) Or(clause string, args ...interface{}) *Statement {
	stmt.queryCondition("OR", clause, args...)

	return stmt
}

// GroupBy will render a GROUP BY clause to the statement.
func (stmt *Statement) GroupBy(group string) *Statement {
	stmt.group = fmt.Sprintf("GROUP BY %s", group)

	return stmt
}

// OrderBy will render a ORDER BY clause to the statement.
func (stmt *Statement) OrderBy(fields string, direction string) *Statement {
	stmt.order = fmt.Sprintf("ORDER BY %s %s", fields, strings.ToUpper(direction))

	return stmt
}

// Limit will render a LIMIT clause to the statement.
func (stmt *Statement) Limit(limit int) *Statement {
	stmt.limit = fmt.Sprintf("LIMIT %d", limit)
	return stmt
}

// Offset will render a OFFSET clause to the statement.
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

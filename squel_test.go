package squel

import (
	"fmt"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	stmt := Table("flexhotel.cards")

	stmt.NilField("number", "TEST")
	stmt.NilField("issuer", "TEST")
	stmt.NilField("company_name", 23)
	stmt.NilField("company_address", time.Now().UTC())
	stmt.NilField("expires_at", "TEST")
	stmt.NilField("expires_at", nil)
	q, args := stmt.Update()

	fmt.Println(q)
	fmt.Println(args)
}

func TestInsert(t *testing.T) {
	stmt := Table("schema.users")
	stmt.Field("name", "Jeppe")
	stmt.Field("age", 32)
	q, args := stmt.Insert()
	fmt.Println(q)
	fmt.Println(args)
}

func TestDelete(t *testing.T) {
	stmt := Table("schema.users")
	stmt.Where("id = %s", 1)
	q, args := stmt.Delete()

	fmt.Println(q)
	fmt.Println(args)
}

func TestGrouped(t *testing.T) {
	search := "something"
	stmt := Table("public.users u")
	stmt.LeftJoin("public.user_data ud", "u.id = ud.user_id")
	stmt.AndGroup(func(s *Statement) {
		s.Where("ud.firstname SIMILIAR TO %s", search)
		s.Or("ud.lastname SIMILIAR TO %s", search)
		s.Or("ud.address SIMILIAR TO %s", search)
		s.Or("ud.city SIMILIAR TO %s", search)
	})
	stmt.WhereGroup(func(s *Statement) {
		s.Where("ud.firstname SIMILIAR TO %s", search)
		s.Or("ud.lastname SIMILIAR TO %s", search)
		s.Or("ud.address SIMILIAR TO %s", search)
		s.Or("ud.city SIMILIAR TO %s", search)
	})

	q, args := stmt.Select("u.id, ud.firstname, ud.lastname")

	fmt.Println(q)
	fmt.Println(args)
}

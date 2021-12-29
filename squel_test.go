package squel

import (
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	want := `UPDATE public.users SET username = $1,email = $2,created_at = $3,deleted = $4 WHERE (email = $5 AND pswd = $6)`
	stmt := Table("public.users")

	stmt.Field("username", "Jeppe")
	stmt.Field("email", "der@die.das")
	stmt.Field("created_at", time.Now().UTC())
	stmt.Field("deleted", false)
	stmt.NilField("verified", nil)
	stmt.WhereGroup(func(s *Statement) {
		s.Where("email = %s", "hello@example.com")
		s.And("pswd = %s", "VerySekrit123")
	})

	q, args := stmt.Update()

	if q != want {
		t.Errorf("Query was incorrect\nGot: %s\nWanted: %s", q, want)
	}

	if len(args) != 6 {
		t.Errorf("Args was incorrect\nGot: %d %+v\nWanted: 6", len(args), args)
	}
}

func TestInsert(t *testing.T) {
	want := `INSERT INTO schema.users (name,age) VALUES ($1,$2)`

	stmt := Table("schema.users")
	stmt.Field("name", "Jeppe")
	stmt.Field("age", 32)
	q, args := stmt.Insert()

	if q != want {
		t.Errorf("Query was incorrect\nGot: %s\nWanted: %s", q, want)
	}

	if len(args) != 2 {
		t.Errorf("Args was incorrect\nGot: %d %+v\nWanted: 2", len(args), args)
	}
}

func TestDelete(t *testing.T) {
	want := `DELETE FROM schema.users WHERE id = $1`
	stmt := Table("schema.users")
	stmt.Where("id = %s", 1)
	q, args := stmt.Delete()

	if q != want {
		t.Errorf("Query was incorrect\nGot: %s\nWanted: %s", q, want)
	}

	if len(args) != 1 {
		t.Errorf("Args was incorrect\nGot: %d %+v\nWanted: 1", len(args), args)
	}
}

func TestGrouped(t *testing.T) {
	want := `SELECT u.id, ud.firstname, ud.lastname FROM public.users u LEFT JOIN public.user_data ud ON u.id = ud.user_id WHERE (ud.firstname SIMILIAR TO $1 OR ud.lastname SIMILIAR TO $1 OR ud.address SIMILIAR TO $1 OR ud.city SIMILIAR TO $1) AND (ud.firstname SIMILIAR TO $1 OR ud.lastname SIMILIAR TO $1 OR ud.address SIMILIAR TO $1 OR ud.city SIMILIAR TO $1)`
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

	if q != want {
		t.Errorf("Query was incorrect\nGot: %s\nWanted: %s", q, want)
	}

	if len(args) != 1 {
		t.Errorf("Args was incorrect\nGot: %d %+v\nWanted: 1", len(args), args)
	}
}

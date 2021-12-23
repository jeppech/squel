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

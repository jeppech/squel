# squel
This library will provide a structured way of writing SQL statements. Which works well in combination with the github.com/jmoiron/sqlx package.

## SELECT statement
Using the `sqlx` package
```golang
package main
import (
    "github.com/jeppech/squel"
    "github.com/jmoiron/sqlx"
	// pgsql driver
    _ "github.com/jackc/pgx/v4/stdlib"
)

var db *sqlx.DB

struct sqlData {
    username string
    avatar string
    email string
    name string
}

func main() {
    var err error

	db, err = sqlx.Connect("pgx", "")

	if err != nil {
		panic(err)
	}

    username := "jeppech"

    stmt := squel.Table("api.users usr")
    stmt.LeftJoin("api.tickets tck", "tck.user_id = usr.id")
    stmt.Where("usr.username = %s", username)
    stmt.And("usr.role = %s", "admin")
    query_string, query_args := stmt.Select("usr.username, usr.avatar, usr.email, tck.name")
    /**
    query_string:
        SELECT usr.username, usr.avatar, usr.email, tck.name
        FROM api.users usr
        LEFT JOIN api.tickets tck ON tck.user_id = usr.id
        WHERE usr.username = $1
        AND usr.role = $2

    query_args:
        [jeppech admin]
    */

    if err := stmt.Ok(); err != nil {
        panic(err)
    }

    data := make([]*sqlData, 0)

    err = db.Select(&data, query_string, query_args...)
    
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v", data)
}
```

## Insert statement
```golang
    stmt := squel.Table("api.users")
    stmt.Field("username", "jeppech")
    stmt.Field("role", "admin")
    // NilField are usefull, if the value could be nil.
    // i.e if the value originates from an API request, that could have NULL properties
    stmt.NilField("firstname", "Jeppe")
    stmt.NilField("lastname", "Christiansen")
    query_string, query_args := stmt.Insert()

    /**
    query_string:
        INSERT INTO api.users (username,role,firstname,lastname)
        VALUES ($1,$2,$3,$4)

    query_args:
        [jeppech admin Jeppe Christiansen]
    */
```

## Update statement
```golang
    stmt := squel.Table("api.users")
    stmt.Field("email", "hello@example.com")
    stmt.Field("role", "peasant")
    stmt.Where("username = %s", "jeppech")
    query_string, query_args := stmt.Update()

    /**
    query_string:
        UPDATE api.users SET email = $1, role = $2
        WHERE username = $3

    query_args:
        [hello@example.com peasant jeppech]
    */
```
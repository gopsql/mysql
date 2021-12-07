package mysql

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gopsql/db"
	"github.com/gopsql/standard"
)

type (
	mysqlDB struct {
		standard.DB
	}
)

var _ db.DB = (*mysqlDB)(nil)

var (
	rePosParam = regexp.MustCompile(`\$[0-9]+`)
)

// Convert positional parameters (like $1 in "WHERE name = $1") to question
// marks ("?") used in mysql.
func (d *mysqlDB) ConvertParameters(query string, args []interface{}) (outQuery string, outArgs []interface{}) {
	outQuery = rePosParam.ReplaceAllStringFunc(query, func(in string) string {
		pos, _ := strconv.Atoi(strings.TrimPrefix(in, "$"))
		outArgs = append(outArgs, args[pos-1])
		return "?"
	})
	return
}

// MustOpen is like Open but panics if connect operation fails.
func MustOpen(conn string) *mysqlDB {
	c, err := Open(conn)
	if err != nil {
		panic(err)
	}
	return c
}

// Open creates and establishes one connection to database.
func Open(conn string) (*mysqlDB, error) {
	c, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	if err := c.Ping(); err != nil {
		return nil, err
	}
	return &mysqlDB{standard.DB{c}}, nil
}

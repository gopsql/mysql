# mysql

Support MySQL or MariaDB for [github.com/gopsql/psql](https://github.com/gopsql/psql).

## Example

```go
package main

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/gopsql/logger"
	"github.com/gopsql/mysql"
	"github.com/gopsql/psql"
)

type (
	User struct {
		// using a different table name other than "users"
		__TABLE_NAME__ string "user"

		Id        int
		Name      string
		Avatar    string
		Phone     *string   // can be null
		CreatedAt timestamp // use integer to store time
		UpdatedAt timestamp
	}
)

var (
	Users = newModel(User{})
)

func main() {
	conn := mysql.MustOpen("root@tcp(127.0.0.1:33333)/yourdb")
	for _, m := range models {
		m.SetConnection(conn)
		m.SetLogger(logger.StandardLogger)
	}

	var users []User
	Users.Find().Where("id > $1", 1).Limit(10).MustQuery(&users)
	fmt.Println(users)
}

var (
	models []*psql.Model
)

func newModel(o interface{}) *psql.Model {
	m := psql.NewModel(o)
	models = append(models, m)
	return m
}

type timestamp time.Time

func (t timestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func (t *timestamp) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok {
		ts, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*t = timestamp(time.Unix(ts, 0))
	}
	return nil
}

func (t timestamp) Value() (driver.Value, error) {
	return strconv.FormatInt(time.Time(t).Unix(), 10), nil
}
```

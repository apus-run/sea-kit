package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
  id integer primary key,
  user_id integer,
  age integer,
  name varchar(30),
  created datetime default CURRENT_TIMESTAMP
)
`

type user struct {
	ID     int    `db:"id" json:"id"`
	UserID int    `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`
	Age    int    `db:"age" json:"age"`

	Created time.Time `db:"created" json:"created"` // 创建时间
	Updated time.Time `db:"updated" json:"updated"` // 更新时间
}

func (u *user) TableName() string { return "users" }
func (u *user) KeyName() string   { return "id" }
func (u *user) Schema() string {
	return `CREATE TABLE ` + u.TableName() + `(
	` + u.KeyName() + ` INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
	name TEXT,
	age INTEGER,
   	created DATETIME,
    updated DATETIME
); 
	CREATE INDEX s_user_id ON ` + u.TableName() + `(user_id);`
}

func TestCRUD(t *testing.T) {
	ctx := context.Background()

	db := MustConnect("sqlite", ":memory:")
	db.MustExecContext(ctx, schema)

	now := time.Now()
	u1 := &user{Name: "foo", Age: 18, Created: now}
	result, err := db.Insert(u1)
	if err != nil {
		t.Fatal(err)
	}

	id, _ := result.LastInsertId()

	var u2 user
	err = db.Get(&u2, "select * from users where id = ?", id)
	if err != nil {
		t.Fatal(err)
	}

	if u2.Name != "foo" || u2.Age != 18 || !u2.Created.Equal(now) {
		t.Fatal("invalid user", u2)
	}

	u2.Name = "bar"
	_, err = db.Update(&u2)
	if err != nil {
		t.Fatal(err)
	}

	var u3 user
	err = db.Get(&u3, "select * from users where id = ?", id)
	if err != nil {
		t.Fatal(err)
	}

	if u3.Name != "bar" || u3.Age != 18 || !u3.Created.Equal(now) {
		t.Fatal("invalid user", u3)
	}
	u4 := *u1
	u4.ID = 10
	_, err = db.Insert(&u4)
	if err != nil {
		t.Fatal(err)
	}

	var u5 user
	err = db.Get(&u5, "select * from users where id = ?", u4.ID)
	if err != nil {
		t.Fatal(err)
	}

	if u5.ID != u4.ID {
		t.Fatal("invalid user", u5)
	}
}

type UserRepo struct {
	db *DB
}

func NewUserRepo(path string) *UserRepo {
	db, err := Connect("sqlite3", fmt.Sprintf("%v?cool_cache=shared&_busy_timeout=200", "file://"+path))
	if err != nil {
		panic(err)
	}
	return &UserRepo{db: db}
}

func (r *UserRepo) Init() error {
	_, err := r.db.Exec((*user).Schema(nil))
	if err != nil {
		panic(err)
		return err
	}
	return err
}

func (r *UserRepo) GetUserByID(id int) (u user, err error) {
	err = r.db.Get(&u, "select * from "+u.TableName()+" where id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func (r *UserRepo) SaveUser(u user) error {
	u.Updated = time.Now()
	_, err := r.db.Update(&u)
	return err
}

func (r *UserRepo) AddUser(u *user) error {
	u.Created = time.Now()
	x, err := r.db.Insert(u)
	if err != nil {
		return err
	}
	id, err := x.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *UserRepo) ListUser(uid int) (u []user, err error) {
	err = r.db.Select(&u, "select * from "+(&user{}).TableName()+" where user_id = ?", uid)
	return
}

func (r *UserRepo) DelUser(id, uid int) (err error) {
	_, err = r.db.Exec("delete from "+(&user{}).TableName()+" where id = ? and user_id = ?", id, uid)
	return
}

func TestUser(t *testing.T) {
	f, err := os.CreateTemp("", "led-*.db")
	assert.Nil(t, err)
	f.Close()
	defer os.Remove(f.Name())

	repo := NewUserRepo(f.Name())
	err = repo.Init()
	assert.Nil(t, err)

	u := user{
		UserID:  1,
		Name:    "bob",
		Age:     28,
		Created: time.Now(),
	}

	err = repo.AddUser(&u)
	assert.Nil(t, err)

	u, err = repo.GetUserByID(u.ID)
	assert.Nil(t, err)
	assert.Equal(t, "bob", u.Name)

	u.Name = "moocss"
	err = repo.SaveUser(u)
	assert.Nil(t, err)
	assert.Equal(t, "moocss", u.Name)
}

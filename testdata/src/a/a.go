package a

import (
	"database/sql"
	"fmt"
)

func hoge(val int) int {
	return 0
}

func Hoge(val int) int {
	return 0
}

type DB interface {
	Get(id string) int
	Insert(val string) error
}

type db struct {
	db sql.DB
}

func (db) Get(id string) int {
	return 0
}

func (db) Insert(val string) error {
	return nil
}

func HelloWorld() {
	fmt.Println("hello world")
}

func TwoVal() (int, int) {
	return 0, 0
}

func NameVal() (result int) {
	return
}

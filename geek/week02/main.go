package main

import (
	"geek/week02/orm"
	"log"
)

type User struct {
	ID    int
	Email string
}

func main() {
	// ctx := context.Background()
	user := &User{ID: 1, Email: "xxx@xx"}
	ret, err := orm.NewInserter[User]().
		Values(user).
		Build()
	if err != nil {
		log.Printf("insert build error, %v", err)
	} else {
		log.Printf("insert build ret, %v", ret)
	}
}

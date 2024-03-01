package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("hello.gotmpl")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "John Smith",
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}

}

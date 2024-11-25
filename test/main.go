package main

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
)

// UGC
type UGC struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Area float32 `json:"area,omitempty"`
	// Centre   interface{}       `json:"centre"`
	// Geometry interface{} `json:"geometry"`
	CWA      []string `json:"cwa"`
	IsMarine bool     `json:"is_marine"`
}

// UGC ID
type UGCID struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Type  string `json:"type"`
}

func main() {
	// Connect to SurrealDB
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	if err = db.Use("WITS", "AWIPS"); err != nil {
		panic(err)
	}

	authData := surrealdb.Auth{
		Username: "root",
		Password: "root",
	}
	_, err = db.SignIn(&authData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to SurrealDB")

	_, err = surrealdb.Query[[]UGC](db, "SELEC * FROM state", map[string]interface{}{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

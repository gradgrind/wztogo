package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

/*
I need a list of database nodes, presumably containing their own DB-key.
The question is how they get built ...
*/

type Node map[string]interface{}

func (n Node) GetId() int {
	return n["Id"].(int)
}

func test1(dbpath string) {
	var data []Node
	m1 := Node{}
	m1["Id"] = 2
	m1["A"] = 1
	m1["L"] = []interface{}{"a", "b", "c", 5}
	data = append(data, m1)
	m2 := Node{}
	m2["Id"] = 4
	m2["A"] = 1
	m2["L"] = []interface{}{"a", "b", "d", 6}
	data = append(data, m2)
	dbwrite(dbpath, data)
}

func dbwrite(dbpath string, nodelist []Node) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
DROP TABLE IF EXISTS NODES;
CREATE TABLE NODES(
	Id INTEGER PRIMARY KEY,
	DATA TEXT NOT NULL
);
`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = "INSERT INTO NODES(Id, DATA) values(?,?)"
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// The primary key will correspond to the node indexes.
	for _, node := range nodelist {
		j, err := json.Marshal(node)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		fmt.Printf(" ??? %+v\n", node)
		_, err = tx.Exec(query, node.GetId(), string(j))
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

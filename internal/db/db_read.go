package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func dbread(dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM NODES")
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var Id int
		var DATA []byte

		err = rows.Scan(&Id, &DATA)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID %d: %s\n", Id, DATA)
		var jsonobj map[string]interface{}
		err = json.Unmarshal([]byte(DATA), &jsonobj)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("  ->: %+v\n", jsonobj)

		ix, ok := jsonobj["X"]
		if ok {
			fx := ix.(float64)
			fmt.Printf("  +++: %d\n", int(fx))
		}

		v := jsonobj["CONSTRAINTS"].(map[string]interface{})["MaxDays"]

		i, err := strconv.Atoi(v.(string))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  +++: %d\n", i)

		m2 := jsonobj["CONSTRAINTS"].(map[string]interface{})
		fmt.Printf("  ***: %d\n", jsoni(&m2, "MinLessonsPerDay"))
	}
}

func jsoni(map0 *map[string]interface{}, field string) int {
	m := *map0
	i, err := strconv.Atoi(m[field].(string))
	if err != nil {
		panic(err)
	}
	return i
}

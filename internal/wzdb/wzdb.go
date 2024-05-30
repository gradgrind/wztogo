package wzdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func dbtest() {
	x := "245"
	y, e := strconv.Atoi(x)
	if e == nil {
		fmt.Printf("%T: %v\n\n", y, y)
	}

	db, err := sql.Open("sqlite3", "../_testdata/db365.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)

	rows, err := db.Query("SELECT * FROM NODES WHERE DB_TABLE=?", "TEACHERS")
	if err != nil {
		//panic(err)
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var DB_TABLE string
		var DATA []byte

		err = rows.Scan(&id, &DB_TABLE, &DATA)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID %d: %s\n", id, DATA)
		var jsonobj map[string]interface{}
		err = json.Unmarshal([]byte(DATA), &jsonobj)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("  ->: %s\n", jsonobj)

		ix, ok := jsonobj["#"]
		if !ok {
			log.Fatal("#")
		}
		fx := ix.(float64)
		fmt.Printf("  +++: %d\n", int(fx))

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

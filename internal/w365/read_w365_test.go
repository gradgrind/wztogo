package w365

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	// "regexp"
	//"gradgrind/wztogo/internal/wzbase"
	//"github.com/RoaringBitmap/roaring"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ncruces/zenity"
)

func TestReadW365(t *testing.T) {
	fmt.Println("\n############## TestReadW365")
	const defaultPath = "../_testdata/*.w365"
	f365, err := zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilter{
			Name:     "Waldorf-365 files",
			Patterns: []string{"*.w365"},
			CaseFold: false,
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n ***** Reading %s *****\n", f365)

	db := ReadW365Raw(f365)
	db.ReadYear(db.ActiveYear)
	for _, yeardata := range db.Years {
		fmt.Printf("\n$$$ %#v\n", yeardata)
	}
	db.read_days()
	fmt.Printf("\n§§NodeList: %#v\n", db.NodeList)
	fmt.Printf("\n§§NodeMap: %#v\n", db.NodeMap)
	fmt.Printf("\n§§TableMap: %#v\n", db.TableMap)
	db.read_hours()
	fmt.Printf("\n§§Config: %#v\n", db.Config)
	db.read_subjects()
	db.read_rooms()
	db.read_absences()
	fmt.Printf("\n§§absences: %#v\n", db.absencemap)
	db.read_categories()
	fmt.Printf("\n§§categories: %#v\n", db.categorymap)
	db.read_teachers()
	db.read_groups()
	fmt.Println("  =======================================================")
	db.read_activities()
	for i, n := range db.NodeList {
		fmt.Printf("\n§node %4d: %#v\n", i, n)
	}
	fmt.Println("\n****************************************************")
	schedules := db.read_lesson_times()
	scheduleNames := []string{}
	scheduleLessons := [][]wzbase.Lesson{}
	for _, xn := range schedules {
		scheduleNames = append(scheduleNames, xn.name)
		c_l := db.read_course_lessons(xn.lessons)
		scheduleLessons = append(scheduleLessons, c_l)
		// At least the initialized activities should be added to the
		// database. Here all activities (including uninitialized ones)
		// are added as a "lesson plan", named as the w365 schedule.
		entry := wzbase.LessonPlan{ID: xn.name, LESSONS: c_l}
		db.add_node("LESSON_PLANS", entry, "")
	}

	// Save data to (new) sqlite file
	dbfile := strings.TrimSuffix(f365, filepath.Ext(f365)) + ".sqlite"
	os.Remove(dbfile)
	dbx, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()

	var version string
	err = dbx.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n(SQLITE_VERSION = %s)\n", version)
	query := `
    CREATE TABLE IF NOT EXISTS NODES(
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        DB_TABLE TEXT NOT NULL,
        DATA TEXT NOT NULL
    );
    `
	_, err = dbx.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n  *** Tables ***")
	query = "INSERT INTO NODES(DB_TABLE, DATA) values(?,?)"
	tx, err := dbx.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// The primary key will correspond to the node indexes.
	for _, wznode := range db.NodeList[1:] {
		j, err := json.Marshal(wznode.Node)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		_, err = tx.Exec(query, wznode.Table, string(j))
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n Database saved to: %s\n", dbfile)

	if len(scheduleNames) == 0 {
		fmt.Println("\n No Schedule")
	} else {
		var plan_name string
		var plan_index int
		if len(scheduleNames) == 1 {
			plan_name = scheduleNames[0]
			plan_index = 0
			err = zenity.Question(
				fmt.Sprintf("Show Schedule '%s'?", plan_name),
				zenity.Title("Question"),
				zenity.QuestionIcon)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			plan_name, err = zenity.ListItems(
				"Select a 'Schedule' (timetable)",
				scheduleNames...)
			if err != nil {
				log.Fatal(err)
			}
			plan_index = slices.Index(scheduleNames, plan_name)
		}
		fmt.Printf("\n Schedule: %s\n", plan_name)
		//fmt.Printf("\n == %s: %+v\n", plan_name, schedules[plan_index].lessons)
		for ci, ll := range scheduleLessons[plan_index] {
			fmt.Printf("\n%4d: %+v\n", ci, ll)
		}
	}
}

/*
func TestMisc(t *testing.T) {
	fmt.Println("\n############## TestMisc")
	fmt.Println(convert_date("24. 12. 2023"))
	fmt.Println(convert_colour("-16777216"))
	fmt.Println(convert_colour("-47834"))
	fmt.Println(convert_colour("-16000000"))

	ndays := 5
	absence_map := make([]([]int), ndays)
	for i := range ndays {
		absence_map[i] = []int{}
	}
	absence_map[1] = append(absence_map[1], 5)
	fmt.Printf("absence_map: %v\n", absence_map)

	var x int = 5
	xp := &x
	*xp += 1
	fmt.Printf("XXX %d\n", x)
	type sx struct {
		a int
		b []int
	}
	x1 := sx{a: 1}
	fmt.Printf("XXX1a: %+v\n", x1)
	x1.b = append(x1.b, 7)
	fmt.Printf("XXX1b: %+v\n", x1.b)
	x1p := &x1.b
	*x1p = append(*x1p, 8)
	fmt.Printf("XXX1c: %+v\n", x1.b)
	x2 := map[int]*sx{}
	x2[10] = &sx{a: 11}
	fmt.Printf("XXX2a: %+v\n", x2[10])
	x2x := *x2[10]
	fmt.Printf("XXX2b: %+v\n", x2x)
	x2x.a += 1
	fmt.Printf("XXX2c: %+v\n", x2x)
	fmt.Printf("XXX2d: %+v\n", x2[10])
	x2p := x2[10]
	(*x2p).a += 2
	fmt.Printf("XXX2e: %+v\n", x2[10])
	x3 := map[int]int{}
	x3[20] = 21
	x3[20]++
	fmt.Printf("XXX3: %+v\n", x3)
	x4 := map[int]*sx{}
	x4[30] = &sx{a: 41}
	x4[30].a++
	fmt.Printf("XXX4a: %+v\n", x4[30])
	x4p := x4[30]
	x4p.a++
	fmt.Printf("XXX4b: %+v\n", x4[30])
	x4pp := &(x4p.a)
	*x4pp++
	fmt.Printf("XXX4c: %+v\n", x4[30])

	fmt.Printf("??? %#v\n", *new(string))

	il := map[string][]int{}
	il["A1"] = append(il["A1"], 1)
	il["A1"] = append(il["A1"], 2)
	il["A2"] = append(il["A2"], 3)
	fmt.Printf("\n il: %+v\n", il)
}
*/

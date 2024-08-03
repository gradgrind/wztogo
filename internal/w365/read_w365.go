// Package w365 provides functions for reading and processing data files
// from "Waldorf 365" (https://waldorf365.de/).
package w365

import (
	"bufio"
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Read a Waldorf 365 data file and divide it into "items" (individual
// table entries).
//
// Note that the XML version of the data should probably be used in
// preference to the basic text dumps read here. However, the latter form
// is available offline in the "Planer" app for non-cloud projects,
// whereas the XML version is (currently?) only available from the cloud.
//
// Return a "W365Data" value containing:
//   - the "SchoolState" item,
//   - a mapping, year-tag -> "Scenario" item,
//   - the year-tag of the "active" year,
//   - a list of all other items.
//
// ------------------------------------------------------------------------
func ReadW365Raw(fpath string) W365Data {
	// open file
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()
	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	//TODO: The default buffer size should be plenty, but the Waldorf 365
	// text editor is very messy and retains no end of junk, including
	// office markup, images, etc. I hope a newer version will perform
	// adequate sanitizing!
	bufsize := 100000
	buffer := make([]byte, bufsize)
	scanner.Buffer(buffer, bufsize)

	var item ItemType
	var schoolstate ItemType
	scenarios := []ItemType{}
	tables := map[string][]ItemType{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			item = nil
			continue
		}
		if strings.HasPrefix(line, "*") {
			if line == "*" {
				break
			}
			line = strings.TrimLeft(line, "*")
			item = make(ItemType)
			if line == w365_Scenario {
				scenarios = append(scenarios, item)
				continue
			}
			if line == w365_SchoolState {
				schoolstate = item
				continue
			}
			tables[line] = append(tables[line], item)
			continue
		}
		if item == nil {
			continue
		}
		k, v, found := strings.Cut(line, "=")
		if !found {
			log.Fatalf("Invalid input line: %s\n", line)
		}
		item[k] = v
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	scenario_map := map[string]YearData{}
	aid := schoolstate[w365_ActiveScenario]
	ayear := ""
	for _, item := range scenarios {
		fp, err := strconv.ParseFloat(item[w365_EpochFactor], 64)
		if err != nil {
			log.Fatal(err)
		}
		tag := item[w365_Name]
		yid := item[w365_Id]
		if yid == aid {
			ayear = tag
		}
		scenario_map[tag] = YearData{
			Tag:         tag,
			Name:        item[w365_Description],
			DATE_Start:  convert_date(item[w365_Start]),
			DATE_End:    convert_date(item[w365_End]),
			EpochFactor: fp,
			W365Id:      yid,
			LastChanged: item[w365_LastChanged],
		}
	}
	return W365Data{
		Schooldata: ItemType{
			"CountryCode": schoolstate["CountryCode"],
			"SchoolName":  schoolstate["SchoolName"],
			"StateCode":   schoolstate["StateCode"],
		},
		Years:            scenario_map,
		ActiveYear:       ayear,
		tables0:          tables,
		group_classgroup: map[int]wzbase.ClassGroup{},
		class_group_div:  map[int]map[int]int{},
	}
}

func (w365data *W365Data) add_node(
	table string,
	node interface{},
	key string,
) int {
	i := len(w365data.NodeList)
	//fmt.Printf("  +++++ %4d %s: %s \n", i, table, key)
	w365data.NodeList = append(w365data.NodeList, wzbase.WZnode{
		Table: table, Node: node,
	})
	if key != "" {
		w365data.NodeMap[key] = i
	}
	w365data.TableMap[table] = append(w365data.TableMap[table], i)
	return i
}

// Collect the data for the selected year.
func (w365data *W365Data) ReadYear(year string) {
	// For the first call it is not necessary to (re)initialize the slice,
	// but if the structure is reused, it needs to be cleared.
	w365data.NodeList = []wzbase.WZnode{}
	// Add a dummy entry at index 0.
	w365data.NodeList = append(w365data.NodeList, wzbase.WZnode{})
	// Maps must be initialized anyway.
	w365data.NodeMap = map[string]int{}
	w365data.TableMap = map[string][]int{}
	w365data.Config = map[string]interface{}{}
	w365data.absencemap = map[string]wzbase.Timeslot{}
	w365data.categorymap = map[string]Category{}
	year_data := w365data.Years[year]
	containerId := year_data.W365Id
	yeartables := map[string][]ItemType{}
	// Filter the items, retain only those for the chosen year.
	for tbl, itemlist := range w365data.tables0 {
		newlist := []ItemType{}
		for _, item := range itemlist {
			if item[w365_ContainerId] == containerId {
				newlist = append(newlist, item)
			}
		}
		// Sort the table according to the "ListPosition" entries.
		slices.SortFunc(newlist, func(a, b ItemType) int {
			af, err := strconv.ParseFloat(a[w365_ListPosition], 64)
			if err != nil {
				log.Fatal(err)
			}
			bf, err := strconv.ParseFloat(b[w365_ListPosition], 64)
			if err != nil {
				log.Fatal(err)
			}
			return cmp.Compare(af, bf)
		})
		yeartables[tbl] = newlist
	}
	w365data.yeartables = yeartables
	w365data.Yeardata = map[string]interface{}{
		"YEAR":            year_data.Tag,
		"SCHOOLYEAR":      year_data.Name,
		"DATE_Start":      year_data.DATE_Start,
		"DATE_End":        year_data.DATE_End,
		"EpochFactor":     year_data.EpochFactor,
		"LastChanged":     year_data.LastChanged,
		"SourceReference": containerId,
	}
	//TODO--
	fmt.Printf("\n §§§§§§§§§§§§§§§§§§§§§§§\n %+v\n §§§§§§§§§§§§§§§§§§§§§§§\n",
		w365data.Yeardata)
}

// Convert an input date. In the Waldorf365 data dumps, "DateOfBirth" fields
// are not ISO dates, but like "22. 01. 2015" (how reliable is this?).
// Return an ISO date.
func convert_date(date string) string {
	// return date     # the dates are correct in the xml dump format
	dmy := strings.Split(date, ". ")
	return fmt.Sprintf("%s-%s-%s", dmy[2], dmy[1], dmy[0])
}

// Convert a colour. In the Waldorf365 data dumps, colours are represented
// as negative integers (how reliable is this?).
// Return a 6-digit hex RGB colour (as "#RRGGBB").
func convert_colour(colour string) string {
	i, err := strconv.ParseInt(colour, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("#%06X", 0x1000000+i)
}

// ************ Tying it together ************

// Read the data from the given file for the "active" year.
// Within various elements there can be references to other elements. These
// are generally saved as integers corresponding to the primary keys of the
// referenced elements within the generated database.
// The internal data structure which is produced uses a (potentially)
// different indexing. All elements ("nodes") are held in a vector, NodeList.
// There is a mapping, IndexMap, which maps the primary-key-references to the
// indexes within NodeList.
// Note that "days" and "hours", which also have nodes within NodeList are
// handled differently. They are referenced by 0-based indexes because of
// the high significance of their ordering. For many purposes, no other
// information about these elements will be required, but if it is, appropriate
// maps can be generated – the entries in the lists TableMap["DAYS"] and
// TableMap["HOURS"] are correctly ordered (the values are normal primary-key
// references).
func ReadW365(w365file string) wzbase.WZdata {
	db365 := ReadW365Raw(w365file)
	//TODO: Might one want to select a different year?
	db365.ReadYear(db365.ActiveYear)
	db365.read_days()
	db365.read_hours()
	db365.read_subjects()
	db365.read_rooms()
	db365.read_absences()
	db365.read_categories()
	db365.read_teachers()
	db365.read_groups()
	db365.read_activities()
	schedules := db365.read_lesson_times()

	// Add all schedules to the database
	for _, xn := range schedules {
		c_l := db365.read_course_lessons(xn.lessons)
		// At least the initialized activities should be added to the
		// database. Here all activities (including uninitialized ones)
		// are added as a "lesson plan", named as the w365 schedule.
		entry := wzbase.LessonPlan{ID: xn.name, LESSONS: c_l}
		db365.add_node("LESSON_PLANS", entry, "")
	}

	// Save data to (new) sqlite file
	dbfile := "../_testdata/db365.sqlite"
	os.Remove(dbfile)
	dbx, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer dbx.Close()
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
	query = "INSERT INTO NODES(DB_TABLE, DATA) values(?,?)"
	tx, err := dbx.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// The primary key will correspond to the node indexes.
	for _, wznode := range db365.NodeList[1:] {
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
	// The index map is in this case an identity mapping ...
	imap := make(map[int]int, len(db365.NodeList))
	for i := range db365.NodeList {
		imap[i] = i
	}
	// Combine the various bits of school & configuration data.
	sdata := map[string]interface{}{}
	for k, v := range db365.Schooldata {
		sdata[k] = v
	}
	for k, v := range db365.Yeardata {
		sdata[k] = v
	}
	for k, v := range db365.Config {
		sdata[k] = v
	}
	// Generate atomic groups for all classes.
	// Include only divisions containing groups which are used in the
	// timetable.
	// Groups not in divisions (empty division Tag) may not be used in the
	// timetable and will be caught here.
	cagdivs := map[int][][]int{}
	ag := wzbase.NewAtomicGroups()
	for _, nc := range db365.TableMap["CLASSES"] {
		node := db365.NodeList[nc].Node.(wzbase.Class)
		//fmt.Printf("\n+++CLASS: %+v\n", node)
		agdivs := [][]int{}
		for _, div := range node.DIVISIONS {
			//fmt.Printf("  ---DIV: %+v\n", div)
			for _, g := range div.Groups {
				//fmt.Printf("    ~~~GROUP: %+v\n", g)
				if db365.ActiveGroups[g] {
					if div.Tag == "" {
						log.Fatalf(
							"Active group (%s) not in class division: %+v",
							wzbase.ClassGroup{
								CIX: nc, GIX: g,
							}.Print(db365), node,
						)
					}
					agdivs = append(agdivs, div.Groups)
					for _, g = range div.Groups {
						if !db365.ActiveGroups[g] {
							log.Printf("Group %s has no activities\n",
								wzbase.ClassGroup{
									CIX: nc, GIX: g,
								}.Print(db365),
							)
						}
					}
					break
				}
			}
		}
		cagdivs[nc] = agdivs
		ag.Add_class_groups2(nc, agdivs)
	}
	/*
		fmt.Println("\n  ****************************************")
		fmt.Printf(" Number of atomic groups: %d\n", ag.X)
		fmt.Printf(" ClassGroups: %+v\n", ag.Class_Groups)
		fmt.Printf(" 12K: %+v //// 12: %+v\n",
			ag.Class_Groups[251], ag.Class_Groups[250],
		)
		fmt.Println("  ++++++++++++++++++++++++++++++++++++++++++")
		fmt.Printf(" For class 12: %+v\n", ag.Group_Atomics[wzbase.ClassGroup{
			CIX: 250, GIX: 0,
		}])
		fmt.Println("  ****************************************")
	*/

	// Build a mapping from the reference (index/primary key) back to the
	// Waldorf 365 node Id, where possible.
	ref2key := make(map[int]string, len(db365.NodeMap))
	for k, v := range db365.NodeMap {
		ref2key[v] = k
	}
	return wzbase.WZdata{
		Schooldata: sdata,
		NodeList:   db365.NodeList,
		IndexMap:   imap,
		TableMap:   db365.TableMap,
		// GroupDiv:     db365.GroupDiv,
		// AtomicGroups: db365.AtomicGroups,
		GroupClassgroup:  db365.group_classgroup,
		ActiveDivisions:  cagdivs,
		AtomicGroups:     ag,
		SourceReferences: ref2key,
	}
}

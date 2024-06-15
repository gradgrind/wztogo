// Package w365 provides functions for reading and processing data files
// from "Waldorf 365" (https://waldorf365.de/).
package w365

import (
	"bufio"
	"cmp"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
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
			if line == "Scenario" {
				scenarios = append(scenarios, item)
				continue
			}
			if line == "SchoolState" {
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
	aid := schoolstate["ActiveScenario"]
	ayear := ""
	for _, item := range scenarios {
		fp, err := strconv.ParseFloat(item["EpochFactor"], 64)
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
			Name:        item["Decription"],
			DATE_Start:  convert_date(item["Start"]),
			DATE_End:    convert_date(item["End"]),
			EpochFactor: fp,
			W365Id:      yid,
			LastChanged: item["LastChanged"],
		}
	}
	return W365Data{
		Schooldata: schoolstate, //TODO: filter
		Years:      scenario_map,
		ActiveYear: ayear,
		tables0:    tables,
	}
}

func (w365data *W365Data) add_node(
	table string,
	node interface{},
	key string,
) int {
	i := len(w365data.NodeList)
	fmt.Printf("  +++++ %4d %s: %s \n", i, table, key)
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

	/* TESTING
	w365data.add_node("TEST1", 100, "$1")
	w365data.add_node("TEST1", map[string]string{"VAL1": "V1"}, "$2")
	w365data.add_node("TEST2", map[string]int{"VAL2": 2}, "$3")
	fmt.Printf("\n &&&&&1 %#v\n", w365data.NodeList)
	fmt.Printf("\n &&&&&2 %#v\n", w365data.NodeMap)
	fmt.Printf("\n &&&&&3 %#v\n", w365data.TableMap)
	nt0_ := w365data.NodeList[w365data.NodeMap["$2"]]
	fmt.Printf("\n &&&&&4 %#v\n", nt0_.Node.(map[string]string)["VAL1"])
	*/

	containerId := w365data.Years[year].W365Id
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

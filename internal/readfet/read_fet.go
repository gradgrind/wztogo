package readfet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"io"
	"log"
	"os"
	"strings"
)

type Day struct {
	//XMLName xml.Name `xml:"Day"`
	Name      string
	Long_Name string
}

type Days_List struct {
	//XMLName xml.Name `xml:"Days_List"`
	Day []Day
}

type Hour struct {
	//XMLName xml.Name `xml:"Hour"`
	Name      string
	Long_Name string
}

type Hours_List struct {
	//XMLName xml.Name `xml:"Hours_List"`
	Hour []Hour
}

type Activity struct {
	//XMLName xml.Name `xml:"Activity"`
	Id       int
	Duration int
	Comments string
}

type Activities_List struct {
	//XMLName  xml.Name `xml:"Activities_List"`
	Activity []Activity
}

type Subject struct {
	//XMLName xml.Name `xml:"Subject"`
	Name      string
	Long_Name string
	Comments  string
}

type Subjects_List struct {
	//XMLName xml.Name `xml:"Subjects_List"`
	Subject []Subject
}

type Teacher struct {
	//XMLName xml.Name `xml:"Teacher"`
	Name      string
	Long_Name string
	Comments  string
}

type Teachers_List struct {
	//XMLName xml.Name `xml:"Teachers_List"`
	Teacher []Teacher
}

// A very special class and group structure is expected here. This form must
// be generated specially, fet cannot do it itself!
// TODO: It may be possible to support the normal fet group structure using
// fet Categories.
type ClassGroup struct {
	XMLName xml.Name `xml:"Group"`
	Name    string
}

type ClassData struct {
	XMLName   xml.Name `xml:"Year"`
	Name      string
	Long_Name string
	// At present the comments field is using to convey the divisions info.
	Comments string
	Group    []ClassGroup
}

const GROUP_SEP = ","
const DIV_SEP = "|"
const CLASS_GROUP_SEP = "."

type Students_List struct {
	//XMLName xml.Name `xml:"Students_List"`
	Year []ClassData
}

type Room struct {
	//XMLName xml.Name `xml:"Room"`
	Name      string
	Long_Name string
	Virtual   bool
	Comments  string
}

type Rooms_List struct {
	//XMLName xml.Name `xml:"Rooms_List"`
	Room []Room
}

type ConstraintActivityPreferredStartingTime struct {
	//XMLName            xml.Name `xml:"ConstraintActivityPreferredStartingTime"`
	Activity_Id        int
	Preferred_Day      string
	Preferred_Hour     string
	Permanently_Locked bool
}

type ConstraintActivityPreferredRoom struct {
	//XMLName     xml.Name `xml:"ConstraintActivityPreferredRoom"`
	Activity_Id int
	Room        string
	Real_Room   []string
	//Permanently_Locked bool
}

type Time_Constraints_List struct {
	//XMLName       xml.Name `xml:"Time_Constraints_List"`
	ConstraintActivityPreferredStartingTime []ConstraintActivityPreferredStartingTime
}

type Space_Constraints_List struct {
	//XMLName    xml.Name `xml:"Space_Constraints_List"`
	ConstraintActivityPreferredRoom []ConstraintActivityPreferredRoom
}

type Result struct {
	XMLName                xml.Name `xml:"fet"`
	Institution_Name       string
	Comments               string
	Days_List              Days_List
	Hours_List             Hours_List
	Activities_List        Activities_List
	Rooms_List             Rooms_List
	Subjects_List          Subjects_List
	Teachers_List          Teachers_List
	Students_List          Students_List
	Time_Constraints_List  Time_Constraints_List
	Space_Constraints_List Space_Constraints_List
}

//TODO: Adapt for wzbase.WZdata
/*

// WZdata represents all the data within the sqlite table "NODES".
// The nodes / db-rows may contain references to other nodes. These
// references are integers (> 0) and are the primary keys of the
// referenced nodes in the database.
// When the database is loaded into memory to produce this structure,
// the contiguous NodeList is produced. IndexMap is built to map
// the node references (primary keys) to the corresponding indexes in
// the NodeList.
// TableMap collects the node references (primary keys) of the entries
// of each "table" ("DB_TABLE" field, not a table within the sqlite
// database).
type WZdata struct {
	Schooldata       map[string]interface{}
	NodeList         []WZnode           // all the db rows
	IndexMap         map[int]int        // map reference to NodeList index
	TableMap         map[string][]int   // map table name to list of references
	GroupClassgroup  map[int]ClassGroup // map group/class index to ClassGroup
	ActiveDivisions  map[int][][]int
	AtomicGroups     AtomicGroups
	SourceReferences map[int]string
	// This one maps subject -> atomic-group -> activities:
	//SbjAGActivities  map[int]map[int]map[int]bool
}
*/

func add_node(
	wzdb *wzbase.WZdata,
	table string,
	node interface{},
	key string,
) int {
	i := len(wzdb.NodeList)
	//fmt.Printf("  +++++ %4d %s: %s \n", i, table, key)
	wzdb.NodeList = append(wzdb.NodeList, wzbase.WZnode{
		Table: table, Node: node,
	})
	//TODO: This is probably not what I need! Indeed this structure seems
	// to serve no really useful purpose in this context.
	if key != "" {
		wzdb.SourceReferences[i] = key
	}
	wzdb.TableMap[table] = append(wzdb.TableMap[table], i)
	return i
}

func ReadFet(fetpath string) string {
	// Open the  XML file
	xmlFile, err := os.Open(fetpath)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer xmlFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(xmlFile)
	log.Printf("*+ Reading: %s\n", fetpath)
	v := Result{}
	err = xml.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("XML error in %s:\n %v\n", fetpath, err)
	}

	// ***** Collect the data
	wzdb := wzbase.WZdata{}
	// Add a dummy entry at index 0.
	wzdb.NodeList = append(wzdb.NodeList, wzbase.WZnode{})
	// Maps must be initialized anyway.
	wzdb.TableMap = map[string][]int{}
	wzdb.Schooldata = map[string]interface{}{}
	wzdb.SourceReferences = map[int]string{}

	for i, d := range v.Days_List.Day {
		n := wzbase.Day{
			ID:   d.Name,
			NAME: d.Long_Name,
			X:    i,
		}
		add_node(&wzdb, "DAYS", n, d.Name)
	}
	for i, h := range v.Hours_List.Hour {
		ln := h.Long_Name
		ln_t := strings.Split(ln, "@")
		var st, et string
		if len(ln_t) == 2 {
			s_e := strings.Split(ln_t[1], "-")
			if len(s_e) == 2 {
				ln = ln_t[0]
				st = s_e[0]
				et = s_e[1]
			} else {
				log.Printf("Invalid Period times: %s\n", h.Long_Name)
			}
		}
		if st == "" {
			log.Printf("Period '%s' has no times.", h.Name)
		}
		n := wzbase.Hour{
			ID:         h.Name,
			NAME:       ln,
			X:          i,
			START_TIME: st,
			END_TIME:   et,
		}
		add_node(&wzdb, "HOURS", n, h.Name)
		fmt.Printf("??? %+v\n", n)
	}

	for i, s := range v.Subjects_List.Subject {
		ln := s.Long_Name
		if ln == "" {
			ln = s.Comments
		}
		n := wzbase.Subject{
			ID:   s.Name,
			NAME: ln,
			X:    i,
		}
		add_node(&wzdb, "SUBJECTS", n, s.Name)
		fmt.Printf("??? %+v\n", n)
	}

	for i, t := range v.Teachers_List.Teacher {
		ln := t.Long_Name
		if ln == "" {
			ln = t.Comments
		}
		n := wzbase.Subject{
			ID:   t.Name,
			NAME: ln,
			X:    i,
		}
		add_node(&wzdb, "TEACHERS", n, t.Name)
		fmt.Printf("??? %+v\n", n)
	}

	for i, r := range v.Rooms_List.Room {
		// Virtual rooms are not needed here.
		if r.Virtual {
			continue
		}
		n := wzbase.Room{
			ID:   r.Name,
			NAME: r.Long_Name,
			X:    i,
		}
		add_node(&wzdb, "ROOMS", n, r.Name)
		fmt.Printf("??? %+v\n", n)
	}

	for i, k := range v.Students_List.Year {
		class := k.Name
		divs := k.Comments
		divlist := []wzbase.DivGroups{}
		if divs != "" {
			for j, div := range strings.Split(divs, DIV_SEP) {
				glist := []int{}
				for _, g := range strings.Split(div, GROUP_SEP) {
					//TODO: Need to add a group node
					gi := add_node(&wzdb,
						"GROUPS",
						wzbase.Group{ID: g},
						fmt.Sprintf("%s%s%s", class, CLASS_GROUP_SEP, g))
					glist = append(glist, gi)
				}
				divlist = append(divlist, wzbase.DivGroups{
					Tag:    fmt.Sprintf("Div-%02d", j),
					Groups: glist,
				})
			}
		}
		n := wzbase.Class{
			ID:        k.Name,
			SORTING:   fmt.Sprintf("%02d:%s", i, k.Name),
			DIVISIONS: divlist,
		}
		//TODO? It may be worth waiting until all groups have been entered,
		// but this stuff won't be saved anywhere, so maybe not ...
		add_node(&wzdb, "CLASSES", n, k.Name)
		fmt.Printf("??? %+v\n", n)
	}

	log.Fatalln("QUITTING")
	return ""

	/*
		//fmt.Printf("*+ Hours: %+v\n", hourmap)
		amap := map[int]Activity{}
		for _, a := range v.Activities_List.Activity {
			amap[a.Id] = a
		}
		roommap := map[string]string{}
		for _, r := range v.Rooms_List.Room {
			roommap[r.Name] = r.Comments
		}
		yearid := v.Comments
		tim := time.Now().Format("2006-01-02T15:04:05")

		room_allocation := map[int]string{}
		for _, rdata := range v.Space_Constraints_List.
			ConstraintActivityPreferredRoom {
			//fmt.Printf("  -- %+v\n", rdata)
			ai := rdata.Activity_Id
			r := rdata.Room
			//rr := rdata.Real_Room
			room_allocation[ai] = roommap[r]
		}

		lessons := []string{}
		lids := []string{}
		for _, tdata := range v.Time_Constraints_List.ConstraintActivityPreferredStartingTime {
			aid := tdata.Activity_Id
			d := daymap[tdata.Preferred_Day]
			h := hourmap[tdata.Preferred_Hour]
			fixed := tdata.Permanently_Locked
			a := amap[aid]
			if a.Comments == "" {
				continue
			}
			for i := range a.Duration {
				lid := uuid.NewString()
				lids = append(lids, lid)
				les := []string{
					"*Lesson",
					fmt.Sprintf("ContainerId=%s", yearid),
					fmt.Sprintf("Course=%s", a.Comments),
					fmt.Sprintf("Day=%d", d),
					fmt.Sprintf("Hour=%d", h+i),
					fmt.Sprintf("Fixed=%t", fixed),
					fmt.Sprintf("Id=%s", lid),
					fmt.Sprintf("LastChanged=%s", tim),
					// W365 seems to accept only single rooms, which may be
					// room-groups, so just use the translation table.
					fmt.Sprintf("LocalRooms=%s", room_allocation[aid]),
					"",
				}
				lessons = append(lessons, les...)
			}
		}
		schedule_name := "fet001"
		list_pos := "100.0"
		schedule := []string{
			"*Schedule",
			fmt.Sprintf("ContainerId=%s", yearid),
			//f"End=",   #01. 03. 2024    # unused?
			fmt.Sprintf("Id=%s", uuid.NewString()),
			fmt.Sprintf("LastChanged=%s", tim), // 2024-03-30T18:59:53
			fmt.Sprintf("Lessons=%s", strings.Join(lids, "#")),
			fmt.Sprintf("ListPosition=%s", list_pos),
			fmt.Sprintf("Name=%s", schedule_name),
			"NumberOfManualChanges=0",
			//f"Start=",  #01. 03. 2024  # unused?
			"YearPercent=1.0",
			"",
		}
		schedule = append(schedule, lessons...)
		return strings.Join(schedule, "\n")

		/*
			//TODO: If single room constraints are used, placements with virtual
			// rooms will have two entries, only one will have the real rooms!

				room_allocation := make([][]string, len(activities))
				for _, rdata := range v.Space_Constraints_List.ConstraintActivityPreferredRoom {
					fmt.Printf("  -- %+v\n", rdata)
					ai := rdata.Activity_Id
					r := rdata.Room
					rr := rdata.Real_Room
					if len(rr) > 0 {
						room_allocation[ai-1] = rr
					} else {
						room_allocation[ai-1] = []string{r}
					}
				}
				for _, r := range room_allocation {
					fmt.Printf(" + %+v\n", r)
				}
	*/
}

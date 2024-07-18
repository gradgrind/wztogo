package wzbase

import (
	"fmt"
	"strings"
)

type Timeslot struct {
	Day  int
	Hour int
}

type Day struct {
	ID   string
	NAME string
	X    int
}

type Hour struct {
	ID         string
	NAME       string
	X          int
	START_TIME string
	END_TIME   string
}

type Teacher struct {
	ID         string
	LASTNAME   string
	FIRSTNAMES string
	// SEX int
	// ...
	CONSTRAINTS   map[string]string
	NOT_AVAILABLE []([]int)
}

type Subject struct {
	ID   string
	NAME string
	X    int
}

type Room struct {
	ID            string
	NAME          string
	X             int
	NOT_AVAILABLE []([]int)
}

// Roomspec describes the room requirements for a course. The total number
// of rooms required is len(Compulory) + len(Choices) + UserInput.
type RoomSpec struct {
	Compulsory []int
	Choices    [][]int
	UserInput  int // or []int ?
}

type Student struct {
	ID         string
	SORTNAME   string
	LASTNAME   string
	FIRSTNAMES string
	FIRSTNAME  string
	GENDER     string
	DATE_BIRTH string
	BIRTHPLACE string
	DATE_ENTRY string
	DATE_EXIT  string
	HOME       string
	POSTCODE   string
	STREET     string
	EMAIL      string
	PHONE      string
}

type Group struct {
	ID       string
	STUDENTS []int
}

type DivGroups struct {
	Tag    string
	Groups []int
}

/*
type DivIndexGroups struct {
	Div    int
	Groups []int
}
*/

type ClassDivGroups struct {
	Class  int
	Div    int
	Groups []int
}

type CourseGroups []ClassDivGroups

/*
	func (cgs CourseGroups) Print(nodelist []WZnode) string {
		gnlist := []string{}
		for _, cdgs := range cgs {
			for _, g := range cdgs.Groups {
				gnlist = append(gnlist, ClassGroup{cdgs.Class, g}.Print(nodelist))
			}
		}
		return strings.Join(gnlist, ",")
	}
*/

func (cgs CourseGroups) Print(ng NodeGetter) string {
	gnlist := []string{}
	for _, cdgs := range cgs {
		for _, g := range cdgs.Groups {
			gnlist = append(gnlist, ClassGroup{cdgs.Class, g}.Print(ng))
		}
	}
	return strings.Join(gnlist, ",")
}

// AddCourseGroups adds the groups from a course to a base CourseGroups
// item, when they are not already contained in the base. Unless the whole
// class is covered, added groups must be in the same division within a class.
func (cg0 *CourseGroups) AddCourseGroups(
	nodelist []WZnode, cg CourseGroups) bool {
	var divs DivGroups
	for _, cdg := range cg {
		for i, cdg0 := range *cg0 {
			if cdg0.Class == cdg.Class {
				if cdg0.Div != -1 {
					if cdg.Div == -1 {
						// Update to full class
						(*cg0)[i] = ClassDivGroups{Class: cdg0.Class, Div: -1}
					} else {
						// Check division compatibility
						if cdg0.Div != cdg.Div {
							return false
						}
						// Add groups which are not already contained
						for _, g := range cdg.Groups {
							for _, g0 := range cdg0.Groups {
								if g0 == g {
									goto skip
								}
							}
							cdg0.Groups = append(cdg0.Groups, g)
							// Test whether now full class
							divs = nodelist[cdg.Class].Node.(Class).DIVISIONS[cdg.Div]
							if len(cdg0.Groups) == len(divs.Groups) {
								// Substitute the whole class
								(*cg0)[i] = ClassDivGroups{Class: cdg0.Class, Div: -1}
							} else {
								(*cg0)[i] = cdg0
							}
						skip:
						}
					}
				}
				goto cfound
			}
		}
		// else: Add class
		*cg0 = append(*cg0, cdg)
	cfound:
	}
	return true
}

type Class struct {
	ID            string // normal short name of class
	SORTING       string // sortable short name of class
	BLOCK_FACTOR  float64
	STUDENTS      []int
	DIVISIONS     []DivGroups
	CONSTRAINTS   map[string]string
	NOT_AVAILABLE []([]int)
}

type ClassGroup struct {
	CIX int
	GIX int
}

/*
func (cg ClassGroup) Print(nodelist []WZnode) string {
	c := nodelist[cg.CIX].Node.(Class)
	if cg.GIX == 0 {
		return c.ID
	}
	g := nodelist[cg.GIX].Node.(Group)
	return fmt.Sprintf("%s.%s", c.ID, g.ID)
}
*/

func (cg ClassGroup) Print(ng NodeGetter) string {
	c := ng.GetNode(cg.CIX).(Class)
	if cg.GIX == 0 {
		return c.ID
	}
	g := ng.GetNode(cg.GIX).(Group)
	return fmt.Sprintf("%s.%s", c.ID, g.ID)
}

func (cg ClassGroup) Printx(ng NodeGetter) string {
	c := ng.GetNode(cg.CIX).(Class)
	if cg.GIX == 0 {
		return c.SORTING
	}
	g := ng.GetNode(cg.GIX).(Group)
	return fmt.Sprintf("%s.%s", c.SORTING, g.ID)
}

type Course struct {
	TEACHERS        []int
	GROUPS          CourseGroups
	SUBJECT         int
	ROOM_WISH       RoomSpec
	WORKLOAD        float64
	WORKLOAD_FACTOR float64
	LESSONS         []int
	BLOCK_UNITS     float64
	FLAGS           map[string]bool
}

/*
func (c Course) Print(nodelist []WZnode) string {
	g := c.GROUPS.Print(nodelist)
	s := nodelist[c.SUBJECT].Node.(Subject).ID
	tl := []string{}
	for _, ti := range c.TEACHERS {
		tl = append(tl, nodelist[ti].Node.(Teacher).ID)
	}
	return fmt.Sprintf("<%s-%s-%s>", g, s, strings.Join(tl, ","))
}
*/

func (c Course) Print(ng NodeGetter) string {
	g := c.GROUPS.Print(ng)
	s := ng.GetNode(c.SUBJECT).(Subject).ID
	tl := []string{}
	for _, ti := range c.TEACHERS {
		tl = append(tl, ng.GetNode(ti).(Teacher).ID)
	}
	return fmt.Sprintf("<%s-%s-%s>", g, s, strings.Join(tl, ","))
}

type Block struct {
	Tag           string
	Base          int
	Components    []int
	BlockGroups   CourseGroups
	BlockTeachers []int
	BlockRooms    RoomSpec
}

type Lesson struct {
	Day    int
	Hour   int
	Length int
	Rooms  []int
	Fixed  bool
	Course int
}

type WZnode struct {
	Table string
	Node  interface{}
}

type WZDB struct { // for the NODES table in the sqlite database
	Id       int    // primary key
	DB_TABLE string // table name
	DATA     string // JSON
}

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
}

func (wzdb WZdata) GetNode(ref int) interface{} {
	return wzdb.NodeList[wzdb.IndexMap[ref]].Node
}

type NodeGetter interface {
	GetNode(ref int) interface{}
}

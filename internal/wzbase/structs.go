package wzbase

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
	SUBROOMS      []int // room group: indexes of component rooms
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

type Course struct {
	TEACHERS        []int
	GROUPS          []int
	SUBJECT         int
	ROOM_WISH       []int
	WORKLOAD        float64
	WORKLOAD_FACTOR float64
	LESSONS         []int
	BLOCK_UNITS     float64
	FLAGS           map[string]bool
}

type Block struct {
	Tag        string
	Base       int
	Components []int
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
	Schooldata   map[string]interface{}
	NodeList     []WZnode         // all the db rows
	IndexMap     map[int]int      // map reference to NodeList index
	TableMap     map[string][]int // map table name to list of references
	AtomicGroups AtomicGroups
}

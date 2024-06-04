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
	ID           string // normal short name of class
	SORTING      string // sortable short name of class
	BLOCK_FACTOR float64
	STUDENTS     []int
	DIVISIONS    []DivGroups
}

type WZnode struct {
	Table string
	Node  interface{}
}

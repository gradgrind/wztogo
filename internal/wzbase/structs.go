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

type WZnode struct {
	Table string
	Node  interface{}
}

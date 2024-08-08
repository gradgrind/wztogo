package readfet

import "encoding/xml"

// The structures used for reading a fet result file

type Day struct {
	//XMLName xml.Name `xml:"Day"`
	Name      string
	Long_Name string
	X         int
}

type Days_List struct {
	//XMLName xml.Name `xml:"Days_List"`
	Day []Day
}

type Hour struct {
	//XMLName xml.Name `xml:"Hour"`
	Name      string
	Long_Name string
	X         int
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
	// Added by constraints:
	Room      string
	RealRooms []string
	Day       int
	Hour      int
	Fixed     bool
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
	X         int
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
	X         int
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

type ClassDivision struct {
	XMLName             xml.Name `xml:"Category"`
	Number_of_Divisions int
	Division            []string
}

type ClassData struct {
	XMLName   xml.Name `xml:"Year"`
	Name      string
	Long_Name string
	// At present the comments field is using to convey the divisions info.
	Comments             string
	Number_of_Categories int
	Separator            string
	Category             []ClassDivision
	Group                []ClassGroup
	X                    int
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
	X         int
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

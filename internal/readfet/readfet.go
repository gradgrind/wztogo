package readfet

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Activity struct {
	//XMLName xml.Name `xml:"Activity"`
	Id       int
	Duration int
	Id365    string `xml:"Comments"`
}

type Activities_List struct {
	//XMLName  xml.Name `xml:"Activities_List"`
	Activity []Activity
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
	YearId365              string   `xml:"Comments"`
	Activities_List        Activities_List
	Time_Constraints_List  Time_Constraints_List
	Space_Constraints_List Space_Constraints_List
}

func to_w365(fetpath string) {
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

	fmt.Printf(" --- Year-Id: %s\n", v.YearId365)
	fmt.Printf(" --- Activities_List:\n%+v\n", v.Activities_List)
	fmt.Printf(" --- Time_Constraints_List:\n%+v\n",
		v.Time_Constraints_List)
	fmt.Printf(" --- Space_Constraints_List:\n%+v\n",
		v.Space_Constraints_List)

	//TODO: Placements with virtual rooms will have two entries, only
	// one will have the real rooms!
}

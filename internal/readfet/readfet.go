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

	fmt.Println("========================================================")
	activities := v.Activities_List.Activity
	fmt.Printf("  Activities: %d\n", len(activities))
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
}

/*
tim := time.Now()
fmt.Printf("Go launched at %s\n", tim.Format("2006-01-02T15:04:05"))

uuid
----
uuid := uuid.NewString()
fmt.Println(uuid)


            lessons.extend([
                "*Lesson",
                f"ContainerId={container_id}",
                f'Course={cx}',
                f'Day={p["Day"]}',
                f'Hour={h}',
                f'Fixed={p["Fixed"]}',
                f"Id={lid}",
                f"LastChanged={date_time}",     # 2024-03-30T18:59:53
                f"ListPosition={lesson_index}",
#TODO
                #f"LocalRooms={}", # 0b5413dc-1420-478f-b266-212fed8d2564
                "",
            ])

*/

/*
func new_lesson() {
	lid := uuid.NewString()
	les := []string{
		"*Lesson",
		fmt.Sprintf("ContainerId=%s", container_id),
		fmt.Sprintf("Course=%s", course_id),
		fmt.Sprintf("Day=%d", day),
		fmt.Sprintf("Hour=%d", hour),
		fmt.Sprintf("Fixed=%t", fixed),
		fmt.Sprintf("Id=%s", lid),
		fmt.Sprintf("LastChanged=%s", date),
		//fmt.Sprintf("LocalRooms=%s", room),
	}

}
*/

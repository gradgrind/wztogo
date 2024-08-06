package readfet

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func to_w365(fetpath string) string {
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

	//fmt.Printf(" --- Year-Id: %s\n", v.Comments)

	daymap := map[string]int{}
	for i, d := range v.Days_List.Day {
		daymap[d.Name] = i
	}
	//fmt.Printf("*+ Days: %+v\n", daymap)
	hourmap := map[string]int{}
	for i, h := range v.Hours_List.Hour {
		hourmap[h.Name] = i
	}
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

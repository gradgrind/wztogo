package readfet

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func MakeW365(fetdata FetResult) string {
	// IMPORTANT: The top-level Comments field must contain the w365 Id of
	// the year ("Container") in question.
	yearid := fetdata.MainComment
	tim := time.Now().Format("2006-01-02T15:04:05")

	lessons := []string{}
	lids := []string{}
	for _, a := range fetdata.Activities {
		courseId := a.Comments
		if courseId == "" {
			// Not a w365 course
			continue
		}
		d := a.Day
		if d < 0 {
			// No placement (actually this shouldn't occur!)
			continue
		}
		h := a.Hour
		fixed := a.Fixed
		r := a.Room
		if r != "" {
			r = fetdata.Rooms[r].Comments
		}
		// Build the w365 lesson(s)
		for i := range a.Duration {
			lid := uuid.NewString()
			lids = append(lids, lid)
			les := []string{
				"*Lesson",
				fmt.Sprintf("ContainerId=%s", yearid),
				fmt.Sprintf("Course=%s", courseId),
				fmt.Sprintf("Day=%d", d),
				fmt.Sprintf("Hour=%d", h+i),
				fmt.Sprintf("Fixed=%t", fixed),
				fmt.Sprintf("Id=%s", lid),
				fmt.Sprintf("LastChanged=%s", tim),
				// W365 seems to accept only single rooms, which may be
				// room-groups, so just use the room's Comments field.
				fmt.Sprintf("LocalRooms=%s", r),
				"",
			}
			lessons = append(lessons, les...)
		}
	}
	// Bundle the lessons up in a w365 "Schedule"
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
		room_allocation := map[int]string{}
		for _, rdata := range v.Space_Constraints_List.
			ConstraintActivityPreferredRoom {
			//fmt.Printf("  -- %+v\n", rdata)
			ai := rdata.Activity_Id
			r := rdata.Room
			//rr := rdata.Real_Room
			room_allocation[ai] = roommap[r]
		}


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

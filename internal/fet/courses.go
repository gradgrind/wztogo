package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
)

type fetActivity struct {
	XMLName           xml.Name `xml:"Activity"`
	Id                int
	Teacher           []string
	Subject           string
	Students          []string
	Active            bool
	Total_Duration    int
	Duration          int
	Activity_Group_Id int
	Comments          string
}

type fetActivitiesList struct {
	XMLName  xml.Name `xml:"Activities_List"`
	Activity []fetActivity
}

// Generate the fet activties.
func getActivities(fetinfo *fetInfo,
	activities []wzbase.Activity,
	course2activities map[int][]int,
) {
	fixed_rooms := []fixedRoom{}
	room_choices := []roomChoice{}
	// Preprocessing because of courses with multiple lessons.
	course_act := map[int]fetActivity{}
	for ci, acts := range course2activities {
		var td int
		var aid int
		activity := activities[acts[0]]
		if len(acts) > 1 {
			td = 0
			for _, l := range acts {
				td += activities[l].Duration
			}
			aid = acts[0] + 1 // fet indexing starts at 1
		} else {
			td = activity.Duration
			aid = 0
		}
		// Teachers
		tlist := []string{}
		for _, ti := range activity.Teachers {
			tlist = append(tlist, fetinfo.ref2fet[ti])
		}
		// Subject
		sbj := fetinfo.ref2fet[activity.Subject]
		// Groups
		glist := []string{}
		for _, cg := range activity.Groups {
			c := fetinfo.ref2fet[cg.CIX]
			if cg.GIX == 0 {
				glist = append(glist, c)
			} else {
				glist = append(glist, c+"."+fetinfo.ref2fet[cg.GIX])
			}
		}
		course_act[ci] = fetActivity{
			//Id:                i + 1, // fet indexing starts at 1
			Teacher:           tlist,
			Subject:           sbj,
			Students:          glist,
			Active:            true,
			Total_Duration:    td,
			Activity_Group_Id: aid,
			Comments:          fetinfo.wzdb.SourceReferences[ci],
		}

		//TODO: This needs quite some work! It is currently not complete.
		// It also needs repeating for multiple activities, with
		// appropriate id field ...
		// And the constraints must be added to the fet file!
		addRoomConstraint(fetinfo,
			&fixed_rooms,
			&room_choices,
			acts[0]+1,
			activity.RoomNeeds,
		)

	}
	// Now generate the full list of fet activities
	starting_times := []startingTime{}
	items := []fetActivity{}
	for i, activity := range activities {
		ci := activity.Course
		fetact := course_act[ci]
		fetact.Id = i + 1
		fetact.Duration = activity.Duration
		items = append(items, fetact)

		// Activity placement
		day := activity.Day
		if day >= 0 {
			hour := activity.Hour
			starting_times = append(starting_times, startingTime{
				Weight_Percentage:  100,
				Activity_Id:        i + 1,
				Preferred_Day:      fetinfo.days[day],
				Preferred_Hour:     fetinfo.hours[hour],
				Permanently_Locked: true,
				Active:             true,
			})
		}

		//TODO: Rooms
		//TODO: subject mappings

	}
	fetinfo.fetdata.Activities_List = fetActivitiesList{
		Activity: items,
	}
	fetinfo.fetdata.Time_Constraints_List.ConstraintActivityPreferredStartingTime = starting_times
}

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
	subject_activities []wzbase.SubjectGroupActivities,
) {
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

func getCourses(fetinfo *fetInfo) {
	type cdata struct {
		lessons      []int
		teachers     []int
		subject      int
		students     wzbase.CourseGroups
		course_index int
	}
	courses := []cdata{}
	// Go through the blocks first.
	block_courses := map[int]bool{} // collect courses connected with blocks
	for _, bi := range fetinfo.wzdb.TableMap["BLOCKS"] {
		bnode := fetinfo.wzdb.GetNode(bi).(wzbase.Block)
		ci := bnode.Base
		block_courses[ci] = true
		n := fetinfo.wzdb.GetNode(ci).(wzbase.Course)
		for _, bci := range bnode.Components {
			block_courses[bci] = true
		}
		if len(n.LESSONS) == 0 {
			continue
		}
		courses = append(courses, cdata{
			lessons:      n.LESSONS,
			teachers:     bnode.BlockTeachers,
			subject:      n.SUBJECT,
			students:     bnode.BlockGroups,
			course_index: ci,
		})
	}
	// Add the other courses.
	for _, ci := range fetinfo.wzdb.TableMap["COURSES"] {
		if block_courses[ci] {
			continue
		} // blocks already handled
		n := fetinfo.wzdb.GetNode(ci).(wzbase.Course)
		if len(n.LESSONS) == 0 {
			continue
		}
		courses = append(courses, cdata{
			lessons:      n.LESSONS,
			teachers:     n.TEACHERS,
			subject:      n.SUBJECT,
			students:     n.GROUPS,
			course_index: ci,
		})
	}
	// Generate the fet activties.
	id := 0 // activity Id (starts at 1)
	items := []fetActivity{}
	for _, course := range courses {
		// Teachers
		tlist := []string{}
		for _, ti := range course.teachers {
			tlist = append(tlist, fetinfo.ref2fet[ti])
		}
		// Subject
		sbj := fetinfo.ref2fet[course.subject]
		// Groups
		glist := []string{}
		for _, cg := range course.students {
			c := fetinfo.ref2fet[cg.Class]
			if cg.Div < 0 {
				glist = append(glist, c)
			} else {
				for _, g := range cg.Groups {
					glist = append(glist,
						c+"."+fetinfo.ref2fet[g],
					)
				}
			}
		}
		// Individual lessons/activities
		ll := 0 // Get total duration
		for _, l := range course.lessons {
			ll += l
		}
		aid := 0
		if len(course.lessons) > 1 {
			aid = id + 1
		}
		source_ref := fetinfo.wzdb.SourceReferences[course.course_index]
		for _, l := range course.lessons {
			id++
			items = append(items, fetActivity{
				Id:                id,
				Teacher:           tlist,
				Subject:           sbj,
				Students:          glist,
				Active:            true,
				Total_Duration:    ll,
				Duration:          l,
				Activity_Group_Id: aid,
				Comments:          source_ref,
			})

			//TODO: Rooms
			//TODO: subject mappings
		}
	}
	fetinfo.fetdata.Activities_List = fetActivitiesList{
		Activity: items,
	}
}

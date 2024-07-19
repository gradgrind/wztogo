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

func getCourses(wzdb *wzbase.WZdata, ref2fet map[int]string) fetActivitiesList {
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
	for _, bi := range wzdb.TableMap["BLOCKS"] {
		bnode := wzdb.GetNode(bi).(wzbase.Block)
		ci := bnode.Base
		block_courses[ci] = true
		n := wzdb.GetNode(ci).(wzbase.Course)
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
	for _, ci := range wzdb.TableMap["COURSES"] {
		if block_courses[ci] {
			continue
		} // blocks already handled
		n := wzdb.GetNode(ci).(wzbase.Course)
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
			tlist = append(tlist, ref2fet[ti])
		}
		// Subject
		sbj := ref2fet[course.subject]
		// Groups
		glist := []string{}
		for _, cg := range course.students {
			c := ref2fet[cg.Class]
			if cg.Div < 0 {
				glist = append(glist, c)
			} else {
				for _, g := range cg.Groups {
					glist = append(glist,
						c+"."+ref2fet[g],
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
		source_ref := wzdb.SourceReferences[course.course_index]
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
	return fetActivitiesList{
		Activity: items,
	}
}

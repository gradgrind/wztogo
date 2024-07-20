package wzbase

import (
	"fmt"
	"log"
)

type Activity struct {
	Day       int          // 0-based day index (-1 => not placed)
	Hour      int          // 0-based period index
	Duration  int          // number of periods occupied by activity
	Fixed     bool         // true => activity may not be moved
	Rooms     []int        // list of actually occupied rooms (node references)
	Course    int          // course node reference
	Teachers  []int        // list of teacher node references
	Subject   int          // subject node reference
	Groups    []ClassGroup // list of class-groups
	RoomNeeds RoomSpec     // Specification of rooms needed for the activity
}

// Build a list of timetable "activities", corresponding to individual
// lessons, etc.
func GetActivities(wzdb *WZdata) []Activity {
	type cdata struct {
		lessons      []int
		teachers     []int
		subject      int
		students     CourseGroups
		rooms        RoomSpec
		course_index int
	}
	courses := []cdata{}
	// *** Go through the blocks first.
	block_courses := map[int]bool{} // collect courses connected with blocks
	for _, bi := range wzdb.TableMap["BLOCKS"] {
		bnode := wzdb.GetNode(bi).(Block)
		ci := bnode.Base
		block_courses[ci] = true
		n := wzdb.GetNode(ci).(Course)
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
			rooms:        bnode.BlockRooms,
			course_index: ci,
		})
	}
	// *** Add the other courses.
	for _, ci := range wzdb.TableMap["COURSES"] {
		if block_courses[ci] {
			continue
		} // blocks already handled
		n := wzdb.GetNode(ci).(Course)
		if len(n.LESSONS) == 0 {
			continue
		}
		courses = append(courses, cdata{
			lessons:      n.LESSONS,
			teachers:     n.TEACHERS,
			subject:      n.SUBJECT,
			students:     n.GROUPS,
			rooms:        n.ROOM_WISH,
			course_index: ci,
		})
	}
	// *** Generate the activities as a list.
	activities := []Activity{}
	// Build a mapping from course (ref) to its activities (list of indexes).
	course_activities := map[int][]int{}
	for _, course := range courses {
		// Groups
		glist := []ClassGroup{}
		for _, cg := range course.students {
			if cg.Div < 0 {
				glist = append(glist, ClassGroup{cg.Class, 0})
			} else {
				for _, g := range cg.Groups {
					glist = append(glist, ClassGroup{cg.Class, g})
				}
			}
		}
		// Individual lessons/activities
		lessons := []int{} // list of activity indexes
		for _, l := range course.lessons {
			ix := len(activities)
			lessons = append(lessons, ix)
			activities = append(activities, Activity{
				Day:       -1, // indicate "not placed"
				Duration:  l,
				Course:    course.course_index,
				Teachers:  course.teachers,
				Subject:   course.subject,
				Groups:    glist,
				RoomNeeds: course.rooms,
			})
		}
		course_activities[course.course_index] = lessons
		// Keep track of lessons in each subject for each group,
		// to support constraints keeping them on separate days.
		AddSubjectMappings(course.subject, glist, lessons)
	}
	return activities
}

// func (smap map[]) AddSubjectMappings(
func AddSubjectMappings(
	subject int,
	groups []ClassGroup,
	lessons []int) {

	//TODO

	fmt.Printf("*** AddSubjectMappings: %d | %+v | %+v\n",
		subject, groups, lessons)
}

func GetSchedules(wzdb *WZdata) []string {
	items := []string{}
	for _, lpi := range wzdb.TableMap["LESSON_PLANS"] {
		lp := wzdb.GetNode(lpi).(LessonPlan)
		items = append(items, lp.ID)
	}
	return items
}

// TODO
func SetLessons(wzdb *WZdata, schedule string, activities []Activity) {
	var lessons []Lesson
	for _, lpi := range wzdb.TableMap["LESSON_PLANS"] {
		lp := wzdb.GetNode(lpi).(LessonPlan)
		if lp.ID == schedule {
			lessons = lp.LESSONS
			goto sfound
		}
	}
	log.Fatalf("Schedule '%s' does not exist\n", schedule)
sfound:
	for _, ll := range lessons {

		//TODO--
		fmt.Printf("* lesson:\n %+v\n", ll)

		if ll.Day < 0 {
			continue // unplaced
		}

		//TODO: Actually, why are there unplaced lessons here? How did they
		// get here???
	}
}

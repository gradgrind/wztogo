package wzbase

import (
	"log"

	"github.com/RoaringBitmap/roaring"
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

type TimetableCourse struct {
	Lessons     []int
	Teachers    []int
	Subject     int
	Students    []ClassGroup
	Rooms       RoomSpec
	CourseIndex int
}

func SplitCourseGroups(groups CourseGroups) []ClassGroup {
	glist := []ClassGroup{}
	for _, cg := range groups {
		if cg.Div < 0 {
			glist = append(glist, ClassGroup{cg.Class, 0})
		} else {
			for _, g := range cg.Groups {
				glist = append(glist, ClassGroup{cg.Class, g})
			}
		}
	}
	return glist
}

// Build a list of timetable "activities", corresponding to individual
// lessons, etc.
func GetActivities(wzdb *WZdata) (
	[]Activity,
	map[int][]int, // course-ref -> list of activity indexes
	map[int][]TimetableCourse, // subject-ref -> list of TimetableCourse
) {
	courses := []TimetableCourse{}
	subject_courses := map[int][]TimetableCourse{}
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
		courses = append(courses, TimetableCourse{
			Lessons:     n.LESSONS,
			Teachers:    bnode.BlockTeachers,
			Subject:     n.SUBJECT,
			Students:    SplitCourseGroups(bnode.BlockGroups),
			Rooms:       bnode.BlockRooms,
			CourseIndex: ci,
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
		courses = append(courses, TimetableCourse{
			Lessons:     n.LESSONS,
			Teachers:    n.TEACHERS,
			Subject:     n.SUBJECT,
			Students:    SplitCourseGroups(n.GROUPS),
			Rooms:       n.ROOM_WISH,
			CourseIndex: ci,
		})
	}
	// *** Generate the activities as a list.
	activities := []Activity{}
	// Build a mapping from course (ref) to its activities (list of indexes).
	course_activities := map[int][]int{}
	for _, course := range courses {
		subject_courses[course.Subject] = append(subject_courses[course.Subject],
			course)
		// Individual lessons/activities
		lessons := []int{} // list of activity indexes
		for _, l := range course.Lessons {
			ix := len(activities)
			lessons = append(lessons, ix)
			activities = append(activities, Activity{
				Day:       -1, // indicate "not placed"
				Duration:  l,
				Course:    course.CourseIndex,
				Teachers:  course.Teachers,
				Subject:   course.Subject,
				Groups:    course.Students,
				RoomNeeds: course.Rooms,
			})
		}
		course_activities[course.CourseIndex] = lessons
		// Keep track of lessons in each subject for each group,
		// to support constraints keeping them on separate days.
		//AddSubjectMappings(wzdb, course.subject, glist, lessons)
	}
	return activities, course_activities, subject_courses
}

// Collect groups of activities which share subject and (some) students.
type SubjectGroupActivities struct {
	Subject    int
	Groups     []ClassGroup
	Activities []int
}

// Build lists of activities which share subject and group(s).
func SubjectActivities(
	wzdb *WZdata,
	// subject-ref -> list of TimetableCourse:
	subject_courses map[int][]TimetableCourse,
	// course-ref -> list of activity indexes:
	course2activities map[int][]int,
) []SubjectGroupActivities {
	type g_lessons struct {
		groups  *roaring.Bitmap
		lessons []int
	}
	sgalist := []SubjectGroupActivities{}
	for sbj, tclist := range subject_courses {
		// Deal with all courses for a single subject.
		gll := []g_lessons{}
		for _, tc := range tclist {
			// When a course has multiple lessons, this lesson-group should
			// be collected.
			lessons := course2activities[tc.CourseIndex]
			if len(lessons) > 1 {
				sgalist = append(sgalist, SubjectGroupActivities{
					Subject:    sbj,
					Groups:     tc.Students,
					Activities: lessons,
				})
			}
			// Seek courses with shared (atomic) groups
			groups := roaring.New()
			for _, cg := range tc.Students {
				groups.Or(wzdb.AtomicGroups.Group_Atomics[cg])
			}
			for _, gl := range gll {
				if gl.groups.Intersects(groups) {
					// Collect a lesson-group for each combination
					for _, l0 := range gl.lessons {
						for _, l1 := range lessons {
							lx := []int{l0, l1}
							sgalist = append(sgalist, SubjectGroupActivities{
								Subject:    sbj,
								Groups:     tc.Students,
								Activities: lx,
							})
						}
					}
				}
			}
			gll = append(gll, g_lessons{groups, lessons})
		}
	}
	return sgalist
}

func GetSchedules(wzdb *WZdata) []string {
	items := []string{}
	for _, lpi := range wzdb.TableMap["LESSON_PLANS"] {
		lp := wzdb.GetNode(lpi).(LessonPlan)
		items = append(items, lp.ID)
	}
	return items
}

func SetLessons(
	wzdb *WZdata,
	schedule string,
	activities []Activity,
	c2alist map[int][]int,
) {
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
lfor:
	for _, ll := range lessons {
		if ll.Day < 0 {
			log.Fatalf("Unset Lesson in schedule %s\n", schedule)
		}
		alist := c2alist[ll.Course]
		for _, ai := range alist {
			a := &activities[ai]
			if a.Day < 0 && a.Duration == ll.Length {
				a.Day = ll.Day
				a.Hour = ll.Hour
				a.Rooms = ll.Rooms
				a.Fixed = ll.Fixed
				continue lfor
			}
		}
		log.Printf("Lesson in schedule %s could not be set: %+v\n",
			schedule, ll)
	}
}

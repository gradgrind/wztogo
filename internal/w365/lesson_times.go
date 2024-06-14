package w365

import (
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"slices"
	"strconv"
	"strings"
)

type xschedule struct {
	sortnum float64
	name    string
	lessons []wzbase.Lesson
}

func (w365data *W365Data) read_lesson_nodes() map[string]wzbase.Lesson {
	lesson_map := map[string]wzbase.Lesson{}
	for _, node := range w365data.yeartables[w365_Lesson] {
		// w365 only has single slot lessons, so a local intermediate
		// form is necessary.
		c := node[w365_Course]
		if c == "" {
			// I am currently not accepting the Waldorf 365 "Epochenschienen".
			//TODO: This should be more forceful, perhaps fatal:
			log.Printf(
				"!!! Lesson without course, Id = %s\n", node[w365_Id],
			)
			continue
		}
		d, err := strconv.Atoi(node[w365_Day])
		if err != nil {
			log.Fatal(err)
		}
		h, err := strconv.Atoi(node[w365_Hour])
		if err != nil {
			log.Fatal(err)
		}
		lnode := wzbase.Lesson{
			Day:    d,
			Hour:   h,
			Fixed:  node[w365_Fixed] == "true",
			Course: w365data.NodeMap[c],
		}
		lr, ok := node[w365_LocalRooms]
		if ok {
			for _, r := range strings.Split(lr, LIST_SEP) {
				lnode.Rooms = append(lnode.Rooms, w365data.NodeMap[r])
			}
		}
		lesson_map[node[w365_Id]] = lnode
	}
	return lesson_map
}

// Get the existing "plans" (W365: schedule).
func (w365data *W365Data) read_lesson_times() []xschedule {
	lesson_map := w365data.read_lesson_nodes()
	schedules := []xschedule{}
	for _, node := range w365data.yeartables[w365_Schedule] {
		//wid := node[w365_Id]
		lidlist := node[w365_Lessons]
		lessons := []wzbase.Lesson{}
		for _, s := range strings.Split(lidlist, LIST_SEP) {
			l, ok := lesson_map[s]
			if ok {
				lessons = append(lessons, l)
			}
		}
		f, err := strconv.ParseFloat(node[w365_ListPosition], 32)
		if err != nil {
			log.Fatal(err)
		}
		schedules = append(schedules, xschedule{
			f, node[w365_Name], lessons,
		})
	}
	// Sort the schedules according to the Waldorf 365 ListPosition
	slices.SortFunc(schedules, func(a, b xschedule) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	return schedules
}

// My current preference is to ignore the W365 Epochen, using tagged
// "normal" courses instead.

//TODO: Need to specify which "Schedule" to use.
// I could consider using the one called "Vorlage" for the moment?
// In any case, I suppose all lessons should be read in, but only those
// with fixed time used for input to the placement algorithm.

func (w365data *W365Data) read_course_lessons(
	lessons []wzbase.Lesson, // the lessons from the chosen "schedule"
) map[int][]wzbase.Lesson {
	// Allocate the lessons in the "schedule" to their courses.
	course_lessons := map[int][]wzbase.Lesson{}
	for _, lesson := range lessons {
		course_lessons[lesson.Course] = append(
			course_lessons[lesson.Course], lesson,
		)
	}
	// Order the timeslots for each course
	for ci, ll := range course_lessons {
		slices.SortFunc(ll, func(a, b wzbase.Lesson) int {
			if a.Day < b.Day {
				return -1
			}
			if a.Day == b.Day && a.Hour < b.Hour {
				return -1
			}
			return 1
		})
		// Amalgamate contiguous lessons, checking against the course needs.
		course := w365data.NodeList[ci].Node.(wzbase.Course)
		nlist := course.LESSONS
		// Because of the way the course lesson lengths are determined
		// (see function "read_activities()"), all lengths are the same,
		// except possibly the last, which can be shorter.
		joined_lessons := []wzbase.Lesson{}
		var day, hour int
		var fixed bool
		var i0 int
		var lesson wzbase.Lesson
		for _, n := range nlist {
			if n == 1 {
				// Take the first lesson.
				if len(ll) > 0 {
					lesson = ll[0]
					ll = ll[1:]
					lesson.Length = 1
				} else {
					// Make an unplaced lesson.
					lesson = wzbase.Lesson{
						Day:    -1,
						Hour:   -1,
						Length: 1,
						Course: ci,
					}
					joined_lessons = append(joined_lessons, lesson)
					continue
				}
			}
			length := 0
			found := false
			for i, lesson := range ll {
				if length > 0 {
					if lesson.Day == day &&
						lesson.Hour == hour+1 &&
						lesson.Fixed == fixed {
						//TODO: What about comparing the room lists???
						// contiguous
						length += 1
						if length == n {
							// match found
							lesson = ll[i0]
							lesson.Length = n
							ll = append(ll[:i0], ll[i+1:]...)
							joined_lessons = append(joined_lessons, lesson)
							found = true
							break
						}
						// not long enough, continue seeking
						hour++
						continue
					}
					// Otherwise start counting again.
				}
				day = lesson.Day
				hour = lesson.Hour
				fixed = lesson.Fixed
				length = 1
				i0 = i
			}
			if !found {
				// Make an unplaced lesson.
				lesson = wzbase.Lesson{
					Day:    -1,
					Hour:   -1,
					Length: n,
					Course: ci,
				}
				joined_lessons = append(joined_lessons, lesson)
			}
		}
		if len(ll) != 0 {
			//TODO: a more helpful error message
			log.Fatalf("Lesson mismatch for course %d\n", ci)
		}
		course_lessons[ci] = joined_lessons
		//fmt.Printf("\n********* %3d: %+v\n", ci, joined_lessons)
	}
	return course_lessons
}

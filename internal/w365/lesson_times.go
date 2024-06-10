package w365

import (
	"log"
	"slices"
	"strconv"
	"strings"
)

type xlesson struct {
	Course int
	Day    int
	Hour   int
	Rooms  []int
}

type xschedule struct {
	sortnum float64
	name    string
	lessons []xlesson
}

func (w365data *W365Data) read_lesson_nodes() map[string]xlesson {
	lesson_map := map[string]xlesson{}
	for _, node := range w365data.yeartables[w365_Lesson] {
		// w365 only has single slot lessons, so a local intermediate
		// form is necessary. Only use lessons with set times.
		if node[w365_Fixed] == "true" {
			c := node[w365_Course]
			if c == "" {
				//TODO: This should be more forceful, perhaps fatal:
				log.Printf(
					"!!! Fixed lesson without course, Id = %s\n", node[w365_Id],
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
			lnode := xlesson{Course: w365data.NodeMap[c], Day: d, Hour: h}
			lr, ok := node[w365_LocalRooms]
			if ok {
				for _, r := range strings.Split(lr, LIST_SEP) {
					lnode.Rooms = append(lnode.Rooms, w365data.NodeMap[r])
				}
			}
			lesson_map[node[w365_Id]] = lnode
		}
	}
	return lesson_map
}

func (w365data *W365Data) read_lesson_times() []xschedule {
	lesson_map := w365data.read_lesson_nodes()
	schedules := []xschedule{}
	for _, node := range w365data.yeartables[w365_Schedule] {
		//wid := node[w365_Id]
		lidlist := node[w365_Lessons]
		lessons := []xlesson{}
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

/*
#TODO: Need to specify which "Schedule" to use
    schedules = [
        (float(node[_ListPosition]), node[_Name], node[_Lessons])
        for node in w365_db.scenario[_Schedule]
    ]
    schedules.sort()
    #for _, n, _ in schedules:
    #    print(" +++", n)

# The "Vorlage" might have only fixed lessons.
# If adding or deleting lessons, the Lessons field of the Schedule
# must be updated, or a new Schedule must be built (maybe better).

#TODO: Assume the last schedule?
    isched = -1
# or the first?
    isched = 0
    lesson_ids = schedules[isched][-1].split(LIST_SEP)
    lesson_set = set(lesson_ids)
# (or maybe rather one with a particular name?)

#TODO?
#    w365_db.w365lessons = lesson_ids

# My current preference is to ignore the W365 Epochen, using tagged
# "normal" courses instead.

    for lid in lesson_ids:
        node = w365_db.idmap[lid]
        try:
            course_id = node[_Course]
        except KeyError:
#TODO
            REPORT_WARNING("Zeiten für Epochenschienen werden nicht berücksichtigt")
            continue
        if node[_Fixed] == "true":
            slot = (node[_Day], node[_Hour])
        else:
            slot = None
        # Add lesson id and time slot (if fixed) to course
        course_lessons[course_id].append((lid, slot))

    # Now deal with the individual lessons
    w365id_nodes.clear()
    for course_id, lslots in course_lessons.items():
        if lslots:
            lesson_times = set()
            for l_id, slot in lslots:
                #print("    ", l_id, slot)
                if slot:
                    lesson_times.add(slot)
            pltimes = process_lesson_times(lesson_times)
            #print(" --c--:", pltimes)
            k = w365_db.id2key[course_id]
            for ll, tlist in pltimes.items():
                for d, p in tlist:
                    xnode =  {
                        "LENGTH": str(ll),
                        "_Course": k,
                        "DAY": str(d),
                        "PERIOD": str(p),
                        "FIXED": "true",
                        #"_Parallel": 0,
                    }
                    w365id_nodes.append(("", xnode))
                    #print("     ++", xnode)
    # Add to database
    w365_db.add_nodes("LESSONS", w365id_nodes)
#TODO: Note that if I am only including "fixed" lessons, I don't need
# them to have a "FIXED" field!


#TODO: Might want to record the ids of non-fixed lessons as these entries
# might get changed? Actually, probably not, because I will probably
# generate a new Schedule.

# Do I need the EpochPlan to discover which teachers are involved in an
# Epoch, or can I get it from the Course entries somehow? No, this is really
# not ideal. There is a tenuous connection between "Epochenschienen" and
# courses only when an "Epochenplan" has been generated: there are then
# lessons which point to the course. Maybe for now I should collect the block
# times associated with the classes (I suppose using the EpochPlan to
# identify the classes is best? – it also supplies the name tag), then
# go through the block courses to find those in a block (test EpochWeeks?)
# and hence any other infos ... especially the teachers, I suppose.

# Für jede Klasse, die an einer Epoche beteiligt ist, gibt es einen Satz
# "Lessons", die identische Zeiten angeben. So entstehen viele überflüssige
# Einträge – es wäre besser, die "Lessons" mit der Epochen zu verknüpfen,
# einmalig.


#TODO: Might want to represent the Epochs as single course items in fet?
# That would be necessary if the teachers are included (but consider also
# the possibility of being involved in other Epochen (e.g. Mittelstufe),
# which might be different ... That's difficult to handle anyway.
# Perhaps it's easier to put no teachers in and block the teachers
# concerned in "Absences"?

*/

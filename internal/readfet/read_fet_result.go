package readfet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/timetable"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type FetResult struct {
	Institution string
	MainComment string
	Days        map[string]Day
	Hours       map[string]Hour
	Rooms       map[string]Room
	Subjects    map[string]Subject
	Teachers    map[string]Teacher
	Students    map[string]ClassData
	Activities  map[int]Activity
}

func ReadFetResult(fetpath string) FetResult {
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

	daymap := map[string]Day{}
	for i, d := range v.Days_List.Day {
		d.X = i
		daymap[d.Name] = d
	}
	//fmt.Printf("*+ Days: %+v\n", daymap)
	hourmap := map[string]Hour{}
	for i, h := range v.Hours_List.Hour {
		h.X = i
		hourmap[h.Name] = h
	}
	//fmt.Printf("*+ Hours: %+v\n", hourmap)
	roommap := map[string]Room{}
	for i, r := range v.Rooms_List.Room {
		r.X = i
		roommap[r.Name] = r
	}
	subjectmap := map[string]Subject{}
	for i, s := range v.Subjects_List.Subject {
		s.X = i
		subjectmap[s.Name] = s
	}
	teachermap := map[string]Teacher{}
	for i, t := range v.Teachers_List.Teacher {
		t.X = i
		teachermap[t.Name] = t
	}
	classmap := map[string]ClassData{}
	for i, sb := range v.Students_List.Year {
		sb.X = i
		classmap[sb.Name] = sb
	}
	setroom := map[int]ConstraintActivityPreferredRoom{}
	for _, rdata := range v.Space_Constraints_List.
		ConstraintActivityPreferredRoom {
		setroom[rdata.Activity_Id] = rdata
	}
	settime := map[int]ConstraintActivityPreferredStartingTime{}
	for _, tdata := range v.Time_Constraints_List.
		ConstraintActivityPreferredStartingTime {
		settime[tdata.Activity_Id] = tdata
	}
	amap := map[int]Activity{}
	for _, a := range v.Activities_List.Activity {
		ai := a.Id
		tdata, ok := settime[ai]
		if ok {
			a.Day = daymap[tdata.Preferred_Day].X
			a.Hour = hourmap[tdata.Preferred_Hour].X
			a.Fixed = tdata.Permanently_Locked
			rdata, ok := setroom[ai]
			if ok {
				a.Room = rdata.Room
				a.RealRooms = rdata.Real_Room
			}
			//fmt.Printf("§ ROOMS: %s - %+v\n", a.Room, a.RealRooms)
		} else {
			a.Day = -1
		}
		amap[ai] = a
		fmt.Printf("§ ACTIVITY: %+v\n", a)
	}
	// Gather all the data together
	return FetResult{
		Institution: v.Institution_Name,
		MainComment: v.Comments,
		Days:        daymap,
		Hours:       hourmap,
		Rooms:       roommap,
		Subjects:    subjectmap,
		Teachers:    teachermap,
		Students:    classmap,
		Activities:  amap,
	}
}

// Get the fet data in a form to pass to the printing functions.
func PrepareFetData(fetdata FetResult) []timetable.LessonData {
	// Class-group infrastructure
	divmap := map[string][][]string{}
	for c, cdata := range fetdata.Students {
		divlist := [][]string{}
		for _, div := range cdata.Category {
			divlist = append(divlist, div.Division)
		}
		divmap[c] = divlist
		//fmt.Printf(" $$$ AD %s: %+v\n", c, divlist)
	}

	lessons := []timetable.LessonData{}
	for _, a := range fetdata.Activities {
		if a.Day < 0 {
			// Unplaced activity, skip it.
			continue
		}
		// Gather the rooms: Use the RealRooms list unless it is empty,
		// in which case use the Room value.
		var rooms []string
		if len(a.RealRooms) == 0 {
			rooms = a.RealRooms
		} else if a.Room != "" {
			rooms = []string{a.Room}
		}
		// Gather the teachers.
		teachers := a.Teacher

		//TODO: Is there any way of associating teachers with particular
		// (sub)groups? Probably not (with the current data structures).

		// Gather student groups, dividing them for the class view.
		classes := map[string][]string{} // mapping: class -> list of groups
		for _, cg := range a.Students {
			cgs := strings.SplitN(cg, CLASS_GROUP_SEP, 2)
			if len(cgs) == 1 {
				classes[cg] = nil
			} else {
				c := cgs[0]
				classes[c] = append(classes[c], cgs[1])
			}
		}
		cgroups := map[string][]timetable.TTGroup{}
		for c, glist := range classes {
			var ttgroups []timetable.TTGroup
			if len(glist) == 0 {
				// whole class
				ttgroups = []timetable.TTGroup{{
					Groups: nil,
					Offset: 0,
					Size:   1,
					Total:  1,
				}}
			} else {
				n := 0
				start := 0
				gs := []string{}
				for _, div := range divmap[c] {
					for i, g := range div {
						if slices.Contains(glist, g) {
							n += 1
							if (start + len(gs)) == i {
								gs = append(gs, g)
								continue
							}
							if len(gs) > 0 {
								ttgroups = append(ttgroups,
									timetable.TTGroup{
										Groups: gs,
										Offset: start,
										Size:   len(gs),
										Total:  len(div),
									})
							}
							gs = []string{g}
							start = i
						}
					}
					if len(gs) > 0 {
						ttgroups = append(ttgroups,
							timetable.TTGroup{
								Groups: gs,
								Offset: start,
								Size:   len(gs),
								Total:  len(div),
							})
					}
					if n != 0 {
						if n != len(glist) {
							log.Fatalf("Groups in activity for class %s"+
								" not in one division: %+v\n", c, glist)
						}
						break
					}
				}
				if n == 0 {
					log.Fatalf("Invalid groups in activity for class %s: %+v\n",
						c, glist)
				}
			}
			cgroups[c] = ttgroups
		}
		lessons = append(lessons, timetable.LessonData{
			Duration:  a.Duration,
			Subject:   a.Subject,
			Teacher:   teachers,
			Students:  cgroups,
			RealRooms: rooms,
			Day:       a.Day,
			Hour:      a.Hour,
		})
	}

	//return lessons

	info := map[string]string{
		"School": fetdata.Institution,
	}
	clist := []string{}
	for _, ci := range fetdata.Students {
		clist = append(clist, ref2id[ci])
	}
	tlmap := map[int]timetable.IdName{}
	xmax := -1
	for t, tdata := range fetdata.Teachers {
		name := tdata.Long_Name
		if name == "" {
			name = tdata.Comments
		}
		x := tdata.X
		tlmap[x] = timetable.IdName{Id: t, Name: name}
		if x > xmax {
			xmax = x
		}
	}
	tlist := []timetable.IdName{}
	for x := 0; x <= xmax; x++ {
		idn, ok := tlmap[x]
		if ok {
			tlist = append(tlist, idn)
		}
	}

	return TimetableData{
		Info:        info,
		ClassList:   clist,
		TeacherList: tlist,
		RoomList:    rlist,
		Lessons:     lessons,
	}

}

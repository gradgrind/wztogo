package readfet

import (
	"encoding/xml"
	"io"
	"log"
	"os"
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
		} else {
			a.Day = -1
		}
		amap[ai] = a
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

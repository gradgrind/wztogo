// Package fet handles interaction with the fet timetabling program.
package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"strings"
)

const fet_version = "6.18.0"

// Function makeXML produces a chunk of pretty-printed XML output from
// the input data.
func makeXML(data interface{}, indent_level int) string {
	const indent = "  "
	prefix := strings.Repeat(indent, indent_level)
	xmlData, err := xml.MarshalIndent(data, prefix, indent)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	return string(xmlData)
}

type fet struct {
	Version          string `xml:"version,attr"`
	Mode             string
	Institution_Name string
	Comments         string // this is a source reference
	Days_List        fetDaysList
	Hours_List       fetHoursList
	Teachers_List    fetTeachersList
	Subjects_List    fetSubjectsList
	Rooms_List       fetRoomsList
	Students_List    fetStudentsList
	//Buildings_List
	Activity_Tags_List     fetActivityTags
	Activities_List        fetActivitiesList
	Time_Constraints_List  timeConstraints
	Space_Constraints_List spaceConstraints
}

type fetInfo struct {
	wzdb             *wzbase.WZdata
	ref2fet          map[int]string
	days             []string
	hours            []string
	fetdata          fet
	fixed_activities []bool
}

type timeConstraints struct {
	XMLName xml.Name `xml:"Time_Constraints_List"`
	//
	ConstraintBasicCompulsoryTime                basicTimeConstraint
	ConstraintStudentsSetNotAvailableTimes       []studentsNotAvailable
	ConstraintTeacherNotAvailableTimes           []teacherNotAvailable
	ConstraintActivityPreferredStartingTime      []startingTime
	ConstraintMinDaysBetweenActivities           []minDaysBetweenActivities
	ConstraintStudentsSetMaxHoursDailyInInterval []lunchBreak
	ConstraintStudentsSetMaxGapsPerWeek          []maxGapsPerWeek
	ConstraintStudentsSetMinHoursDaily           []minLessonsPerDay
}

type basicTimeConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsoryTime"`
	Weight_Percentage int
	Active            bool
}

type spaceConstraints struct {
	XMLName                          xml.Name `xml:"Space_Constraints_List"`
	ConstraintBasicCompulsorySpace   basicSpaceConstraint
	ConstraintActivityPreferredRoom  []fixedRoom
	ConstraintActivityPreferredRooms []roomChoice
}

type basicSpaceConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsorySpace"`
	Weight_Percentage int
	Active            bool
}

func make_fet_file(wzdb *wzbase.WZdata,
	activities []wzbase.Activity,
	course2activities map[int][]int,
	subject_activities []wzbase.SubjectGroupActivities,
) string {
	//TODO--
	fmt.Printf("\n????? %+v\n", wzdb.Schooldata)

	// Build ref-index -> fet-key mapping
	ref2fet := wzdb.Ref2IdMap()

	fetinfo := fetInfo{
		wzdb:    wzdb,
		ref2fet: ref2fet,
		fetdata: fet{
			Version:          fet_version,
			Mode:             "Official",
			Institution_Name: wzdb.Schooldata["SchoolName"].(string),
			Comments:         wzdb.Schooldata["SourceReference"].(string),
			Time_Constraints_List: timeConstraints{
				ConstraintBasicCompulsoryTime: basicTimeConstraint{
					Weight_Percentage: 100, Active: true},
			},
			Space_Constraints_List: spaceConstraints{
				ConstraintBasicCompulsorySpace: basicSpaceConstraint{
					Weight_Percentage: 100, Active: true},
			},
		},
	}

	getDays(&fetinfo)
	getHours(&fetinfo)
	getTeachers(&fetinfo)
	getSubjects(&fetinfo)
	getRooms(&fetinfo)
	getClasses(&fetinfo)
	getActivities(&fetinfo, activities, course2activities)
	gap_subject_activities(&fetinfo, subject_activities)

	return xml.Header + makeXML(fetinfo.fetdata, 0)
}

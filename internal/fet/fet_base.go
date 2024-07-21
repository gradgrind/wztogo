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
	Students_List    fetStudentsList
	//Buildings_List
	//TODO:
	//	Rooms_List fetRoomsList
	Activities_List        fetActivitiesList
	Time_Constraints_List  timeConstraints
	Space_Constraints_List spaceConstraints
}

type fetInfo struct {
	wzdb    *wzbase.WZdata
	ref2fet map[int]string
	days    []string
	hours   []string
	fetdata fet
}

type timeConstraints struct {
	XMLName                                 xml.Name `xml:"Time_Constraints_List"`
	ConstraintBasicCompulsoryTime           basicTimeConstraint
	ConstraintStudentsSetNotAvailableTimes  []studentsNotAvailable
	ConstraintTeacherNotAvailableTimes      []teacherNotAvailable
	ConstraintActivityPreferredStartingTime []startingTime
}

type basicTimeConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsoryTime"`
	Weight_Percentage int
	Active            bool
}

type spaceConstraints struct {
	XMLName                        xml.Name `xml:"Space_Constraints_List"`
	ConstraintBasicCompulsorySpace basicSpaceConstraint
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
	ref2fet := map[int]string{}
	for ref := range wzdb.IndexMap {
		item := wzdb.NodeList[wzdb.IndexMap[ref]]
		node := item.Node
		var v string
		switch item.Table {
		case "DAYS":
			v = node.(wzbase.Day).ID
		case "HOURS":
			v = node.(wzbase.Hour).ID
		case "SUBJECTS":
			v = node.(wzbase.Subject).ID
		case "TEACHERS":
			v = node.(wzbase.Teacher).ID
		case "ROOMS":
			v = node.(wzbase.Room).ID
		case "CLASSES":
			v = node.(wzbase.Class).ID
		case "GROUPS":
			v = node.(wzbase.Group).ID
		default:
			continue
		}
		ref2fet[ref] = v
	}
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

	//	fetdata.Time_Constraints_List.constraints = append(
	//		fetdata.Time_Constraints_List.constraints,
	//		basicTimeConstraint{Weight_Percentage: 100, Active: true},
	//	)
	getDays(&fetinfo)
	getHours(&fetinfo)
	getTeachers(&fetinfo)
	getSubjects(&fetinfo)
	getClasses(&fetinfo)
	//getCourses(&fetinfo)
	getActivities(&fetinfo, activities, course2activities, subject_activities)
	return xml.Header + makeXML(fetinfo.fetdata, 0)
}

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
	Activities_List fetActivitiesList
	//TODO ...
	/*
			<Time_Constraints_List>
		    <ConstraintBasicCompulsoryTime>
		      <Weight_Percentage>100</Weight_Percentage>
		      <Active>true</Active>
		      <Comments></Comments>
		    </ConstraintBasicCompulsoryTime>

			...

			</Time_Constraints_List>
		    <Space_Constraints_List>
		    <ConstraintBasicCompulsorySpace>
		      <Weight_Percentage>100</Weight_Percentage>
		      <Active>true</Active>
		      <Comments></Comments>
		    </ConstraintBasicCompulsorySpace>

			...

		    </Space_Constraints_List>
	*/
}

func make_fet_file(wzdb wzbase.WZdata) string {
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
	fetdata := fet{
		Version:          fet_version,
		Mode:             "Official",
		Institution_Name: wzdb.Schooldata["SchoolName"].(string),
		Comments:         wzdb.Schooldata["SourceReference"].(string),
		Days_List:        getDays(ref2fet, wzdb.TableMap["DAYS"]),
		Hours_List:       getHours(ref2fet, wzdb.TableMap["HOURS"]),
		Teachers_List:    getTeachers(&wzdb),
		Subjects_List:    getSubjects(&wzdb),
		Students_List:    getClasses(&wzdb, ref2fet),
		Activities_List:  getCourses(&wzdb, ref2fet),
	}
	return xml.Header + makeXML(fetdata, 0)
}

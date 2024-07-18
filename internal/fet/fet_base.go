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
	fmt.Printf("\n????? %+v\n", wzdb.Schooldata)
	fetdata := fet{
		Version:          fet_version,
		Mode:             "Official",
		Institution_Name: wzdb.Schooldata["SchoolName"].(string),
		Comments:         wzdb.Schooldata["SourceReference"].(string),
		Days_List:        getDays(&wzdb),
		Hours_List:       getHours(&wzdb),
		Teachers_List:    getTeachers(&wzdb),
		Subjects_List:    getSubjects(&wzdb),
		Students_List:    getClasses(&wzdb),
		Activities_List:  getCourses(&wzdb),
	}
	return xml.Header + makeXML(fetdata, 0)
}

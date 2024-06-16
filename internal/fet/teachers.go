package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
)

type fetTeacher struct {
	XMLName  xml.Name `xml:"Teacher"`
	Name     string
	Comments string
}

type fetTeachersList struct {
	XMLName xml.Name `xml:"Teachers_List"`
	Teacher []fetTeacher
}

func getTeachers(wzdb *wzbase.WZdata) string {
	trefs := wzdb.TableMap["TEACHERS"]
	items := []fetTeacher{}
	for _, ti := range trefs {
		n := wzdb.NodeList[wzdb.IndexMap[ti]].Node
		items = append(items, fetTeacher{
			Name: n.(wzbase.Teacher).ID,
			Comments: fmt.Sprintf("%s %s",
				n.(wzbase.Teacher).FIRSTNAMES,
				n.(wzbase.Teacher).LASTNAME,
			),
			//<Target_Number_of_Hours>0</Target_Number_of_Hours>
			//<Qualified_Subjects></Qualified_Subjects>
		})
	}
	data := fetTeachersList{
		Teacher: items,
	}
	return string(makeXML(data, 0))
}

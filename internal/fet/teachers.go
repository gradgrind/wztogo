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

func getTeachers(wzdb *wzbase.WZdata) fetTeachersList {
	trefs := wzdb.TableMap["TEACHERS"]
	items := []fetTeacher{}
	for _, ti := range trefs {
		n := wzdb.GetNode(ti).(wzbase.Teacher)
		items = append(items, fetTeacher{
			Name: n.ID,
			Comments: fmt.Sprintf("%s %s",
				n.FIRSTNAMES,
				n.LASTNAME,
			),
			//<Target_Number_of_Hours>0</Target_Number_of_Hours>
			//<Qualified_Subjects></Qualified_Subjects>
		})
	}
	return fetTeachersList{
		Teacher: items,
	}
}

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

type teacherNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintTeacherNotAvailableTimes"`
	Weight_Percentage             int
	Teacher                       string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}

func getTeachers(fetinfo *fetInfo) {
	trefs := fetinfo.wzdb.TableMap["TEACHERS"]
	items := []fetTeacher{}
	natimes := []teacherNotAvailable{}
	for _, ti := range trefs {
		n := fetinfo.wzdb.GetNode(ti).(wzbase.Teacher)
		items = append(items, fetTeacher{
			Name: n.ID,
			Comments: fmt.Sprintf("%s %s",
				n.FIRSTNAMES,
				n.LASTNAME,
			),
			//<Target_Number_of_Hours>0</Target_Number_of_Hours>
			//<Qualified_Subjects></Qualified_Subjects>
		})

		// "Not available" times
		nats := []notAvailableTime{}
		for d, dna := range n.NOT_AVAILABLE {
			for _, h := range dna {
				nats = append(nats,
					notAvailableTime{
						Day: fetinfo.days[d], Hour: fetinfo.hours[h]})
			}
		}
		if len(nats) > 0 {
			natimes = append(natimes,
				teacherNotAvailable{
					Weight_Percentage:             100,
					Teacher:                       n.ID,
					Number_of_Not_Available_Times: len(nats),
					Not_Available_Time:            nats,
					Active:                        true,
				})
		}

	}
	fetinfo.fetdata.Teachers_List = fetTeachersList{
		Teacher: items,
	}
	fetinfo.fetdata.Time_Constraints_List.ConstraintTeacherNotAvailableTimes = natimes
}

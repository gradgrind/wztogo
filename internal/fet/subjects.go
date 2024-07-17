package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
)

type fetSubject struct {
	XMLName  xml.Name `xml:"Subject"`
	Name     string
	Comments string
}

type fetSubjectsList struct {
	XMLName xml.Name `xml:"Subjects_List"`
	Subject []fetSubject
}

func getSubjects(wzdb *wzbase.WZdata) string {
	trefs := wzdb.TableMap["SUBJECTS"]
	items := []fetSubject{}
	for _, ti := range trefs {
		n := wzdb.GetNode(ti).(wzbase.Subject)
		items = append(items, fetSubject{
			Name:     n.ID,
			Comments: n.NAME,
		})
	}
	data := fetSubjectsList{
		Subject: items,
	}
	return string(makeXML(data, 0))
}

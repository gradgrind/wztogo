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
		n := wzdb.NodeList[wzdb.IndexMap[ti]].Node
		items = append(items, fetSubject{
			Name:     n.(wzbase.Subject).ID,
			Comments: n.(wzbase.Subject).NAME,
		})
	}
	data := fetSubjectsList{
		Subject: items,
	}
	return string(makeXML(data, 0))
}

package fet

import (
	"encoding/xml"
	"gradgrind/wztogo/internal/wzbase"
)

// TODO: At present this is used for adding activity-tags to activities
// for these subjects (see courses.go). This is just a temporary hack ...
var tagged_subjects = []string{"Eu", "Sp"}

type fetSubject struct {
	XMLName  xml.Name `xml:"Subject"`
	Name     string
	Comments string
}

type fetSubjectsList struct {
	XMLName xml.Name `xml:"Subjects_List"`
	Subject []fetSubject
}

func getSubjects(fetinfo *fetInfo) {
	trefs := fetinfo.wzdb.TableMap["SUBJECTS"]
	items := []fetSubject{}
	for _, ti := range trefs {
		n := fetinfo.wzdb.GetNode(ti).(wzbase.Subject)
		items = append(items, fetSubject{
			Name:     n.ID,
			Comments: n.NAME,
		})
	}
	fetinfo.fetdata.Subjects_List = fetSubjectsList{
		Subject: items,
	}
}

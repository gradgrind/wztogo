package fet

import (
	"cmp"
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

//func TestDays(t *testing.T) {
//	readDays()
//}

func TestFet(t *testing.T) {
	// w365file := "../_testdata/fms.w365"
	w365file := "../_testdata/test.w365"
	wzdb := w365.ReadW365(w365file)

	fmt.Println("\n *******************************************")

	alist, course2activities, subject_courses := wzbase.GetActivities(&wzdb)
	fmt.Println("\n -------------------------------------------")
	for _, a := range alist {
		fmt.Printf(" >>> %+v\n", a)
	}

	fmt.Println("\n +++++++++++++++++++++++++++++++++++++++++++")
	wzbase.SetLessons(&wzdb, "Vorlage", alist, course2activities)
	for _, a := range alist {
		fmt.Printf("+++ %+v\n", a)
	}

	fmt.Println("\n +++++ SubjectActivities +++++")
	sgalist := wzbase.SubjectActivities(&wzdb,
		subject_courses, course2activities)
	type sgapr struct {
		subject    string
		groups     string
		activities []int
	}
	sgaprl := []sgapr{}
	for _, sga := range sgalist {
		s := wzdb.GetNode(sga.Subject).(wzbase.Subject).ID
		gg := []string{}
		for _, cg := range sga.Groups {
			gg = append(gg, cg.Print(wzdb))
		}
		sgaprl = append(sgaprl,
			sgapr{s, strings.Join(gg, ","), sga.Activities})
	}
	slices.SortStableFunc(sgaprl,
		func(a, b sgapr) int {
			if n := cmp.Compare(a.groups, b.groups); n != 0 {
				return n
			}
			// If names are equal, order by age
			return cmp.Compare(a.subject, b.subject)
		})
	for _, sga := range sgaprl {
		fmt.Printf("XXX %s / %s: %+v\n", sga.groups, sga.subject, sga.activities)
	}

	// ********** Build the fet file **********
	xmlitem := make_fet_file(&wzdb, alist, course2activities, sgalist)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	fetfile0 := "../_testdata/test.fet"
	fetfile, err := filepath.Abs(fetfile0)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", fetfile0)
	}
	f, err := os.Create(fetfile)
	if err != nil {
		log.Fatalf("Couldn't open output file: %s\n", fetfile)
	}
	defer f.Close()
	_, err = f.WriteString(xmlitem)
	if err != nil {
		log.Fatalf("Couldn't write fet output to: %s\n", fetfile)
	}
	log.Printf("\nFET file written to: %s\n", fetfile)

	/*
		cg0 := wzbase.CourseGroups{}
		cg := wzbase.CourseGroups{}
		cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{308}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: -1})
		cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{308, 328}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: 1, Groups: []int{10, 13}})
		//cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 0, Groups: []int{311, 328}})
		cg = append(cg, wzbase.ClassDivGroups{Class: 2, Div: 2, Groups: []int{14, 15}})
		//cg = append(cg, wzbase.ClassDivGroups{Class: 386, Div: 1, Groups: []int{18}})
		if !cg0.AddCourseGroups(wzdb.NodeList, cg) {
			log.Fatalln("INCOMPATIBLE GROUP")
		}
		log.Printf("\n --> %+v\n", cg0)
	*/
}

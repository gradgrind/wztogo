package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"path/filepath"
	"testing"
)

//func TestDays(t *testing.T) {
//	readDays()
//}

func TestFet(t *testing.T) {
	// w365file := "../_testdata/fms.w365"
	w365file := "../_testdata/test.w365"
	abspath, err := filepath.Abs(w365file)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", abspath)
	}
	wzdb := w365.ReadW365(abspath)

	fmt.Println("\n +++++ GetActivities +++++")
	alist, course2activities, subject_courses := wzbase.GetActivities(&wzdb)
	fmt.Println("\n -------------------------------------------")
	/*
		for _, a := range alist {
			fmt.Printf(" >>> %+v\n", a)
		}
	*/

	fmt.Println("\n +++++ SetLessons +++++")
	wzbase.SetLessons(&wzdb, "Vorlage", alist, course2activities)
	/*
		for _, a := range alist {
			fmt.Printf("+++ %+v\n", a)
		}
	*/

	fmt.Println("\n +++++ SubjectActivities +++++")
	sgalist := wzbase.SubjectActivities(&wzdb,
		subject_courses, course2activities)

	/*
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
	*/

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
}

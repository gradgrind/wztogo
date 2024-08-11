package readfet

import (
	"fmt"
	"gradgrind/wztogo/internal/timetable"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToW365(t *testing.T) {
	fpath := "../_testdata/test_data_and_timetable.fet"
	//fpath := "../_testdata/v003/test_data_and_timetable.fet"
	abspath, err := filepath.Abs(fpath)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", fpath)
	}
	fetdata := ReadFetResult(abspath)

	w365 := MakeW365(fetdata)
	//	w365 := to_w365(abspath)
	fmt.Println("====================================================")
	//fmt.Println(w365)

	fout := strings.TrimSuffix(abspath, filepath.Ext(abspath)) + ".schedule"
	f, err := os.Create(fout)
	if err != nil {
		fmt.Println(err)
	}
	// close the file with defer
	defer f.Close()
	f.WriteString(w365)
	fmt.Printf("Saved to %s\n", fout)

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
	datadir, err := filepath.Abs("../data/")
	if err != nil {
		log.Fatal(err)
	}
	lessons := PrepareFetData(fetdata)
	timetable.PrintClassTimetables(lessons, "fet", datadir,
		strings.TrimSuffix(abspath, filepath.Ext(abspath))+"_Klassen.pdf")
	timetable.PrintTeacherTimetables(lessons, "fet", datadir,
		strings.TrimSuffix(abspath, filepath.Ext(abspath))+"_Lehrer.pdf")
	timetable.PrintRoomTimetables(lessons, "fet", datadir,
		strings.TrimSuffix(abspath, filepath.Ext(abspath))+"_Räume.pdf")
}

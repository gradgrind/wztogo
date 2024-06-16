package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"testing"
)

//func TestDays(t *testing.T) {
//	readDays()
//}

func TestFet(t *testing.T) {
	w365file := "../_testdata/fms.w365"
	// w365file := "../_testdata/test.w365"
	wzdb := w365.ReadW365(w365file)
	//fmt.Printf("\nINPUT: %+v\n", wzdb)
	xmlitem := getDays(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getHours(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getSubjects(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
	xmlitem = getTeachers(&wzdb)
	fmt.Printf("\n*** fet:\n%v\n", xmlitem)
}

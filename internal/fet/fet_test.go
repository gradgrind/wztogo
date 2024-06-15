package fet

import (
	"fmt"
	"gradgrind/wztogo/internal/w365"
	"testing"
)

func TestDays(t *testing.T) {
	readDays()
}

func TestFet(t *testing.T) {
	w365file := "../_testdata/fms.w365"
	// w365file := "../_testdata/test.w365"
	wzdb := w365.ReadW365(w365file)
	fmt.Printf("\nINPUT: %+v\n", wzdb)
}

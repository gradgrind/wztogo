package readfet

import (
	"log"
	"path/filepath"
	"testing"
)

func TestToW365(t *testing.T) {
	//fpath := "../_testdata/testx_data_and_timetable.fet"
	fpath := "../_testdata/test_data_and_timetable.fet"
	abspath, err := filepath.Abs(fpath)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", fpath)
	}
	to_w365(abspath)
}

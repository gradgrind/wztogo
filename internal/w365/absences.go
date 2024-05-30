package w365

import (
	"log"
	"slices"
	"strconv"
	"strings"

	"gradgrind/wztogo/internal/wzbase"
)

func (w365data *W365Data) read_absences() {
	for _, node := range w365data.yeartables[w365_Absence] {
		d, err := strconv.Atoi(node[w365_day])
		if err != nil {
			log.Fatal(err)
		}
		h, err := strconv.Atoi(node[w365_hour])
		if err != nil {
			log.Fatal(err)
		}
		w365data.absencemap[node[w365_Id]] = wzbase.Timeslot{
			Day: d, Hour: h,
		}
	}
}

func (w365data *W365Data) absences(item ItemType) []([]int) {
	ndays := len(w365data.TableMap["DAYS"])
	absence_map := make([]([]int), ndays)
	for i := range ndays {
		absence_map[i] = []int{}
	}
	a0 := item[w365_Absences]
	if a0 != "" {
		alist := []wzbase.Timeslot{}
		for _, wid := range strings.Split(a0, LIST_SEP) {
			alist = append(alist, w365data.absencemap[wid])
		}
		// Sort the timeslots
		slices.SortFunc(alist, func(a, b wzbase.Timeslot) int {
			if a.Day < b.Day {
				return -1
			}
			if a.Day == b.Day {
				if a.Hour < b.Hour {
					return -1
				}
				if a.Hour == b.Hour {
					return 0
				}
			}
			return 1
		})
		for _, ts := range alist {
			absence_map[ts.Day] = append(absence_map[ts.Day], ts.Hour)
		}
	}
	return absence_map
}

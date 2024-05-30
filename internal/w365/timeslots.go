package w365

import (
	"log"
	"slices"
	"strconv"

	"gradgrind/wztogo/internal/wzbase"
)

type xday struct {
	sortnum float64
	wid     string
	tag     string
	name    string
}

func (w365data *W365Data) read_days() {
	xnodes := []xday{}
	for _, node := range w365data.yeartables[w365_Day] {
		sf, err := strconv.ParseFloat(node[w365_ListPosition], 64)
		if err != nil {
			log.Fatal(err)
		}
		xnodes = append(xnodes, xday{
			sortnum: sf,
			wid:     node[w365_Id],
			tag:     node[w365_Shortcut],
			name:    node[w365_Name],
		})
	}
	// Sort the days
	slices.SortFunc(xnodes, func(a, b xday) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	for i, xd := range xnodes {
		var id string
		if xd.tag == "" {
			id = strconv.Itoa(i)
		} else {
			id = xd.tag
		}
		n := wzbase.Day{
			ID:   id,
			NAME: xd.name,
			X:    i,
		}
		w365data.add_node("DAYS", n, xd.wid)
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type xhour struct {
	sortnum    float64
	wid        string
	tag        string
	name       string
	start_time string
	end_time   string
	_lb        bool
	_pm        bool
}

func (w365data *W365Data) read_hours() {
	xnodes := []xhour{}
	for _, node := range w365data.yeartables[w365_Period] {
		sf, err := strconv.ParseFloat(node[w365_ListPosition], 64)
		if err != nil {
			log.Fatal(err)
		}
		xnodes = append(xnodes, xhour{
			sortnum:    sf,
			wid:        node[w365_Id],
			tag:        node[w365_Shortcut],
			name:       node[w365_Name],
			start_time: node[w365_Start],
			end_time:   node[w365_End],
			_lb:        node[w365_MiddayBreak] == "true",
			_pm:        node[w365_FirstAfternoonHour] == "true",
		})
	}
	// Sort the hours
	slices.SortFunc(xnodes, func(a, b xhour) int {
		if a.sortnum <= b.sortnum {
			return -1
		}
		return 1
	})
	lb := []int{}
	pm := 0
	for i, xh := range xnodes {
		var id string
		if xh.tag == "" {
			id = strconv.Itoa(i)
		} else {
			id = xh.tag
		}
		n := wzbase.Hour{
			ID:         id,
			NAME:       xh.name,
			X:          i,
			START_TIME: xh.start_time,
			END_TIME:   xh.end_time,
		}
		w365data.add_node("HOURS", n, xh.wid)
		if xh._lb {
			lb = append(lb, i)
		}
		if xh._pm && pm == 0 {
			pm = i
		}
	}
	w365data.Config["LUNCHBREAK"] = lb
	w365data.Config["AFTERNOON_START_PERIOD"] = pm
}

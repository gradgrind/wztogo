package w365

import (
	"log"
	"strconv"
	"strings"
)

type Category struct {
	Role           map[string]bool
	WorkloadFactor float64
	NotColliding   bool
	roleproto      int

	// Do I need this?
	shortcut string

	// The following are fields that I have added, their values being
	// based on special values in the Shortcut. I am not sure whether
	// they are really needed.

	// Courses only (course tag, alternative to Epochenschienen?):
	Block string
	// Teachers only (not quite sure what it should do):
	MaxLunchDays int
}

func (w365data *W365Data) read_categories() {
	for _, node := range w365data.yeartables[w365_Category] {
		role, err := strconv.Atoi(node["Role"])
		if err != nil {
			log.Fatal(err)
		}
		wf_ := node["WorkloadFactor"]
		wf := 1.0
		if wf_ != "1.0" {
			wf, err = strconv.ParseFloat(wf_, 64)
			if err != nil {
				log.Fatal(err)
			}
		}
		// TODO: Where multiple categories are used is it right to assume or
		// assert that only one will have a WorkloadFactor other than "1.0"?
		sc := node["Shortcut"]
		//wid := node[w365_Id]
		block := ""
		if strings.HasPrefix(sc, "_") {
			// only relevant for courses
			block = sc
		}
		maxLunchDays := -1
		_, n, ok := strings.Cut(sc, "*") // only relevant for teachers
		if ok {
			maxLunchDays, err = strconv.Atoi(n)
			if err != nil {
				log.Fatal(err)
			}
		}
		w365data.categorymap[node[w365_Id]] = Category{
			WorkloadFactor: wf,
			NotColliding:   node["Colliding"] == "false",
			MaxLunchDays:   maxLunchDays,
			Block:          block,
			roleproto:      role,
			shortcut:       sc,
		}
	}
}

func (w365data *W365Data) categories(item ItemType) Category {
	c := item[w365_Categories]
	wf := 1.0
	role := 0
	notcolliding := false
	block := ""
	maxlunchdays := -1
	if c != "" {
		ctag := false
		for _, wid := range strings.Split(c, LIST_SEP) {
			cat := w365data.categorymap[wid]
			role |= cat.roleproto
			if cat.WorkloadFactor != 1.0 {
				if wf == 1.0 {
					wf = cat.WorkloadFactor
				} else {
					log.Fatal("Category: Multiple WorkloadFactors != 1.0")
				}
			}
			if cat.NotColliding {
				notcolliding = true
			}
			if cat.Block != "" { // only relevant for courses
				if ctag {
					log.Fatal("Category: Unexpected block marker")
				}
				ctag = true
				block = cat.Block
			}
			if cat.MaxLunchDays >= 0 { // only relevant for teachers
				if ctag {
					log.Fatal("Category: Unexpected MaxLunchDays")
				}
				ctag = true
				maxlunchdays = cat.MaxLunchDays
			}
		}
	}
	return Category{
		Role: map[string]bool{
			// mark a special subject as representing "available for substitutions":
			"CanSubstitute": (role & 1) != 0,
			// course with no report:
			"NoReport": (role & 2) != 0,
			// course with no register (Klassenbuch) entry:
			"NotRegistered": (role & 4) != 0,
			// for "Epoch" definitions:
			"WholeDayBlock": (role & 8) != 0,
			// to mark a "dummy" teacher:
			"NoTeacher": (role & 16) != 0,
		},
		WorkloadFactor: wf,
		NotColliding:   notcolliding,
		Block:          block,
		MaxLunchDays:   maxlunchdays,
		roleproto:      role,
	}
}

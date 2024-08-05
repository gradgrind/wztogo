package timetable

import (
	"encoding/json"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type Tile struct {
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Duration int    `json:"duration"`
	Fraction int    `json:"fraction"`
	Offset   int    `json:"offset"`
	Total    int    `json:"total"`
	Centre   string `json:"centre"`
	TL       string `json:"tl"`
	TR       string `json:"tr"`
	BR       string `json:"br"`
	BL       string `json:"bl"`
}

type Timetable struct {
	Title  string
	School string
	Plan   string
	Pages  [][]interface{}
}

func PrintClasses(wzdb *wzbase.WZdata,
	plan_name string,
	activities []wzbase.Activity,
	datadir string,
	outpath string, // full path to output pdf
) {
	pages := [][]interface{}{}
	//

	ref2id := wzdb.Ref2IdMap()
	// Get the rooms contained in room-groups
	room_groups := map[int][]string{}
	for _, ri := range wzdb.TableMap["ROOMS"] {
		rg := wzdb.GetNode(ri).(wzbase.Room).SUBROOMS
		if len(rg) != 0 {
			rglist := []string{}
			for _, r := range rg {
				rglist = append(rglist, ref2id[r])
			}
			slices.Sort(rglist)
			room_groups[ri] = rglist
		}
	}

	// Class-group infrastructure
	divmap := map[string][][]string{}
	for c, ad := range wzdb.ActiveDivisions {
		divlist := [][]string{}
		for _, div := range ad {
			gs := []string{}
			for _, g := range div {
				gs = append(gs, ref2id[g])
			}
			divlist = append(divlist, gs)
		}
		divmap[ref2id[c]] = divlist
		fmt.Printf(" $$$ AD %s: %+v\n", ref2id[c], divlist)
	}

	type chip struct {
		class  string
		groups []string
		offset int
		parts  int
		total  int
	}

	// Generate the tiles.
	classTiles := map[string][]Tile{}
	for _, a := range activities {
		if a.Day < 0 {
			// Unplaced activity, skip it.
			continue
		}
		// Gather the rooms.
		rooms := []string{}
		if len(a.Rooms) == 0 {
			// Check whether there are compulsory rooms (possible with
			// undeclared room-group).
			for _, r := range a.RoomNeeds.Compulsory {
				rooms = append(rooms, ref2id[r])
			}
			if len(rooms) > 1 {
				slices.Sort(rooms)
			}
		} else {
			for _, r := range a.Rooms {
				rg, ok := room_groups[r]
				if ok {
					rooms = append(rooms, rg...)
				} else {
					rooms = append(rooms, ref2id[r])
				}
			}
		}
		// Limit the length of the room list.
		var room string
		if len(rooms) > 6 {
			room = strings.Join(rooms[:5], ",") + "..."
		} else {
			room = strings.Join(rooms, ",")
		}
		// Gather the teachers.
		teachers := []string{}
		for _, t := range a.Teachers {
			teachers = append(teachers, ref2id[t])
		}
		var teacher string
		if len(teachers) > 6 {
			teacher = strings.Join(teachers[:5], ",") + "..."
		} else {
			teacher = strings.Join(teachers, ",")
		}

		//TODO: Is there any way of associating teachers with particular
		// (sub)groups? Probably not.

		// Gather student groups.
		classes := map[string][]string{} // mapping: class -> list of groups
		for _, cg := range a.Groups {
			c := ref2id[cg.CIX]
			g := ref2id[cg.GIX]
			// Assume the groups are valid
			if g == "" {
				classes[c] = nil
			} else {
				classes[c] = append(classes[c], g)
			}
			/*
				gl, ok := classes[c]
				if ok {
					// Assume the groups are valid
					classes[c] = append(gl, g)
				} else if g == "" {
					classes[c] = nil
				} else {
					classes[c] = []string{g}
				}
			*/
		}

		chips := []chip{}
		for c, glist := range classes {
			if len(glist) == 0 {
				// whole class
				chips = append(chips, chip{c, nil, 0, 1, 1})
				continue
			}
			n := 0
			start := 0
			gs := []string{}
			for _, div := range divmap[c] {
				for i, g := range div {
					if slices.Contains(glist, g) {
						n += 1
						if (start + len(gs)) == i {
							gs = append(gs, g)
							continue
						}
						if len(gs) > 0 {
							chips = append(chips, chip{c, gs, start, len(gs), len(div)})
						}
						gs = []string{g}
						start = i
					}
				}
				if len(gs) > 0 {
					chips = append(chips, chip{c, gs, start, len(gs), len(div)})
				}
				if n != 0 {
					if n != len(glist) {
						log.Fatalf("Groups in activity for class %s"+
							" not in one division: %+v\n", c, glist)
					}
					break
				}
			}
			if n == 0 {
				log.Fatalf("Invalid groups in activity for class %s: %+v\n",
					c, glist)
			}
		}

		fmt.Printf("*** GROUPS (%s / %s): %+v\n", teacher, room, chips)

		for _, c := range chips {
			tile := Tile{
				Day:      a.Day,
				Hour:     a.Hour,
				Duration: a.Duration,
				Fraction: c.parts,
				Offset:   c.offset,
				Total:    c.total,
				Centre:   ref2id[a.Subject],
				TL:       teacher,
				TR:       strings.Join(c.groups, ","),
				BR:       room,
			}
			classTiles[c.class] = append(classTiles[c.class], tile)
		}

	}
	// Assume the classes table is sorted!
	for _, ci := range wzdb.TableMap["CLASSES"] {
		c := ref2id[ci]
		ctiles, ok := classTiles[c]
		if !ok {
			continue
		}
		pages = append(pages, []interface{}{
			fmt.Sprintf("Klasse %s", c),
			ctiles,
		})

	}

	tt := Timetable{
		Title:  "Stundenpläne der Klassen",
		School: wzdb.Schooldata["SchoolName"].(string),
		Plan:   plan_name,
		Pages:  pages,
	}
	b, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	//os.Stdout.Write(b)
	jsonfile := filepath.Join("out", "tmp.json")
	jsonpath := filepath.Join(datadir, jsonfile)
	err = os.WriteFile(jsonpath, b, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote json to: %s\n", jsonpath)
	cmd := exec.Command("typst", "compile",
		"--root", datadir,
		"--input", "ifile="+filepath.Join("..", jsonfile),
		filepath.Join(datadir, "resources", "print_timetable.typ"),
		outpath)
	fmt.Printf(" ::: %s\n", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

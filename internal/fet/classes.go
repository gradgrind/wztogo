package fet

import (
	"encoding/xml"
	"fmt"
	"gradgrind/wztogo/internal/wzbase"
)

type fetCategory struct {
	//XMLName             xml.Name `xml:"Category"`
	Number_of_Divisions int
	Division            []string
}

type fetSubgroup struct {
	Name string // 13.m.MaE
	//Number_of_Students int // 0
	//Comments string // ""
}

type fetGroup struct {
	Name string // 13.K
	//Number_of_Students int // 0
	//Comments string // ""
	Subgroup []fetSubgroup
}

type fetClass struct {
	//XMLName  xml.Name `xml:"Year"`
	Name     string
	Comments string
	//Number_of_Students int (=0)
	// The information regarding categories, divisions of each category,
	// and separator is only used in the dialog to divide the year
	// automatically by categories.
	Number_of_Categories int    // 0 or 1
	Separator            string // "."
	Category             []fetCategory
	Group                []fetGroup
}

/*
	<Name>13</Name>
	<Number_of_Students>0</Number_of_Students>
	<Comments>13. Klasse</Comments>
	<!-- The information regarding categories, divisions of each category, and separator is only used in the divide year automatically by categories dialog. -->
	<Number_of_Categories>1</Number_of_Categories>
	<Category>
		<Number_of_Divisions>6</Number_of_Divisions>
		<Division>k.MaE</Division>
		<Division>k.MaG</Division>
		<Division>m.MaE</Division>
		<Division>m.MaG</Division>
		<Division>s.MaE</Division>
		<Division>s.MaG</Division>
	</Category>
	<Separator>.</Separator>
	<Group>
		<Name>13.K</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.m.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.M</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.k.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.k.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.MaE</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.k.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.MaG</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.k.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.S</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.k.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.k.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.k</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.k.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.k.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.m</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.m.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.m.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
	<Group>
		<Name>13.s</Name>
		<Number_of_Students>0</Number_of_Students>
		<Comments></Comments>
		<Subgroup>
			<Name>13.s.MaE</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
		<Subgroup>
			<Name>13.s.MaG</Name>
			<Number_of_Students>0</Number_of_Students>
			<Comments></Comments>
		</Subgroup>
	</Group>
*/

type fetStudentsList struct {
	XMLName xml.Name `xml:"Students_List"`
	Year    []fetClass
}

// TODO: Handle the groups ...
// Note that there may well be "superfluous" divisions – ones with no
// actual lessons associated. These should be stripped out for fet!
// That might be an argument for not generating the atomic groups until
// it is clear in which form they are needed.
func getClasses(wzdb *wzbase.WZdata) string {
	trefs := wzdb.TableMap["CLASSES"]
	items := []fetClass{}
	for _, ti := range trefs {
		cl := wzdb.NodeList[wzdb.IndexMap[ti]].Node.(wzbase.Class)
		divs := cl.DIVISIONS
		nc := 0
		if len(divs) > 0 {
			nc = 1

		}
		items = append(items, fetClass{
			Name:                 cl.SORTING, //?
			Comments:             cl.ID,      //?
			Number_of_Categories: nc,
			Separator:            ".",
		})
		fmt.Printf("\nCLASS %s: %+v\n", cl.SORTING, cl.DIVISIONS)
	}
	data := fetStudentsList{
		Year: items,
	}
	return string(makeXML(data, 0))
}

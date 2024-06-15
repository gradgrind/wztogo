package wzbase

import (
	"github.com/RoaringBitmap/roaring"
)

type AtomicGroups struct {
	x             uint32
	Class_Groups  map[int][]ClassGroup
	Group_Atomics map[ClassGroup]*roaring.Bitmap
}

func NewAtomicGroups() AtomicGroups {
	return AtomicGroups{
		x:             0,
		Class_Groups:  map[int][]ClassGroup{},
		Group_Atomics: map[ClassGroup]*roaring.Bitmap{},
	}
}

func (ag *AtomicGroups) Add_class_groups(cix int, cdata Class) {

	cg2rbm := ag.Group_Atomics
	//?
	type grbm struct {
		g  int
		bm *roaring.Bitmap
	}
	c2groups := map[int][]grbm{}

	//fmt.Printf("\n********* %s:\n", cdata.ID)
	glist0, cg := gen_class_groups(cdata.DIVISIONS)
	g2bm := map[int]*roaring.Bitmap{}
	for _, g := range glist0 {
		rbm := roaring.New()
		g2bm[g] = rbm
		cg2rbm[ClassGroup{CIX: cix, GIX: g}] = rbm
		c2groups[cix] = append(c2groups[cix], grbm{g, rbm})
	}
	var rbm *roaring.Bitmap
	rbm0 := roaring.New()
	if len(cg) == 0 {
		ag.x++
		rbm0 = roaring.BitmapOf(ag.x)
	} else {
		for _, glist := range cg {
			ag.x++
			rbm = roaring.BitmapOf(ag.x)
			for _, g := range glist {
				g2bm[g].Or(rbm)
				rbm0.Or(rbm)
			}
		}
	}
	c2groups[cix] = append(c2groups[cix], grbm{0, rbm0})
	cg2rbm[ClassGroup{CIX: cix, GIX: 0}] = rbm0
	//for _, cgr := range c2groups[cix] {
	//	fmt.Printf("\n +++ %d: %v", cgr.g, cgr.bm)
	//}
}

func gen_class_groups(class_divisions []DivGroups) ([]int, [][]int) {
	// If there is any ordering, it must be visible in the input data, no
	// sorting is done here.
	lists := [][]int{}
	glist0 := []int{}
	if len(class_divisions) != 0 {
		for _, divgroups := range class_divisions {
			lists2 := [][]int{}
			if divgroups.Tag == "" {
				continue
			}
			for _, g := range divgroups.Groups {
				glist0 = append(glist0, g)
				if len(lists) == 0 {
					lists2 = append(lists2, []int{g})
				} else {
					for j := range len(lists) {
						lists2 = append(lists2, append(lists[j], g))
					}
				}
			}
			lists = lists2
		}
	}
	//fmt.Printf("** gen_class_groups: %+v\n  %+v\n", glist0, lists)
	return glist0, lists
}

/* Here is code to build a simple bitmap for a class. The number of groups
would have to be limited to the word size, so it would need checking and
could not be used to encompass all classes (except in a very simple
situation).
The code above using "roaring" bitmap should make coding simpler, as it
can include all classes and uses integer indexes. It may be slower, though.

	g2ag := map[int]int{}
	i := 1
	if len(lists) == 0 {
		g2ag[0] = i
	} else {
		for _, glist := range lists {
			for _, g := range glist {
				g2ag[g] |= i
			}
			i <<= 1
		}
		g2ag[0] = i - 1
	}
	for g, i := range g2ag {
		fmt.Printf("  [%d: %b]", g, i)
	}
	fmt.Println()
*/

package wzbase

import (
	"github.com/RoaringBitmap/roaring"
)

// AtomicGroups uses "roaring bitmaps" to represent the atomic groups of
// the classes and their groups. To the outside this uses simple consecutive
// integers to represent the atomic groups, but allows set operations on the
// collections of atomic groups associated with the classes and groups.
type AtomicGroups struct {
	// X counts the atomic groups. The first one has value 1. When all have
	// been read, X has the value of the last one, i.e. the total count.
	X uint32
	//TODO: Class_Groups maps a class reference to a list of its ClassGroup
	// elements.
	Class_Groups map[int][]ClassGroup
	// Group_Atomics maps a ClassGroup to its atomic groups.
	Group_Atomics map[ClassGroup]*roaring.Bitmap
}

func NewAtomicGroups() AtomicGroups {
	return AtomicGroups{
		X:             0,
		Class_Groups:  map[int][]ClassGroup{},
		Group_Atomics: map[ClassGroup]*roaring.Bitmap{},
	}
}

func (ag *AtomicGroups) Add_class_groups(cix int, cdata Class) {

	cg2rbm := ag.Group_Atomics
	c2cg := ag.Class_Groups
	/*
		type grbm struct {
			g  int
			bm *roaring.Bitmap
		}
		c2groups := map[int][]grbm{}
	*/
	//fmt.Printf("\n********* %s:\n", cdata.ID)
	glist0, cg := gen_class_groups(cdata.DIVISIONS)
	g2bm := map[int]*roaring.Bitmap{}
	for _, g := range glist0 {
		rbm := roaring.New()
		g2bm[g] = rbm
		c_g := ClassGroup{CIX: cix, GIX: g}
		cg2rbm[c_g] = rbm
		c2cg[cix] = append(c2cg[cix], c_g)
		//c2groups[cix] = append(c2groups[cix], grbm{g, rbm})
	}
	var rbm *roaring.Bitmap
	rbm0 := roaring.New()
	if len(cg) == 0 {
		ag.X++
		rbm0 = roaring.BitmapOf(ag.X)
	} else {
		for _, glist := range cg {
			ag.X++
			rbm = roaring.BitmapOf(ag.X)
			for _, g := range glist {
				g2bm[g].Or(rbm)
				rbm0.Or(rbm)
			}
		}
	}
	//c2groups[cix] = append(c2groups[cix], grbm{0, rbm0})
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

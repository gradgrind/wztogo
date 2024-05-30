package wzbase

//import "github.com/wk8/go-ordered-map/v2"
//import "github.com/elliotchance/orderedmap/v2"

type KVtuple[T any] struct {
	key string
	val T
}

// Because of the list representation, entry deletion will be expensive.
type Omap[T any] struct {
	imap  map[string]int
	vlist []KVtuple[T]
}

func get[T any](o Omap[T], x string) (T, bool) {
	i, ok := o.imap[x]
	if ok {
		return o.vlist[i].val, true
	}
	return *new(T), false
}

func set[T any](o Omap[T], x string, v T) {
	i, ok := o.imap[x]
	if ok {
		o.vlist[i].val = v
	} else {
		i = len(o.vlist)
		o.vlist = append(o.vlist, KVtuple[T]{x, v})
		o.imap[x] = i
	}
}

func items[T any](o Omap[T]) []KVtuple[T] {
	return o.vlist
}

package walker

import (
	"go/types"
	"strings"
	"sync"

	"1pkg/gopium"
)

// maven defines visiting helper
// that aggregates has and enum facilities
type maven struct {
	exposer gopium.Exposer
	idfunc  gopium.IDFunc
	store   sync.Map
	backref *struct {
		vals    map[string]int64
		signals map[string]chan struct{}
		mutex   sync.Mutex
	}
}

func newm(exposer gopium.Exposer, idfunc gopium.IDFunc, backref bool) *maven {
	m := &maven{
		exposer: exposer,
		idfunc:  idfunc,
		store:   sync.Map{},
	}
	if backref {
		m.backref = new(struct {
			vals    map[string]int64
			signals map[string]chan struct{}
			mutex   sync.Mutex
		})
	}
	return m
}

// has defines struct store id helper
// that uses gopium.IDFunc to build id
// for a structure and check that
// builded id has not been stored already
func (m *maven) has(tn *types.TypeName) (string, bool) {
	// build id for the structure
	id := m.idfunc(tn.Pos())
	// in case id of structure
	// has been already stored
	if _, ok := m.store.Load(id); ok {
		return id, true
	}
	// mark id of structure as stored
	m.store.Store(id, struct{}{})
	return id, false
}

// enum defines struct enumerating converting helper
// that goes through all structure fields
// and uses gopium.Exposer to expose gopium.Field DTO
// for each field and puts them back
// to resulted gopium.Struct object
func (m *maven) enum(name string, st *types.Struct) (r gopium.Struct) {
	// set structure name
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     m.exposer.Name(f.Type()),
			Size:     m.size(f.Type()),
			Align:    m.exposer.Align(f.Type()),
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}

func (m *maven) size(t types.Type) int64 {
	if m.backref == nil {
		// just use default exposer size
		return m.exposer.Size(t)
	}
	// for refsize only structures
	// and arrays should be calculated
	// not with default exposer size
	switch tp := t.(type) {
	case *types.Array:
		// ignore not struct arrays
		if _, ok := tp.Elem().(*types.Struct); !ok {
			break
		}
		// note: copied from `go/types/sizes.go`
		n := tp.Len()
		if n <= 0 {
			return 0
		}
		// n > 0
		a := m.exposer.Align(tp.Elem())
		z := m.size(tp.Elem())
		return gopium.Align(z, a)*(n-1) + z
	case *types.Struct:
		name := tp.String()
		// ignore structs from different pkg
		if strings.Contains(name, ".") {
			break
		}
		// ignore anonymus structs
		if strings.Contains(name, "struct") {
			break
		}
		// get size of the structure from ref
		m.backref.mutex.Lock()
		ch, ok := m.backref.signals[name]
		m.backref.mutex.Unlock()
		if !ok {
			break
		}
		<-ch
		defer m.backref.mutex.Unlock()
		m.backref.mutex.Lock()
		if val, ok := m.backref.vals[name]; ok {
			return val
		}
	}
	// just use default exposer size
	return m.exposer.Size(t)
}

func (m *maven) stref(st gopium.Struct) {
	var size int64
	for _, f := range st.Fields {
		size += f.Size
	}
	key := st.Name

	defer m.backref.mutex.Unlock()
	m.backref.mutex.Lock()
	if ch, ok := m.backref.signals[key]; ok {
		select {
		case <-ch:
		default:
			m.backref.vals[key] = size
			close(ch)
		}
	}
}

func (m *maven) link(key string) {
	defer m.backref.mutex.Unlock()
	m.backref.mutex.Lock()
	if _, ok := m.backref.signals[key]; !ok {
		m.backref.signals[key] = make(chan struct{})
	}
}

func (m *maven) prune() {
	m.store = sync.Map{}
	defer m.backref.mutex.Unlock()
	m.backref.mutex.Lock()
	m.backref.signals = nil
	for _, ch := range m.backref.signals {
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
}

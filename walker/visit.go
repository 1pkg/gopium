package walker

import (
	"context"
	"go/types"
	"regexp"
	"sync"

	"1pkg/gopium"
)

// applied encapsulates visited by strategy
// structs results: id, origin, result structs and error
type applied struct {
	ID             string
	Origin, Result gopium.Struct
	Error          error
}

// appliedCh defines abstraction that helpes
// keep applied stream results
type appliedCh chan applied

// visitFunc defines abstraction that helpes
// visit and filtered structures on the scope
type govisit func(context.Context, *types.Scope)

// Visit helps to implement Walker VisitTop and VisitDeep methods
// depends on deep flag (different tree levels)
// it creates visitFunc instance that
// goes through all struct decls inside the scope
// convert them to inner gopium format
// and applies the strategy if struct name matches regex
// then it push result of the strategy to the chan
func visit(
	regex *regexp.Regexp,
	stg gopium.Strategy,
	exposer gopium.Exposer,
	idfunc gopium.IDFunc,
	ch appliedCh,
	deep bool,
) govisit {
	// determinate which function
	// should be applied depends on
	// deep flag
	var v func(
		ctx context.Context,
		scope *types.Scope,
		regex *regexp.Regexp,
		stg gopium.Strategy,
		exposer gopium.Exposer,
		idfunc gopium.IDFunc,
		visited *sync.Map,
		ch appliedCh,
	)
	if deep {
		v = vdeep
	} else {
		v = vscope
	}
	// return govisit func applied
	// visiting implementation
	// that goes through all structures
	// with names that match regex
	// and applies strategy to them
	return func(ctx context.Context, scope *types.Scope) {
		v(
			ctx,
			scope,
			regex,
			stg,
			exposer,
			idfunc,
			&sync.Map{},
			ch,
		)
	}
}

// vdeep defines deep visiting helper
// that goes through all structures on all scopes concurently
// with names that match regex and applies strategy to them
func vdeep(
	ctx context.Context,
	scope *types.Scope,
	regex *regexp.Regexp,
	stg gopium.Strategy,
	exposer gopium.Exposer,
	idfunc gopium.IDFunc,
	visited *sync.Map,
	ch appliedCh,
) {
	// wait group visits counter
	var wg sync.WaitGroup
	// after deep visiting is done
	// wait until all visits finished
	// and then close the channel
	defer func() {
		wg.Wait()
		close(ch)
	}()
	// indeep defines recursive inner
	// visitig helper that visits
	// all scope one by one
	// and runs vscope on them
	var indeep govisit
	indeep = func(ctx context.Context, scope *types.Scope) {
		// manage context actions
		// in case of cancelation
		// break from futher traverse
		select {
		case <-ctx.Done():
			return
		default:
		}
		// increment wait group visits counter
		wg.Add(1)
		// concurently visit current scope
		go func() {
			// decrement wait group visits counter
			// after scope visiting is done
			defer wg.Done()
			// run vscope on current scope
			nch := make(appliedCh)
			go vscope(
				ctx,
				scope,
				regex,
				stg,
				exposer,
				idfunc,
				visited,
				nch,
			)
			// redirect results of vscope
			// to final applied chanel
			for applied := range nch {
				ch <- applied
			}
		}()
		// traverse through children scopes
		for i := 0; i < scope.NumChildren(); i++ {
			// visit children scopes iteratively
			// using child context and scope
			go indeep(ctx, scope.Child(i))
		}
	}
	// start indeep chain
	indeep(ctx, scope)
}

// vscope defines top visiting helper
// that goes through structures on the scope
// with names that match regex and applies strategy to them
func vscope(
	ctx context.Context,
	scope *types.Scope,
	regex *regexp.Regexp,
	stg gopium.Strategy,
	exposer gopium.Exposer,
	idfunc gopium.IDFunc,
	visited *sync.Map,
	ch appliedCh,
) {
	// wait group visits counter
	var wg sync.WaitGroup
	// after visiting is done
	// wait until all visits finished
	// and then close the channel
	defer func() {
		wg.Wait()
		close(ch)
	}()
loop:
	// go through all names inside the package scope
	for _, name := range scope.Names() {
		// check if object name doesn't matches regex
		if !regex.MatchString(name) {
			continue
		}
		// in case it does and object is
		// a type name and it's not an alias for struct
		// then apply strategy to it
		if tn, ok := scope.Lookup(name).(*types.TypeName); ok && !tn.IsAlias() {
			// if underlying type is struct
			if st, ok := tn.Type().Underlying().(*types.Struct); ok {
				// build id for structure
				id := idfunc(tn.Pos())
				// in case id of structure
				// has been already visited
				if _, ok := visited.Load(id); ok {
					continue
				}
				// mark hierarchy name of structure to visited
				visited.Store(id, struct{}{})
				// manage context actions
				// in case of cancelation break from
				// futher traverse
				select {
				case <-ctx.Done():
					break loop
				default:
				}
				// increment wait group visits counter
				wg.Add(1)
				// concurently visit the structure
				// and apply strategy to it
				go func(id, name string, st *types.Struct) {
					// decrement wait group visits counter
					defer wg.Done()
					// convert original struct
					// to inner gopium format
					o := enum(exposer, name, st)
					// apply provided strategy
					r, err := stg.Apply(ctx, o)
					// and push results to the chan
					ch <- applied{
						ID:     id,
						Origin: o,
						Result: r,
						Error:  err,
					}
				}(id, name, st)
			}
		}
	}
}

// enum defines struct enumerating visit converting helper
// that goes through all structure fields and uses gopium.Exposer
// to expose gopium.Field DTO for each field
// and puts it back to resulted gopium.Struct object
func enum(exposer gopium.Exposer, name string, st *types.Struct) (r gopium.Struct) {
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
			Type:     exposer.Name(f.Type()),
			Size:     exposer.Size(f.Type()),
			Align:    exposer.Align(f.Type()),
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return
}

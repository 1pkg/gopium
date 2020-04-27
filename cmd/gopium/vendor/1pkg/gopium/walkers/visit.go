package walkers

import (
	"context"
	"go/types"
	"regexp"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/collections"
)

// applied encapsulates visited by strategy
// structs results: id, loc, origin, result structs and error
type applied struct {
	ID     string
	Loc    string
	Origin gopium.Struct
	Result gopium.Struct
	Error  error
}

// appliedCh defines abstraction that helps
// keep applied stream results
type appliedCh chan applied

// govisit defines abstraction that helps
// visit filtered structures on the scope
type govisit func(context.Context, *types.Scope)

// prepare defines abstraction that helps
// setup visiting maven for future visit action
type prepare func() (*maven, context.CancelFunc)

// with helps to create prepare func
// with exposer, locator and backref
func with(exp gopium.Exposer, loc gopium.Locator, bref bool) prepare {
	return func() (*maven, context.CancelFunc) {
		// create visiting maven with reference
		// and return it back,
		// with ref prune cancelation func
		ref := collections.NewReference(!bref)
		return &maven{exp: exp, loc: loc, ref: ref}, ref.Prune
	}
}

// visit helps to implement Walker Visit method
// depends on deep flag (different tree levels),
// it creates govisit func instance that
// goes through all struct decls inside the scope,
// converts them to inner gopium format
// and applies the strategy if struct name matches regex,
// then it push result of the strategy to the provided chan
func (p prepare) visit(regex *regexp.Regexp, stg gopium.Strategy, ch appliedCh, deep bool) govisit {
	// return govisit func applied
	// visiting implementation
	// that goes through all structures
	// with names that match regex
	// and applies strategy to them
	return func(ctx context.Context, s *types.Scope) {
		// prepare visiting maven
		m, cancel := p()
		defer cancel()
		// determinate which function
		// should be applied for visiting
		// depends on deep flag
		if deep {
			vdeep(ctx, s, regex, stg, m, ch)
		} else {
			vscope(ctx, s, regex, stg, m, ch)
		}
	}
}

// vdeep defines deep visiting helper
// that goes through all structures on all scopes concurently,
// if their names match regex then applies strategy to them
// uses vscope helper to visit single scope
func vdeep(ctx context.Context, s *types.Scope, r *regexp.Regexp, stg gopium.Strategy, m *maven, ch appliedCh) {
	// wait until all visits finished
	// and then close the channel
	defer close(ch)
	// wait group visits counter
	var wg sync.WaitGroup
	// indeep defines recursive inner
	// visitig helper that visits
	// all scope one by one
	// and runs vscope on them
	var indeep govisit
	indeep = func(ctx context.Context, s *types.Scope) {
		// manage context actions
		// in case of cancelation
		// break from further traverse
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
			// setup chan for current scope
			nch := make(appliedCh)
			// run vscope on current scope
			go vscope(ctx, s, r, stg, m, nch)
			// redirect results of vscope
			// to final applied chanel
			for applied := range nch {
				ch <- applied
			}
		}()
		// traverse through children scopes
		nc := s.NumChildren()
		for i := 0; i < nc; i++ {
			// visit children scopes iteratively
			// using child context and scope
			go indeep(ctx, s.Child(i))
		}
	}
	// start indeep chain
	indeep(ctx, s)
	// sync all visiting to finish
	// the same time
	wg.Wait()
}

// vscope defines visiting helper
// that goes through structures on the single scope concurently,
// if their names match regex then applies strategy to them
func vscope(ctx context.Context, s *types.Scope, r *regexp.Regexp, stg gopium.Strategy, m *maven, ch appliedCh) {
	// wait until all visits finished
	// and then close the channel
	defer close(ch)
	// wait group visits counter
	var wg sync.WaitGroup
loop:
	// go through all names inside the package scope
	for _, name := range s.Names() {
		// check if object name doesn't matches regex
		if !r.MatchString(name) {
			continue
		}
		// in case it does and object is
		// a type name and it's not an alias for struct
		// then apply strategy to it
		if tn, ok := s.Lookup(name).(*types.TypeName); ok && !tn.IsAlias() {
			// if underlying type is struct
			if st, ok := tn.Type().Underlying().(*types.Struct); ok {
				// structure id
				var id, loc string
				// in case id of structure
				// has been already visited
				if id, loc, ok = m.has(tn); ok {
					continue
				}
				// manage context actions
				// in case of cancelation break from
				// further traverse
				select {
				case <-ctx.Done():
					// push error to the chan
					ch <- applied{Error: ctx.Err()}
					break loop
				default:
				}
				// increment wait group visits counter
				wg.Add(1)
				// prepare struct ref notifier
				notif := m.refst(id)
				// concurently visit the structure
				// and apply strategy to it
				go func(id, loc, name string, st *types.Struct, notif func(gopium.Struct)) {
					// decrement wait group visits counter
					defer wg.Done()
					// convert original struct
					// to inner gopium format
					o := m.enum(name, st)
					// apply provided strategy
					r, err := stg.Apply(ctx, o)
					// notify ref with result structure
					notif(r)
					// and push results to the chan
					ch <- applied{
						ID:     id,
						Loc:    loc,
						Origin: o,
						Result: r,
						Error:  err,
					}
				}(id, loc, name, st, notif)
			}
		}
	}
	// wait until all visits finished
	wg.Wait()
}

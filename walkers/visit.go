package walkers

import (
	"context"
	"go/types"
	"regexp"
	"sync"

	"1pkg/gopium"
	"1pkg/gopium/walkers/ref"
)

// applied encapsulates visited by strategy
// structs results: id, origin, result structs and error
type applied struct {
	ID, Loc        string
	Origin, Result gopium.Struct
	Error          error
}

// appliedCh defines abstraction that helps
// keep applied stream results
type appliedCh chan applied

// visitFunc defines abstraction that helps
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
	loc gopium.Locator,
	ch appliedCh,
	deep, backref bool,
) govisit {
	// return govisit func applied
	// visiting implementation
	// that goes through all structures
	// with names that match regex
	// and applies strategy to them
	return func(ctx context.Context, scope *types.Scope) {
		// setup visiting maven and reference
		m := &maven{exposer: exposer, locator: loc}
		ref := ref.NewRef(backref)
		defer ref.Prune()
		// determinate which function
		// should be applied for visiting
		// depends on deep flag
		if deep {
			vdeep(ctx, scope, regex, stg, m, ref, ch)
		} else {
			vscope(ctx, scope, regex, stg, m, ref, ch)
		}
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
	maven *maven,
	ref *ref.Ref,
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
			go vscope(
				ctx,
				scope,
				regex,
				stg,
				maven,
				ref,
				nch,
			)
			// redirect results of vscope
			// to final applied chanel
			for applied := range nch {
				ch <- applied
			}
		}()
		// for child visiting
		// create separate child context and
		// wait until all visits finished
		// and then cancel the context
		nctx, cancel := context.WithCancel(ctx)
		defer func() {
			wg.Wait()
			cancel()
		}()
		// traverse through children scopes
		for i := 0; i < scope.NumChildren(); i++ {
			// visit children scopes iteratively
			// using child context and scope
			go indeep(nctx, scope.Child(i))
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
	maven *maven,
	ref *ref.Ref,
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
				// structure id
				var id, loc string
				// in case id of structure
				// has been already visited
				if id, loc, ok = maven.has(tn); ok {
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
				notif := ref.StRef(id)
				// concurently visit the structure
				// and apply strategy to it
				go func(id, loc, name string, st *types.Struct, notif func(gopium.Struct)) {
					// decrement wait group visits counter
					defer wg.Done()
					// convert original struct
					// to inner gopium format
					o := maven.enum(name, st, ref)
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
}

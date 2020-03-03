package gopium

import (
	"context"
	"fmt"
	"go/types"
	"regexp"
	"sync"
)

// VisitedStructCh encapsulates visited by strategy
// structs results: origin, result structs and error
type VisitedStructCh chan struct {
	Origin, Result Struct
	Error          error
}

// VisitFunc defines abstraction that helpes
// visit filtered structures in the scope
type VisitFunc func(context.Context, *types.Scope)

// Visit helps to implement Walker VisitTop and VisitDeep methods
// depends on deep flag (different tree levels)
// it creates VisitFunc instance that
// goes through all struct decls inside the scope
// and applies the strategy if struct name matches regex
// then it push result of the strategy to the StructError chan
func Visit(regex *regexp.Regexp, stg Strategy, ch VisitedStructCh, deep bool) (f VisitFunc) {
	// wait group visits counter
	var wg sync.WaitGroup
	// govisit defines shallow function
	// that goes through structures on the scope
	// with names that match regex and applies strategy to them
	//nolint
	var govisit VisitFunc
	// visited holds visited structure
	// hierarchy names list
	// should be shared between govisit funcs
	visited := sync.Map{}
	govisit = func(ctx context.Context, scope *types.Scope) {
		// after visiting is done
		// wait until all visits finished
		// and then close the channel
		// still will close channel gracefully
		// even in case of context cancelation
		defer func() {
			// in case of deep visiting
			// do nothing as godeep
			// will close channel itself
			if !deep {
				wg.Wait()
				close(ch)
			}
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
					// build full hierarchy name of structure
					name = fmt.Sprintf("%s/%s", name, st)
					// in case hierarchy name of structure
					// has been already visited
					if _, ok := visited.Load(name); ok {
						continue
					}
					// mark hierarchy name of structure to visited
					visited.Store(name, struct{}{})
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
					go func(name string, st *types.Struct) {
						// decrement wait group visits counter
						defer wg.Done()
						// apply provided strategy
						o, r, err := stg.Apply(ctx, name, st)
						// and push results to the chan
						ch <- struct {
							Origin, Result Struct
							Error          error
						}{
							Origin: o,
							Result: r,
							Error:  err,
						}
					}(name, st)
				}
			}
		}
	}
	// assign result func
	f = govisit
	// in case of deep visit
	if deep {
		// deep wait group visits counter
		var dwg sync.WaitGroup
		// deep once channel close helper
		var donce sync.Once
		// godeep defines recursive function
		// that goes through all nested scopes with govisit
		var godeep VisitFunc
		godeep = func(ctx context.Context, scope *types.Scope) {
			// create child context here
			chctx, cancel := context.WithCancel(ctx)
			// after deep visiting is done
			// wait until all visits finished
			// and then close the channel
			// still will close channel gracefully
			// even in case of context cancelation
			defer func() {
				dwg.Wait()
				// should be called only once
				// from top level godeep
				donce.Do(func() {
					cancel()
					close(ch)
				})
			}()
			// manage parent context actions
			// in case of cancelation
			// break from futher traverse
			select {
			case <-ctx.Done():
				return
			default:
			}
			// increment deep wait group visits counter
			dwg.Add(1)
			// concurently visit current scope
			go func() {
				// decrement deep wait group visits counter
				defer dwg.Done()
				// run govisit on current scope
				govisit(ctx, scope)
				// wait until scope wait group is resolved
				wg.Wait()
			}()
			// traverse children scopes
			for i := 0; i < scope.NumChildren(); i++ {
				// visit children scopes iteratively
				// using child context and scope
				go godeep(chctx, scope.Child(i))
			}
		}
		// assign result func
		f = godeep
	}
	return
}

package gopium

import (
	"context"
	"go/token"
	"go/types"
	"regexp"
	"sync"
)

// applied encapsulates visited by strategy
// structs results: id, origin, result structs and error
type applied struct {
	ID             string
	Origin, Result Struct
	Error          error
}

// VisitedStructCh defines abstraction that helpes
// keep applied stream results
type VisitedStructCh chan applied

// IDFunc defines abstraction that helpes
// create unique identifier by token.Pos
type IDFunc func(token.Pos) string

// VisitFunc defines abstraction that helpes
// visit filtered structures in the scope
type VisitFunc func(context.Context, *types.Scope)

// Visit helps to implement Walker VisitTop and VisitDeep methods
// depends on deep flag (different tree levels)
// it creates VisitFunc instance that
// goes through all struct decls inside the scope
// and applies the strategy if struct name matches regex
// then it push result of the strategy to the chan
func Visit(regex *regexp.Regexp, stg Strategy, idfunc IDFunc, ch VisitedStructCh, deep bool) (f VisitFunc) {
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
						// apply provided strategy
						o, r, err := stg.Apply(ctx, name, st)
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
	// assign result func
	f = govisit
	// in case of deep visit
	if deep {
		// deep wait group visits counter
		var dwg sync.WaitGroup
		// godeep defines recursive function
		// that goes through all nested scopes with govisit
		//nolint
		var godeep VisitFunc
		godeep = func(ctx context.Context, scope *types.Scope) {
			// after deep visiting is done
			// wait until all visits finished
			// and then close the channel
			// still will close channel gracefully
			// even in case of context cancelation
			defer func() {
				// wait for deep wait group
				// and close chan
				dwg.Wait()
				close(ch)
			}()
			var ingodeep VisitFunc
			ingodeep = func(ctx context.Context, scope *types.Scope) {
				// create child context here
				nctx, cancel := context.WithCancel(ctx)
				// after deep visiting is done
				// wait until all visits finished
				// and then cancel child context
				defer func() {
					// wait for deep wait group
					// and cancel child context
					dwg.Wait()
					cancel()
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
					go ingodeep(nctx, scope.Child(i))
				}
			}
			// run ingodeep chain
			ingodeep(ctx, scope)
		}
		// assign result func
		f = godeep
	}
	return
}

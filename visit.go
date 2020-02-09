package gopium

import (
	"context"
	"go/types"
	"regexp"
	"sync"
)

// VisitFunc defines abstraction that helpes
// visit filtered structures on the scope
type VisitFunc func(context.Context, *types.Scope)

// Visit helps to implement Walker VisitTop and VisitDeep methods
// depends on deep flag (different tree levels)
// it creates VisitFunc instance that
// goes through all struct decls inside the scope
// and applies the strategy if struct name matches regex
// then it push result of the strategy to the StructError chan
func Visit(regex *regexp.Regexp, stg Strategy, ch chan<- StructError, deep bool) (f VisitFunc) {
	// deep wait group visits counter
	var wg sync.WaitGroup
	// govisit defines shallow function
	// that goes through structures on the scope
	// with names that match regex and applies strategy to them
	var govisit VisitFunc
	govisit = func(ctx context.Context, scope *types.Scope) {
		// after visiting is done
		// wait until all visits finished
		// and then close the channel
		defer func() {
			// in case of deep visiting
			// do nothing as godeep
			// will close channel itself
			if !deep {
				wg.Wait()
				close(ch)
			}
		}()
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
					// increment wait group visits counter
					wg.Add(1)
					// concurently visit the structure
					// and apply strategy to it
					go func() {
						// apply strategy
						// and push result to the chan
						ch <- stg.Apply(ctx, name, st)
						// decrement wait group visits counter
						wg.Done()
					}()

				}
			}
		}
	}
	// assign result func
	f = govisit

	if deep {
		// deep wait group visits counter
		var dwg sync.WaitGroup
		// godeep defines recursive function
		// that goes through all nested scopes with govisit
		var godeep VisitFunc
		godeep = func(ctx context.Context, scope *types.Scope) {
			// after deep visiting is done
			// wait until all visits finished
			// and then close the channel
			defer func() {
				dwg.Wait()
				close(ch)
			}()
			// increment deep wait group visits counter
			dwg.Add(1)
			// concurently visit top scope
			go func() {
				// govisit visit scope
				govisit(ctx, scope)
				// decrement deep wait group visits counter
				dwg.Done()
			}()
			// traverse children scopes
			for i := 0; i < scope.NumChildren(); i++ {
				chs := scope.Child(i)
				// visit them iteratively
				godeep(ctx, chs)
			}
		}
		// assign result func
		f = godeep
	}

	return
}

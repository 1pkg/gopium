package strategy

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"1pkg/gopium"

	"golang.org/x/sync/errgroup"
)

// list of group presets
var (
	grp = group{}
)

// group defines strategy implementation
// that uses fields tags annotation
// in order to execute different set of strategies
// on different groups and then combine it
// in single struct result
type group struct {
	b Builder
}

// container carries all relevant group data
type container struct {
	g    string
	o, r gopium.Struct
	stg  gopium.Strategy
}

// Apply group implementation
func (stg group) Apply(ctx context.Context, o gopium.Struct) (r gopium.Struct, err error) {
	// copy original structure to result
	r = o
	// parse tag annotation
	// into containers groups
	containers, err := stg.parse(r)
	// in case of any error
	// just return error back
	if err != nil {
		return o, err
	}
	// create sync error group
	// with cancelation context
	group, _ := errgroup.WithContext(ctx)
	// go through all containers and apply
	// all strategies concurently on inner structs
	for i := range containers {
		container := &containers[i]
		group.Go(func() error {
			// apply strategy on struct
			res, err := container.stg.Apply(ctx, container.o)
			// in case of any error
			// just return error back
			if err != nil {
				return err
			}
			// in case of success
			// update result on container
			container.r = res
			// if we faced default group
			// update result comment and doc
			if container.g == "" {
				r = res
			}
			return nil
		})
	}
	// wait until all strategies
	// have been applied and resolved
	if err = group.Wait(); err != nil {
		return o, nil
	}
	// sort result containers lexicographicaly
	sort.SliceStable(containers, func(i, j int) bool {
		return containers[i].g < containers[j].g
	})
	// combine all results to single result struct
	r.Fields = nil
	for _, container := range containers {
		r.Fields = append(r.Fields, container.r.Fields...)
	}
	return
}

// parse helps to parse structure fields tags
// into groups container or throw parse error
// - `gopium:"stg,stg,stg"` parsed to `default` group
// - `gopium:"group:def;stg,stg,stg"` parsed to named group
// - otherwise throw a parse error
func (stg group) parse(st gopium.Struct) ([]container, error) {
	// setup temporary groups maps
	// for fields and strategies
	gfields := make(map[string][]gopium.Field)
	gstrategies := make(map[string]string)
	// go through all struct fields
	for _, f := range st.Fields {
		// grab the field tag
		tag, ok := reflect.StructTag(f.Tag).Lookup(gopium.TAG)
		// in case tag is empty
		// or marked as skipped
		if !ok || tag == "-" {
			gfields["-"] = append(gfields["-"], f)
			continue
		}
		// otherwise parse the tag
		tokens := strings.Split(tag, ";")
		if len(tokens) == 1 {
			stgs := tokens[0]
			// check that strategies list is consistent
			if gstg, ok := gstrategies[""]; ok && gstg != stgs {
				return nil, fmt.Errorf("inconsistent strategies list %q for field %q", stgs, f.Name)
			}
			// collect strategies and fields
			gstrategies[""] = stgs
			gfields[""] = append(gfields[""], f)
		} else if len(tokens) == 2 {
			group := tokens[0]
			stgs := tokens[1]
			// check that tag contains group anchor
			if !strings.Contains(group, "group:") {
				return nil, fmt.Errorf("tag %q can't be parsed named group `group:` anchor wasn't found", tag)
			}
			// remove group anchor
			group = strings.Replace(group, "group:", "", 1)
			// check that strategies list is consistent
			if gstg, ok := gstrategies[group]; ok && gstg != stgs {
				return nil, fmt.Errorf("inconsistent strategies list %q for field %q", stgs, f.Name)
			}
			// collect strategies and fields
			gstrategies[group] = stgs
			gfields[group] = append(gfields[group], f)
		} else {
			// return parsing error msg
			return nil, fmt.Errorf("tag %q can't be parsed neither as `default` nor named group", tag)
		}
	}
	// setup result containers
	containers := make([]container, len(gfields))
	// go through all collected group fields
	for g, gfs := range gfields {
		// prepare new empty group container
		var cnt container
		// set container group
		cnt.g = g
		// set container original
		// struct and its fields
		cnt.o = st
		cnt.o.Fields = gfs
		// if group has strategy
		// then build it
		// otherwise set nil strategy
		if gstgs, ok := gstrategies[g]; ok {
			// prepare strategy pipe
			p := pipe{}
			names := strings.Split(gstgs, ",")
			// go through list of strategy name
			for _, name := range names {
				// try to build new strategy by name
				stg, err := stg.b.Build(gopium.StrategyName(name))
				// in case of any error
				// just return it back
				if err != nil {
					return nil, err
				}
				// otherwise append strategy to pipe
				p = append(p, stg)
			}
			cnt.stg = p
		} else {
			cnt.stg = nl
		}
		// append current container to result
		containers = append(containers, cnt)
	}
	// return result containers
	return containers, nil
}

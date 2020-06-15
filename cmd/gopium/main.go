package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/runners"
	"github.com/spf13/cobra"
)

// list of all cli vars
var (
	// cli command iteself
	cli *cobra.Command
	// target platform vars
	tcompiler string
	tarch     string
	tcpulines []int
	// package parser vars
	ppath   string
	pbenvs  []string
	pbflags []string
	// gopium walker vars
	wregex   string
	wdeep    bool
	wbackref bool
	// gopium printer vars
	pindent   int
	ptabwidth int
	pusespace bool
	pusegofmt bool
	// gopium global vars
	timeout int
)

// init cli command runner
// and global context
func init() {
	// set root cli command app
	cli = &cobra.Command{
		Use:     "gopium -flag_0 -flag_n walker package strategy_1 strategy_2 strategy_3 ...",
		Short:   gopium.STAMP,
		Version: gopium.VERSION,
		Example: "gopium -r ^A go_std 1pkg/gopium filter_pads memory_pack separate_padding_cpu_l1_top separate_padding_cpu_l1_bottom",
		Long: `
Gopium is the tool that was designed to automate and simplify some common performance transformations for structs, such as:
- cpu cache alignment
- memory packing
- false sharing guarding
- auto annotation
- generic fields management, etc.

In order to use gopium cli you need to provide at least package name (full package name is expected),
list of strategies which is applied one by one and single walker.
Outcome of execution is fully defined by list of strategies and walker combination.
List of strategies modifies structs inside the package, walker facilitates and insures,
that outcome is formatted and written to one of provided destinations.

Gopium provides next walkers:

 - ast_go (directly syncs result as go code to orinal file)
 - ast_go_tree (directly syncs result as go code to copy package)
 - ast_std (prints result as go code to stdout)
 - ast_gopium (directly syncs result as go code to copy gopium files)
 - json_file (prints json encoded results to single file inside package directory)
 - xml_file (prints xml encoded results to single file inside package directory)
 - csv_file (prints csv encoded results to single file inside package directory)
 - md_table_file (prints markdown table encoded results to single file inside package directory)
 - size_align_md_table_file (prints markdown encoded table of sizes and aligns difference for results to single file
	inside package directory)
 - fields_html_table_file (prints html encoded table of fields difference for results to single file
	inside package directory)

Gopium provides next strategies:

 - process_tag_group (uses gopium fields tags annotation in order to process different set of strategies
	on different groups and then combine results in single struct result)
 - memory_pack (rearranges structure fields to obtain optimal memory utilization)
 - memory_unpack (rearranges structure field list to obtain inflated memory utilization)
 - cache_rounding_cpu_l1 (fits structure into cpu cache line #1 by adding bottom partial rounding cpu cache padding)
 - cache_rounding_cpu_l2 (fits structure into cpu cache line #2 by adding bottom partial rounding cpu cache padding)
 - cache_rounding_cpu_l3 (fits structure into cpu cache line #3 by adding bottom partial rounding cpu cache padding)
 - full_cache_rounding_cpu_l1 (fits structure into full cpu cache line #1 by adding bottom rounding cpu cache padding)
 - full_cache_rounding_cpu_l2 (fits structure into full cpu cache line #2 by adding bottom rounding cpu cache padding)
 - full_cache_rounding_cpu_l3 (fits structure into full cpu cache line #3 by adding bottom rounding cpu cache padding)
 - false_sharing_cpu_l1 (guards structure from false sharing by adding extra cpu cache line #1 paddings
	for each structure field)
 - false_sharing_cpu_l2 (guards structure from false sharing by adding extra cpu cache line #1 paddings
	for each structure field)
 - false_sharing_cpu_l3 (guards structure from false sharing by adding extra cpu cache line #1 paddings
	for each structure field)
 - separate_padding_system_alignment_top (separates structure with extra system alignment padding by adding
	the padding at the top)
 - separate_padding_cpu_l1_top (separates structure with extra cpu cache line #1 padding by adding
	the padding at the top)
 - separate_padding_cpu_l2_top (separates structure with extra cpu cache line #2 padding by adding
	the padding at the top)
 - separate_padding_cpu_l3_top (separates structure with extra cpu cache line #3 padding by adding
	the padding at the top)
 - separate_padding_system_alignment_bottom (separates structure with extra system alignment padding by adding
	the padding at the bottom)
 - separate_padding_cpu_l1_bottom (separates structure with extra cpu cache line #1 padding by adding
	the padding at the bottom)
 - separate_padding_cpu_l2_bottom (separates structure with extra cpu cache line #2 padding by adding
	the padding at the bottom)
 - separate_padding_cpu_l3_bottom (separates structure with extra cpu cache line #3 padding by adding
	the padding at the bottom)
 - explicit_padings_system_alignment (explicitly aligns each structure field to system alignment padding by adding
	missing paddings for each field)
 - explicit_padings_type_natural (explicitly aligns each structure field to max type alignment padding by adding
	missing paddings for each field)
 - add_tag_group_soft (adds gopium fields tags annotation if no previous annotation found)
 - add_tag_group_force (adds gopium fields tags annotation if previous annotation found overwrites it)
 - add_tag_group_discrete (discretely adds gopium fields tags annotation if no previous annotation found)
 - add_tag_group_force_discrete (discretely adds gopium fields tags annotation if previous annotation found overwrites it)
 - remove_tag_group (removes gopium fields tags annotation)
 - doc_fields_annotate (adds align and size doc annotation for each structure field)
 - comment_fields_annotate adds align and size comment annotation for each structure field)
 - doc_struct_annotate (adds aggregated align and size doc annotation for structure)
 - comment_struct_annotate (adds aggregated align and size comment annotation for structure)
 - name_lexicographical_ascending (sorts fields accordingly to their names in ascending order)
 - name_lexicographical_descending (sorts fields accordingly to their names descending order)
 - type_lexicographical_ascending (sorts fields accordingly to their types in ascending order)
 - type_lexicographical_descending (sorts fields accordingly to their types in descending order)
 - filter_pads (filters out all structure padding fields)
 - ignore (does nothing by returning original structure)

Notes:
 - it might be useful to use filter_pads in pipes with other strategies to clean paddings first.
 - process_tag_group currently supports only next fields tags annotation formats:
  - gopium:"stg,stg,stg" processed as default group
  - gopium:"group:def;stg,stg,stg" processed as named group
 - by specifying tag_type you can automatically generate fields tags annotation suitable for process_tag_group.
 - add_tag_* strategies just add list of applied transformations to structure fields tags and NOT change results of
	other strategies, you can execute process_tag_group strategy afterwards to reuse saved strategies list.
		`,
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// create cli app instance
			cli, err := runners.NewCli(
				// target platform vars
				tcompiler,
				tarch,
				tcpulines,
				// package parser vars
				args[1], // package name
				ppath,
				pbenvs,
				pbflags,
				// gopium walker vars
				args[0], // single walker
				wregex,
				wdeep,
				wbackref,
				args[2:], // strategies slice
				// gopium printer vars
				pindent,
				ptabwidth,
				pusespace,
				pusegofmt,
				// gopium global vars
				timeout,
			)
			if err != nil {
				return err
			}
			// execute app
			return cli.Run(cmd.Context())
		},
	}
	// set target_compiler flag
	cli.Flags().StringVarP(
		&tcompiler,
		"target_compiler",
		"c",
		"gc",
		"Gopium target platform compiler, possible values are: gc or gccgo.",
	)
	// set target_architecture flag
	cli.Flags().StringVarP(
		&tarch,
		"target_architecture",
		"a",
		"amd64",
		"Gopium target platform architecture, possible values are: 386, arm, arm64, amd64, mips, etc.",
	)
	// set target_cpu_cache_lines_sizes flag
	cli.Flags().IntSliceVarP(
		&tcpulines,
		"target_cpu_cache_lines_sizes",
		"l",
		[]int{64, 64, 64},
		`
Gopium target platform CPU cache line sizes in bytes, cache line size is set one by one l1,l2,l3,...
For now only 3 lines of cache are supported by strategies.
		`,
	)
	// set package_path flag
	cli.Flags().StringVarP(
		&ppath,
		"package_path",
		"p",
		filepath.Join("src", "{{package}}"),
		`
Gopium go package path, either relative or absolute path to root of the package is expected.
To obtain full path from relative, package path is concatenated with current GOPATH env var.
Template {{package}} part is replaced with package name.
		`,
	)
	// set package_build_envs flag
	cli.Flags().StringSliceVarP(
		&pbenvs,
		"package_build_envs",
		"e",
		[]string{},
		"Gopium go package build envs, additional list of building envs is expected.",
	)
	// set package_build_flags flag
	cli.Flags().StringSliceVarP(
		&pbflags,
		"package_build_flags",
		"f",
		[]string{},
		"Gopium go package build flags, additional list of building flags is expected.",
	)
	// set walker_regexp flag
	cli.Flags().StringVarP(
		&wregex,
		"walker_regexp",
		"r",
		".*",
		`
Gopium walker regexp, regexp that defines which structures are subjects for visiting.
Visiting is done only if structure name matches the regexp.
		`,
	)
	// set walker_deep flag
	cli.Flags().BoolVarP(
		&wdeep,
		"walker_deep",
		"d",
		true,
		`
Gopium walker deep flag, flag that defines type of nested scopes visiting.
By default it visits all nested scopes.
		`,
	)
	// set walker_backref flag
	cli.Flags().BoolVarP(
		&wbackref,
		"walker_backref",
		"b",
		true,
		`
Gopium walker backref flag, flag that defines type of names referencing.
By default any previous visited types have affect on future relevant visits.
		`,
	)
	// set printer_indent flag
	cli.Flags().IntVarP(
		&pindent,
		"printer_indent",
		"i",
		0,
		"Gopium printer width of tab, defines the least code indent.",
	)
	// set printer_tab_width flag
	cli.Flags().IntVarP(
		&ptabwidth,
		"printer_tab_width",
		"w",
		8,
		"Gopium printer width of tab, defines width of tab in spaces for printer.",
	)
	// set printer_use_space flag
	cli.Flags().BoolVarP(
		&pusespace,
		"printer_use_space",
		"s",
		false,
		"Gopium printer use space flag, flag that defines if all formatting should be done by spaces.",
	)
	// set printer_use_gofmt flag
	cli.Flags().BoolVarP(
		&pusegofmt,
		"printer_use_gofmt",
		"g",
		true,
		`
Gopium printer use gofmt flag, flag that defines if canonical gofmt tool should be used for formatting.
By default it is used and overrides other printer formatting parameters.
`,
	)
	// set timeout flag
	cli.Flags().IntVarP(
		&timeout,
		"timeout",
		"t",
		0,
		"Gopium global timeout of cli command in seconds, considered only if value greater than 0.",
	)
}

// signals creates context with cancelation
// which listens to provided list of signals
func signals(ctx context.Context, sigs ...os.Signal) (context.Context, context.CancelFunc) {
	// prepare global context
	// with cancelation
	// on system signals
	ctx, cancel := context.WithCancel(ctx)
	// run separate listener goroutine
	go func() {
		defer cancel()
		// prepare signal chan for
		// global context cancelation
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, sigs...)
		// on signal or cancelation
		// stop the goroutine
		select {
		case <-ctx.Done():
		case <-sig:
		}
	}()
	return ctx, cancel
}

// main gopium cli entry point
func main() {
	// explicitly set number of threads
	// to number of logical cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	// prepare context with signals cancelation
	ctx, cancel := signals(
		context.Background(),
		os.Interrupt,
		os.Kill,
	)
	defer cancel()
	// execute cobra cli command
	// with ctx and log error if any
	if err := cli.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"

	"1pkg/gopium"
	"1pkg/gopium/runners"

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
	// walker vars
	wregex   string
	wdeep    bool
	wbackref bool
	// printer vars
	pindent   int
	ptabwidth int
	pusespace bool
	// global vars
	timeout int
	// global context
	gctx    context.Context
	gcancel func()
)

// init cli command runner
// and global context
func init() {
	// set root cli command app
	cli = &cobra.Command{
		Use:     "gopium walker package strategy#1 strategy#2 strategy#3 ...",
		Short:   gopium.STAMP,
		Version: gopium.VERSION,
		Example: "gopium -E -r ^A json_std 1pkg/gopium filter_pads memory_pack separate_padding_cpu_l1_top separate_padding_cpu_l1_bottom",
		Long: `
Gopium is the tool which was designed to automate and simplify non trivial actions on structs, like:
 - cpu cache alignment
 - memory packing
 - false sharing guarding
 - auto annotation
 - generic fields management
 - other relevant activities

In order to use gopium cli you need to provide at least package name (full package name is expected),
list of strategies which is applied one by one and single walker.
Outcome of execution is fully defined by list of strategies and walker combination.
List of strategies modifies structs inside the package, walker facilitates and insures,
that outcome is written to one of supported destinations.

Gopium supports next walkers: 
 - json_std (prints json encoded result to stdout)
 - xml_std (prints xml encoded result to stdout)
 - csv_std (prints csv encoded result to stdout)
 - json_files (prints json encoded result to files inside package directory)
 - xml_files (prints xml encoded result to files inside package directory)
 - csv_files (prints csv encoded result to files inside package directory)
 - ast_std (prints result as go code to stdout)
 - ast_go (directly syncs result as go code in orinal file)
 - ast_gopium (directly syncs result as go code in copy package)

Gopium supports next strategies: 
 - process_tag_group (uses gopium fields tags annotation in order to process different set of strategies
	on different groups and then combine results in single struct result)

 - memory_pack (rearranges structure fields to obtain optimal memory utilization)
 - memory_unpack (rearranges structure field list to obtain inflated memory utilization)
	
 - cache_rounding_cpu_l1 (fits structure into cpu cache line #1 by adding bottom rounding cpu cache padding)
 - cache_rounding_cpu_l2 (fits structure into cpu cache line #2 by adding bottom rounding cpu cache padding)
 - cache_rounding_cpu_l3 (fits structure into cpu cache line #3 by adding bottom rounding cpu cache padding)

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
 - doc_struct_annotate (adds aggregated align and size doc annotation for whole structure)
 - comment_struct_annotate (adds aggregated align and size comment annotation for whole structure)
 - doc_struct_stamp (adds doc gopium stamp to structure)
 - comment_struct_stamp (adds comment gopium stamp to structure)

 - name_lexicographical_ascending (sorts fields accordingly to their names in ascending order)
 - name_lexicographical_descending (sorts fields accordingly to their names descending order)
 - name_length_ascending (sorts fields accordingly to their names length in ascending order)
 - name_length_descending (sorts fields accordingly to their names length in descending order)
 - type_lexicographical_ascending (sorts fields accordingly to their types in ascending order)
 - type_lexicographical_descending (sorts fields accordingly to their types in descending order)
 - type_length_ascending (sorts fields accordingly to their types length in ascending order)
 - type_length_descending (sorts fields accordingly to their types length in descending order)

 - embedded_ascending (sorts fields accordingly to their embedded flag in ascending order)
 - embedded_descending (sorts fields accordingly to their embedded flag in descending order)
 - exported_ascending (sorts fields accordingly to their exported flag in ascending order)
 - exported_descending (sorts fields accordingly to their exported flag in descending order)

 - filter_pads (filters out all structure padding fields)
 - filter_embedded (filters out all structure embedded fields)
 - filter_not_embedded (filters out all structure not embedded fields)
 - filter_exported (filters out all structure exported fields)
 - filter_not_exported (filters out all structure not exported fields)

 - nope (does nothing by returning original structure)
 - void (does nothing by returning void struct)

Notes:
 - it might be useful to use filter_pads in pipes with other strategies to clean paddings first
 - process_tag_group currently supports only next fields tags annotation formats:
  - gopium:"stg,stg,stg" processed as default group
  - gopium:"group:def;stg,stg,stg" processed as named group
- by specifying tag_type you can automatically generate fields tags annotation suitable for process_tag_group
		`,
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// create app instance
			app, err := runners.NewCliApp(
				// target platform vars
				tcompiler,
				tarch,
				tcpulines,
				// package parser vars
				args[1], // package name
				ppath,
				pbenvs,
				pbflags,
				// walker vars
				args[0], // walker
				wregex,
				wdeep,
				wbackref,
				args[2:], // strategies slice
				// printer vars
				pindent,
				ptabwidth,
				pusespace,
				// global vars
				timeout,
			)
			if err != nil {
				return err
			}
			// execute app
			return app.Run(cmd.Context())
		},
	}
	// set target_compiler flag
	cli.Flags().StringVarP(
		&tcompiler,
		"target_compiler",
		"c",
		"gc",
		"Target platform compiler, possible values are: gc or gccgo.",
	)
	// set target_architecture flag
	cli.Flags().StringVarP(
		&tarch,
		"target_architecture",
		"a",
		"amd64",
		"Target platform architecture, possible values are: 386, arm, arm64, amd64, mips, etc.",
	)
	// set target_cpu_cache_line_sizes flag
	cli.Flags().IntSliceVarP(
		&tcpulines,
		"target_cpu_cache_line_sizes",
		"l",
		[]int{64, 64, 64},
		`
Target platform CPU cache line sizes in bytes, cache line size is set one by one l1,l2,l3,...
For now only 3 lines of cache are supported by strategies.
		`,
	)
	// set package_path flag
	cli.Flags().StringVarP(
		&ppath,
		"package_path",
		"p",
		"src/{{package}}",
		`
Go package path, relative path to root of the package is expected.
To obtain fill path, package path is concatenated with current GOPATH env var.
Template {{package}} part is replaced with package name.
		`,
	)
	// set package_build_envs flag
	cli.Flags().StringSliceVarP(
		&pbenvs,
		"package_build_envs",
		"e",
		[]string{},
		"Go package build envs, additional list of building envs is expected.",
	)
	// set package_build_flags flag
	cli.Flags().StringSliceVarP(
		&pbflags,
		"package_build_flags",
		"f",
		[]string{},
		"Go package build flags, additional list of building flags is expected.",
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
	// set timeout flag
	cli.Flags().IntVarP(
		&timeout,
		"timeout",
		"t",
		0,
		"Global timeout of cli command in seconds, considered only if value > 0.",
	)
	// prepare global context
	// with cancelation
	// on system signals
	gctx, gcancel = context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		select {
		case <-gctx.Done():
		case <-sig:
			gcancel()
		}
	}()
}

// main gopium cli entry point
func main() {
	// explicitly set number of threads
	// to number of logical cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	// execute cobra cli command
	// on global running context
	// and log error if any
	defer gcancel()
	if err := cli.ExecuteContext(gctx); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

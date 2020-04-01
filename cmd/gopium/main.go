package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"1pkg/gopium/runners"

	"github.com/spf13/cobra"
)

// list of all cli vars
var (
	// cli command iteself
	cli *cobra.Command
	// target platform vars
	tcompiler, tarch string
	tcpulines        []int
	// package parser vars
	pname, ppath    string
	pbenvs, pbflags []string
	// walker strategies vars
	wname, wregex   string
	wdeep, wbackref bool
	snames          []string
	tagtype         string
	// global vars
	timeout int
)

// init cli command runner
func init() {
	// set root cli command app
	cli = &cobra.Command{
		Use:   "gopium",
		Short: "",
		Long:  ``,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := runners.NewCliApp(
				tcompiler,
				tarch,
				tcpulines,
				pname,
				ppath,
				pbenvs,
				pbflags,
				wname,
				wregex,
				wdeep,
				wbackref,
				snames,
				tagtype,
			)
			ctx := cmd.Context()
			if timeout > 0 {
				nctx, cancel := context.WithTimeout(
					ctx,
					time.Duration(timeout)*time.Second,
				)
				defer cancel()
				ctx = nctx
			}
			return app.Run(ctx)
		},
	}
	// set target_compiler flag
	cli.Flags().StringVarP(
		&tcompiler,
		"target_compiler",
		"c",
		"gc",
		`
		Target compiler name (default gc), possible values are:
			- gc
			- gccgo
		for more info, please check go official list of supported compilers.
		`,
	)
	// set target_architecture flag
	cli.Flags().StringVarP(
		&tarch,
		"target_architecture",
		"a",
		"amd64",
		`
		Target architecture name (default amd64), possible values are: 
			- 386
			- arm
			- arm64
			- amd64
			- mips
			- etc.
		for more info, please check go official list of supported architectures.
		`,
	)
	// set target_cpu_cache_line_sizes flag
	cli.Flags().IntSliceVarP(
		&tcpulines,
		"target_cpu_cache_line_sizes",
		"l",
		[]int{64, 64, 64},
		`
		Target CPU cache line sizes (default [64,64,64]),
		cache size is set one by one l1,l2,l3,...
		Now maximum 3 lines of cache is supported by all strategies,
		by default typical cache size 64 is used.
		`,
	)
	// set required package_name flag
	cli.Flags().StringVarP(
		&pname,
		"package_name",
		"n",
		"",
		`
		Package name (required),
		only full package names are accepted.
		`,
	)
	cli.MarkFlagRequired("package_name")
	// set package_path flag
	cli.Flags().StringVarP(
		&ppath,
		"package_path",
		"p",
		"",
		`
		Package path (default ""),
		relative path to package root.
		`,
	)
	// set package_build_envs flag
	cli.Flags().StringSliceVarP(
		&pbenvs,
		"package_build_envs",
		"v",
		[]string{},
		`
		Package build envs (default []),
		list of building envs for types parser.
		`,
	)
	// set package_build_flags flag
	cli.Flags().StringSliceVarP(
		&pbflags,
		"package_build_flags",
		"g",
		[]string{},
		`
		Package build flags (default []),
		list of building flags for types parser.
		`,
	)
	// set required walker_name flag
	cli.Flags().StringVarP(
		&wname,
		"walker_name",
		"w",
		"",
		`
		Walker name (required), possible values are: 
			- json_std (print json encoded result to stdout)
			- xml_std (print xml encoded result to stdout)
			- csv_std (print csv encoded result to stdout)
			- json_files (print json encoded result to set of files in pkg dirs)
			- xml_files (print xml encoded result to set of files in pkg dirs)
			- csv_files (print csv encoded result to set of files in pkg dirs)
			- sync_ast (directly sync print result as go code to orinal file)
		`,
	)
	cli.MarkFlagRequired("walker_name")
	// set walker_regex flag
	cli.Flags().StringVarP(
		&wregex,
		"walker_regex",
		"r",
		".*",
		`
		Walker regex (default ".*" visit all),
		regex that filters struct names for visiting. 
		`,
	)
	// set walker_deep flag
	cli.Flags().BoolVarP(
		&wdeep,
		"walker_deep",
		"d",
		true,
		`
		Walker deep flag (default true),
		type of scopes visiting. 
		`,
	)
	// set walker_backref flag
	cli.Flags().BoolVarP(
		&wbackref,
		"walker_backref",
		"b",
		true,
		`
		Walker backref flag (default true),
		type of names referencing. 
		`,
	)
	// set required strategies_names flag
	cli.Flags().StringSliceVarP(
		&snames,
		"strategies_names",
		"s",
		[]string{},
		`
		Strategies names list (required), possible values are: 
			- nil
			- comment_fields_annotate
			- comment_struct_stamp
			- group_tag
			- lexicographical_ascending
			- lexicographical_descending
			- length_ascending
			- length_descending
			- memory_pack
			- memory_unpack
			- explicit_padings_system_alignment
			- explicit_padings_type_natural
			- false_sharing_cpu_l1
			- false_sharing_cpu_l2
			- false_sharing_cpu_l2
			- cache_rounding_cpu_l1
			- cache_rounding_cpu_l2
			- cache_rounding_cpu_l3
			- separate_padding_system_alignment
			- separate_padding_cpu_l1
			- separate_padding_cpu_l2
			- separate_padding_cpu_l3
		`,
	)
	cli.MarkFlagRequired("strategies_names")
	// set tag_type flag
	cli.Flags().StringVarP(
		&tagtype,
		"tag_type",
		"e",
		"",
		`
		Tag type, possible values are: 
			- none
			- soft
			- force
		`,
	)
	// set timeout flag
	cli.Flags().IntVarP(
		&timeout,
		"timeout",
		"t",
		0,
		`
		Gopium global cli timeout (default no timeout),
		timeout specified in seconds. 
		`,
	)
}

// main gopium cobra cli entry point
func main() {
	// prepare context with cancelation
	// on system signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		select {
		case <-ctx.Done():
		case <-sig:
			cancel()
		}
	}()
	// execute cobra cli command
	// and log error if any
	if err := cli.ExecuteContext(ctx); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"github.com/drognisep/wms/data"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

func main() {
	var (
		minorFlag       bool
		largeOnlyFlag   bool
		debugFlag       bool
		listFlag        bool
		resultLimitFlag int
	)
	if len(os.Args) <= 1 {
		printSpace()
	}

	flags := flag.NewFlagSet("flags", flag.ExitOnError)
	flags.IntVarP(&resultLimitFlag, "limit", "n", 0, "Specifies a result limit.")
	flags.BoolVarP(&listFlag, "list", "l", false, "Lists files and sizes without doing outlier analysis.")
	flags.BoolVarP(&minorFlag, "minor", "m", false, "Specifies that even minor outliers should be returned.")
	flags.BoolVarP(&largeOnlyFlag, "major", "M", false, "Specifies that only major outliers should be returned. Good for a first pass.")
	flags.BoolVar(&debugFlag, "debug", false, "Turns on debug logging.")
	flags.Usage = func() {
		fmt.Printf(`wms (where my space?) will walk your filesystem from its root to determine where a majority of space is used.
This can be useful to determine if some data may be removed.

Usage: wms
       wms FLAGS [DIR]

The first form will gather space metrics for the entire (current) drive and return.
The second form will gather space metrics for the specified directory, or the current working directory if none is specified.

%s
`, flags.FlagUsages())
	}
	if err := flags.Parse(os.Args[1:]); err != nil {
		return
	}

	if flags.NArg() > 1 {
		exit(1, "Unexpected number of arguments: %d", flags.NArg())
	}
	var target string
	if flags.NArg() == 0 {
		target = "."
	} else {
		target = flags.Arg(0)
	}
	fmt.Println("Walking directory...")
	dir, err := walkDir(debugFlag, target)
	errExit(err, "Failed to walk directory '%s': %v", flags.Arg(0), err)
	fmt.Println("Done")

	if listFlag {
		var ranks ranking
		ranks = ranks.AddFiles(dir.Files...)
		ranks = ranks.AddDirectories(dir.Directories...)
		if resultLimitFlag > 0 {
			ranks[:resultLimitFlag].PrintList()
		} else {
			ranks.PrintList()
		}
		success("Done listing files")
	}

	var (
		ranks ranking
		files []*data.File
		dirs  []*data.Directory
	)
	if largeOnlyFlag {
		files, dirs = dir.LargeOutliers(debugFlag)
	} else {
		files, dirs = dir.Outliers(debugFlag)
	}
	ranks = ranks.AddFiles(files...)
	ranks = ranks.AddDirectories(dirs...)
	if len(ranks) == 0 {
		success("No outliers found")
	}
	if resultLimitFlag > 0 {
		ranks[:resultLimitFlag].PrintList()
	} else {
		ranks.PrintList()
	}
}

func printSpace() {
	total, free, err := data.GetSpaceFree()
	errExit(err, "Failed to get disk space attributes: %v", err)
	success("Total: %s, Free: %s", total, free)
}

func errExit(err error, msg string, args ...any) {
	if err != nil {
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Printf(msg, args...)
		os.Exit(1)
	}
}

func exit(code int, msg string, args ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf(msg, args...)
	os.Exit(code)
}

func success(msg string, args ...any) {
	printFlush(msg, args...)
	os.Exit(0)
}

func printFlush(msg string, args ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf(msg, args...)
	_ = os.Stdout.Sync()
}

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"runtime/pprof"
	"strings"

	"bitbucket.org/pancakeio/delta/delta"
)

func main() {
	open := flag.Bool("open", false, "open the file in the gui")
	html := flag.Bool("html", false, "print out html")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	pathFrom := flag.Arg(0)
	pathTo := flag.Arg(1)
	pathBase := os.Getenv("BASE")

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *open {
		openDiff(pathBase, pathFrom, pathTo)
	} else {
		printDiff(pathFrom, pathTo, *html)
	}
}

// openDiffs diffs the given files and writes the result to a tempfile,
// then opens it in the gui.
func openDiff(pathBase, pathFrom, pathTo string) {
	d, err := diff(pathFrom, pathTo)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	f, err := ioutil.TempFile("", "delta-diff")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	io.WriteString(f, d.HTML())

	dir, _ := os.Getwd()
	u, _ := url.Parse("delta://openset")
	v := url.Values{}

	tmpFrom := strings.HasPrefix(pathFrom, os.TempDir())
	tmpTo := strings.HasPrefix(pathTo, os.TempDir())
	if tmpFrom && !tmpTo {
		pathFrom = pathTo
	} else if !tmpFrom && tmpTo {
		pathTo = pathFrom
	}

	v.Add("wd", dir)
	v.Add("base", pathBase)
	v.Add("left", pathFrom)
	v.Add("right", pathTo)
	v.Add("diff", f.Name())
	u.RawQuery = v.Encode()
	exec.Command("open", u.String()).Run()
}

// diff reads in files in pathFrom and pathTo, and returns a diff
func diff(pathFrom, pathTo string) (*delta.DiffSolution, error) {
	from, err := ioutil.ReadFile(pathFrom)
	if err != nil {
		return nil, fmt.Errorf("error reading %q: %v", pathFrom, err)
	}
	to, err := ioutil.ReadFile(pathTo)
	if err != nil {
		return nil, fmt.Errorf("error reading %q: %v", pathTo, err)
	}
	return delta.DiffHistogram(string(from), string(to)), nil
	// return delta.Diff(string(from), string(to)), nil
}

func printDiff(pathFrom, pathTo string, html bool) {
	d, err := diff(pathFrom, pathTo)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	if html {
		fmt.Println(d.HTML())
	} else {
		for _, l := range d.Raw() {
			if l[2] == "=" && l[0] == l[1] {
				// fmt.Printf("%d %s = %s \n", i, l[2], l[0])
				fmt.Printf(" %s \n", l[0])
				continue
			}
			if l[0] != "" {
				// fmt.Printf("\x1b[31m%d %s < %s\x1b[0m\n", i, l[2], l[0])
				fmt.Printf("\x1b[31m-%s\x1b[0m\n", l[0])
			}
			if l[1] != "" {
				// fmt.Printf("\x1b[32m%d %s > %s\x1b[0m\n", i, l[2], l[1])
				fmt.Printf("\x1b[32m+%s\x1b[0m\n", l[1])
			}
		}
	}
}
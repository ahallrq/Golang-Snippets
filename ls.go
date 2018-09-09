package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

func main() {
	progArgs := os.Args
	fmt.Println(progArgs)

	dirArg := progArgs[1]
	fileList, err := ioutil.ReadDir(dirArg)
	if (err != nil) {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, f := range fileList {
		//println(f.Name(), FSZtoString(f.Size()), FModetoString(f.Mode()))

		fmt.Printf("%s%s %d %d %5s %s\n",
			FModetoString(f.Mode()),
			f.Mode().Perm().String()[1:],
			f.Sys().(*syscall.Stat_t).Uid, // Expect this and the following line
			f.Sys().(*syscall.Stat_t).Gid, // to break on anything other than Linux
			FSZtoString(f.Size()),
			f.Name(),
		)
	}
}

func FSZtoString(size int64) string {
	switch {
	case size < 1024: return strconv.FormatInt(size, 10)
	case size < 1048576: return strconv.FormatInt(size/1024, 10)+"K"
	case size < 1073741824: return strconv.FormatInt(size/1048576, 10)+"M"
	case size < 1099511627776: return strconv.FormatInt(size/1073741824, 10)+"G"
	default: return strconv.FormatInt(size/1125899906842624, 10)+"T"
	}
}

func FModetoString(mode os.FileMode) string {
	switch {
	case mode.IsRegular(): return "-"
	case mode.IsDir(): return "d"
	default: return "?"
	// TODO: add lines for links and stuff
	}
}
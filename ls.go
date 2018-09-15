package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "<dir>")
	fmt.Printf("%s", os.ModeCharDevice)
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Fprintln(os.Stderr, "This program is only supported on Linux.")
		os.Exit(1)
	}

	progArgs := os.Args
	if len(progArgs) < 2 {
		usage()
		os.Exit(1)
	}

	dirArg := progArgs[1]
	fileList, err := ioutil.ReadDir(dirArg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var dirListStr = make([][]string, len(fileList))
	var colUserLen, colGroupLen, colFsizeLen = 0, 0, 0

	for i, f := range fileList {
		//println(f.Name(), FSZtoString(f.Size()), FModetoString(f.Mode()))

		// The following two functions only work on Linux
		pwdUser := strconv.FormatUint(uint64(f.Sys().(*syscall.Stat_t).Uid), 10)
		pwdGroup := strconv.FormatUint(uint64(f.Sys().(*syscall.Stat_t).Gid), 10)
		pwdULookup, err := user.LookupId(pwdUser)
		if err == nil { pwdUser = pwdULookup.Username }
		pwdGLookup, err := user.LookupGroupId(pwdGroup)
		if err == nil { pwdGroup = pwdGLookup.Name }

		fileSize := FSZtoString(f.Size())

		dirListStr[i] = []string{
			FModetoString(f),
			f.Mode().Perm().String()[1:],
			pwdUser, pwdGroup,
			fileSize,
			f.Name(),
		}

		if len(pwdUser) > colUserLen { colUserLen = len(pwdUser) }
		if len(pwdGroup) > colGroupLen { colGroupLen = len(pwdGroup) }
		if len(fileSize) > colFsizeLen { colFsizeLen = len(fileSize) }
	}

	for _, s := range dirListStr {
		fmt.Printf("%s%s %s %s %s %s\n",
			s[0], s[1],
			Rightpad(s[2], colUserLen), Rightpad(s[3], colGroupLen),
			Leftpad(s[4], colFsizeLen), s[5],
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

func FModetoString(f os.FileInfo) string {
	mode := f.Mode()

	switch {
	case mode.IsRegular(): return "-"
	case mode.IsDir(): return "d"
	case mode & os.ModeSymlink != 0: return "l"
	case mode & os.ModeDevice != 0 && mode & os.ModeCharDevice == 0: return "b"
	case mode & os.ModeCharDevice != 0: return "c"
	default: return "?"
	// TODO: add lines for links and stuff
	}
}
package main

import (
    "fmt"
    "flag"
    "os"
)

// represents a file to be listed
type ListFile struct {
    ModifiedTime    time.Time
    IsLink          bool
    IsDir           bool
    LinksTo         string
    Size            int64
    Name            string
    Children        []*ListFile
}

func main() {

    var showHelp bool
    var path string
    var recursive bool
    var outFmt string

    flag.BoolVar(&showHelp, "help", false, "show this message and exit")
    flag.StringVar(&path, "path", "", "path to folder")
    flag.BoolVar(&recursive, "recursive", false, "list files recursively")
    flag.StringVar(&outFmt, "output", "text", "output format - json|yaml|text")
    flag.Parse()

    if showHelp {
        flag.Usage()
        os.Exit(0)
    }
    fmt.Println("we made it")
}

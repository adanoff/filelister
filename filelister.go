package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
    "log"
    "time"
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

// converts an os.FileInfo into a ListFile
func toListFile(fi os.FileInfo, parentPath string, recursive bool, logger *log.Logger) (lf *ListFile) {

    mtime := fi.ModTime()
    fMode := fi.Mode()
    isLink := fMode & os.ModeSymlink == os.ModeSymlink
    fName := fi.Name()
    fPath, _ := filepath.Abs(filepath.Join(parentPath, fName))
    isDir := fi.IsDir()
    fSize := fi.Size()

    linkPath := ""

    if isLink {
        var err error
        linkPath, err = os.Readlink(fPath)
        if err != nil {
            logger.Fatal(err)
        }
        linkPath, _ = filepath.Abs(linkPath)
    }

    children := []*ListFile{}

    // add the children
    if recursive && isDir {
        listing, err := ioutil.ReadDir(fPath)
        if err != nil {
            logger.Fatal(err)
        }

        for _, l := range(listing) {
            children = append(children, toListFile(l, fPath, recursive, logger))
        }
    }

    lf = &ListFile{
        ModifiedTime: mtime,
        IsLink: isLink,
        IsDir: isDir,
        LinksTo: linkPath,
        Size: fSize,
        Name: fName,
        Children: children,
    }

    return

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

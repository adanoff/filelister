package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
    "encoding/json"
    "strings"
    "io/ioutil"
    "log"
    "time"
    "bytes"
    "gopkg.in/yaml.v2"
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

// print a ListFile as text
func (lf *ListFile) TextPrint(level int) {

    for i := 0; i < level; i++ {
        fmt.Print("\t")
    }

    displayName := lf.Name

    if lf.IsLink {
        displayName += "*  ->  " + lf.LinksTo
    } else if lf.IsDir {
        displayName += "/"
    }

    fmt.Println(displayName)

    for _, child := range(lf.Children) {
        child.TextPrint(level + 1)
    }

}

// print a ListFile as json
func (lf *ListFile) JSONPrint(logger *log.Logger) {

    jsonBytes, err := json.Marshal(lf)
    if err != nil {
        logger.Fatal(err)
    }

    var out bytes.Buffer
    json.Indent(&out, jsonBytes, "", "    ")
    out.WriteTo(os.Stdout)
    fmt.Println()

}

// print a ListFile as YAML
func (lf *ListFile) YAMLPrint(logger *log.Logger) {

    yamlBytes, err := yaml.Marshal(lf)
    if err != nil {
        logger.Fatal(err)
    }

    fmt.Println(string(yamlBytes))

}

type WalkFunc func(string, []*ListFile, *log.Logger)

// list files in listing (at path) in text format
func textWalk(path string, listing []*ListFile, logger *log.Logger) {

    if !strings.HasSuffix(path, "/") {
        path += "/"
    }

    fmt.Println(path)

    for _, lf := range(listing) {
        lf.TextPrint(1)
    }

}

// list files in listing (at path) in JSON format
func jsonWalk(path string, listing []*ListFile, logger *log.Logger) {

    jsonBytes, err := json.Marshal(listing)
    if err != nil {
        logger.Fatal(err)
    }

    var out bytes.Buffer
    json.Indent(&out, jsonBytes, "", "  ")
    out.WriteTo(os.Stdout)
    fmt.Println()

}

// list files in listing (at path) in YAML format
func yamlWalk(path string, listing []*ListFile, logger *log.Logger) {

    yamlBytes, err := yaml.Marshal(listing)
    if err != nil {
        logger.Fatal(err)
    }

    fmt.Println(string(yamlBytes))

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

    outFmt = strings.ToLower(outFmt)

    logger := log.New(os.Stderr, "", 0)

    walkFuncs := make(map[string]WalkFunc)
    walkFuncs["text"] = textWalk
    walkFuncs["yaml"] = yamlWalk
    walkFuncs["json"] = jsonWalk

    // show help and quit
    if showHelp {
        flag.Usage()
        os.Exit(0)
    }

    // empty path, quit
    if path == "" {
        logger.Fatal("path cannot be empty\n")
    }

    fInfos, err := ioutil.ReadDir(path)

    // make sure we could read the directory
    if err != nil {
        logger.Fatal(path, " - ", err)
    }

    listing := []*ListFile{}

    // convert the directory contents into ListFiles
    for _, fInfo := range(fInfos) {
        listing = append(listing, toListFile(fInfo, path, recursive, logger))
    }

    // invalid output format
    if _, ok := walkFuncs[outFmt]; !ok {
        logger.Fatal("output format must be one of text|json|yaml\n");
    }

    walkFuncs[outFmt](path, listing, logger)

}

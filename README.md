# filelister

List directories using different formats (text, JSON, YAML).

## Installation

Run

```
go get -u github.com/adanoff/filelister
```

to install `filelister` into your `$GOPATH`.

## Usage

```
Usage of filelister:
  -help
      show this message and exit
  -output string
      output format - json|yaml|text (default "text")
  -path string
      path to folder
  -recursive
      list files recursively
```

**NOTE:** When specifying options, you can use either one or two `-`s (e.g.
`--help` is the same as `-help`)

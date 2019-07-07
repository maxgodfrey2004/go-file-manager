[![Build Status](https://travis-ci.com/maxgodfrey2004/go-file-manager.svg?branch=master)](https://travis-ci.com/maxgodfrey2004/go-file-manager)
[![Coverage Status](https://coveralls.io/repos/github/maxgodfrey2004/go-file-manager/badge.svg?branch=master)](https://coveralls.io/github/maxgodfrey2004/go-file-manager?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/maxgodfrey2004/go-file-manager)](https://goreportcard.com/report/github.com/maxgodfrey2004/go-file-manager)
[![GolangCI](https://golangci.com/badges/github.com/maxgodfrey2004/go-file-manager.svg)](https://golangci.com/r/github.com/maxgodfrey2004/go-file-manager)

# go-file-manager
A simple text based file manager written in Go.

Before running the application, have a read of the table under the "Controls" heading below. It is probably good to know which key does what before throwing yourself in the deep end!

## Contents

  * [Controls](#controls) (Read me!)
  * [Installation and Building](#installation-and-building)
    * [Installing Go](#installing-go)
    * [Installing Dependencies](#installing-dependencies)
    * [Building the Application](#building-the-application)
      * [Getting to the right directory](#getting-to-the-right-directory)
      * [Installing local packages](#installing-local-packages)
      * [Building the app](#building-the-app)

## Controls

Note that keys separated by a comma do not have to be pressed together to activate pertinent functionality. The comma means that pressing any of the listed keys will activate the functionality listed in the next column.

| Key(s)                  | Functionality                               |
| ----------------------- | ------------------------------------------- |
| `Arrow Up`              | Move the caret to the file/directory above  |
| `Arrow Right`, `Return` | Move to the current selected directory      |
| `Arrow Down`            | Move the caret to the file/directory below  |
| `A`, `a`                | Toggle listing all files                    |
| `Q`, `q`                | Quit the application                        |

## Installation and Building

### Installing Go

To install Go, have a read of the [official installation page](https://golang.org/doc/install); and follow the instructions.

### Installing Dependencies

This project depends on [termbox-go](github.com/nsf/termbox-go). To install it, paste the following command into your terminal:

```bash
go get github.com/nsf/termbox-go
```

Then, install the application by running:

```bash
go get github.com/maxgodfrey2004/go-file-manager
```

### Building the application

#### Getting to the right directory

To build the application, you must first be located in the application's root directory. This can be found in your `GOPATH`. If you are not sure of where this is, you will be after running the command:

```bash
go env GOPATH
```

#### Installing local packages

You are almost there! Now we just have to install the local packages which `main.go` depends on. Do this by running the command:

```bash
go get -t ./...
```

#### Building the app

Build the application by running the following command in your terminal:

```bash
go build -o gfm.exe main.go
```

Note that there are other forms of executable which you may produce, for example you may execute:

```bash
go build -o gfm main.go
```

Then, run the executable and revel in all of your new text-based file managing powers! I hope that you find this application useful!

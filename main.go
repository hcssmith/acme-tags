package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	)

type MetaTags struct {
	MetaTag string
	Option string
	Comment string
}

type FunctionTags struct {
	Name string
	File string
	Line int
	Kind string
	Signature string
}

type VariableTags struct {
	Name string
	File string
	Line int
	Kind string
	Signature string
}

type Ctags struct {
	Files     []string
	Meta      []MetaTags
	Functions []FunctionTags
	Variables []VariableTags
}

func u_add_if_missing(slice []string, entry string) (result []string) {
	var b bool
	b = false
	result = slice
	for _, value := range slice {
		if value == entry {
			b = true
		}
	}
	if b == false {
		result = append(result, entry)
	}
	return
}

func find_tags_file() (path string, timestamp time.Time) {
	cwd, _ := os.Getwd()
	files, _ := ioutil.ReadDir(cwd)
	for _, value := range files {
		switch value.Name() {
			case "tags":
				path = cwd+"/tags"
			case "TAGS":
				path = cwd+"/TAGS"
			case "ctags":
				path = cwd+"/ctags"
			default:
				continue
		}
		if path != "" {
			timestamp = value.ModTime()
			return
		}
		path = cwd+"/tags"
	}
	return
}

func parse_ctags_file(path string) (ctags Ctags) {
	contents, _ := ioutil.ReadFile(path)
	var files     []string
	var metatags  []MetaTags
	var functions []FunctionTags
	var variables []VariableTags
	length := len(strings.Split(string(contents[:]), "\n"))
	for key, value := range strings.Split(string(contents[:]), "\n") {
		if key == length - 1 {
			break
		}
		baseline := strings.Split(value, "\t")
		if baseline[0][0] == '!' {
			var meta MetaTags
			meta.MetaTag = baseline[0]
			meta.Option  = baseline[1]
			meta.Comment = baseline[2]
			metatags = append(metatags, meta)
		}
		if len(baseline) >= 5 {
			if baseline[3] == "f" {
				var fun FunctionTags
				fun.Name = baseline[0]
				fun.File = baseline[1]
				files = u_add_if_missing(files, baseline[1])
				fun.Line, _ = strconv.Atoi(strings.Split(baseline[2], ";")[0])
				fun.Kind = baseline[3]
				fun.Signature = baseline[4]
				functions = append(functions, fun)
			} else if baseline[3] == "v" {
				var v VariableTags
				v.Name = baseline[0]
				v.File = baseline[1]
				files = u_add_if_missing(files, baseline[1])
				v.Line, _ = strconv.Atoi(strings.Split(baseline[2], ";")[0])
				v.Kind = baseline[3]
				v.Signature = baseline[4]
				variables = append(variables, v)
			}
		}
	}
	ctags.Files     = files
	ctags.Meta      = metatags
	ctags.Functions = functions
	ctags.Variables = variables
	return
}

func regen_ctags_file_timestamp(tagsfile string, tagstime time.Time,  tags Ctags) (ctags Ctags, timestamp time.Time) {
	cwd, _ := os.Getwd()
	files, _ := ioutil.ReadDir(cwd)
	regen := false
	for _, file := range files {
		for _, name := range tags.Files {
			if name == file.Name() && tagstime < file.ModTime() {
				regen = true
			}
		}
	}
	if regen == true {
		// regen cmd (os.exec?)
		// get timestamp of tags
		// parse
	} else {
		ctags = tags
	}
	return
}


// Find tags file in cwd
// regen if needed (timestamps?)
// check plumber is running
// Parse to acme filelinks
// Start ev_loop
// Cleanup then exit on close
func main() {
	tagsfile, _ := find_tags_file()
	tags := parse_ctags_file(tagsfile)
	fmt.Printf("%s\n", tags.Files)
	fmt.Printf("%s=> %s:%d\n",tags.Functions[0].Name, tags.Functions[0].File, tags.Functions[0].Line)

}
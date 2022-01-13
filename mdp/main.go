package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
{{ .Body }}
  </body>
</html>
`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	file := flag.String("file", "", "markdown file to preview")
	skip := flag.Bool("skip", false, "skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	if *file == "" {
		fmt.Println("app exit: argument not passed")
		fmt.Printf("usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(*file, *tFname, os.Stdout, *skip); err != nil {
		fmt.Fprintln(os.Stderr, "err")
		os.Exit(1)
	}
}

func run(file string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	htmlData, err := parseInput(input, tFname)
	if err != nil {
		return err
	}

	temp, err := ioutil.TempFile("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)
	err = saveHTML(outName, htmlData)
	if err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)
	return preview(outName)
}

func saveHTML(name string, data []byte) error {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func parseInput(input []byte, tFname string) ([]byte, error) {
	var buffer bytes.Buffer
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	if tFname != "" {
		log.Println("I am here")
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Body:  template.HTML(body),
		Title: "Markdown preview tool",
	}
	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("%v: platform not supported", runtime.GOOS)
	}

	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return err
}

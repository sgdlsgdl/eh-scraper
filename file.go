package main

import (
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
	template2 "html/template"
	"log"
	"os"
	"text/template"
)

const (
	configPath = "./config.yaml"
	basePath   = "./output/"
	csvPath    = basePath + "total.csv"
	mdPath     = basePath + "new.md"
	htmlPath   = basePath + "new.html"
)

const (
	mdTemplateStr = `
# emm
| Gallery | Image | Date | Name | Link |
|---|---|---|---|---|
{{- range .}}
| {{.Gallery}} | <img src="{{.Image}}" width="50"> | {{.Date}} | {{.Name}} | [Link]({{.Link}}) |
{{- end}}
`

	htmlTemplateStr = `
<!DOCTYPE html>
<html>
<head>
<title>emm</title>
</head>
<body>
<table>
    <thead>
        <tr>
            <th>Gallery</th>
            <th>Image</th>
            <th>Date</th>
            <th>Link</th>
        </tr>
    </thead>
    <tbody>
        {{range .}}
        <tr>
            <td>{{.Gallery}}</td>
            <td><img src="{{.Image}}" alt="{{.Name}}"></td>
            <td>{{.Date}}</td>
            <td><a href="{{.Link}}">{{.Name}}</a></td>
        </tr>
        {{end}}
    </tbody>
</table>
</body>
</html>
`
)

func readConfig(path string) Config {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read Yaml file: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		log.Fatalf("Failed to decode Yaml: %v", err)
	}
	return config
}

func readCsv(path string) (res ItemList) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, &res); err != nil {
		log.Printf("gocsv.UnmarshalFile error %v", err)
		return
	}
	return
}

func saveCsv(path string, items ItemList) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("os.OpenFile error %v", err)
		return
	}
	defer file.Close()

	if err := gocsv.MarshalFile(items, file); err != nil {
		log.Printf("gocsv.MarshalFile error %v", err)
		return
	}
	return
}

func saveMd(path string, items ItemList) {
	tmpl, err := template.New("itemList").Parse(mdTemplateStr)
	if err != nil {
		log.Printf("template.New error %v", err)
		return
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("os.OpenFile error %v", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, items)
	if err != nil {
		log.Printf("tmpl.Execute error %v", err)
		return
	}
}

func saveHtml(path string, items ItemList) {
	tmpl, err := template2.New("itemList").Parse(htmlTemplateStr)
	if err != nil {
		log.Printf("template.New error %v", err)
		return
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("os.OpenFile error %v", err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, items)
	if err != nil {
		log.Printf("tmpl.Execute error %v", err)
		return
	}
}

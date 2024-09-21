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
# data
| Gallery | Image | Date | Name | Link | Key |
|---|---|---|---|---|---|
{{- range .}}
| {{.Gallery}} | <img src="{{.Image}}" width="50"> | {{.Date}} | {{.Name}} | [Link]({{.Link}}) | {{.Key}} |
{{- end}}
`

	htmlTemplateStr = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Data</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <style>
        body {
            margin: 20px;
        }
        table {
            width: 100%;
        }
        th, td {
            text-align: center;
            vertical-align: middle;
        }
        img {
            max-width: 1000px;
            height: auto;
        }
    </style>
</head>
<body>
<div class="container">
    <h2 class="text-center my-4">Data Gallery</h2>
    <table class="table table-striped table-bordered">
        <thead class="thead-dark">
            <tr>
                <th>Gallery</th>
                <th>Image</th>
                <th>Date</th>
                <th>Link</th>
                <th>Key</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.Gallery}}</td>
                <td><img src="{{.Image}}" alt="{{.Name}}"></td>
                <td>{{.Date}}</td>
                <td><a href="{{.Link}}" class="btn btn-primary">{{.Name}}</a></td>
                <td>{{.Key}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</div>
<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.3/dist/umd/popper.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js"></script>
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
	if config.Concurrency <= 0 {
		config.Concurrency = 5
	}
	if config.MaxDay <= 0 {
		config.MaxDay = 90
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

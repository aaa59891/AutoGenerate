package main

import (
	"io/ioutil"
	"log"
	"strings"
	"reflect"
	"bitbucket.org/mositech/AutoGenerate/utils"
	"time"
)

var models = []string{
	"Test",
	"Chong",
	"Person",
}

type IdRule struct {
	Id          uint   `gorm:"primary key"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	DigitNumber int    `json:"digitNumber"`
	IsDateShow  bool   `json:"isDateShow"`
	DateFormat  string
	Category    int

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

var AutoSlice = []interface{}{
	IdRule{},
}

func main() {
	controllerTemplate, err := ioutil.ReadFile("./templates/controllers/ControllerTemplate.txt")
	if err != nil{
		log.Panic(err)
	}

	listHtmlTemplate, err := ioutil.ReadFile("./templates/views/listView.html")
	if err != nil{
		log.Panic(err)
	}

	editHtmlTemplate, err := ioutil.ReadFile("./templates/views/editView.html")
	if err != nil{
		log.Panic(err)
	}


	for _, model := range AutoSlice{

		listHtml, listFilename := CreateHtmlList(model, listHtmlTemplate)
		if err := ioutil.WriteFile("./export/html/" + listFilename, []byte(listHtml), 0755); err != nil{
			log.Panic(err)
		}

		editHtml, editFilename := CreateHtmlEdit(model, editHtmlTemplate)
		if err := ioutil.WriteFile("./export/html/" + editFilename, []byte(editHtml), 0755); err != nil{
			log.Panic(err)
		}

		contData, contFilename := CreateGoController(model, controllerTemplate)
		if err := ioutil.WriteFile("./export/controller/" + contFilename, []byte(contData), 0755); err != nil{
			log.Panic(err)
		}

	}
}

func CreateGoController(model interface{}, template []byte)(fileStr, fileName string){
	t := reflect.TypeOf(model)
	result := string(template)
	modelName := t.Name()
	result = strings.Replace(result, "MODEL_NAME", modelName, -1)

	return result, modelName + "Controller.go"
}

func CreateHtmlList(model interface{}, template []byte) (html, fileName string){
	result := string(template)

	t := reflect.TypeOf(model)

	modelName := utils.ToLowerFirst(t.Name())

	var tHeader, tData string

	for i := 0; i < t.NumField(); i++{
		fieldName := t.Field(i).Name
		firstLowerName := utils.ToLowerFirst(fieldName)
		tHeader += "<th>" + fieldName + "</th>\n"
		tData += `<td class="` + firstLowerName + `">{{.` + fieldName + `}}</td>` + "\n"
	}

	result = strings.Replace(result, "-MODEL_NAME-", modelName, -1)
	result = strings.Replace(result, "-TABLE_HEADER-", tHeader, -1)
	result = strings.Replace(result, "-TABLE_DATA-", tData, -1)

	return result, t.Name() + "List.html"
}

func CreateHtmlEdit(model interface{}, template []byte)(html, fileName string){
	result := string(template)

	t := reflect.TypeOf(model)

	modelName := utils.ToLowerFirst(t.Name())

	var formData string

	for i := 0; i < t.NumField(); i++{
		fieldName := t.Field(i).Name
		firstLowerName := utils.ToLowerFirst(fieldName)
		formData += `
		<div class="form-group">
			<label>` + fieldName + `</label>
		`
		if fieldName == "Id"{
			formData += `<input class="form-control" required name="id" {{if .data.Id}} readonly {{end}} value="{{.data.Id}}">`
		}else{
			switch t.Field(i).Type.Kind(){
			case reflect.String:
				formData += `<input class="form-control" name="` + firstLowerName + `" value="{{.data.` + fieldName + `}}">`
			case reflect.Int:
				formData += `<input class="form-control" type="number" name="` + firstLowerName + `" value="{{.data.` + fieldName + `}}">`
			case reflect.Bool:
				formData +=
					`<br><label class="radio-inline"><input name="` + firstLowerName + `" type="radio" value="1" {{if .data.` + fieldName + `}}checked{{end}}>Yes</label>
					<label class="radio-inline"><input name="` + firstLowerName + `" type="radio" value="0" {{if not .data.` + fieldName + `}}checked{{end}}>No</label>`
			}
		}

		formData += "\n"

		formData += "</div>"
	}


	result = strings.Replace(result, "-MODEL_NAME-", modelName, -1)
	result = strings.Replace(result, "-FORM_DATA-", formData, -1)

	return result, t.Name() + "Edit.html"
}


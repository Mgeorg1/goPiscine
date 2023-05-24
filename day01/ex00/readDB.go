package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Ingredients struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type Cake struct {
	Name       string        `xml:"name" json:"name"`
	Time       string        `xml:"stovetime" json:"time"`
	Ingredient []Ingredients `xml:"ingredients>item" json:"ingredients"`
}

type Recipes struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Cakes   []Cake   `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read(file []byte) (Recipes, error)
}

type Json Recipes
type XML Recipes

func (ptr *Json) Read(file []byte) (Recipes, error) {
	err := json.Unmarshal(file, ptr)
	return Recipes(*ptr), err

}
func (ptr *XML) Read(file []byte) (Recipes, error) {
	err := xml.Unmarshal(file, ptr)
	return Recipes(*ptr), err
}

func Rewrite(reader DBReader, file []byte) ([]byte, error) {
	recipes, err := reader.Read(file)
	if err != nil {
		return nil, err
	}
	var result []byte
	switch reader.(type) {
	case *Json:
		result, err = xml.MarshalIndent(recipes, "", "    ")
	case *XML:
		result, err = json.MarshalIndent(recipes, "", "    ")
	default:
		break
	}
	return result, err

}

func RemoveComments(file *[]byte) {
	sFile := string(*file)
	var x, y int

	runes := []rune(sFile)
	var buf []rune
	for i := range runes {
		if i > 0 && runes[i] == '/' && runes[i-1] == '/' {
			x = i
		}
		if x > 0 && runes[i] == '\n' {
			y = i
			buf = append(runes[0:x-1], runes[y:]...)
			x = 0
		}
	}

	if len(buf) > 0 {
		newS := string(buf)
		*file = []byte(newS)
	}
}

func main() {
	var f bool

	flag.BoolVar(&f, "f", false, "FilePath")
	flag.Parse()
	if !f {
		fmt.Println("Specify the filepath: -f filePath")
		os.Exit(1)
	}
	if len(os.Args) < 3 {
		fmt.Println("Specify the filepath: -f filePath")
		os.Exit(1)
	}
	filePath := os.Args[2]
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lenght := len(filePath)
	var isXml bool
	if lenght >= 4 && filePath[lenght-4:] == ".xml" {
		isXml = true
	} else if lenght >= 5 && filePath[lenght-5:] == ".json" {
		isXml = false
	} else {
		fmt.Println("Unsuported file format")
		os.Exit(1)
	}
	RemoveComments(&file)
	var res []byte
	if isXml {
		recipes := new(XML)
		res, err = Rewrite(recipes, file)
	} else {
		recipes := new(Json)
		res, err = Rewrite(recipes, file)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s\n", res)
}

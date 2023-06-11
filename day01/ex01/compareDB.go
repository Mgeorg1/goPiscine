package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/r3labs/diff"
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

func pairs(p []string, r *Recipes) string {
	pairs := make([]string, len(p)/2+len(p)%2)
	var a, b int
	for a = len(pairs) - 1; b < len(p)-1; b, a = b+2, a-1 {
		idx, err := strconv.Atoi(p[b+1])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if p[b] == "Cakes" {
			pairs[a] = fmt.Sprintf("%s %s", p[b], r.Cakes[idx].Name)
		} else {
			idx1, err := strconv.Atoi(p[1])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			pairs[a] = fmt.Sprintf("%s %s", p[b], r.Cakes[idx1].Ingredient[idx1].Name)
		}
	}
	if a == 0 {
		pairs[a] = p[b]
	}
	return strings.Join(pairs, " for ")
}

func CompareRecipes(old *Recipes, new *Recipes) {
	differ, err := diff.NewDiffer(diff.DisableStructValues(), diff.SliceOrdering(false))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log, err := differ.Diff(old, new)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, change := range log {
		a := change.Path
		if a[0] == "XMLName" {
			continue
		}
		switch change.Type {
		case diff.CREATE:
			fmt.Printf("ADDED %s\n", pairs(a, new))
		case diff.UPDATE:
			fmt.Printf("CHANGED %s - %s instead of %s\n", pairs(a, new), change.To, change.From)
		case diff.DELETE:
			switch n := len(a) - 1; a[n] {
			case "unit":
				a = append(a, change.From.(string))
			case "ingredient":
				a = a[:n]
			}
			fmt.Printf("REMOVED %s\n", pairs(a, old))
		}

	}
}

func Parse(reader DBReader, file []byte) (Recipes, error) {
	recipes, err := reader.Read(file)
	return recipes, err
}

func main() {
	var f1, f2 string
	flag.StringVar(&f1, "f1", "", "FilePath of old file")
	flag.StringVar(&f2, "f2", "", "FilePath of the old file")
	if len(os.Args) < 5 {
		fmt.Println("Specify the filepath: -f1 filePathOld -f2 filePathNew")
		os.Exit(1)
	}
	flag.Parse()

	if len(f1) == 0 || len(f2) == 0 {
		fmt.Println("Specify the filepath: -f1 and -f2 filePaths")
		os.Exit(1)
	}
	filePath := os.Args[2]
	file, err := os.ReadFile(filePath)
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
	var res1 Recipes
	if isXml {
		recipes := new(XML)
		res1, err = Parse(recipes, file)
	} else {
		recipes := new(Json)
		res1, err = Parse(recipes, file)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	filePath = os.Args[4]
	file, err = os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lenght = len(filePath)
	if lenght >= 4 && filePath[lenght-4:] == ".xml" {
		isXml = true
	} else if lenght >= 5 && filePath[lenght-5:] == ".json" {
		isXml = false
	} else {
		fmt.Println("Unsuported file format")
		os.Exit(1)
	}
	RemoveComments(&file)
	var res2 Recipes
	if isXml {
		recipes := new(XML)
		res2, err = Parse(recipes, file)
	} else {
		recipes := new(Json)
		res2, err = Parse(recipes, file)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	CompareRecipes(&res1, &res2)
}

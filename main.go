package main;

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/xml"
)

const (
	XML_HEADER  = `<?xml version="1.0"?>` + "\n"
	WOOORDHUNT_SITE = "http://wooordhunt.ru"
)

type Tip struct {
	Word string `json:"w"`
	Translate string `json:"t"`
}

type Tips struct {
	Tips []Tip `json:"tips"`
}

type ItemIcon struct {
	Path string `json:"path" xml:"path,attr"`
}

type Item struct {
	Title string `xml:"title" json:"title"`
	SubTitle string `xml:"subtitle" json:"subtitle"`
	Arg string `xml:"arg,attr" json:"arg"`
	Icon ItemIcon `xml:"icon" json:"icon"`
}

type Items struct {
	XMLName xml.Name `xml:"items" json:"-"`
	Item []Item `xml:"item" json:"items"`
}

func main() {
	var query *string = flag.String("q", "", "Search word")
	var format *string = flag.String("f", "xml", "Output type")
	flag.Parse()

	uri := WOOORDHUNT_SITE + "/get_tips.php?abc=" + *query

	resp, _ := http.Get(uri)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonResponse Tips
	json.Unmarshal(body, &jsonResponse)

	var items []Item
	for _, v := range jsonResponse.Tips {
		items = append(items, Item{
			Title: v.Word,
			SubTitle: v.Translate,
			Arg: WOOORDHUNT_SITE + "/word/" + v.Word,
			Icon: ItemIcon{Path: "icon.png"},
		})
	}

	output := ""
	if *format == "json" {
		result, _ := json.Marshal(Items{Item: items})
		output = string(result)
	} else if *format == "xml" {
		result, _ := xml.Marshal(Items{Item: items})
		output = string(result)
	}
	fmt.Println(output)
}
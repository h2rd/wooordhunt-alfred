package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	XML_HEADER      = `<?xml version="1.0"?>` + "\n"
	WOOORDHUNT_SITE = "https://wooordhunt.ru"
)

type Tip struct {
	Word      string `json:"w"`
	Translate string `json:"t"`
}

type Tips struct {
	Tips []Tip `json:"tips"`
}

type ItemIcon struct {
	Path string `json:"path" xml:"path,attr"`
}

type Item struct {
	Title    string   `xml:"title" json:"title"`
	SubTitle string   `xml:"subtitle" json:"subtitle"`
	Arg      string   `xml:"arg,attr" json:"arg"`
	Icon     ItemIcon `xml:"icon" json:"icon"`
}

type Items struct {
	XMLName xml.Name `xml:"items" json:"-"`
	Item    []Item   `xml:"item" json:"items"`
}

func main() {
	var targetQuery *string = flag.String("q", "", "Search word")
	var format *string = flag.String("f", "xml", "Output type")
	flag.Parse()

	uri := fmt.Sprintf(
		"%s/openscripts/forjs/get_tips.php?abc=%s",
		WOOORDHUNT_SITE,
		strings.ToLower(*targetQuery),
	)

	resp, _ := http.Get(uri)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var jsonResponse Tips
	json.Unmarshal(body, &jsonResponse)

	var items []Item
	for _, v := range jsonResponse.Tips {
		query := strings.ToLower(v.Word)
		items = append(items, Item{
			Title:    query,
			SubTitle: v.Translate,
			Arg:      WOOORDHUNT_SITE + "/word/" + query,
			Icon:     ItemIcon{Path: "icon.png"},
		})
	}

	if len(items) == 0 {
		items = append(items, Item{
			Title:    strings.ToLower(*targetQuery),
			SubTitle: "The requested word was not found.",
			Arg:      WOOORDHUNT_SITE,
			Icon:     ItemIcon{Path: "icon.png"},
		})
	}

	output := ""
	if *format == "json" {
		result, _ := json.Marshal(Items{Item: items})
		output = string(result)
	} else if *format == "xml" {
		result, _ := xml.Marshal(Items{Item: items})
		output = XML_HEADER + string(result)
	}

	fmt.Println(output)
}

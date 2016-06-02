package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"

	h "github.com/sunreaver/gotools/http"
)

var (
	reg = regexp.MustCompile(`<td id="LC[0-9]+" class="blob-code blob-code-inner js-file-line">(.*)</td>`)
)

func main() {
	c, e := h.Get("https://github.com/gfwlist/gfwlist/blob/master/gfwlist.txt")
	if e != nil {
		log.Println("Get Err")
		return
	}

	pacs := reg.FindAllStringSubmatch(c, -1)

	var b64 string
	var result []string
	for _, item := range pacs {
		b64 += item[1]
	}

	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println("Base64 Err:", err)
		return
	}
	strs := strings.Split(string(data), "\n")
	for _, item := range strs {
		if strings.HasPrefix(item, "!") ||
			len(item) == 0 {
			continue
		}
		result = append(result, item)
	}

	for n, item := range result {
		if n == len(result)-1 {
			fmt.Printf("  \"%s\"\n", item)
		} else {
			fmt.Printf("  \"%s\",\n", item)
		}
	}
}

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

	h "github.com/sunreaver/gotools/http"
	"github.com/sunreaver/gotools/system"
)

var (
	reg = regexp.MustCompile(`<td id="LC[0-9]+" class="blob-code blob-code-inner js-file-line">(.*)</td>`)
)

func main() {
	result := getGfwlist("https://github.com/gfwlist/gfwlist/blob/master/gfwlist.txt")
	if len(result) == 0 {
		return
	}

	fileName := makeJsFile(result)
	if len(fileName) == 0 {
		return
	}

	e := moveFile2ShadowsocksX(fileName)
	if e != nil {
		log.Println("MoveFile Err:", e)
		return
	}

	log.Println("Update OK")
}

func getGfwlist(url string) (gfwlist []string) {
	c, e := h.Get(url)
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
			strings.HasPrefix(item, "[") ||
			len(item) == 0 {
			continue
		}
		result = append(result, item)
	}
	gfwlist = result
	return
}

func makeJsFile(gfwlist []string) (fileName string) {
	fileName = system.CurPath() + system.SystemSep() + time.Now().Format("2006_01_02_15_04_05.js")
	if system.IsFileExists(fileName) {
		os.Remove(fileName)
	}

	f, err2 := os.Create(fileName)
	if err2 != nil {
		log.Println("CreateFile Err:", err2)
		fileName = ""
		return
	}

	tmpl, err1 := template.ParseFiles(system.CurPath() + system.SystemSep() + "gfwlist.tmpl")
	if err1 != nil {
		log.Println("Tmpl Err:", err1)
		fileName = ""
		return
	}
	tmpl.Execute(f, gfwlist)
	f.Close()
	return
}

func moveFile2ShadowsocksX(source string) (e error) {
	path := strings.Split(system.CurPath(), system.SystemSep())
	desPath := fmt.Sprintf("/%s/%s/.ShadowsocksX/gfwlist.js", path[1], path[2])

	cmd := exec.Command("mv", source, desPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	e = cmd.Run()
	if e != nil {
		return
	}
	return
}

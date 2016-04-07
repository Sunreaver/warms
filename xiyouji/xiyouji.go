package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	sys "github.com/sunreaver/goTools/system"
	"github.com/sunreaver/mahonia"
)

func main() {
	for i := 0; i <= 120; i++ {
		url := "http://www.ziyexing.com/book/shuihu/shuihu_index.htm"
		if i != 0 {
			url = fmt.Sprintf("http://www.ziyexing.com/book/shuihu/shuihu_%03d.htm", i)
		}

		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		r, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			continue
		}
		s := mahonia.NewDecoder("gb2312").ConvertString(string(r))
		s = strings.Replace(s, "www.ziyexing.com", "23.83.239.152", -1)
		s = strings.Replace(s, "charset=gb2312", "charset=utf-8", -1)
		s = strings.Replace(s, "子夜星网站", "sunreaver", -1)
		s = strings.Replace(s, "Midnight Star", "sunreaver", -1)
		s = strings.Replace(s, "../../images/", "../images/", -1)

		dir := sys.CurPath()
		fileName := dir + sys.SystemSep() + fmt.Sprintf("shuihu_%03d.htm", i)
		if i == 0 {
			fileName = dir + sys.SystemSep() + "shuihu_index.htm"
		}

		f, e1 := os.Create(fileName)
		if e1 != nil {
			f, e = os.Open(fileName)
			if e != nil {
				log.Println("Error with Create : ", fileName)
				continue
			}
		}

		s = mahonia.NewEncoder("utf-8").ConvertString(s)
		n, e2 := f.WriteString(s)
		if e2 != nil {
			log.Println("Error with Write : ", fileName)
		} else {
			log.Println("Write file : ", fileName)
			log.Println("Write size : ", n)
		}
		f.Close()

		//下载资源图
		getPic(s)
	}
}

func getPic(s string) {
	rg := regexp.MustCompile(`(?U)src="(.*).gif"`)
	matchs := rg.FindAllStringSubmatch(s, -1)

	for _, item := range matchs {
		if strings.Contains(item[1], "./") {
			continue
		}

		dir := sys.CurPath()
		fileName := dir + sys.SystemSep() + item[1] + ".gif"

		if sys.IsFileExists(fileName) {
			continue
		}
		f, e := os.Create(fileName)
		if e != nil {
			log.Println("Error Create: ", fileName)
			continue
		}

		resp, e1 := http.Get(fmt.Sprintf("http://www.ziyexing.com/book/shuihu/%s.gif", item[1]))
		if e1 != nil {
			log.Println("Error Download: ", item[1], ".gif")
			continue
		}
		defer resp.Body.Close()
		r, e2 := ioutil.ReadAll(resp.Body)
		if e2 != nil {
			log.Println("Error Read: ", item[1], ".gif")
			continue
		}
		n, e3 := f.Write(r)
		if e3 != nil {
			log.Println("Error Write: ", item[1], ".gif")
		} else {
			log.Println("Download: ", item[1], ".gif size: ", n)
		}
		f.Close()
	}
}

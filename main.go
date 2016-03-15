package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tanweirush/mahonia"
)

var (
	StockDataRegexp = regexp.MustCompile(`="(.*)";`)
)

func main() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	stocks, e := readFile(filepath.Dir(path) + systemSep() + "stock.json")
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	format := "%s:\t>>>>>昨收:%s\t今收:%s\t今开:%s\t幅度:%s<<<<<\r\n"
	outStr := ""

	for _, v := range stocks {
		resp, err := http.Get("http://hq.sinajs.cn/list=" + v)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			continue
		}

		str := string(body)
		enc := mahonia.NewDecoder("gbk")
		matchs := StockDataRegexp.FindAllStringSubmatch(enc.ConvertString(str), -1)
		if len(matchs) == 0 || len(matchs[0]) < 2 {
			continue
		}

		numerical := strings.Split(matchs[0][1], ",")
		if len(numerical) < 5 {
			continue
		}

		upDown := "未知"
		yestoday, e1 := strconv.ParseFloat(numerical[2], 64)
		today, e2 := strconv.ParseFloat(numerical[3], 64)
		if e1 == nil && e2 == nil {
			upDown = strconv.FormatFloat(today-yestoday, 'f', -1, 64)[0:6]
		}
		outStr = outStr + fmt.Sprintf(format, numerical[0], numerical[2], numerical[3], numerical[1], upDown)
	}

	outStr = outStr + "\r\nHappy day!\r\n"

	e = SendMail(outStr, []string{"tanwei.rush@gmail.com"})
	if e != nil {
		fmt.Println(e.Error())
		fmt.Println(outStr)
	}
}

func readFile(fileName string) ([]string, error) {
	b, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil, e
	}
	var stocks []string
	e = json.Unmarshal(b, &stocks)
	if e != nil {
		return nil, e
	}
	return stocks, nil
}

//获取系统分隔符
func systemSep() (path string) {

	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	return path
}

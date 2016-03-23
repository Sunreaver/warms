package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	tanweiTools "github.com/sunreaver/goTools/system"
)

var (
	ptnIndexItem    = regexp.MustCompile(`app\.page\["pins"\] = (.*);\napp\.page\["ads"\]`)
	ptnContentRough = regexp.MustCompile(`(?s).*<div class="artcontent">(.*)<div id="zhanwei">.*`)
	ptnBrTag        = regexp.MustCompile(`<br>`)
	ptnHTMLTag      = regexp.MustCompile(`(?s)</?.*?>`)
	ptnSpace        = regexp.MustCompile(`(^\s+)|( )`)

	warmUrl = `http://huaban.com/favorite/beauty/`
	imgPath = `http://img.hb.aicdn.com/`

	FileHadExist = errors.New("文件已经存在")
)

type HuaBan struct {
	Board struct {
		BoardID     int    `json:"board_id"`
		CategoryID  string `json:"category_id"`
		CreatedAt   int    `json:"created_at"`
		Deleting    int    `json:"deleting"`
		Description string `json:"description"`
		Extra       string `json:"extra"`
		FollowCount int    `json:"follow_count"`
		IsPrivate   int    `json:"is_private"`
		LikeCount   int    `json:"like_count"`
		PinCount    int    `json:"pin_count"`
		Seq         int    `json:"seq"`
		Title       string `json:"title"`
		UpdatedAt   int    `json:"updated_at"`
		UserID      int    `json:"user_id"`
	} `json:"board"`
	BoardID      int `json:"board_id"`
	CommentCount int `json:"comment_count"`
	CreatedAt    int `json:"created_at"`
	File         struct {
		Bucket string `json:"bucket"`
		Farm   string `json:"farm"`
		Frames int    `json:"frames"`
		Height int    `json:"height"`
		Key    string `json:"key"`
		Theme  string `json:"theme"`
		Type   string `json:"type"`
		Width  int    `json:"width"`
	} `json:"file"`
	FileID     int      `json:"file_id"`
	IsPrivate  int      `json:"is_private"`
	LikeCount  int      `json:"like_count"`
	Link       string   `json:"link"`
	MediaType  int      `json:"media_type"`
	OrigSource string   `json:"orig_source"`
	Original   int      `json:"original"`
	PinID      int      `json:"pin_id"`
	RawText    string   `json:"raw_text"`
	RepinCount int      `json:"repin_count"`
	Source     string   `json:"source"`
	TextMeta   struct{} `json:"text_meta"`
	User       struct {
		Avatar struct {
			Bucket string `json:"bucket"`
			Farm   string `json:"farm"`
			Frames int    `json:"frames"`
			Height int    `json:"height"`
			ID     int    `json:"id"`
			Key    string `json:"key"`
			Theme  string `json:"theme"`
			Type   string `json:"type"`
			Width  int    `json:"width"`
		} `json:"avatar"`
		CreatedAt int    `json:"created_at"`
		Urlname   string `json:"urlname"`
		UserID    int    `json:"user_id"`
		Username  string `json:"username"`
	} `json:"user"`
	UserID  int `json:"user_id"`
	Via     int `json:"via"`
	ViaUser struct {
		Avatar struct {
			Bucket string `json:"bucket"`
			Farm   string `json:"farm"`
			Frames int    `json:"frames"`
			Height int    `json:"height"`
			ID     int    `json:"id"`
			Key    string `json:"key"`
			Type   string `json:"type"`
			Width  int    `json:"width"`
		} `json:"avatar"`
		CreatedAt int    `json:"created_at"`
		Urlname   string `json:"urlname"`
		UserID    int    `json:"user_id"`
		Username  string `json:"username"`
	} `json:"via_user"`
	ViaUserID int `json:"via_user_id"`
}

func Get(url string) (content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		err = err2
		return
	}
	content = string(data)
	return
}

func findIndex(content string) (hb []HuaBan, err error) {
	matches := ptnIndexItem.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, errors.New("查找json出错")
	}
	json.Unmarshal([]byte(matches[0][1]), &hb)

	return hb, nil
}

//创建目录
func makeDirWithToday(dirName string) error {

	dir := tanweiTools.CurPath()

	fullPath := dir + tanweiTools.SystemSep() + dirName
	if tanweiTools.IsDirExists(fullPath) { //目录已经存在
		return nil
	}

	err := os.Mkdir(fullPath, os.ModePerm) //在当前目录下生成新目录
	if err != nil {
		return err
	}
	return nil
}

func readContent(hb HuaBan) error {

	if !strings.HasPrefix(hb.File.Type, "image/") {
		return errors.New("不是图片")
	}

	fileType := hb.File.Type[len("image/"):]

	dir := tanweiTools.CurPath() //当前的目录
	dirName := fmt.Sprintf("huaban_%s", time.Now().Format("2006-01-02"))
	makeDirWithToday(dirName)
	filename := dir + tanweiTools.SystemSep() + dirName + tanweiTools.SystemSep() + fmt.Sprintf("%s_%d", hb.User.Username, hb.FileID) + "." + fileType

	if tanweiTools.IsFileExists(filename) {
		return FileHadExist
	}

	res, err := http.Get(imgPath + hb.File.Key)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	data, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return e
	}

	go func(h HuaBan, d []byte) {

		file, e1 := os.Create(filename)
		if e1 != nil {
			fmt.Printf("Error Create File %d.\n", h.FileID)
			return
		}
		defer file.Close()

		_, e2 := file.Write(d)
		if e2 != nil {
			fmt.Printf("Error Write File %d.\n", h.FileID)
			return
		}
	}(hb, data)

	return nil
}

func main() {
	var n = 1
	for {
		con, err := Get(warmUrl)
		if err != nil {
			fmt.Println("\r\n链接出错,3分钟后再试: " + time.Now().Format("06/01/02-15:04"))
			time.Sleep(3 * time.Minute)
			continue
		}

		index, _ := findIndex(con)

		var exist, errCount = 0, 0
		n = 0
		for _, item := range index {
			// fmt.Printf("Get content %s from %s and write to file.\n", item.title, item.url)
			n++
			e := readContent(item)
			if e == FileHadExist {
				exist++
			} else if e != nil {
				errCount++
			}
		}

		fmt.Printf("\r\n总\t保存\t已存在\t出错\r\n%d\t%d\t%d\t%d\r\n", n, n-exist-errCount, exist, errCount)

		if n-exist-errCount <= 0 {
			n = 3
		} else {
			n = (n - exist - errCount) * 5
		}

		//总休眠时间
		var sleepTime int = int(100.0 / n)
		//显示倒计时
		go func() {
			fmt.Print(sleepTime, ".")
			for i := 1; i < sleepTime; i++ {
				time.Sleep(1 * time.Minute)
				fmt.Print(".")
			}
			fmt.Println(time.Now().Format("06/01/02-15:04"))
			fmt.Printf("\r\n")
		}()
		time.Sleep(time.Duration(sleepTime) * time.Minute)
	}
}

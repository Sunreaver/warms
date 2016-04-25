package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	tanweiTools "github.com/sunreaver/gotools/system"
)

var (
	ptnIndexItem    = regexp.MustCompile(`app\.page\["pins"\] = (.*);\napp\.page\["ads"\]`)
	ptnContentRough = regexp.MustCompile(`(?s).*<div class="artcontent">(.*)<div id="zhanwei">.*`)
	ptnBrTag        = regexp.MustCompile(`<br>`)
	ptnHTMLTag      = regexp.MustCompile(`(?s)</?.*?>`)
	ptnSpace        = regexp.MustCompile(`(^\s+)|( )`)

	warmUrl = `http://huaban.com/favorite/beauty/`
	// warmUrl = `http://huaban.com/boards/28266958`
	imgPath = `http://img.hb.aicdn.com/`

	FileHadExist = errors.New("文件已经存在")

	wMinMax = [2]int{450, 1000}
	hMinMax = [2]int{600, 2000}
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
	if fileType != "png" && fileType != "jpeg" {
		return errors.New("不是png或者jpeg")
	}
	// else if hb.File.Width < wMinMax[0] || hb.File.Width > wMinMax[1] ||
	// 	hb.File.Height < hMinMax[0] || hb.File.Height > hMinMax[1] {
	// 	return errors.New("尺寸不匹配: " + strconv.Itoa(hb.File.Width) + "x" + strconv.Itoa(hb.File.Height))
	// }

	dir := tanweiTools.CurPath() //当前的目录
	dirName := fmt.Sprintf("huaban_%s", time.Now().Format("2006-01"))
	makeDirWithToday(dirName)
	ttitle := strings.Replace(hb.Board.Title, " ", "-", -1)
	filename := dir + tanweiTools.SystemSep() + dirName + tanweiTools.SystemSep() + fmt.Sprintf("%s_%d", ttitle, hb.FileID) + "." + fileType

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
		r := bytes.NewReader(d)

		//cf 用来检测图像真实尺寸
		var cf image.Config
		var e0 error
		if h.File.Type[len("image/"):] == "png" {
			cf, e0 = png.DecodeConfig(r)
		} else {
			cf, e0 = jpeg.DecodeConfig(r)
		}
		if e0 != nil {
			log.Printf("Error []byte2io.reader %d.\n", h.FileID)
			return
		} else if (cf.Height > hMinMax[1] || cf.Height < hMinMax[0]) ||
			(cf.Width < wMinMax[0] || cf.Width > wMinMax[1]) {
			log.Printf("Width||Height No Match %d. %dx%d\n", h.FileID, cf.Width, cf.Height)
			return
		} else if len(d) > 2*1024*1024 || len(d) < 50*1024 {
			log.Printf("size no match %d.\n", h.FileID)
			return
		}
		file, e1 := os.Create(filename)
		defer file.Close()
		if e1 != nil {
			log.Printf("Error Create File %d.\n", h.FileID)
			return
		}

		_, e2 := file.Write(d)
		if e2 != nil {
			log.Printf("Error Write File %d.\n", h.FileID)
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
			log.Println("\r\n链接出错,3分钟后再试: " + time.Now().Format("06/01/02-15:04"))
			time.Sleep(3 * time.Minute)
			continue
		}

		index, _ := findIndex(con)

		var exist, errCount = 0, 0
		n = 0
		for _, item := range index {
			// log.Printf("Get content %s from %s and write to file.\n", item.title, item.url)
			n++
			e := readContent(item)
			if e == FileHadExist {
				exist++
			} else if e != nil {
				errCount++
				log.Println(e)
			}
		}

		log.Printf("\r\n总\t保存\t已存在\t出错\r\n%d\t%d\t%d\t%d\r\n", n, n-exist-errCount, exist, errCount)

		if n-exist-errCount <= 0 {
			n = 3
		} else {
			n = (n - exist - errCount) * 5
		}

		//总休眠时间
		sleepTime := 30 - n
		if sleepTime <= 0 {
			sleepTime = 1
		}
		//显示倒计时
		go func() {
			log.Print(sleepTime, ".")
			for i := 1; i < sleepTime; i++ {
				time.Sleep(1 * time.Minute)
				log.Print(".")
			}
			log.Println(time.Now().Format("06/01/02-15:04"))
			log.Printf("\r\n")
		}()
		time.Sleep(time.Duration(sleepTime) * time.Minute)
	}
}

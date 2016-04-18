package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/sunreaver/gotools/system"
)

var (
	_W    = sort.IntSlice{}
	_H    = sort.IntSlice{}
	_Size = sort.IntSlice{}
)

func main() {

	if len(os.Args) < 2 {
		getMiddleSize(system.CurPath())
	} else {
		for i := 1; i < len(os.Args); i++ {
			getMiddleSize(os.Args[i])
		}
	}

	sort.Sort(_W)
	sort.Sort(_H)

	if len(_W) > 0 && len(_H) > 0 {
		log.Println("Count    : ", len(_W))
		log.Println("Middle W : ", _W[len(_W)/2])
		log.Println("Middle H : ", _H[len(_H)/2])
		log.Println("Middle Size : ", _Size[len(_Size)/2])
	} else {
		log.Println("No pic")
	}
}

func getMiddleSize(dir string) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		log.Println(e)
		return
	}

	var file *os.File
	defer file.Close()
	for _, item := range files {
		file.Close()
		if item.IsDir() {
			log.Println("Dir : ", item.Name())
			if item.Name() == ".git" {
				continue
			}
			getMiddleSize(dir + system.SystemSep() + item.Name())
			continue
		}
		if strings.HasSuffix(item.Name(), ".jpeg") ||
			strings.HasSuffix(item.Name(), ".jpg") ||
			strings.HasSuffix(item.Name(), ".png") {

			var e1 error
			file, e1 = os.Open(dir + system.SystemSep() + item.Name())
			if e1 != nil {
				continue
			}
			var cf image.Config
			var e2 error
			if strings.HasSuffix(item.Name(), ".png") {
				cf, e2 = png.DecodeConfig(file)
			} else {
				cf, e2 = jpeg.DecodeConfig(file)
			}

			if e2 != nil {
				continue
			}
			_W = append(_W, cf.Width)
			_H = append(_H, cf.Height)
			fileInfo, err := os.Stat(dir + system.SystemSep() + item.Name())
			if err != nil {
				panic(err)
			} else {
				_Size = append(_Size, int(fileInfo.Size()/int64(1024)))
			}
		}
	}
}

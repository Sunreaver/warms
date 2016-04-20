package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sunreaver/gotools/system"
)

var (
	maxw int
	minw int
	maxh int
	minh int
	maxs int
	mins int
)

func main() {
	flag.IntVar(&maxw, "aw", 900, "最大的宽度")
	flag.IntVar(&minw, "iw", 400, "最小的宽度")
	flag.IntVar(&maxh, "ah", 3000, "最大的高度")
	flag.IntVar(&minh, "ih", 600, "最小的高度")
	flag.IntVar(&maxs, "as", 5, "最大的size（单位MB）")
	flag.IntVar(&mins, "is", 50, "最小的size（单位KB）")

	flag.Parse()
	log.Println(flag.Args())
	if len(flag.Args()) == 0 {
		removeFileWithDir(system.CurPath())
	} else {
		for i := 0; i < len(flag.Args()); i++ {
			removeFileWithDir(os.Args[i])
		}
	}
}

func removeFileWithDir(dir string) {
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
			removeFileWithDir(dir + system.SystemSep() + item.Name())
			continue
		}
		if strings.HasSuffix(item.Name(), ".jpeg") ||
			strings.HasSuffix(item.Name(), ".jpg") ||
			strings.HasSuffix(item.Name(), ".png") {

			var e1 error
			file, e1 = os.Open(dir + system.SystemSep() + item.Name())
			if e1 != nil {
				removeFile(dir + system.SystemSep() + item.Name())
				continue
			}
			var cf image.Config
			var e2 error
			if strings.HasSuffix(item.Name(), ".png") {
				cf, e2 = png.DecodeConfig(file)
			} else {
				cf, e2 = jpeg.DecodeConfig(file)
			}

			fileInfo, e3 := os.Stat(dir + system.SystemSep() + item.Name())
			if e2 != nil ||
				cf.Height < minh || cf.Width < minw ||
				cf.Height > maxh || cf.Width > maxw {
				removeFile(dir + system.SystemSep() + item.Name())
			} else if e3 != nil || fileInfo.Size()/1024/1024 > int64(maxs) || fileInfo.Size()/1024 < int64(mins) {
				removeFile(dir + system.SystemSep() + item.Name())
			}
		} else if strings.HasSuffix(item.Name(), ".gif") {
			removeFile(dir + system.SystemSep() + item.Name())
		}
	}
}

func removeFile(file string) {
	if e := os.Remove(file); e != nil {
		log.Println("remove error : ", e)
	} else {
		log.Println("removed : ", file)
	}
}

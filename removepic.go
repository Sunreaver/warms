package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sunreaver/gotools/system"
)

func main() {

	imgPath := system.CurPath()
	removeFileWithDir(imgPath)
}

func removeFileWithDir(dir string) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		log.Println(e)
		panic(e)
	}

	var file *os.File
	for _, item := range files {
		file.Close()
		if item.IsDir() {
			log.Println("Dir : ", item.Name())
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
			if e2 != nil {
				removeFile(dir + system.SystemSep() + item.Name())
				continue
			}
			if cf.Height < 700 || cf.Width < 600 {
				removeFile(dir + system.SystemSep() + item.Name())
			}
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

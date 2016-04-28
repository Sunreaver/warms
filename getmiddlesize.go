package main

import (
	"crypto/md5"
	"fmt"
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

	// md5文件
	md5file := map[string][]string{}
	if len(os.Args) < 2 {
		getMD5(system.CurPath(), md5file)
	} else {
		for i := 1; i < len(os.Args); i++ {
			getMD5(os.Args[i], md5file)
		}
	}

	for key, item := range md5file {
		if len(item) == 1 {
			delete(md5file, key)
		}
	}

	log.Println("The Same Pic Count: ", len(md5file))
	for _, item := range md5file {
		log.Println(item)
		log.Print("\r\n")

		for i := 0; i < len(item)-1; i++ {
			os.Remove(item[i])
		}
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

func getMD5(dir string, outMD5OfFile map[string][]string) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		log.Println(e)
		return
	}

	var f *os.File
	defer f.Close()
	for _, item := range files {
		f.Close()
		if item.IsDir() {
			log.Println("Dir : ", item.Name())
			if item.Name() == ".git" {
				continue
			}
			getMD5(dir+system.SystemSep()+item.Name(), outMD5OfFile)
			continue
		}
		if strings.HasSuffix(item.Name(), ".jpeg") ||
			strings.HasSuffix(item.Name(), ".jpg") ||
			strings.HasSuffix(item.Name(), ".png") {

			var e1 error
			f, e1 = os.Open(dir + system.SystemSep() + item.Name())
			if e1 != nil {
				continue
			}

			var b []byte
			b, e1 = ioutil.ReadAll(f)
			if e1 != nil {
				continue
			}
			key := fmt.Sprintf("%x", md5.Sum(b))
			if _, ok := outMD5OfFile[key]; ok {
				outMD5OfFile[key] = append(outMD5OfFile[key], dir+system.SystemSep()+item.Name())

			} else {
				outMD5OfFile[key] = []string{dir + system.SystemSep() + item.Name()}
			}
		}
	}
}

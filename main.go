package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/nfnt/resize"
)

const dirPath string = "./picture"

func main() {
	// open "test.jpg"
	var smallImagePath string
	if runtime.GOOS == "windows" {
		smallImagePath = dirPath + `\\small`
	} else {
		smallImagePath = dirPath + "/small"
	}
	err := os.MkdirAll(smallImagePath, os.ModePerm)
	if err != nil {
		fmt.Printf("创建小图文件夹失败,err:%v\n", err)
		return
	}
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("read dir failed,err:%v", err)
		return
	}
	for i := 0; i < len(dirs); i++ {
		if dirs[i].IsDir() {
			fmt.Printf("skip dir:%v\n", dirs[i].Name())
			continue
		}
		info, err := dirs[i].Info()
		if err != nil {
			fmt.Printf("get file info failed,err:%v\n", err)
			continue
		}
		name := dirs[i].Name()
		if !strings.Contains(name, "/") {
			name = dirPath + "/" + name
		}
		img, err := readImg(name)
		if err != nil {
			continue
		}
		m := beSmall(img)
		err = saveImg(smallImagePath, getFileName(info.Name()), m)
		if err != nil {
			fmt.Printf("save image failed,err:%v\n", err)
		}
	}
}

func beSmall(img image.Image) image.Image {
	var pcr float64
	var maxSize = img.Bounds().Dx()
	if img.Bounds().Dy() > maxSize {
		maxSize = img.Bounds().Dy()
	}
	if maxSize <= 100 {
		return img
	}
	pcr = float64(maxSize) / 1600
	outputX := float64(img.Bounds().Dx()) / pcr
	outputY := float64(img.Bounds().Dy()) / pcr
	//fmt.Printf("Dx:%v,Dy:%v\n", img.Bounds().Dx(), img.Bounds().Dy())
	//fmt.Printf("pcr:%v\n", pcr)
	//fmt.Printf("output x:%v ,y:%v\n", outputX, outputY)
	m := resize.Resize(uint(outputX), uint(outputY), img, resize.Lanczos3)
	return m
}

func readImg(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("open file(path:%v) failed,err:%v\n", path, err)
		return nil, err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		//fmt.Printf("decode jpg failed,err:%v\n", err)
		return nil, err
	}
	err = file.Close()
	if err != nil {
		fmt.Printf("close file failed,err:%v\n", err)
		return nil, err
	}
	return img, nil
}

func saveImg(basePath string, name string, m image.Image) error {
	sp := "/"
	if runtime.GOOS == "windows" {
		sp = `\\`
	}
	savePwd := basePath + sp + name
	f, err := os.Create(savePwd)
	if err != nil {
		return err
	}
	defer f.Close()

	// write new image to file
	err = jpeg.Encode(f, m, nil)
	if err != nil {
		return err
	}
	fmt.Printf("写入数据成功,文件:%v\n", basePath)
	return nil
}

func getFileName(str string) string {
	_, name := path.Split(str)
	return name
}

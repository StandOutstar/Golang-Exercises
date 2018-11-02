package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	help bool
	imgOutput string
	imgPath string
)

// usage flag usage replacement
func usage() {
	usageString := `
cut 9 image 

version: 0.0.1
Usage: command [-h] [-f file_path] [-o output_name]
        
Options:
`
	fmt.Fprintf(os.Stderr, usageString)
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&help, "h", false, "show help message")
	flag.StringVar(&imgOutput, "o", "output", "output dir name")
	flag.StringVar(&imgPath, "f", "", "target img file path")
	flag.Usage = usage
}

// fillImage fillImage 填充图片为正方形
func fillImage(img image.Image) image.Image {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	log.Printf("width: %v height: %v", width, height)

	lengthf := math.Max(float64(width), float64(height))
	length := int(lengthf)
	dst := imaging.New(length, length, color.NRGBA{0, 0, 0, 0})
	dst = imaging.PasteCenter(dst, img)

	return dst
}

// cutImage cutImage 裁剪图片为9个
func cutImage(img image.Image) []image.Image {
	width := img.Bounds().Dx()
	itemWidth := int(width/3)
	boxList := make([]image.Rectangle, 0)

	for i := 0; i < 3; i++ {
	    for j := 0; j < 3; j++ {
	        box := image.Rectangle{Min:image.Point{X:j*itemWidth, Y:i*itemWidth}, Max:image.Point{X:(j+1)*itemWidth, Y:(i+1)*itemWidth}}
			boxList = append(boxList, box)
	    }
	}
	log.Println("boxs:", len(boxList))
	log.Println("boxs:", boxList)
	imgList := make([]image.Image, 0)
	for i := 0; i < 9; i++ {
	    imgList = append(imgList, imaging.Crop(img, boxList[i]))
	}
	return imgList
}


// saveImage saveImage 将裁剪后的9个图像进行保存
func saveImage(imgList []image.Image) {
	for index, img := range imgList {
		tempDirPath := GetCurrentDirectory() + "/" + imgOutput + "/"
		exists, err := PathExists(tempDirPath)
		if err != nil {
			log.Fatalln("check path error:", err)
		}
		if !exists {
			err = os.Mkdir(tempDirPath, 0700)
			if err != nil {
				log.Fatalln("mkdir error:", err)
			}
		}
		tempPath := tempDirPath + strconv.Itoa(index) + ".png"
		log.Println(tempPath)
		err = imaging.Save(img, tempPath)
		if err != nil {
			log.Println("save tempPath error:", err)
		}
	}
}

//GetCurrentDirectory GetCurrentDirectory 获取当前路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// PathExists PathExists 判断指定路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	flag.Parse()
	if help {    // 设置 -h
		flag.Usage()
		os.Exit(0)
	}

	if imgPath == "" {
		log.Fatalln("not have a image file, -f imageFilePath")
	}

	src, err := imaging.Open(imgPath)
	if err != nil {
		log.Fatalln("open image error", err)
	}

	img := fillImage(src)
	imgs := cutImage(img)
	saveImage(imgs)
}

package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"log"
	"os"
)

var (
	help bool
	imgOutput string
	imgWidth int
	imgHeight int
	imgPath string
	asciiChar = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "
)

func init() {
	flag.BoolVar(&help, "h", false, "show help message")
	flag.StringVar(&imgOutput, "o", "output.txt", "output file name")
	flag.IntVar(&imgWidth, "width", 80, "image width") // if 0 all layers
	flag.IntVar(&imgHeight, "height", 80, "image height") // if 0 all layers
	flag.StringVar(&imgPath, "f", "", "target img file path")
	flag.Usage = usage
}

// usage flag usage replacement
func usage() {
	usageString := `
convert image to ascii 

version: 0.0.1
Usage: ascii [-h] [-f file_path] [-o output_name] [-width resize width] [-height resize height]
        
Options:
`
	fmt.Fprintf(os.Stderr, usageString)
	flag.PrintDefaults()
}


// getChar getChar将256灰度转换到70个字符上
func getChar(r, g, b, alpha uint8) (c string) {
	if alpha == 0 {
		return " "
	}
	length := len(asciiChar)
	gray := int(0.299 * float32(r) + 0.587 * float32(g) + 0.114 * float32(b))
	uintT := 256.0/length + 1
	fmt.Println("length:", length,"gray:", gray, "uintT:", uintT, "index:", int(gray/uintT))
	char := asciiChar[int(gray/uintT)]
	return string(char)
}

func main() {
	flag.Parse() // 解析flag
	if help {    // 设置 -h
		flag.Usage()
		os.Exit(0)
	}

	// Open a image.
	src, err := imaging.Open(imgPath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Resize the cropped image preserving the aspect ratio.
	dstSrc := imaging.Resize(src, imgWidth, imgHeight, imaging.Lanczos)

	txt := ""
	for i := 0; i < imgHeight; i++ {
		for j := 0; j < imgWidth; j++ {

		    txt += getChar(dstSrc.NRGBAAt(j, i).R, dstSrc.NRGBAAt(j, i).G, dstSrc.NRGBAAt(j, i).B, dstSrc.NRGBAAt(j, i).A)
		}
		txt += "\n"
	}

	print(txt)

	// 字符画输出到文件
	outputFile, err := os.OpenFile(imgOutput, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	outputFile.Write([]byte(txt))
}

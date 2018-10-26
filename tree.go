package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	help            bool   // -h : show help message
	displayDatetime bool   // -D : show change datetime
	maxLayer        int    // -L : 显示的层级数
	targetPath      string // -p : 指定目录
	rootPath        string
)

// init
func init() {
	flag.BoolVar(&help, "h", false, "show help message")
	flag.BoolVar(&displayDatetime, "D", false, "show change datetime")
	flag.IntVar(&maxLayer, "L", 0, "max layers to show") // if 0 all layers
	flag.StringVar(&targetPath, "p", "", "show taget path")
	flag.Usage = usage
	rootPath = GetCurrentDirectory()
}

// flag usage replacement
func usage() {
	fmt.Fprintf(os.Stderr, "Current Path: %s\n", rootPath)

	usageString := `
tree version: tree/0.0.0
Usage: tree [-hDL]
        
Options:
`
	fmt.Fprintf(os.Stderr, usageString)
	flag.PrintDefaults()
}

//GetCurrentDirectory 获取当前路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 遍历打印文件
func printPath(dir string, cn int) error {
	var prefix = " |__"
	if dir == rootPath {
		fmt.Println("./")
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, info := range infos {

		fmt.Print(strings.Repeat(" |  ", cn)) // 打印占位字符

		if displayDatetime { // 通过 displayDatetime 控制显示修改时间
			fmt.Println(prefix, info.Name(), info.ModTime().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println(prefix, info.Name())
		}

		// is dir
		if info.IsDir() {
			if (cn+1 < maxLayer) || maxLayer == 0 { // 通过 maxLayer 控制显示的层级
				err := printPath(dir+"/"+info.Name(), cn+1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	flag.Parse() // 解析flag
	if help {    // 设置 -h
		flag.Usage()
		os.Exit(0)
	}

	if targetPath != "" { // 设置 -p
		rootPath = targetPath
	}

	err := printPath(rootPath, 0)
	if err != nil {
		fmt.Println(err)
	}
}

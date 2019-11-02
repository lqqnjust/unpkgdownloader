package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	// "container/list"
)

// 判断文件是否存在
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 下载页面内容
func getHtml(url string, downloadDirTrue string){
	resp, err := http.Get(url)
    if err != nil {
        fmt.Println("http get error", err)
        return
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("read error", err)
        return
	}
	reg := regexp.MustCompile("<table(.*?)</table>")
	data := reg.Find(body)
	reg2 := regexp.MustCompile("href=\"(.*?)\"")
	hrefs := reg2.FindAllStringSubmatch(string(data), -1)
	
	for _, href := range hrefs {
		
		link_name := href[1]
		if link_name == "../" {
			continue
		}
		newUrl := url + link_name
		fmt.Println(newUrl)
	
		if (strings.HasSuffix(link_name, "/")) {
			newDir := filepath.Join(downloadDirTrue, link_name)
			if !IsFileExist(newDir){
		
				err := os.Mkdir(newDir, os.ModeDir)
				if err != nil {
					fmt.Println("err")
				}
			}
			
			getHtml(newUrl, newDir)

		}else {
			
			res, err := http.Get(newUrl)
			if err != nil {
				panic(err)
			}
			f, err := os.Create(filepath.Join(downloadDirTrue, link_name))
			if err != nil {
				panic(err)
			}
			io.Copy(f, res.Body)
		}
		
    }
}

func main(){
	pkgname := flag.String("n", "", "Package name to download")
	version := flag.String("v", "", "Package version to download")
	downloadDir := flag.String("d", "", "Package save directory")
	
	flag.Parse()
	fmt.Println("pkgname:", *pkgname)
	fmt.Println("version:", *version)
	fmt.Println("downloadDir:", *downloadDir)
	downloadDirTrue := filepath.Join(*downloadDir, *pkgname+"@"+ *version)
	fmt.Println("save dir:", downloadDirTrue)
	if !IsFileExist(downloadDirTrue){
		
		err := os.Mkdir(downloadDirTrue, os.ModeDir)
		if err != nil {
			fmt.Println("err")
		}
	}

	fmt.Println("Starting")
	url := "https://unpkg.com/"+ *pkgname + "@" + *version + "/"
	fmt.Println("url:", url)

	getHtml(url, downloadDirTrue)

	
}
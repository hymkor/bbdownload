package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

func download(url string) ([]byte, error) {
	log.Printf("GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(resp.Status)
	return bin, err
}

var rxAnchor = regexp.MustCompile(`(?si)<a[^>]+href="([^"]+)"[^>]*>(.*?)</a>`)

func mains(args []string) error {
	for _, arg1 := range args {
		baseUrl, err := url.Parse(arg1)
		if err != nil {
			return err
		}

		baseHtmlBin, err := download(arg1)
		if err != nil {
			return err
		}

		baseHtml := string(baseHtmlBin)

		m := rxAnchor.FindAllStringSubmatch(baseHtml, -1)
		downloadCount := 0
		if m != nil {
			for _, m1 := range m {
				subUrl, err := baseUrl.Parse(m1[1])
				if err != nil {
					return err
				}

				subUrlString := subUrl.String()
				switch path.Ext(subUrl.Path) {
				case ".zip", ".bz2", ".gz":
					arcbin, err := download(subUrlString)
					if err != nil {
						return err
					}

					fname := path.Base(subUrl.Path)
					ioutil.WriteFile(fname, arcbin, 0666)
					fmt.Println("done.")
					downloadCount++
				}
			}
		}
		if downloadCount <= 0 {
			log.Println(baseHtml)
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

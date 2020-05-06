package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
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
					fmt.Printf("%s ...", subUrlString)

					arcbin, err := download(subUrlString)
					if err != nil {
						return err
					}

					fname := path.Base(subUrl.Path)
					ioutil.WriteFile(fname, arcbin, 0666)
					fmt.Println("done.")
					downloadCount++
					break
				default:
					fmt.Printf("ignore %s\n", subUrlString)
					break
				}
			}
		}
		if downloadCount <= 0 {
			fmt.Println(baseHtml)
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

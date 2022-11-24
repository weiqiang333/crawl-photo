package wallpaperscraft_com

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Wallpaperscraft struct {
	UrlMain   string
	Url       string
	FileName  string
	File      *os.File
	PhotoUrls []string
	Doc       *goquery.Document
}

func NewWallpaperscraft(url string, file string) *Wallpaperscraft {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Failed open file err: ", err.Error())
		os.Exit(2)
	}
	return &Wallpaperscraft{
		UrlMain:   url,
		Url:       url,
		FileName:  file,
		File:      f,
		PhotoUrls: []string{},
	}
}

// CrawlAllPhoto 获取所有图片
func (w *Wallpaperscraft) CrawlAllPhoto() {
	for {
		if err := w.GetHtml(); err != nil {
			time.Sleep(time.Second * 2)
			return
		}
		w.getPhoto()
		fmt.Println(w.PhotoUrls)
		w.fileWrite()
		w.getNextPageUrl()
		if w.Url == "" {
			log.Println("W 不存在下一页, 它很可能是最后一页.")
			break
		}
	}
}

func (w *Wallpaperscraft) GetHtml() error {
	res, err := http.Get(w.Url)
	if err != nil {
		log.Println("GetHtml err: ", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("GetHtml status code error: ", res.StatusCode, res.Status)
		time.Sleep(time.Second * 10)
		return fmt.Errorf("GetHtml status code error: %v %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("GetHtml NewDocumentFromReader error: ", err.Error())
		return err
	}
	w.Doc = doc
	return nil
}

func (w *Wallpaperscraft) getNextPageUrl() {
	w.Url = ""
	w.Doc.Find("a.pager__link").Each(func(i int, s *goquery.Selection) {
		href, IsExist := s.Attr("href")
		if IsExist && s.Text() == "→" {
			w.Url = w.UrlMain + href
			log.Println(s.Text(), w.Url)
			return
		}
	})
}

func (w *Wallpaperscraft) getPhoto() {
	w.Doc.Find("li.wallpapers__item a.wallpapers__link").Each(func(i int, s *goquery.Selection) {
		href, IsExist := s.Attr("href")
		if IsExist {
			err, url := getPhotoUrl(w.UrlMain + href)
			if err != nil {
				log.Println("Failed get page photo url err: ", err.Error())
			} else {
				w.PhotoUrls = append(w.PhotoUrls, url)
			}
		}
	})
}

func getPhotoUrl(url string) (error, string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("GetHtml err: ", err)
		return err, ""
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("getPhotoUrl status code error: ", res.StatusCode, res.Status)
		time.Sleep(time.Second * 10)
		return fmt.Errorf("getPhotoUrl status code error: %v %s", res.StatusCode, res.Status), ""
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	url, isExist := doc.Find("img.wallpaper__image").Attr("src")
	if isExist {
		return nil, url
	}
	return fmt.Errorf("not url"), ""
}

func (w *Wallpaperscraft) fileWrite() {
	for _, url := range w.PhotoUrls {
		if url == "" {
			continue
		}
		_, err := w.File.WriteString(fmt.Sprintf("%s\n", url))
		if err != nil {
			log.Println("Failed file write err: ", err.Error())
			continue
		}
	}
	w.PhotoUrls = []string{}
}

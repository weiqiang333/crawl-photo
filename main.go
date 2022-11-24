package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	wallpaperscraft_com "crawl-photo/pkg/wallpaperscraft.com"
)

func init() {
	pflag.String("site_name", "", fmt.Sprintf("需要抓取的站点名称. 可选值: \n"+
		"- wallpaperscraft"))
}

func main() {
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("Failed in BindPFlags")
	}
	siteName := viper.GetString("site_name")

	if siteName == "wallpaperscraft" {
		w := wallpaperscraft_com.NewWallpaperscraft("https://wallpaperscraft.com", "./data/wallpaperscraft_photo_urls.text")
		w.CrawlAllPhoto()
	}
}

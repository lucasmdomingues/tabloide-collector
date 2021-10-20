package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"

	"github.com/gocolly/colly"
)

func init() {
	if _, err := os.Stat("storage"); err != nil {
		if err := os.Mkdir("storage", 0774); err != nil {
			log.Fatal(err)
		}
	} else {
		dir, err := ioutil.ReadDir("storage")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range dir {
			os.RemoveAll(path.Join([]string{"storage", file.Name()}...))
		}
	}
}

func main() {
	c := colly.NewCollector()

	c.OnHTML("a[title*='Jornal de Ofertas']", func(e *colly.HTMLElement) {
		c.OnHTML("img", func(e *colly.HTMLElement) {
			if err := downloadImage(e.Attr("src")); err != nil {
				log.Fatal("error on download image:", err)
			}
		})

		href := e.Attr("href")

		if err := c.Visit(href); err != nil {
			log.Fatal(fmt.Sprintf("error on visit link %s: %s", href, err))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal(fmt.Sprintf("request URL: %s failed with status code: %d Error: %s", r.Request.URL, r.StatusCode, err.Error()))
	})

	if err := c.Visit("http://federzonisupermercados.com.br/web/"); err != nil {
		log.Fatal("error on visit website:", err)
	}
}

func downloadImage(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("storage/%d.jpg", rand.Int())

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return err
	}

	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("image '%s' has storagered with success", fileName))

	return nil
}

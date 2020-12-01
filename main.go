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
		os.Mkdir("storage", 0774)
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
			err := downloadImage(e.Attr("src"))
			if err != nil {
				log.Fatal(err)
			}
		})

		err := c.Visit(e.Attr("href"))
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal(fmt.Sprintf("Request URL: %s failed with status code: %d Error: %s", r.Request.URL, r.StatusCode, err.Error()))
	})

	err := c.Visit("http://federzonisupermercados.com.br/web/")
	if err != nil {
		log.Fatal(err)
	}
}

func downloadImage(link string) error {
	fileName := fmt.Sprintf("storage/%d.jpg", rand.Int())

	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return err
	}

	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Image '%s' has storagered with success", fileName))

	return nil
}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

func main() {
	if _, err := os.Stat("tabloides"); err != nil {
		os.Mkdir("tabloides", 0774)
	} else {
		dir, err := ioutil.ReadDir("tabloides")
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, d := range dir {
			os.RemoveAll(path.Join([]string{"tabloides", d.Name()}...))
		}
	}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		link := e.Attr("src")

		err := downloadImage(link)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	err := c.Visit("http://www.federzonisupermercados.com.br/tabloide/tabloide2.html")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func downloadImage(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("tabloides/%s.png", uuid.New().String())

	defer resp.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Image %s downloaded", fileName))

	return nil
}

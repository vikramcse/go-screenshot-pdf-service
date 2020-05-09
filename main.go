// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	host := flag.String("host", "localhost", "the host address")
	port := flag.String("port", "8000", "the port number")
	flag.Parse()

	addr := *host + ":" + *port

	log.Println("server listening on address", addr)
	http.HandleFunc("/", pdfHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("got error while running server %v", err)
	}
}

func getPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().Do(ctx)
			if err != nil {
				log.Printf("something went wrong while in chromedp %v", err)
				return err
			}

			*res = buf
			return nil
		}),
	}
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	rurl := queryParams.Get("url")
	log.Println("request received for url ", rurl)

	if rurl == "" {
		fmt.Fprint(w, "URL can not be empty")
		return
	}

	_, err := url.ParseRequestURI(rurl)
	if err != nil {
		fmt.Fprintf(w, "invalid url format, eg. https://www.google.com")
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// to store data of pdf content
	var buf []byte

	if err := chromedp.Run(ctx, getPDF(rurl, &buf)); err != nil {
		log.Fatalf("error while processing request %v", err)
	}

	tempFile, err := ioutil.TempFile("", "upload-*.pdf")
	if err != nil {
		log.Fatalf("error while creating temp file %v", err)
	}
	log.Println("created temp file", tempFile.Name())
	defer os.Remove(tempFile.Name()) // clean up
	// defer tempFile.Close()

	if _, err := tempFile.Write(buf); err != nil {
		log.Fatalf("error while writing file to temp file %v", err)
	}

	fileStat, _ := tempFile.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)
	log.Println("the size of the file is", fileSize)

	name := strings.Split(tempFile.Name(), "/")[2]
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Type", "pdf")
	w.Header().Set("Content-Length", fileSize)

	tempFile.Seek(0, 0)
	io.Copy(w, tempFile)
	log.Println("Downloading pdf file ....")
}

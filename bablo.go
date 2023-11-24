package main

import (
    "flag"
    "fmt"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/gocolly/colly/v2"
)

func main() {
    customTextDone := make(chan bool)

    go displayCustomText(customTextDone)

    <-customTextDone

    threadCount := flag.Int("t", 1, "Number of threads")
    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Println("Usage: go run main.go -t <threadCount> <URL>")
        os.Exit(1)
    }

    startURL := flag.Arg(0)

    var visitedURLs sync.Map
    var wg sync.WaitGroup
    var totalURLs int
    urlQueue := make(chan string, 100000)

    maxURLLength := 200

    for i := 0; i < *threadCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c := colly.NewCollector()

            c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0"

            c.OnHTML("a[href]", func(e *colly.HTMLElement) {
                href := e.Attr("href")
                parsedURL, err := url.Parse(href)
                if err != nil {
                    return
                }
                absURL := e.Request.AbsoluteURL(parsedURL.String())

                // Check the URL length before processing
                if len(absURL) <= maxURLLength {
                    if strings.Contains(absURL, startURL) {
                        _, loaded := visitedURLs.LoadOrStore(absURL, true)
                        if !loaded {
                            urlQueue <- absURL
                        }
                    }
                }
            })

            for url := range urlQueue {
                fmt.Println(url)
                c.Visit(url)
                if totalURLs++; totalURLs >= 1000000 {
                    close(urlQueue)
                    return
                }
            }
        }()
    }

    urlQueue <- startURL

    wg.Wait()
}

func displayCustomText(customTextDone chan<- bool) {
    customText := "Made with   ˗ˋˏ ♡ ˎˊ˗  by GrozdniyAndy of XSS.is"

    for i := 0; i <= len(customText); i++ {
        fmt.Print("\r" + customText[:i] + "_")
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Print("\r" + customText)
    time.Sleep(1000 * time.Millisecond)

    customTextDone <- true

    fmt.Println("")
}

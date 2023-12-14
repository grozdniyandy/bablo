package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "net/http"
    "net/url"
    "strings"
    "sync"
    "time"

    "golang.org/x/net/html"
    "golang.org/x/net/proxy"
)

var baseDomain string
var visited = make(map[string]bool)
var mutex = &sync.Mutex{}
var urlCount int

func crawl(startUrl string, minDuration time.Duration, numWorkers int, headers []string, disableRedirects bool, outputFile string, proxyURL string) {
    parsedUrl, err := url.Parse(startUrl)
    if err != nil {
        fmt.Println("Error parsing base URL:", err)
        return
    }
    baseDomain = parsedUrl.Host

    queue := make(chan string, 100)
    var queueWg sync.WaitGroup
    var workersWg sync.WaitGroup
    var file io.WriteCloser
    var outputToFile bool
    if outputFile != "" {
        file, err = os.Create(outputFile)
        if err != nil {
            fmt.Println("Error creating output file:", err)
            return
        }
        defer file.Close()
        outputToFile = true
    }

    queueWg.Add(1)
    go func() {
        queue <- startUrl
        queueWg.Done()
    }()

    for i := 0; i < numWorkers; i++ {
        workersWg.Add(1)
        go func() {
            for uri := range queue {
                if !isVisited(uri) && isSameDomain(uri) {
                    visit(uri, queue, &queueWg, headers, disableRedirects, file, outputToFile, proxyURL)
                }
            }
            workersWg.Done()
        }()
    }

    done := make(chan bool)
    if outputToFile {
        ticker := time.NewTicker(5 * time.Second) // Update statistics every 5 seconds

        go func() {
            for {
                select {
                case <-ticker.C:
                    mutex.Lock()
                    fmt.Printf("\rTotal URLs Crawled: %d", urlCount)
                    mutex.Unlock()
                case <-done:
                    ticker.Stop()
                    return
                }
            }
        }()
    }

    go func() {
        time.Sleep(minDuration)
        done <- true
        close(done)
    }()

    go func() {
        select {
        case <-done:
            queueWg.Wait()
            close(queue)
        }
    }()

    workersWg.Wait()
    fmt.Printf("\rTotal URLs Crawled: %d\n", urlCount)
}

func visit(uri string, queue chan string, queueWg *sync.WaitGroup, headers []string, disableRedirects bool, file io.Writer, outputToFile bool, proxyURL string) {
    markVisited(uri)
    mutex.Lock()
    urlCount++
    if !outputToFile {
        fmt.Println("Crawled:", uri)
    }
    mutex.Unlock()

    if file != nil {
        _, err := file.Write([]byte(uri + "\n"))
        if err != nil {
            fmt.Println("Error writing to output file:", err)
        }
    }

    var client *http.Client
    if proxyURL != "" {
        proxyURI, err := url.Parse(proxyURL)
        if err != nil {
            fmt.Println("Error parsing proxy URL:", err)
            return
        }

        switch proxyURI.Scheme {
        case "http", "https":
            httpTransport := &http.Transport{Proxy: http.ProxyURL(proxyURI)}
            client = &http.Client{Transport: httpTransport}
        case "socks4", "socks5":
            dialer, err := proxy.FromURL(proxyURI, proxy.Direct)
            if err != nil {
                fmt.Println("Error setting up proxy dialer:", err)
                return
            }
            httpTransport := &http.Transport{}
            httpTransport.Dial = dialer.Dial
            client = &http.Client{Transport: httpTransport}
        default:
            fmt.Println("Unsupported proxy scheme:", proxyURI.Scheme)
            return
        }
    } else {
        client = &http.Client{}
    }

    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if disableRedirects {
            return http.ErrUseLastResponse
        }
        return nil
    }

    req, err := http.NewRequest("GET", uri, nil)
    if err != nil {
        fmt.Println("Error: Failed to create request for \"" + uri + "\"")
        return
    }

    for _, h := range headers {
        parts := strings.SplitN(h, ":", 2)
        if len(parts) == 2 {
            req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
        }
    }

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error: Failed to crawl \"" + uri + "\"")
        return
    }
    defer resp.Body.Close()

    links := getLinks(resp)
    for _, link := range links {
        absolute := fixUrl(link, uri)
        if absolute != "" && !isVisited(absolute) && isSameDomain(absolute) {
            queueWg.Add(1)
            go func(u string) {
                queue <- u
                queueWg.Done()
            }(absolute)
        }
    }
}

func getLinks(resp *http.Response) []string {
    var links []string
    z := html.NewTokenizer(resp.Body)
    for {
        tt := z.Next()
        switch {
        case tt == html.ErrorToken:
            return links
        case tt == html.StartTagToken, tt == html.EndTagToken:
            t := z.Token()
            if t.Data == "a" {
                for _, a := range t.Attr {
                    if a.Key == "href" {
                        links = append(links, a.Val)
                        break
                    }
                }
            }
        }
    }
}

func fixUrl(href, base string) string {
    uri, err := url.Parse(href)
    if err != nil {
        return ""
    }
    baseUrl, err := url.Parse(base)
    if err != nil {
        return ""
    }
    return baseUrl.ResolveReference(uri).String()
}

func isVisited(uri string) bool {
    mutex.Lock()
    defer mutex.Unlock()
    return visited[uri]
}

func markVisited(uri string) {
    mutex.Lock()
    defer mutex.Unlock()
    visited[uri] = true
}

func isSameDomain(uri string) bool {
    parsedUri, err := url.Parse(uri)
    if err != nil {
        return false
    }
    host := parsedUri.Host
    return host == baseDomain || strings.HasSuffix(host, "."+baseDomain)
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

func main() {
    customTextDone := make(chan bool)

    go displayCustomText(customTextDone)

    <-customTextDone
    threadsPtr := flag.Int("t", 10, "number of concurrent workers")
    disableRedirectsPtr := flag.Bool("dr", false, "disable following redirects")
    outputPtr := flag.String("o", "", "output file to write crawled URLs")
    proxyPtr := flag.String("p", "", "proxy URL (e.g., 'socks5://127.0.0.1:9050')")
    var headers []string
    flag.Func("H", "Add header (can be used multiple times)", func(s string) error {
        headers = append(headers, s)
        return nil
    })
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing URL argument")
        return
    }
    numWorkers := *threadsPtr
    disableRedirects := *disableRedirectsPtr
    outputFile := *outputPtr
    proxyURL := *proxyPtr
    minDuration := time.Minute
    crawl(args[0], minDuration, numWorkers, headers, disableRedirects, outputFile, proxyURL)
}

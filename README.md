# Bablo (Бабло)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Bablo (Бабло) is a web scraping tool designed to explore and collect URLs from a specified starting URL. It utilizes the "github.com/gocolly/colly/v2" package.

## Table of Contents
- [Features](https://github.com/grozdniyandy/bablo#features)
- [Usage](https://github.com/grozdniyandy/bablo#usage)
- [Installation](https://github.com/grozdniyandy/bablo#installation)
- [Dependencies](https://github.com/grozdniyandy/bablo#dependencies)
- [License](https://github.com/grozdniyandy/bablo#license)
- [Author](https://github.com/grozdniyandy/bablo#author)
- [Contributing](https://github.com/grozdniyandy/bablo#contributing)

## Features
- Recursive crawling
- Multithreading

## Usage
1. **Clone or Download:** Clone this repository or download the code to your local machine.
2. **Install Golang:** Install the latest Golang from https://go.dev/dl/. For example:
    ```
    wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    go version
    ```
3. **Install colly:** Run the following commands:
   ```
   go mod init bablo
   go get golang.org/x/net/html && go get golang.org/x/net/proxy
   ```
4. **Run the tool:** Run the tool using the following command:
   ```
   go run main.go -t <thread-count> https://example.com
   ```
   You can use Flags for additional features.
   ```
   Flags:
   -t int       Number of concurrent workers (default 10)
   -dr          Disable following redirects
   -o string    Output file to write crawled URLs
   -p string    Proxy URL (e.g., 'http://127.0.0.1:8080', 'socks5://127.0.0.1:9050')
   -H string    Add header (can be used multiple times)
   ```
5. **Automation:** For automating this tool with others check - https://github.com/grozdniyandy/hackerlines

## Installation
You can either check the "Usage" or download already compiled code from "releases".

## Dependencies
This code uses the Go standard library and one external dependency github.com/gocolly/colly/v2.

## License
This code is released under the [MIT License](LICENSE).

## Author
Bablo is developed by GrozdniyAndy of [XSS.is](https://xss.is).

## Contributing
Feel free to contribute, report issues, or suggest improvements by creating pull requests or issues in the GitHub repository. Enjoy using this crawler!

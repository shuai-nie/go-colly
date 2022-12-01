package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	dd := largestPalindrome(1)
	fmt.Println(dd)
	os.Exit(1)

	urlstr := "https://news.baidu.com"
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector()

	c.SetRequestTimeout(100 * time.Second)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"

	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", urlstr)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("title:", e.Text)
	})

	c.OnResponse(func(resp *colly.Response) {
		fmt.Println("response received", resp.StatusCode)

		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))

		if err != nil {
			log.Fatal(err)
		}

		htmlDoc.Find(".hotnews a").Each(func(i int, s *goquery.Selection) {
			band, _ := s.Attr("href")
			title := s.Text()
			fmt.Printf("热点新闻 %d: %s - $s\n", i, title, band)
			c.Visit(band)
		})
	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(urlstr)

}

var Client http.Client

type Result struct{}

func Spider(url string, i int) {
	reqSpider, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	reqSpider.Header.Set("", "")
	reqSpider.Header.Set("", "")
	reqSpider.Header.Set("", "")
	respSpider, err := Client.Do(reqSpider) //Client.Do(reqSpider)
	if err != nil {
		log.Fatal(err)
	}
	bodyTest, _ := ioutil.ReadAll(respSpider.Body)
	var result Result
	_ = json.Unmarshal(bodyTest, &result)
	fmt.Println(i, result)

}

// 最大回文数乘积
// 输入：n = 2 输出：987
// 解释：99 x 91 = 9009, 9009 % 1337 = 987
func largestPalindrome(n int) int {
	if n == 1 {
		return 9
	}

	key := int(math.Pow10(n)) - 1
	for num := key; num > 0; num-- {
		t := num
		for side := t; side > 0; side /= 10 {
			t = t*10 + side%10
		}

		for test := key; test*test > t; test-- {
			if t/test <= key && t%test == 0 {
				return t % 1337
			}
		}
	}
	return 0
}

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
	"unicode"
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

// 给定一个段落 (paragraph) 和一个禁用单词列表 (banned)。返回出现次数最多，同时不在禁用列表中的单词。
// https://blog.csdn.net/weixin_52001449/article/details/124233408
func mostCommonWord(paragraph string, banned []string) string {
	ban := map[string]bool{}
	for _, s := range banned {
		ban[s] = true
	}
	freq := map[string]int{}
	maxFreq := 0
	var word []byte
	for i, n := 0, len(paragraph); i <= n; i++ {
		if i < n && unicode.IsLetter(rune(paragraph[i])) {
			word = append(word, byte(unicode.ToLower(rune(paragraph[i]))))
		} else if word != nil {
			s := string(word)
			if !ban[s] {
				freq[s]++
				maxFreq = max(maxFreq, freq[s])
			}
			word = nil
		}
	}

	for s, f := range freq {
		if f == maxFreq {
			return s
		}
	}
	return ""
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// https://blog.csdn.net/weixin_52001449/article/details/124257115
// 给你一个整数 n ，按字典序返回范围 [1, n] 内所有整数。
// 输入：n = 13
// 输出：[1,10,11,12,13,2,3,4,5,6,7,8,9]
func lexicalOrder(n int) []int {
	var list = make([]int, n)
	last := 1
	for i := 0; i < n; i++ {
		list[i] = last
		if last*10 <= n {
			last = 10 * last
			continue
		} else {
			if last == n {
				last = last / 10
			}

			for last%10 == 9 {
				last = last / 10
			}
		}
		last++
	}
	return list
}

// 给你一个字符串 s 和一个字符 c ，且 c 是 s 中出现过的字符。
//返回一个整数数组 answer ，其中 answer.length == s.length 且 answer[i] 是 s 中从下标 i 到离它 最近 的字符 c 的 距离 。
//两个下标 i 和 j 之间的 距离 为 abs(i - j) ，其中 abs 是绝对值函数。
// https://blog.csdn.net/weixin_52001449/article/details/124285363
func shortestToChar(s string, c byte) []int {
	n := len(s)
	answer := make([]int, n)

	index := -n
	for i, ch := range s {
		if byte(ch) == c {
			index = i
		}
		answer[i] = i - index
	}

	index = 2 * n
	for i := n - 1; i >= 0; i-- {
		if s[i] == c {
			index = i
		}
		if answer[i] > index-i {
			answer[i] = index - i
		}
	}
	return answer
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SendAPost() {
	cli := &http.Client{}

	data := make(url.Values)
	data.Set("type", "text")
	data.Set("text", "再测试一条，打扰了")
	u, _ := url.Parse("https://api.pkuhollow.com/v3/send/post")

	params := url.Values{}
	params.Set("device", "0")
	params.Set("v", "v3.0.6-459040")
	u.RawQuery = params.Encode()

	req, _ := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	req.Header.Add("token", "TODOTODO")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	//req.Header.Add("origin", "https://web.pkuhollow.com")
	//req.Header.Add("referer", "https://web.pkuhollow.com/")
	//req.Header.Add("accept", "*/*")
	//req.Header.Add("accept-encoding", "gzip, deflate, br")
	//req.Header.Add("accept-language", "zh,en-GB;q=0.9,en-US;q=0.8,en;q=0.7")

	//req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\"")
	//req.Header.Add("sec-ch-ua-mobile", "?0")
	//req.Header.Add("sec-ch-ua-platform", "macOS")
	//req.Header.Add("sec-fetch-dest", "empty")
	//req.Header.Add("sec-fetch-mode", "cors")
	//req.Header.Add("sec-fetch-site", "same-site")
	//req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	resp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Fatalf("error status_code=%v, resp=%v", resp.StatusCode, string(body))
	}

	logger.Println(string(body))
}

func main() {
	SendAPost()
}

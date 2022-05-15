package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

const (
	ProtocolDomain = "https://pkuhelper.pku.edu.cn"
	GetlistUri     = "/services/pkuhole/api.php"
	TOKEN          = "TODOTODO"
)

var (
	cli *http.Client
)

func init() {
	cli = &http.Client{}
}

type GetListResp struct {
	Code int     `json:"code"`
	Data []*Post `json:"data"`
}

func GetLists(numOfPages int) ([]*Post, error) {
	var result []*Post
	for p := 1; p <= numOfPages; p++ {
		pageResult, err := GetList(p)
		if err != nil {
			return result, err
		}
		result = append(pageResult, result...)
	}

	// Remove dup
	var resultAfterDup []*Post
	dupMap := make(map[string]bool)
	for _, p := range result {
		if _, ok := dupMap[p.Pid]; !ok {
			resultAfterDup = append(resultAfterDup, p)
			dupMap[p.Pid] = true
		}
	}

	return resultAfterDup, nil
}

func GetList(page int) ([]*Post, error) {
	params := url.Values{}
	GetListURL, _ := url.Parse(ProtocolDomain + GetlistUri)
	params.Set("action", "getlist")
	params.Set("PKUHelperAPI", "3.0")
	params.Set("user_token", TOKEN)
	params.Set("p", strconv.Itoa(page))
	GetListURL.RawQuery = params.Encode()
	GetListURLRaw := GetListURL.String()

	req, err := http.NewRequest("GET", GetListURLRaw, nil)
	if err != nil {
		return []*Post{}, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return []*Post{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*Post{}, err
	}

	if resp.StatusCode != 200 {
		logger.Printf("error status_code=%v, resp=%v\n", resp.StatusCode, string(body))
		return []*Post{}, errors.New("error status code")
	}
	var getListResp GetListResp
	if err := json.Unmarshal(body, &getListResp); err != nil {
		panic(err)
	}

	sort.Sort(PostList(getListResp.Data))
	logger.Printf("Received posts from pid=%v, len=%v\n", getListResp.Data[0].Pid, len(getListResp.Data))

	return getListResp.Data, nil
}

type GetCommentResp struct {
	Data []*Comment `json:"data"`
}

func GetComment(pid string) ([]*Comment, error) {
	params := url.Values{}
	GetCommentURL, _ := url.Parse(ProtocolDomain + GetlistUri)
	params.Set("action", "getcomment")
	params.Set("PKUHelperAPI", "3.0")
	params.Set("user_token", TOKEN)
	params.Set("pid", pid)
	GetCommentURL.RawQuery = params.Encode()
	GetCommentURLRaw := GetCommentURL.String()

	req, err := http.NewRequest("GET", GetCommentURLRaw, nil)
	if err != nil {
		return []*Comment{}, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return []*Comment{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*Comment{}, err
	}

	if resp.StatusCode != 200 {
		logger.Printf("error status_code=%v, resp=%v\n", resp.StatusCode, string(body))
		return []*Comment{}, errors.New("error status code in get comment")
	}

	var getCommentResp GetCommentResp
	if err := json.Unmarshal(body, &getCommentResp); err != nil {
		panic(err)
	}

	logger.Printf("Received comments of pid=%v, num of comments=%v\n", pid, len(getCommentResp.Data))
	return getCommentResp.Data, nil
}

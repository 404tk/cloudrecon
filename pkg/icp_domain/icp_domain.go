package icp_domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/404tk/cloudrecon/common"
)

var IcpNotForRecord = errors.New("域名未备案")
var timeout = 15 * time.Second

type Icp struct {
	token  string
	cookie string
	ip     string
}

func (i *Icp) Query(input string) ([]*DomainInfo, error) {
	i.ip = common.RandIp()
	i.cookie = getCookie()
	if err := i.auth(); err != nil {
		return nil, err
	}
	resp, err := i.query(input, "20")
	if err != nil {
		return nil, err
	}
	if resp.LastPage > 1 {
		resp, err = i.query(input, strconv.Itoa(resp.Total))
		if err != nil {
			return nil, err
		}
	}

	return resp.List, nil
}

func (i *Icp) query(input, pageSize string) (*QueryParams, error) {
	queryRequest, _ := json.Marshal(&QueryRequest{
		PageSize: pageSize,
		UnitName: input,
	})

	result := &IcpResponse{Params: &QueryParams{}}
	err := i.post("icpAbbreviateInfo/queryByCondition", bytes.NewReader(queryRequest), "application/json;charset=UTF-8", i.token, result)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("查询：%s", result.Msg)
	}

	queryParams := result.Params.(*QueryParams)
	if len(queryParams.List) == 0 {
		return nil, IcpNotForRecord
	}

	return queryParams, nil
}

func getCookie() string {
	req, err := http.NewRequest(http.MethodGet, "https://beian.miit.gov.cn/", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.42")

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "__jsluid_s" {
			return cookie.Value
		}
	}
	return ""
}

func (i *Icp) auth() error {
	timestamp := time.Now().Unix()
	authKey := common.Md5(fmt.Sprintf("testtest%d", timestamp))
	authBody := fmt.Sprintf("authKey=%s&timeStamp=%d", authKey, timestamp)

	result := &IcpResponse{Params: &AuthParams{}}
	err := i.post("auth", bytes.NewReader([]byte(authBody)), "application/x-www-form-urlencoded;charset=UTF-8", "0", result)
	if err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("获取token失败：%s", result.Msg)
	}

	authParams := result.Params.(*AuthParams)
	i.token = authParams.Bussiness

	return nil
}

func (i *Icp) post(url string, data io.Reader, contentType, token string, result interface{}) error {
	postUrl := fmt.Sprintf("https://hlwicpfwc.miit.gov.cn/icpproject_query/api/%s", url)
	queryReq, err := http.NewRequest(http.MethodPost, postUrl, data)
	if err != nil {
		return err
	}

	queryReq.Header.Set("Content-Type", contentType)
	queryReq.Header.Set("Cookie", "__jsluid_s="+i.cookie)
	queryReq.Header.Set("Origin", "https://beian.miit.gov.cn/")
	queryReq.Header.Set("Referer", "https://beian.miit.gov.cn/")
	queryReq.Header.Set("token", token)
	queryReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.42")
	queryReq.Header.Set("CLIENT_IP", i.ip)
	queryReq.Header.Set("X-FORWARDED-FOR", i.ip)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(queryReq)
	return GetHTTPResponse(resp, postUrl, err, result)
}

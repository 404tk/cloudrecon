package asset_check

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/404tk/cloudrecon/common"
	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Wg      *sync.WaitGroup
	Ch      chan struct{}
	Bar     *progressbar.ProgressBar
	Results []string
}

func (cfg *Config) Check(url, method, region string, exp []int) {
	cfg.Wg.Add(1)
	defer cfg.Wg.Done()
	defer cfg.Bar.Add(1)
	client := newClient()
	req := requestBuilder(url, method)
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println(err)
		<-cfg.Ch
		return
	}
	if common.IntContain(exp, resp.StatusCode) {
		cfg.extractRealRegion(resp, url, region)
	}
	resp.Body.Close()
	<-cfg.Ch
}

func (cfg *Config) extractRealRegion(resp *http.Response, url, region string) {
	if resp.StatusCode == 301 && strings.Contains(resp.Request.Host, "amazonaws.com") {
		_region := resp.Header.Get("X-Amz-Bucket-Region")
		url = strings.Replace(url, region, _region, -1)
	} else if resp.StatusCode == 403 && strings.Contains(resp.Request.Host, "aliyuncs.com") {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			r, _ := regexp.Compile(`<Bucket>(?P<name>.*?)</Bucket>\s+<Endpoint>(?P<endpoint>.*?)</Endpoint>`)
			result := r.FindAllStringSubmatch(string(body), -1)
			if len(result) > 0 && len(r.SubexpNames()) > 2 {
				url = fmt.Sprintf("https://%s.%s", result[0][1], result[0][2])
			}
		}
	}
	cfg.Bar.Clear()
	fmt.Printf("[%d] %s\n", resp.StatusCode, url)
	cfg.Results = append(cfg.Results, url)
}

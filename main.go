package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/404tk/cloudrecon/common"
	"github.com/404tk/cloudrecon/pkg/asset_check"
	"github.com/404tk/cloudrecon/pkg/icp_domain"
	"github.com/schollz/progressbar/v3"
)

var (
	threads int
	domain  string
)

func init() {
	flag.IntVar(&threads, "t", 25, "threads")
	flag.StringVar(&domain, "d", "", "domain")
	flag.Parse()
}

func main() {
	start := time.Now()
	if domain == "" {
		fmt.Println("[-] 请指定域名，-d example.com")
		return
	}
	fmt.Println("[*] 当前输入域名：", domain)
	var words []string
	f, err := asset_check.ParseDomain(domain)
	if err != nil || !f.ICANN {
		fmt.Println("[-] 域名格式解析错误。")
		return
	}
	words = append(words, f.Domain)

	icp := &icp_domain.Icp{}
	domainInfo, err := icp.Query(domain)
	if err != nil {
		fmt.Println("[-] ", err)
	} else {
		fmt.Println("[+] ICP备案/许可证号：", domainInfo[0].ServiceLicence)
		fmt.Println("[+] 主办单位名称：", domainInfo[0].UnitName)
		records, err := icp.Query(domainInfo[0].MainLicence)
		if err != nil {
			fmt.Println("[-] 查询备案号关联域名出错：", err)
		}
		fmt.Printf("\n关联域名\n")
		fmt.Println("-------")
		for _, info := range records {
			if info.Domain != domain {
				fmt.Println(info.Domain)
			}
			f, err := asset_check.ParseDomain(info.Domain)
			// 判断域名格式
			if err != nil || !f.ICANN {
				continue
			}
			// 去重 & 过滤中文
			if common.StrContain(words, f.Domain) || common.UnicodeCheck(f.Domain) {
				continue
			}
			words = append(words, f.Domain)
		}
		fmt.Println()
	}

	// 加载yaml规则
	asset_check.InitDefaultPoc()
	// 解析规则payload
	payloads, num := asset_check.ParsePayloads(words)
	// 进度条展示
	bar := progressbar.NewOptions(num,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("[reset]Enumerating [yellow]%d [reset]urls ...", num)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	cfg := &asset_check.Config{
		Wg:  &sync.WaitGroup{},
		Ch:  make(chan struct{}, threads),
		Bar: bar,
	}
	// 开始枚举
	for _, d := range payloads {
		for url, region := range d.Urls {
			cfg.Ch <- struct{}{}
			go cfg.Check(url, d.Method, region, d.Exp)
		}
	}
	cfg.Wg.Wait()

	fmt.Printf("\n\n疑似关联的云上资产地址\n")
	fmt.Println("------------------")
	for _, url := range cfg.Results {
		fmt.Println(url)
	}
	fmt.Println()

	elapsed := time.Since(start)
	fmt.Println(fmt.Sprintf("[+] Done. Used %f seconds.", elapsed.Seconds()))
}

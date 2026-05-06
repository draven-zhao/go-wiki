package prometheus

import (
	"encoding/json"
	"fmt"
	"go-wiki/config"
	"go-wiki/services"
	"net"
	"net/http"
	"sync"
	"time"
)

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
}
type PromResp struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []Alert `json:"alerts"`
	} `json:"data"`
}

var (
	Lock sync.RWMutex
)

// 全局HTTP Client
var HttpClient = &http.Client{
	Timeout: time.Second * 10,
	Transport: &http.Transport{
		// 最大空闲连接
		MaxIdleConns: 100,
		// 每个host最大空闲连接
		MaxIdleConnsPerHost: 20,
		// 空闲连接保持时间
		IdleConnTimeout: time.Second * 90,
		// TLS握手超时
		TLSHandshakeTimeout: time.Second * 5,
		//ExpectContinue超时
		ExpectContinueTimeout: time.Second * 1,
		//连接建立超时
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		//开启连接复用
		DisableKeepAlives: true,
	},
}

// 后台监控
func StartScheduler() {
	ticker := time.NewTicker(300 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			alerts, err := FetchActiveAlerts()
			var results []map[string]interface{}
			collector := make(map[string]interface{})
			if err != nil {
				fmt.Println("抓取失败", err)
				time.Sleep(5 * time.Second)
				continue
			}
			Lock.Lock()
			for _, alert := range alerts {
				alertName, instance := BuildQuery(alert)
				fmt.Println("========", alert.Annotations["description"])
				//if alertName != "磁盘使用率过高" {
				//	continue
				//}
				//搜索wiki
				search, err := services.Search(alert.Annotations["description"])
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				//调用LLM整理
				organize, err := services.Organize(alert.Annotations["description"], search)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Println("-------", organize)
				// 汇总结果
				collector["告警名称"] = alertName
				collector["实例"] = instance
				collector["处理方法"] = organize
				results = append(results, collector)

				fmt.Println("当前告警", collector)
				//todo 将结果写入钉钉或者邮件
			}
			Lock.Unlock()
		}
	}()
}
func FetchActiveAlerts() ([]Alert, error) {
	c := config.Init()
	req, err := http.NewRequest("GET", c.Prometheus.Tg.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.SetBasicAuth(c.Prometheus.Tg.Name, c.Prometheus.Tg.Password)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	var result PromResp
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %v", err)
	}
	var actiive []Alert
	for _, a := range result.Data.Alerts {
		if a.State == "firing" || a.State == "pending" {
			actiive = append(actiive, a)
		}
	}
	return actiive, nil
}

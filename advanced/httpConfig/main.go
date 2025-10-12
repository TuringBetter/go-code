package main

import (
	"net"
	"net/http"
	"time"
)

func createProductionClient() *http.Client {
	transport := &http.Transport{
		// 连接池配置
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 25, // 并发请求数的1/3到1/2
		MaxConnsPerHost:     50, // 峰值并发数
		IdleConnTimeout:     90 * time.Second,

		// 超时配置
		// DialContext是一个方法
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // 建立TCP连接超时
			KeepAlive: 30 * time.Second, // TCP保活周期
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,  // TLS握手超时
		ResponseHeaderTimeout: 10 * time.Second, // 读取响应头超时
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 整体请求超时
	}
}

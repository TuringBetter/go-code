package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 数据模型
type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

//go:embed index.html
var html string

func main() {
	// 设置路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/stream/sse", sseHandler)
	http.HandleFunc("/stream/text", textStreamHandler)
	http.HandleFunc("/stream/json", jsonStreamHandler)

	// 启动服务器
	fmt.Println("流式输出服务器启动在 http://localhost:8080")
	fmt.Println("可用的端点:")
	fmt.Println("  - http://localhost:8080/ (主页)")
	fmt.Println("  - http://localhost:8080/stream/sse (SSE流式输出)")
	fmt.Println("  - http://localhost:8080/stream/text (文本流式输出)")
	fmt.Println("  - http://localhost:8080/stream/json (JSON流式输出)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 主页处理器
func indexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// SSE (Server-Sent Events) 流式输出处理器
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 2. 获取Flusher接口
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 3. 发送初始数据
	fmt.Fprintf(w, "data: %s\n\n", "SSE连接已建立")
	flusher.Flush()

	// 4. 模拟数据流
	for i := 1; i <= 10; i++ {
		message := Message{
			ID:      i,
			Content: fmt.Sprintf("这是第 %d 条SSE消息", i),
			Time:    time.Now().Format("15:04:05"),
		}

		// 将消息转换为JSON
		jsonData, _ := json.Marshal(message)

		// 发送SSE格式的数据
		fmt.Fprintf(w, "data: %s\n\n", string(jsonData))
		flusher.Flush()

		// 模拟处理延迟
		time.Sleep(1 * time.Second)
	}

	// 5. 发送结束信号
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// 文本流式输出处理器
func textStreamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 设置文本流响应头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 2. 获取Flusher接口
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 3. 模拟文本数据流
	textChunks := []string{
		"开始文本流式输出...\n",
		"这是第一行文本\n",
		"这是第二行文本\n",
		"正在处理数据...\n",
		"数据1: 处理完成\n",
		"数据2: 处理完成\n",
		"数据3: 处理完成\n",
		"所有数据处理完毕\n",
		"文本流输出结束\n",
	}

	for _, chunk := range textChunks {
		fmt.Fprint(w, chunk)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

// JSON流式输出处理器
func jsonStreamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 设置JSON流响应头
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 2. 获取Flusher接口
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 3. 开始JSON数组
	fmt.Fprint(w, "[\n")
	flusher.Flush()

	// 4. 模拟JSON数据流
	for i := 1; i <= 5; i++ {
		message := Message{
			ID:      i,
			Content: fmt.Sprintf("JSON流消息 %d", i),
			Time:    time.Now().Format("15:04:05"),
		}

		jsonData, _ := json.MarshalIndent(message, "  ", "  ")

		// 添加逗号（除了第一个元素）
		if i > 1 {
			fmt.Fprint(w, ",\n")
		}

		fmt.Fprintf(w, "  %s", string(jsonData))
		flusher.Flush()

		time.Sleep(1 * time.Second)
	}

	// 5. 结束JSON数组
	fmt.Fprint(w, "\n]")
	flusher.Flush()
}

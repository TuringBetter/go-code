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
	http.HandleFunc("/stream/pipeline", pipelineHandler) // 新增：通道解耦示例

	// 启动服务器
	fmt.Println("流式输出服务器启动在 http://localhost:8080")
	fmt.Println("可用的端点:")
	fmt.Println("  - http://localhost:8080/ (主页)")
	fmt.Println("  - http://localhost:8080/stream/sse (SSE流式输出)")
	fmt.Println("  - http://localhost:8080/stream/text (文本流式输出)")
	fmt.Println("  - http://localhost:8080/stream/json (JSON流式输出)")
	fmt.Println("  - http://localhost:8080/stream/pipeline (通道解耦示例)")

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

	ctx := r.Context()
	// 3. 发送初始数据
	fmt.Fprintf(w, "data: %s\n\n", "SSE连接已建立")
	flusher.Flush()

	// 4. 模拟数据流
	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			log.Printf("[SSE] ⚠️ 客户端断开连接（Context取消，在第 %d/10 条消息时）", i)
			return
		default:
		}
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

	ctx := r.Context()
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

	for i, chunk := range textChunks {
		select {
		case <-ctx.Done():
			log.Printf("[Text] ⚠️ 客户端断开连接（Context取消，在第 %d/%d 块时）", i+1, len(textChunks))
			return
		default:
		}
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

	ctx := r.Context()
	// 3. 开始JSON数组
	fmt.Fprint(w, "[\n")
	flusher.Flush()

	// 4. 模拟JSON数据流
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			log.Printf("[JSON] ⚠️ 客户端断开连接（Context取消，在第 %d/5 条消息时）", i)
			return
		default:
		}
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

// ============ 通道解耦：生产与传输分离示例 ============

// generateWithPipeline 模拟大模型逐token生成
// 💡 关键点：返回只读通道 (<-chan string)，调用者只能接收数据
func generateWithPipeline(prompt string) <-chan string {
	ch := make(chan string, 5) // 带缓冲的通道，生产者不会因为消费者慢而阻塞

	// 在独立的 goroutine 中生成数据（生产者）
	go func() {
		defer close(ch) // 确保生成完成后关闭通道
		log.Printf("[Pipeline-生产者] 开始生成，提示词: %s", prompt)

		// 模拟大模型逐token生成（如 OpenAI/Claude streaming API）
		tokens := []string{
			"你好", "！", "我", "是", "AI", "助手", "。\n",
			"根据", "你的", "提示", "「", prompt, "」", "，\n",
			"我", "将", "逐步", "生成", "回答", "内容", "。\n",
			"这", "展示", "了", "通道", "解耦", "的", "威力", "！",
		}

		for i, token := range tokens {
			// 模拟大模型API的延迟（生成延迟）
			time.Sleep(100 * time.Millisecond)

			// 发送到通道
			ch <- token
			log.Printf("[Pipeline-生产者] ✓ 生成token %d/%d: %q", i+1, len(tokens), token)
		}

		log.Printf("[Pipeline-生产者] ✓ 生成完成，通道已关闭")
	}()

	return ch // 立即返回通道，不等待生成完成
}

// pipelineHandler 演示通道解耦的流式输出处理器
func pipelineHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 2. 获取 Flusher 接口
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 3. 获取查询参数作为提示词
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		prompt = "通道解耦示例"
	}

	ctx := r.Context()
	log.Printf("[Pipeline-消费者] 客户端连接: %s, 提示词: %s", r.RemoteAddr, prompt)

	// 4. 启动生产者（立即返回通道）
	tokenCh := generateWithPipeline(prompt)

	fmt.Fprintf(w, "=== 通道解耦流式输出示例 ===\n")
	fmt.Fprintf(w, "提示词: %s\n", prompt)
	fmt.Fprintf(w, "开始接收生成的token...\n\n")
	flusher.Flush()

	// 5. 消费者：从通道读取并传输（传输过程）
	tokenCount := 0
	for {
		select {
		case <-ctx.Done():
			// 客户端断开连接
			log.Printf("[Pipeline-消费者] ⚠️ 客户端断开连接（已接收 %d 个token）", tokenCount)
			return

		case token, ok := <-tokenCh:
			if !ok {
				// 通道已关闭，生产者完成
				fmt.Fprintf(w, "\n\n=== 生成完成 ===\n")
				fmt.Fprintf(w, "共接收到 %d 个token\n", tokenCount)
				flusher.Flush()
				log.Printf("[Pipeline-消费者] ✓ 传输完成，共发送 %d 个token", tokenCount)
				return
			}

			// 发送token给客户端
			tokenCount++
			fmt.Fprint(w, token)
			flusher.Flush()
			log.Printf("[Pipeline-消费者] → 发送token %d: %q", tokenCount, token)

			// 模拟网络传输延迟（可选）
			// 注意：即使这里延迟，也不会阻塞生产者的生成
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// ============ 对比：无通道解耦的传统方式 ============
// 
// 传统方式的问题：
// func traditionalHandler(w http.ResponseWriter, r *http.Request) {
//     for i := 0; i < 10; i++ {
//         token := generateToken()        // 生成（阻塞）
//         fmt.Fprint(w, token)           // 传输（阻塞）
//         flusher.Flush()
//         // 问题：生成完一个才能传输一个，串行执行
//     }
// }
//
// 通道解耦的优势：
// 1. 生产者（大模型生成）在独立 goroutine 中运行，不被传输阻塞
// 2. 消费者（网络传输）在主 goroutine 中运行，不被生成阻塞
// 3. 通过带缓冲的通道，允许生产者"提前"生成多个token
// 4. 两者并行执行，总体响应更快

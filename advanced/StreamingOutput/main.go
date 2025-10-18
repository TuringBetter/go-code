package main

import (
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
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>流式输出Demo</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        button { padding: 10px 20px; margin: 10px; font-size: 16px; }
        #output { border: 1px solid #ccc; padding: 20px; height: 300px; overflow-y: auto; background: #f9f9f9; }
        .message { margin: 5px 0; padding: 5px; background: white; border-left: 3px solid #007cba; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Go 流式输出 Demo</h1>
        <p>选择一种流式输出类型进行测试：</p>
        
        <button onclick="testSSE()">测试 SSE (Server-Sent Events)</button>
        <button onclick="testTextStream()">测试文本流</button>
        <button onclick="testJSONStream()">测试 JSON 流</button>
        <button onclick="clearOutput()">清空输出</button>
        
        <h3>输出区域：</h3>
        <div id="output"></div>
    </div>

    <script>
        function clearOutput() {
            document.getElementById('output').innerHTML = '';
        }
        
        function addMessage(content) {
            const output = document.getElementById('output');
            const div = document.createElement('div');
            div.className = 'message';
            div.innerHTML = content;
            output.appendChild(div);
            output.scrollTop = output.scrollHeight;
        }
        
        function testSSE() {
            clearOutput();
            addMessage('开始 SSE 流式输出...');
            
            const eventSource = new EventSource('/stream/sse');
            
            eventSource.onmessage = function(event) {
                addMessage('SSE 消息: ' + event.data);
            };
            
            eventSource.onerror = function(event) {
                addMessage('SSE 连接错误');
                eventSource.close();
            };
            
            // 10秒后关闭连接
            setTimeout(() => {
                eventSource.close();
                addMessage('SSE 连接已关闭');
            }, 10000);
        }
        
        function testTextStream() {
            clearOutput();
            addMessage('开始文本流式输出...');
            
            fetch('/stream/text')
                .then(response => {
                    const reader = response.body.getReader();
                    const decoder = new TextDecoder();
                    
                    function readStream() {
                        return reader.read().then(({ done, value }) => {
                            if (done) {
                                addMessage('文本流结束');
                                return;
                            }
                            
                            const chunk = decoder.decode(value, { stream: true });
                            addMessage('文本块: ' + chunk);
                            
                            return readStream();
                        });
                    }
                    
                    return readStream();
                })
                .catch(error => {
                    addMessage('文本流错误: ' + error);
                });
        }
        
        function testJSONStream() {
            clearOutput();
            addMessage('开始 JSON 流式输出...');
            
            fetch('/stream/json')
                .then(response => {
                    const reader = response.body.getReader();
                    const decoder = new TextDecoder();
                    
                    function readStream() {
                        return reader.read().then(({ done, value }) => {
                            if (done) {
                                addMessage('JSON 流结束');
                                return;
                            }
                            
                            const chunk = decoder.decode(value, { stream: true });
                            addMessage('JSON 块: ' + chunk);
                            
                            return readStream();
                        });
                    }
                    
                    return readStream();
                })
                .catch(error => {
                    addMessage('JSON 流错误: ' + error);
                });
        }
    </script>
</body>
</html>`

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

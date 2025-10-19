package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // 导入pprof用于性能分析
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// 模拟一个消耗资源的任务
func leakyTask(ctx context.Context, id int) {
	// 分配一些内存（模拟资源占用）
	buffer := make([]byte, 1024*1024) // 1MB
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		// 清理资源
		buffer = nil
		fmt.Printf("[Task %d] 资源已清理\n", id)
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Task %d] 收到退出信号，正在清理...\n", id)
			return
		case <-ticker.C:
			// 模拟工作（访问buffer防止被优化掉）
			buffer[0] = byte(id)
		}
	}
}

// 不断创建goroutine但不取消（会泄漏）
func demoWithLeak() {
	fmt.Println("========== 泄漏模式：不断创建goroutine但不调用cancel ==========")
	fmt.Println("提示：使用 Ctrl+C 退出程序")
	fmt.Println()

	counter := 0
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			fmt.Println("\n收到退出信号，程序即将退出...")
			fmt.Println("注意：所有泄漏的goroutine将被强制终止")
			return

		case <-ticker.C:
			counter++
			// 创建Context但不调用cancel - 这会导致泄漏！
			ctx, cancel := context.WithCancel(context.Background())
			_ = cancel // 故意不调用，模拟泄漏

			go leakyTask(ctx, counter)

			// 打印统计信息
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("已启动 %d 个goroutine | 当前goroutine数: %d | 内存使用: %.2f MB\n",
				counter,
				runtime.NumGoroutine(),
				float64(m.Alloc)/1024/1024)
		}
	}
}

// 正确的做法：调用cancel清理资源
func demoWithoutLeak() {
	fmt.Println("========== 正常模式：正确调用cancel，资源会被清理 ==========")
	fmt.Println("提示：使用 Ctrl+C 退出程序")
	fmt.Println()

	counter := 0
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 存储cancel函数，以便后续清理
	cancels := make([]context.CancelFunc, 0)

	// 监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 定期清理旧的goroutine
	cleanupTicker := time.NewTicker(10 * time.Second)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n收到退出信号，正在优雅退出...")
			// 取消所有goroutine
			for i, cancel := range cancels {
				fmt.Printf("正在停止 Task %d...\n", i+1)
				cancel()
			}
			time.Sleep(500 * time.Millisecond) // 等待清理完成
			fmt.Println("所有资源已清理完毕")
			return

		case <-cleanupTicker.C:
			// 每10秒清理前面的goroutine
			if len(cancels) > 3 {
				fmt.Println("\n--- 开始清理旧的goroutine ---")
				for i := 0; i < 2 && i < len(cancels); i++ {
					cancels[i]()
				}
				cancels = cancels[2:]
				time.Sleep(200 * time.Millisecond)
				fmt.Println("--- 清理完成 ---\n")
			}

		case <-ticker.C:
			counter++
			ctx, cancel := context.WithCancel(context.Background())
			cancels = append(cancels, cancel) // 保存cancel函数以便后续调用

			go leakyTask(ctx, counter)

			// 打印统计信息
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("已启动 %d 个goroutine | 当前goroutine数: %d | 内存使用: %.2f MB | 活跃任务: %d\n",
				counter,
				runtime.NumGoroutine(),
				float64(m.Alloc)/1024/1024,
				len(cancels))
		}
	}
}

func main() {
	// 启动pprof服务器，用于性能分析
	go func() {
		fmt.Println("pprof服务已启动: http://localhost:6060/debug/pprof/")
		fmt.Println()
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			fmt.Printf("pprof服务启动失败: %v\n", err)
		}
	}()

	// 等待pprof服务启动
	time.Sleep(500 * time.Millisecond)

	if len(os.Args) > 1 && os.Args[1] == "leak" {
		demoWithLeak()
	} else if len(os.Args) > 1 && os.Args[1] == "normal" {
		demoWithoutLeak()
	} else {
		fmt.Println("用法: go run leak_demo.go [leak|normal]")
		fmt.Println()
		fmt.Println("  leak   - 演示资源泄漏（不调用cancel）")
		fmt.Println("  normal - 演示正常情况（调用cancel）")
		fmt.Println()
		fmt.Println("监控命令：")
		fmt.Println("  1. 实时查看goroutine数量:")
		fmt.Println("     watch -n 1 'curl -s http://localhost:6060/debug/pprof/goroutine?debug=1 | grep \"goroutine profile\"'")
		fmt.Println()
		fmt.Println("  2. 查看进程内存:")
		fmt.Println("     watch -n 1 'ps aux | grep leak_demo'")
		fmt.Println()
		fmt.Println("  3. 使用pprof分析:")
		fmt.Println("     go tool pprof http://localhost:6060/debug/pprof/goroutine")
		fmt.Println("     go tool pprof http://localhost:6060/debug/pprof/heap")
	}
}


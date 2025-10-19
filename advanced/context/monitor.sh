#!/bin/bash

# 监控脚本 - 用于观察goroutine和内存泄漏

echo "================================"
echo "Go Context 资源泄漏监控脚本"
echo "================================"
echo ""

# 检查程序是否运行
check_running() {
    pgrep -f "leak_demo" > /dev/null
    return $?
}

if ! check_running; then
    echo "错误: leak_demo 程序未运行"
    echo "请先运行: go run leak_demo.go leak"
    exit 1
fi

# 获取进程ID
PID=$(pgrep -f "leak_demo" | head -1)
echo "监控进程 PID: $PID"
echo ""

# 显示监控选项
echo "请选择监控方式:"
echo "1. 实时监控 goroutine 数量和内存（推荐）"
echo "2. 查看 goroutine 堆栈详情"
echo "3. 查看内存分配详情"
echo "4. 生成 goroutine 分析报告"
echo "5. 生成内存分析报告"
echo ""
read -p "请输入选项 (1-5): " choice

case $choice in
    1)
        echo ""
        echo "开始实时监控（按 Ctrl+C 停止）..."
        echo "----------------------------------------"
        while true; do
            clear
            echo "=== 实时监控 ==="
            echo "时间: $(date '+%H:%M:%S')"
            echo ""
            
            # 显示进程信息
            echo "--- 进程资源 ---"
            ps -p $PID -o pid,vsz,rss,%mem,%cpu,etime,cmd 2>/dev/null || { echo "进程已退出"; exit 1; }
            echo ""
            
            # 获取goroutine数量
            echo "--- Goroutine 统计 ---"
            GOROUTINE_COUNT=$(curl -s http://localhost:6060/debug/pprof/goroutine?debug=1 | grep -c "^goroutine")
            echo "当前 Goroutine 数量: $GOROUTINE_COUNT"
            echo ""
            
            # 获取内存统计
            echo "--- 内存统计 (pprof) ---"
            curl -s http://localhost:6060/debug/pprof/heap?debug=1 | grep -A 5 "# runtime.MemStats"
            
            sleep 2
        done
        ;;
        
    2)
        echo ""
        echo "获取 goroutine 堆栈详情..."
        curl -s http://localhost:6060/debug/pprof/goroutine?debug=2 | less
        ;;
        
    3)
        echo ""
        echo "获取内存分配详情..."
        curl -s http://localhost:6060/debug/pprof/heap?debug=1 | less
        ;;
        
    4)
        echo ""
        echo "生成 goroutine 分析报告..."
        FILENAME="goroutine_$(date +%Y%m%d_%H%M%S).txt"
        curl -s http://localhost:6060/debug/pprof/goroutine?debug=2 > "$FILENAME"
        echo "报告已保存到: $FILENAME"
        echo ""
        echo "统计信息:"
        grep -c "^goroutine" "$FILENAME" | xargs echo "Goroutine 总数:"
        echo ""
        echo "按函数分组统计:"
        grep "^goroutine" "$FILENAME" | awk '{print $3}' | sort | uniq -c | sort -rn | head -10
        ;;
        
    5)
        echo ""
        echo "生成内存分析报告（需要安装 go tool pprof）..."
        go tool pprof -text http://localhost:6060/debug/pprof/heap
        ;;
        
    *)
        echo "无效选项"
        exit 1
        ;;
esac


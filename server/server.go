package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// HeartRateServer 心率HTTP服务器结构体
type HeartRateServer struct {
	FilePath   string
	ListenAddr string
	server     *http.Server
	lastHR     int
	mutex      sync.Mutex
	IsRunning  bool
}

// NewHeartRateServer 创建一个新的心率服务器
func NewHeartRateServer(filePath, listenAddr string) *HeartRateServer {
	return &HeartRateServer{
		FilePath:   filePath,
		ListenAddr: listenAddr,
		lastHR:     0,
		IsRunning:  false,
	}
}

// Start 启动心率服务器
func (hrs *HeartRateServer) Start() error {
	if hrs.IsRunning {
		return fmt.Errorf("服务器已经在运行")
	}

	// 初始化心率
	hr, err := hrs.ReadHRFromFile()
	if err != nil {
		log.Printf("读取心率文件失败: %v, 使用默认值0", err)
	} else {
		hrs.lastHR = hr
	}

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", hrs.handleRequest)

	// 创建服务器
	hrs.server = &http.Server{
		Addr:    hrs.ListenAddr,
		Handler: mux,
	}

	// 启动服务器
	go func() {
		log.Printf("心率服务器开始运行于 %s\n", hrs.ListenAddr)
		if err := hrs.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("心率服务器错误: %v\n", err)
		}
	}()

	hrs.IsRunning = true
	return nil
}

// Stop 停止心率服务器
func (hrs *HeartRateServer) Stop() error {
	if !hrs.IsRunning || hrs.server == nil {
		return nil
	}

	err := hrs.server.Close()
	if err != nil {
		return fmt.Errorf("关闭服务器失败: %v", err)
	}

	hrs.IsRunning = false
	return nil
}

// 处理HTTP请求
func (hrs *HeartRateServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hrs.handleGetRequest(w, r)
	case http.MethodPost:
		hrs.handlePostRequest(w, r)
	default:
		http.Error(w, "不支持的请求方法", http.StatusMethodNotAllowed)
	}
}

// 处理GET请求
func (hrs *HeartRateServer) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	hrs.mutex.Lock()
	defer hrs.mutex.Unlock()

	// 读取最新心率
	hr, err := hrs.ReadHRFromFile()
	if err != nil {
		log.Printf("读取心率文件失败: %v, 使用上次记录的值", err)
		hr = hrs.lastHR
	} else {
		hrs.lastHR = hr
	}

	// 返回JSON格式的心率数据
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"bpm": %d}`, hr)
}

// 处理POST请求
func (hrs *HeartRateServer) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	hrs.mutex.Lock()
	defer hrs.mutex.Unlock()

	// 从请求头中获取心率
	bpmHeader := r.Header.Get("bpm")
	if bpmHeader != "" {
		if hr, err := strconv.Atoi(bpmHeader); err == nil {
			hrs.lastHR = hr
			log.Printf("收到新的心率: %d\n", hrs.lastHR)
			
			// 更新心率到文件
			if err := hrs.UpdateHRToFile(hr); err != nil {
				log.Printf("更新心率到文件失败: %v\n", err)
			}
		}
	}

	// 读取请求体 - 但不使用内容，只是消耗掉请求体
	_, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("读取请求体失败: %v\n", err)
	}
	defer r.Body.Close()

	// 响应
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "POST请求已接收")
}

// ReadHRFromFile 从文件读取心率
func (hrs *HeartRateServer) ReadHRFromFile() (int, error) {
	data, err := os.ReadFile(hrs.FilePath)
	if err != nil {
		return 0, fmt.Errorf("读取文件失败: %v", err)
	}

	hrStr := strings.TrimSpace(string(data))
	hr, err := strconv.Atoi(hrStr)
	if err != nil {
		return 0, fmt.Errorf("解析心率值失败: %v", err)
	}

	return hr, nil
}

// UpdateHRToFile 更新心率到文件
func (hrs *HeartRateServer) UpdateHRToFile(hr int) error {
	return os.WriteFile(hrs.FilePath, []byte(strconv.Itoa(hr)), 0644)
} 
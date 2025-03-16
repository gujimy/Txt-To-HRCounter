package main

import (
	"go-hr-counter/ui"
)

func main() {
	// 创建GUI实例
	heartRateUI := ui.NewHeartRateMonitorUI()
	
	// 设置GUI界面
	heartRateUI.Setup()
	
	// 运行应用
	heartRateUI.Run()
} 
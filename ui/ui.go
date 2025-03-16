package ui

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	fyredialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	
	nativedialogs "github.com/sqweek/dialog"

	"go-hr-counter/config"
	"go-hr-counter/server"
)

// HeartRateMonitorUI GUI界面结构体
type HeartRateMonitorUI struct {
	app            fyne.App
	window         fyne.Window
	server         *server.HeartRateServer
	config         *config.Config
	filePathEntry  *widget.Entry
	listeningAddr  *widget.Entry
	listeningPort  *widget.Entry
	serverStatus   *widget.Label
	currentHR      *widget.Label
	startStopBtn   *widget.Button
	selectFileBtn  *widget.Button
}

// NewHeartRateMonitorUI 创建一个新的GUI界面
func NewHeartRateMonitorUI() *HeartRateMonitorUI {
	// 创建应用和窗口
	a := app.New()
	w := a.NewWindow("心率监测服务")
	w.Resize(fyne.NewSize(600, 400))  // 调整窗口大小

	// 创建GUI结构体
	ui := &HeartRateMonitorUI{
		app:    a,
		window: w,
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("加载配置失败: %v, 使用默认配置", err)
		cfg = config.DefaultConfig()
	}
	ui.config = cfg

	// 设置默认值
	ui.server = server.NewHeartRateServer(ui.config.FilePath, ui.config.GetFullListenAddr())

	return ui
}

// Setup 设置GUI界面
func (ui *HeartRateMonitorUI) Setup() {
	// 创建文件路径输入框
	ui.filePathEntry = widget.NewEntry()
	ui.filePathEntry.SetText(ui.config.FilePath)
	ui.filePathEntry.SetPlaceHolder("请输入心率文件的完整路径")
	fileEntryWrapped := container.NewHScroll(ui.filePathEntry)
	fileEntryWrapped.Resize(fyne.NewSize(300, 36)) // 调整大小
	ui.filePathEntry.OnChanged = func(value string) {
		ui.config.FilePath = value
		ui.server.FilePath = value
		ui.saveConfig()
	}

	// 创建选择文件按钮
	ui.selectFileBtn = widget.NewButtonWithIcon("选择文件", theme.FileIcon(), func() {
		// 设置初始目录
		initialDir := ui.config.LastSaveDir
		if initialDir == "" {
			// 如果没有上次保存的目录，尝试使用当前文件路径的目录
			if ui.config.FilePath != "" {
				initialDir = ui.config.FilePath
			}
		}

		// 使用系统原生文件选择对话框
		filePath, err := nativedialogs.File().Filter("文本文件", "txt").Title("选择心率文件").SetStartDir(initialDir).Load()
		if err != nil {
			if err != nativedialogs.ErrCancelled {
				// 只在非取消操作时显示错误
				fyredialog.ShowError(err, ui.window)
			}
			return
		}
		
		// 保存选择的目录
		ui.config.LastSaveDir = filePath
		
		// 设置选中的文件路径
		ui.filePathEntry.SetText(filePath)
		ui.config.FilePath = filePath
		ui.server.FilePath = filePath
		ui.saveConfig()
	})

	// 创建监听地址和端口输入框
	ui.listeningAddr = widget.NewEntry()
	ui.listeningAddr.SetText(ui.config.ListenAddr)
	
	ui.listeningPort = widget.NewEntry()
	ui.listeningPort.SetText(ui.config.ListenPort)

	// 当地址或端口发生变化时更新服务器设置
	updateListeningAddr := func(s string) {
		addr := ui.listeningAddr.Text
		port := ui.listeningPort.Text
		
		ui.config.ListenAddr = addr
		ui.config.ListenPort = port
		ui.server.ListenAddr = fmt.Sprintf("%s:%s", addr, port)
		ui.saveConfig()
	}

	ui.listeningAddr.OnChanged = updateListeningAddr
	ui.listeningPort.OnChanged = updateListeningAddr

	// 创建状态标签，使用粗体显示状态
	ui.serverStatus = widget.NewLabel("服务器未启动")
	ui.serverStatus.TextStyle = fyne.TextStyle{Bold: true}
	
	ui.currentHR = widget.NewLabel("当前心率: 未知")
	ui.currentHR.TextStyle = fyne.TextStyle{Bold: true}

	// 创建启动/停止按钮，增加按钮大小
	ui.startStopBtn = widget.NewButton("启动服务器", func() {
		if ui.server.IsRunning {
			err := ui.server.Stop()
			if err != nil {
				log.Printf("停止服务器失败: %v\n", err)
				fyredialog.ShowError(fmt.Errorf("停止服务器失败: %v", err), ui.window)
				return
			}
			ui.serverStatus.SetText("服务器已停止")
			ui.startStopBtn.SetText("启动服务器")
		} else {
			err := ui.server.Start()
			if err != nil {
				log.Printf("启动服务器失败: %v\n", err)
				fyredialog.ShowError(fmt.Errorf("启动服务器失败: %v", err), ui.window)
				return
			}
			ui.serverStatus.SetText("服务器运行中 - " + ui.server.ListenAddr)
			ui.startStopBtn.SetText("停止服务器")

			// 读取初始心率并显示
			hr, _ := ui.server.ReadHRFromFile()
			ui.currentHR.SetText("当前心率: " + strconv.Itoa(hr))
		}
	})
	ui.startStopBtn.Importance = widget.HighImportance

	// 创建刷新心率按钮
	refreshHRBtn := widget.NewButton("刷新心率", func() {
		hr, err := ui.server.ReadHRFromFile()
		if err != nil {
			ui.currentHR.SetText("心率读取失败")
		} else {
			ui.currentHR.SetText("当前心率: " + strconv.Itoa(hr))
		}
	})

	// 创建配置说明标签
	configInstructions := widget.NewLabel("Beat Saber配置：修改 UserData/HRCounter.json 设置")
	configSettings := widget.NewLabel("\"DataSource\": \"WebRequest\",\n\"FeedLink\": \"http://" + ui.server.ListenAddr + "/\"")
	configSettings.TextStyle = fyne.TextStyle{Monospace: true}

	// 当地址变化时更新配置说明
	updateConfigSettings := func(s string) {
		configSettings.SetText("\"DataSource\": \"WebRequest\",\n\"FeedLink\": \"http://" + ui.server.ListenAddr + "/\"")
	}
	ui.listeningAddr.OnChanged = func(s string) {
		updateListeningAddr(s)
		updateConfigSettings(s)
	}
	ui.listeningPort.OnChanged = func(s string) {
		updateListeningAddr(s)
		updateConfigSettings(s)
	}

	// 创建带有说明标签的文件路径区域
	fileLabel := widget.NewLabel("心率文件路径:")
	fileLabel.Alignment = fyne.TextAlignLeading

	// 布局UI组件，使用更好的布局方式
	fileBox := container.NewGridWithColumns(3,
		fileLabel,
		fileEntryWrapped,
		ui.selectFileBtn,
	)

	// 监听地址和端口区域
	addrLabel := widget.NewLabel("监听地址:")
	addrLabel.Alignment = fyne.TextAlignLeading
	
	portLabel := widget.NewLabel("端口:")
	portLabel.Alignment = fyne.TextAlignLeading

	addrPortBox := container.NewGridWithColumns(4,
		addrLabel,
		ui.listeningAddr,
		portLabel,
		ui.listeningPort,
	)

	// 状态和按钮区域
	statusBox := container.NewHBox(
		ui.serverStatus,
		layout.NewSpacer(),
		ui.currentHR,
	)

	buttonBox := container.NewHBox(
		ui.startStopBtn,
		layout.NewSpacer(),
		refreshHRBtn,
	)

	// 配置说明区域
	configBox := container.NewVBox(
		container.NewPadded(configInstructions),
		container.NewPadded(configSettings),
	)

	// 创建主布局，增加边距
	content := container.NewVBox(
		container.NewPadded(fileBox),
		container.NewPadded(addrPortBox),
		widget.NewSeparator(),
		container.NewPadded(statusBox),
		container.NewPadded(buttonBox),
		widget.NewSeparator(),
		container.NewPadded(configBox),
	)

	// 设置窗口内容
	ui.window.SetContent(content)
	
	// 设置窗口关闭事件，保存配置
	ui.window.SetOnClosed(func() {
		ui.saveConfig()
	})
}

// saveConfig 保存配置
func (ui *HeartRateMonitorUI) saveConfig() {
	err := config.SaveConfig(ui.config)
	if err != nil {
		log.Printf("保存配置失败: %v", err)
	}
}

// Run 运行GUI应用
func (ui *HeartRateMonitorUI) Run() {
	ui.window.ShowAndRun()
} 
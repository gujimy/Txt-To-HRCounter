# 心率显示服务器 (Go版本)

这是一个为小米手环（Mi Band 2/3/4）设计的心率数据转发服务器，用于节奏光剑（Beat Saber）游戏。本程序将心率文本文件中的数据通过HTTP服务转发给游戏的心率显示mod。

![界面预览](1.png)

## 功能特点

- **简单易用的GUI界面**：直观的界面设计，一目了然的操作流程
- **完整的配置保存**：自动保存所有设置，下次启动自动恢复
- **文件路径自定义**：支持自由选择心率文本文件路径，使用系统原生文件选择器
- **监听地址配置**：可自定义HTTP服务监听地址和端口
- **实时状态显示**：提供实时心率显示和服务状态指示
- **一键操作**：支持一键启动和停止服务
- **完整配置指引**：界面内置Beat Saber配置说明，简化设置过程

## 系统要求

- Windows 系统
- 小米手环 2/3/4 及其蓝牙连接
- 需要安装 [Mi Band Heartrate](https://github.com/Eryux/miband-heartrate) 软件获取心率数据
- 节奏光剑需要安装 [HRCounter](https://github.com/qe201020335/HRCounter) mod

## 使用方法

### 第一步：获取心率数据
1. 安装并运行 [Mi Band Heartrate](https://github.com/Eryux/miband-heartrate) 软件
2. 连接小米手环并开始记录心率数据到文本文件

### 第二步：配置服务器
1. 下载并运行本程序 `heart-rate-monitor.exe`
2. 在GUI界面上设置心率文本文件路径（可使用"选择文件"按钮）
3. 设置监听地址和端口（默认为`localhost:2548`）
4. 点击"启动服务器"按钮开启服务

### 第三步：配置Beat Saber
1. 安装 [HRCounter](https://github.com/qe201020335/HRCounter) mod
2. 修改 `Beat Saber\UserData\HRCounter.json` 文件
3. 按照程序界面下方显示的配置信息设置：
   ```json
   "DataSource": "WebRequest",
   "FeedLink": "http://localhost:2548/"
   ```
4. 保存配置文件并启动游戏

## 关于杀毒软件误报

一些杀毒软件可能会将本程序误报为病毒或木马。这是因为：

1. **技术原因**：
   - Go语言编译的程序包含完整运行时，生成较大可执行文件
   - 程序包含网络服务器代码和文件读写功能
   - 使用了系统API和窗口隐藏技术

2. **解决方法**：
   - 在杀毒软件中将程序添加到白名单/排除列表
   - 您可以查看源代码并自行编译程序
   - 本程序是开源的，不含任何恶意代码

3. **安全保证**：
   - 程序只读取指定的心率文本文件
   - 只提供HTTP服务，不发送任何数据到互联网
   - 所有配置保存在程序目录下的config.json文件中

## 技术实现

- 使用Go语言开发，提供跨平台支持
- GUI界面使用Fyne框架实现
- 原生文件选择对话框使用sqweek/dialog库
- HTTP服务使用Go标准库实现
- 配置保存使用JSON格式

## 常见问题

**Q: 为什么心率数据无法显示在游戏中？**

A: 请检查以下几点：
- 确认心率文件正确生成并包含有效数据
- 确认服务器已成功启动（状态显示"服务器运行中"）
- 确认Beat Saber的HRCounter.json配置已正确设置
- 尝试点击"刷新心率"按钮检查当前读取的心率值

**Q: 程序启动后保存在哪里？**

A: 所有设置都保存在程序目录下的`config.json`文件中，程序下次启动时会自动加载。

**Q: 如何修改默认设置？**

A: 可以直接编辑`config.json`文件，或在程序界面中修改后自动保存。

## 许可证

本程序基于MIT许可证开源，您可以自由使用、修改和分发本程序，详情请参阅LICENSE文件。

## 鸣谢

- [Eryux/miband-heartrate](https://github.com/Eryux/miband-heartrate) - 小米手环心率读取
- [qe201020335/HRCounter](https://github.com/qe201020335/HRCounter) - Beat Saber心率显示mod 
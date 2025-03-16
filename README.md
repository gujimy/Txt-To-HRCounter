# 心率显示服务器 (Go版本)

这是一个为小米手环（Mi Band 2/3/4）设计的心率数据转发服务器，用于节奏光剑（Beat Saber）游戏。本程序将心率文本文件中的数据通过HTTP服务转发给游戏的心率显示mod。



* 电脑需要有蓝牙连接手环
* 使用 Mi Band Heartrate([download](https://github.com/Eryux/miband-heartrate#mi-band-heartrate))) 获得txt心率文件
* 节奏光剑心率mod HRCounter ([download](https://github.com/qe201020335/HRCounter)))

## 功能特点

- 使用Go语言重构实现，提供直观的GUI界面
- 支持自由选择心率文本文件路径
- 可自定义HTTP服务监听地址和端口
- 提供实时心率显示和服务状态指示
- 支持一键启动和停止服务

## 使用方法

1. 下载并运行本程序
2. 在GUI界面上设置心率文本文件路径（默认为`D:\heartrate\heartrate.txt`）
3. 设置监听地址和端口（默认为`localhost:2548`）
4. 点击"启动服务器"按钮开启服务
5. 按照界面上的提示，修改Beat Saber的HRCounter.json配置文件
* mod部分需要修改的部分找到\Beat Saber\UserData\HRCounter.json 打开
  修改部分为：

```
"DataSource": "WebRequest",

"FeedLink": "http://localhost:2548/",
```
![image](https://github.com/gujimy/Txt-To-HRCounter/assets/40573598/bdb7ef17-66d4-4d8b-bf0a-fe870c1e9bce)

## 系统要求

- Windows 系统
- 小米手环 2/3/4 及其蓝牙连接
- 需要安装 [Mi Band Heartrate](https://github.com/Eryux/miband-heartrate) 软件获取心率数据
- 节奏光剑需要安装 [HRCounter](https://github.com/qe201020335/HRCounter) mod

## 编译方法

```bash
# 克隆仓库
git clone https://github.com/your-username/go-hr-counter.git
cd go-hr-counter

# 安装依赖
go mod tidy

# 编译
go build
```

## 项目结构

- `server/`: 包含心率HTTP服务器的实现
- `ui/`: 包含GUI界面的实现
- `main.go`: 程序入口点 
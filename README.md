# Txt-To-HRCounter
Mi Band 2 3 4 Heartrate for Beat Saber 
给小米手环2 3 4 提供的一个节奏光剑 心率显示方案

使用 https://github.com/Eryux/miband-heartrate 获得txt心率文件

节奏光剑心率mod为 https://github.com/qe201020335/HRCounter

心率txt文件默认设定的存放路径为 D:\heartrate\heartrate.txt 

mod部分需要修改的部分找到\Beat Saber\UserData\HRCounter.json 打开
修改部分为：

"DataSource": "WebRequest",

"FeedLink": "http://localhost:2548/",

![image](https://github.com/gujimy/Txt-To-HRCounter/assets/40573598/bdb7ef17-66d4-4d8b-bf0a-fe870c1e9bce)

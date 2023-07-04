using System;
using System.IO;
using System.Net;
using System.Text;
using System.Linq;

namespace HeartRateServer
{
    class Program
    {
        static int lastHr = 0;

        static void Main(string[] args)
        {
            string filePath = @"D:\heartrate\heartrate.txt"; // 文本文件路径

            // 读取初始的心率值
            lastHr = ReadHrFromFile(filePath);

            // 创建HTTP服务器
            HttpListener listener = new HttpListener();
            listener.Prefixes.Add("http://localhost:2548/"); // 监听的URL
            listener.Start();

            Console.WriteLine("Server started. Listening for requests...");

            // 处理请求
            while (true)
            {
                HttpListenerContext context = listener.GetContext();
                HttpListenerRequest request = context.Request;
                HttpListenerResponse response = context.Response;

                if (request.HttpMethod == "GET")
                {
                    // 在每次GET请求中读取最新的心率值
                    int hr = ReadHrFromFile(filePath);

                    byte[] buffer = Encoding.UTF8.GetBytes($"{{\"bpm\": {hr}}}");
                    response.StatusCode = 200;
                    response.ContentType = "application/json";
                    response.ContentLength64 = buffer.Length;
                    response.OutputStream.Write(buffer, 0, buffer.Length);
                    response.Close();
                }
                else if (request.HttpMethod == "POST")
                {
                    using (StreamReader sr = new StreamReader(request.InputStream))
                    {
                        string body = sr.ReadToEnd();
                        if (request.Headers.AllKeys.Contains("bpm"))
                        {
                            int hr;
                            if (Int32.TryParse(request.Headers["bpm"], out hr))
                            {
                                lastHr = hr;
                                Console.WriteLine("Received new heart rate: " + lastHr);
                            }
                        }

                        byte[] buffer = Encoding.UTF8.GetBytes("POST request received");
                        response.StatusCode = 200;
                        response.ContentType = "text/plain";
                        response.ContentLength64 = buffer.Length;
                        response.OutputStream.Write(buffer, 0, buffer.Length);
                        response.Close();
                    }

                    // 在POST请求处理完成后，更新心率值到文本文件
                    UpdateHrToFile(lastHr, filePath);
                }
            }
        }

        static int ReadHrFromFile(string filePath)
        {
            try
            {
                string hrText = File.ReadAllText(filePath);
                int hr;
                if (Int32.TryParse(hrText.Trim(), out hr))
                {
                    return hr;
                }
            }
            catch (Exception e)
            {
                Console.WriteLine("Error reading heart rate from file: " + e.Message);
            }

            return 0;
        }

        static void UpdateHrToFile(int hr, string filePath)
        {
            try
            {
                File.WriteAllText(filePath, hr.ToString());
            }
            catch (Exception e)
            {
                Console.WriteLine("Error updating heart rate to file: " + e.Message);
            }
        }
    }
}

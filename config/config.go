package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 表示程序的配置信息
type Config struct {
	FilePath    string `json:"file_path"`    // 心率文件路径
	ListenAddr  string `json:"listen_addr"`  // 监听地址
	ListenPort  string `json:"listen_port"`  // 监听端口
	LastSaveDir string `json:"last_save_dir"` // 上次保存的目录
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		FilePath:    "D:\\heartrate\\heartrate.txt",
		ListenAddr:  "localhost",
		ListenPort:  "2548",
		LastSaveDir: "",
	}
}

// GetFullListenAddr 获取完整的监听地址
func (c *Config) GetFullListenAddr() string {
	return fmt.Sprintf("%s:%s", c.ListenAddr, c.ListenPort)
}

// SetFullListenAddr 设置完整的监听地址
func (c *Config) SetFullListenAddr(addr string) {
	// 解析地址和端口
	parts := strings.Split(addr, ":")
	if len(parts) == 2 {
		c.ListenAddr = parts[0]
		c.ListenPort = parts[1]
	} else {
		// 格式错误，使用默认值
		c.ListenAddr = "localhost"
		c.ListenPort = "2548"
	}
}

// ConfigFilePath 返回配置文件的路径
func ConfigFilePath() (string, error) {
	// 获取程序所在目录
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取程序路径失败: %v", err)
	}
	
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "config.json"), nil
}

// LoadConfig 加载配置信息
func LoadConfig() (*Config, error) {
	configPath, err := ConfigFilePath()
	if err != nil {
		return DefaultConfig(), err
	}
	
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 文件不存在，使用默认配置
		config := DefaultConfig()
		// 尝试保存默认配置
		_ = SaveConfig(config)
		return config, nil
	}
	
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("读取配置文件失败: %v", err)
	}
	
	// 解析JSON
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return DefaultConfig(), fmt.Errorf("解析配置文件失败: %v", err)
	}
	
	return config, nil
}

// SaveConfig 保存配置信息
func SaveConfig(config *Config) error {
	configPath, err := ConfigFilePath()
	if err != nil {
		return err
	}
	
	// 序列化配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	
	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	
	return nil
} 
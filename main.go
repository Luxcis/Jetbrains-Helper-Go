package main

import (
	"Jetbrain-Helper/config"
	"Jetbrain-Helper/helper"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"time"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}
	err = viper.Unmarshal(&config.Conf)
	if err != nil {
		log.Fatalf("配置文件类型转换失败: %v", err)
	}
	// 监听配置文件修改
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.Unmarshal(&config.Conf)
		if err != nil {
			log.Printf("配置文件类型转换失败: %v\n", err)
			return
		}
		log.Printf("配置文件修改(%s): %s\n", time.Now().Format("2006-01-01 15:04:05"), e.Name)
	})
	helper.InitProducts()
	helper.InitPlugins()
	helper.InitCertificate()
	helper.InitAgent()
}

func main() {
	gin.SetMode(config.Conf.Mode)
	r := gin.Default()
	r.Static("/assets", "./static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", index)
	r.GET("/search", search)
	r.POST("/generateLicense", generateLicense)
	r.GET("/ja-netfilter", download)
	startCron()
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func startCron() {
	c := cron.New()
	_, err := c.AddFunc("0 0 12 * *", helper.RefreshJsonFile)
	if err != nil {
		log.Fatalf("定时任务添加失败: %v", err)
	} else {
		log.Println("定时重启启动成功")
	}
	c.Start()
}

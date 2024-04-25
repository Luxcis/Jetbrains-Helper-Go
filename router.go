package main

import (
	"Jetbrain-Helper/config"
	"Jetbrain-Helper/helper"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"products": helper.Products,
		"plugins":  helper.Plugins,
		"defaults": config.Conf.License,
	})
}

func search(c *gin.Context) {
	var filteredProductList []helper.ProductCache
	var filteredPluginList []helper.PluginCache
	search := c.Query("search")
	if search != "" {
		for _, productCache := range helper.Products {
			if strings.Contains(strings.ToLower(productCache.Name), strings.ToLower(search)) {
				filteredProductList = append(filteredProductList, productCache)
			}
		}
		for _, pluginCache := range helper.Plugins {
			if strings.Contains(strings.ToLower(pluginCache.Name), strings.ToLower(search)) {
				filteredPluginList = append(filteredPluginList, pluginCache)
			}
		}
	} else {
		filteredProductList = helper.Products
		filteredPluginList = helper.Plugins
	}
	c.HTML(http.StatusOK, "articles", gin.H{
		"products": filteredProductList,
		"plugins":  filteredPluginList,
		"defaults": config.Conf.License,
	})
}

func generateLicense(c *gin.Context) {
	var req helper.GenerateLicenseReqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var productCodeSet []string
	if req.ProductCode == "" {
		var productCodeList []string
		for _, product := range helper.Products {
			if product.ProductCode != "" {
				productCodeList = append(productCodeList, strings.Split(product.ProductCode, ",")...)
			}
		}
		for _, plugin := range helper.Plugins {
			if plugin.ProductCode != "" {
				productCodeList = append(productCodeList, plugin.ProductCode)
			}
		}
	} else {
		productCodeSet = append(productCodeSet, strings.Split(req.ProductCode, ",")...)
	}
	license, err := helper.GenerateLicense(
		req.LicenseName,
		req.AssigneeName,
		req.ExpiryDate,
		productCodeSet,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, license)
	}
}

func download(c *gin.Context) {
	file := helper.OpenFile(helper.JaNetfilterFilePath + ".zip")
	fileStat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=ja-netfilter.zip")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	// 将文件内容复制到响应主体
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}

	c.Status(http.StatusOK)
}

package helper

import (
	"log"
)

var Products []ProductCache
var Plugins []PluginCache

func InitProducts() {
	log.Println("ProductCache context init loading...")
	err := ReadJson(products, &Products)
	if err != nil {
		log.Fatalf("Product初始化失败: %v", err)
		return
	}
}

func InitPlugins() {
	log.Println("PluginCache context init loading...")
	err := ReadJson(plugins, &Plugins)
	if err != nil {
		log.Fatalf("Plugin初始化失败: %v", err)
		return
	}
}

const (
	products = "external/data/product.json"
	plugins  = "external/data/plugin.json"
)

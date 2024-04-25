package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	PluginBasicUrl = "https://plugins.jetbrains.com"
	PluginListUrl  = "/api/searchPlugins?max=10000&offset=0&orderBy=name"
	PluginInfoUrl  = "/api/plugins/"
)

func RefreshJsonFile() {
	log.Println("Init or Refresh plugin context from 'JetBrains.com' loading...")
	plugins, err := pluginList()
	if err != nil {
		log.Printf("Plugin context init or refresh failed: %v\n", err)
		return
	}

	filteredPlugins := pluginListFilter(plugins)
	pluginCaches := pluginConversion(filteredPlugins)

	overrideJsonFile(pluginCaches)

	log.Println("Init or Refresh plugin context success !")
}

func overrideJsonFile(pluginCaches []PluginCache) {
	Plugins = append(Plugins, pluginCaches...)

	jsonStr, err := json.Marshal(Plugins)
	if err != nil {
		log.Fatalf("%s File write failed: %v\n", plugins, err)
		return
	}

	err = os.WriteFile(plugins, jsonStr, 0644)
	if err != nil {
		log.Fatalf("%s File write failed: %v\n", plugins, err)
		return
	}
}

func pluginList() (PluginList, error) {
	resp, err := http.Get(PluginListUrl)
	if err != nil {
		return PluginList{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PluginList{}, fmt.Errorf("The request failed = %d", resp.StatusCode)
	}

	var pluginList PluginList
	err = json.NewDecoder(resp.Body).Decode(&pluginList)
	if err != nil {
		return PluginList{}, err
	}

	return pluginList, nil
}

func pluginListFilter(pluginList PluginList) []Plugin {
	filteredPlugins := []Plugin{}
	for _, plugin := range pluginList.Plugins {
		if !contains(plugin.Id) && plugin.PricingModel != "FREE" {
			filteredPlugins = append(filteredPlugins, plugin)
		}
	}
	return filteredPlugins
}

func contains(pluginId int) bool {
	for _, plugin := range Plugins {
		if plugin.Id == pluginId {
			return true
		}
	}
	return false
}

func pluginConversion(pluginList []Plugin) []PluginCache {
	pluginCaches := []PluginCache{}
	for _, plugin := range pluginList {
		info, err := pluginInfo(plugin.Id)
		if err != nil {
			log.Printf("Failed to get plugin info: %v\n", err)
			continue
		}
		cache := PluginCache{
			Id:           plugin.Id,
			ProductCode:  info.PurchaseInfo.ProductCode,
			Name:         plugin.Name,
			PricingModel: plugin.PricingModel,
			Icon:         PluginBasicUrl + plugin.Icon,
		}
		pluginCaches = append(pluginCaches, cache)
	}
	return pluginCaches
}

func pluginInfo(pluginId int) (PluginInfo, error) {
	resp, err := http.Get(PluginInfoUrl + fmt.Sprint(pluginId))
	if err != nil {
		return PluginInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PluginInfo{}, fmt.Errorf("The request failed = %d", resp.StatusCode)
	}

	var pluginInfo PluginInfo
	err = json.NewDecoder(resp.Body).Decode(&pluginInfo)
	if err != nil {
		return PluginInfo{}, err
	}

	return pluginInfo, nil
}

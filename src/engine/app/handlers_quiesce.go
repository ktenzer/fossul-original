package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"log"
	"os"
	"strings"
)

func PreQuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PreAppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing pre quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func QuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.AppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func Quiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/app/" + config.AppPlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "\nERROR: App plugin does not exist: " + plugin
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "app", plugin, "--action", "quiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostQuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PostAppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing post quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
package main

import (
	"encoding/json"
	"engine/util"
	"engine/client"
	"net/http"
	"strings"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)

	var sendTrapErrorCmdResult util.Result
	var results []util.Result

	commentResult := util.SetResultMessage(0,"COMMENT","Welcome to Fossil Backup Framework, Performing Backup Workflow")
	results = append(results, commentResult)

	var preQuiesceCmdResult util.Result
	preQuiesceCmdResult = client.PreQuiesceCmd(config)
	results = append(results, preQuiesceCmdResult)

	if preQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Application Quiesce")
	results = append(results, commentResult)

	var quiesceCmdResult util.Result
	quiesceCmdResult = client.QuiesceCmd(config)	
	results = append(results, quiesceCmdResult)

	if quiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}
	
	var quiesceResult util.Result
	quiesceResult = client.Quiesce(config)
	results = append(results, quiesceResult)

	if quiesceResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	var postQuiesceCmdResult util.Result
	postQuiesceCmdResult = client.PostQuiesceCmd(config)
	results = append(results, postQuiesceCmdResult)

	if postQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Backup")
	results = append(results, commentResult)

	var backupResult util.Result
	backupResult = client.Backup(config)
	results = append(results, backupResult)

	if backupResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	var preUnquiesceCmdResult util.Result
	preUnquiesceCmdResult = client.PreUnquiesceCmd(config)
	results = append(results, preUnquiesceCmdResult)

	if preUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Application Unquiesce")
	results = append(results, commentResult)

	var unquiesceCmdResult util.Result
	unquiesceCmdResult = client.UnquiesceCmd(config)	
	results = append(results, unquiesceCmdResult)

	if unquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	var unquiesceResult util.Result
	unquiesceResult = client.Unquiesce(config)
	results = append(results, unquiesceResult)

	if unquiesceResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	var postUnquiesceCmdResult util.Result
	postUnquiesceCmdResult = client.PostUnquiesceCmd(config)
	results = append(results, postUnquiesceCmdResult)

	if postUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)	
	}

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Backup Retention")
	results = append(results, commentResult)

	var backupDeleteResult util.Result
	backupDeleteResult = client.BackupDelete(config)
	results = append(results, backupDeleteResult)

	if backupDeleteResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	var sendTrapSuccessCmdResult util.Result
	sendTrapSuccessCmdResult = client.SendTrapSuccessCmd(config)	
	results = append(results, sendTrapSuccessCmdResult)

	if sendTrapSuccessCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
		sendError(w,r,results)
	}

	commentResult = util.SetResultMessage(0,"INFO","Backup Completed Successfully")
	results = append(results, commentResult)
	
	results = append(results, sendTrapErrorCmdResult)

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func sendError(w http.ResponseWriter, r *http.Request, results []util.Result) {
	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result

	var config util.Config = util.GetConfig(w,r)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapSuccessCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap success command")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func SendTrapErrorCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result

	var config util.Config = util.GetConfig(w,r)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapErrorCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap error command")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
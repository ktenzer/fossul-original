package main

import (
	"fmt"
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/util"
	"strings"
)

type appPlugin string

var AppPlugin appPlugin

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var result util.Result

	return result
}

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MysqlDb"]

	var dataFilePaths []string
	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId
	dataFilePaths = append(dataFilePaths, dumpPath)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is ["+dumpPath+"]")
	messages = append(messages, msg)

	discoverList = append(discoverList, discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func (a appPlugin) Quiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	var args []string
	var mkdirArgs []string

	podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId

	//create tmp directory for storing dump
	mkdirArgs = append(mkdirArgs, "mkdir")
	mkdirArgs = append(mkdirArgs, "-p")
	mkdirArgs = append(mkdirArgs, dumpPath)

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], mkdirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	//execute database dump
	args = append(args, "/bin/sh")
	args = append(args, "-c")

	if config.AppPluginParameters["MysqlPassword"] != "" {
		args = append(args, config.AppPluginParameters["MysqlDumpCmd"]+" -h "+config.AppPluginParameters["MysqlHost"]+
			" -P "+config.AppPluginParameters["MysqlPort"]+" -u "+config.AppPluginParameters["MysqlUser"]+" -p "+
			config.AppPluginParameters["MysqlPassword"]+" "+config.AppPluginParameters["MysqlDb"]+" "+dumpPath+"/mysql.sql")
	} else {
		args = append(args, config.AppPluginParameters["MysqlDumpCmd"]+" -h "+config.AppPluginParameters["MysqlHost"]+
			" -P "+config.AppPluginParameters["MysqlPort"]+" -u "+config.AppPluginParameters["MysqlUser"]+
			" "+config.AppPluginParameters["MysqlDb"]+" >"+dumpPath+"/mysql.sql")
	}

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Unquiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId

	var args []string
	args = append(args, "rm")
	args = append(args, "-rf")
	args = append(args, dumpPath)

	podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PreRestore(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message

	podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	//create tmp directory for storing dump
	var mkdirArgs []string
	restoreDir := "/tmp/" + util.IntToString(config.SelectedWorkflowId)
	mkdirArgs = append(mkdirArgs, "mkdir")
	mkdirArgs = append(mkdirArgs, "-p")
	mkdirArgs = append(mkdirArgs, restoreDir)

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], mkdirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PostRestore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var lsDirArgs []string
	lsDirArgs = append(lsDirArgs, "ls")
	lsDirArgs = append(lsDirArgs, "/tmp/"+util.IntToString(config.SelectedWorkflowId))

	cmdResult, restoreDir := k8s.ExecuteCommandWithStdout(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], lsDirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	restorePath := "/tmp/" + util.IntToString(config.SelectedWorkflowId) + "/" + strings.TrimSpace(restoreDir) + "/" + util.IntToString(config.SelectedWorkflowId) + "/mysql.sql"

	// execute dump using pg_dump (requires ld_library_path)
	var restoreArgs []string
	restoreArgs = append(restoreArgs, "/bin/sh")
	restoreArgs = append(restoreArgs, "-c")

	if config.AppPluginParameters["MysqlPassword"] != "" {
		restoreArgs = append(restoreArgs, config.AppPluginParameters["MysqlRestoreCmd"]+" -h "+config.AppPluginParameters["MysqlHost"]+
			" -P "+config.AppPluginParameters["MysqlPort"]+" -u "+config.AppPluginParameters["MysqlUser"]+" -p "+
			config.AppPluginParameters["MysqlPassword"]+" "+config.AppPluginParameters["MysqlDb"]+" "+restorePath)
	} else {
		restoreArgs = append(restoreArgs, config.AppPluginParameters["MysqlRestoreCmd"]+" -h "+config.AppPluginParameters["MysqlHost"]+
			" -P "+config.AppPluginParameters["MysqlPort"]+" -u "+config.AppPluginParameters["MysqlUser"]+
			" "+config.AppPluginParameters["MysqlDb"]+" <"+restorePath)
	}

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], restoreArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	var rmDirArgs []string
	restoreTmpDir := "/tmp/" + util.IntToString(config.SelectedWorkflowId)
	rmDirArgs = append(rmDirArgs, "rm")
	rmDirArgs = append(rmDirArgs, "-rf")
	rmDirArgs = append(rmDirArgs, restoreTmpDir)

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], rmDirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mariadb-dump"
	plugin.Description = "MariaDB plugin for backing up MySql or MariaDB using mysqldump utility"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var preRestoreCap util.Capability
	preRestoreCap.Name = "preRestore"

	var postRestoreCap util.Capability
	postRestoreCap.Name = "postRestore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, preRestoreCap, postRestoreCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func checkErr(err error) {
	fmt.Println("error handling")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {}

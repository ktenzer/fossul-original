package util

import (
	"fmt"
	"plugin"
	"errors"
)

type AppPlugin interface {
	SetEnv(Config) (Result)
	Quiesce() (Result)
	Unquiesce() (Result)
	Info() (Plugin)
}

type StoragePlugin interface {
	SetEnv(Config) (Result)
	Backup() (Result)
	BackupDelete() (Result)
	BackupList() ([]Backup)
	Info() (Plugin)
}

type ArchivePlugin interface {
	SetEnv(Config) (Result)
	Archive() (Result)
	ArchiveDelete() (Result)
	ArchiveList() ([]Archive)
	Info() (Plugin)
}

func GetPluginPath(pluginName string) string {
	var path string
	switch pluginName {
	case "mariadb.so":
		path = "./plugins/app/mariadb.so"
	case "container-basic.so":
		path = "./plugins/storage/container-basic.so"		
	case "sample-app.so":
		path = "./plugins/app/sample-app.so"	
	case "sample-storage.so":
		path = "./plugins/storage/sample-storage.so"	
	case "sample-archive.so":
		path = "./plugins/archive/sample-archive.so"						
	default:
		fmt.Println("Built-in plugin [" + pluginName + "] does not exist")
		path = ""
	}

	return path
}

func GetAppInterface(path string) (AppPlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("AppPlugin")
	if err != nil {
		return nil,err
	}
	type Plugin interface {
	
	}
	var appPlugin AppPlugin
	appPlugin, ok := symPlugin.(AppPlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], check plugin type, should be AppPlugin")
	}

	return appPlugin,nil
}

func GetStorageInterface(path string) (StoragePlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("StoragePlugin")
	if err != nil {
		return nil,err
	}

	var storagePlugin StoragePlugin
	storagePlugin, ok := symPlugin.(StoragePlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], check plugin type, should be StoragePlugin")
	}

	return storagePlugin,nil
}

func GetArchiveInterface(path string) (ArchivePlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("ArchivePlugin")
	if err != nil {
		return nil,err
	}

	var archivePlugin ArchivePlugin
	archivePlugin, ok := symPlugin.(ArchivePlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], check plugin type, should be ArchivePlugin and methods")
	}

	return archivePlugin,nil
}
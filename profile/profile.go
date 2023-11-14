package cherryProfile

import (
	"os"
	"path/filepath"

	cerror "github.com/cherry-game/cherry/error"
	cfile "github.com/cherry-game/cherry/extend/file"
	cfacade "github.com/cherry-game/cherry/facade"
	"github.com/spf13/viper"
)

var (
	cfg = &struct {
		profilePath string  // profile root dir
		profileName string  // profile name
		jsonConfig  *Config // profile-x.json parse to json object
		env         string  // env name
		debug       bool    // debug default is true
		printLevel  string  // cherry log print level
	}{}
)

func Path() string {
	return cfg.profilePath
}

func Name() string {
	return cfg.profileName
}

func Env() string {
	return cfg.env
}

func Debug() bool {
	return cfg.debug
}

func PrintLevel() string {
	return cfg.printLevel
}

func Init(filePath, nodeId string) (cfacade.INode, error) {
	if filePath == "" {
		return nil, cerror.Error("Profile file path is nil.")
	}

	if nodeId == "" {
		return nil, cerror.Error("NodeId is nil.")
	}

	judgePath, ok := cfile.JudgeFile(filePath)
	if !ok {
		return nil, cerror.Errorf("File path error. filePath = %s", filePath)
	}

	p, f := filepath.Split(judgePath)
	jsonConfig, err := loadFile(p, f)
	if err != nil {
		return nil, cerror.Errorf("Load profile file error. [err = %v]", err)
	}

	node, err := GetNodeWithConfig(jsonConfig, nodeId)
	if err != nil {
		return nil, cerror.Errorf("Failed to get node config from profile file. [err = %v]", err)
	}

	// init cfg
	cfg.profilePath = p
	cfg.profileName = f
	cfg.jsonConfig = jsonConfig
	cfg.env = jsonConfig.GetString("env", "default")
	cfg.debug = jsonConfig.GetBool("debug", true)
	cfg.printLevel = jsonConfig.GetString("print_level", "debug")

	return node, nil
}

func GetConfig(path ...interface{}) cfacade.ProfileCfg {
	return cfg.jsonConfig.GetConfig(path...)
}

func loadFile(filePath, fileName string) (*Config, error) {

	// read master json file
	fileNamePath := filepath.Join(filePath, fileName)
	viper.SetConfigFile(fileNamePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	paths := viper.GetStringSlice("include")

	for _, p := range paths {
		includePath := filepath.Join(filePath, p)
		includeFile, err := os.Open(includePath)
		if err != nil {
			panic(err)
		}
		err = viper.MergeConfig(includeFile)
		if err != nil {
			panic(err)
		}
	}

	return Wrap(viper.GetViper().AllSettings()), nil
}

//func judgeNameList(path, name string) ([]string, error) {
//	var list []string
//
//	if name != "" {
//		fileName := mergeProfileName(name)
//		list = append(list, fileName)
//
//	} else {
//		// find path
//		filesPath, err := cfile.ReadDir(path, "profile-", ".json")
//		if err != nil {
//			return nil, err
//		}
//
//		if len(filesPath) < 1 {
//			return nil, cerror.Errorf("[path = %s] profile file not found.", path)
//		}
//
//		for _, fp := range filesPath {
//			list = append(list, fp)
//		}
//	}
//
//	return list, nil
//}

//func mergeProfileName(name string) string {
//	return fmt.Sprintf("%s%s%s", profilePrefix, name, profileSuffix)
//}

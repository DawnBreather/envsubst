package cicd_envsubst

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/file"
	pth "cicd_envsubst/utils/path"
	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
)

//func main(){
//	ExecuteSetSubst()
//}

var setSubstOpts struct {
	Sets []string `short:"s" long:"set" description:"Variables set name" required:"true"`
}

func ExecuteSetSubst() {

	//os.Setenv("APPSETTINGS_JSON", playground.TestEnvVar)

	targetPaths := readCliOptions()

	_logger.Infof("Replacing environment variable sets { %v } in the following locations: { %v }", strings.Join(setSubstOpts.Sets, ", "), strings.Join(targetPaths, ", "))
	_logger.Infof("PROCESSING")

	//if len(os.Args) < 2 {
	//	_logger.Errorf("Not enough arguments:\n(1) - name of varibales set\n(2..) - target files")
	//}

	jointVariablesSet := readVariablesFromSets()

	substituteVariablesInTargetPaths(targetPaths, jointVariablesSet)

	_logger.Infof("DONE")
}

func readCliOptions() (args []string) {
	args, err := flags.Parse(&setSubstOpts)

	if err != nil {
		_logger.Fatalf("Unable to read cli options: %v", err)
	}

	return
}

func readVariablesFromSets() (res map[string]string) {

	res = map[string]string{}

	for _, setName := range setSubstOpts.Sets {

		variablesSetYaml := ""

		ev := env_var.EnvVar{}
		if ev.
			SetName(setName).
			IsExist() {
			variablesSetYaml = ev.Value()
		} else {
			_logger.Fatalf("Unable to locate variables set { %s }", setName)
			os.Exit(1)
		}

		tmpObj := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(variablesSetYaml), &tmpObj)
		if err != nil {
			_logger.Fatalf("Unable to unmarshal YAML content of variables set { %s }: %v", setName, err)
		}

		for varName, varValue := range tmpObj {
			res[varName] = varValue.(string)
		}
	}

	return
}

func substituteVariablesInTargetPaths(targetPaths []string, variablesSet map[string]string) {
	for _, targetPath := range targetPaths {
		var p = pth.Path{}
		p.SetPath(targetPath)

		if p.Exists() {
			var f = file.File{}
			if p.IsFile() {
				f.SetPath(targetPath)
				_logger.Infof("-> %s", f.GetPath())
				f.ReadContent().ReplaceVarsSetPlaceholder(variablesSet, "{{", "}}").Save()
			}
			if p.IsDirectory() {
				files := file.FindFilesRecursively(targetPath)
				for _, file := range files {
					var name = file.GetPath()
					f.SetPath(name)
					_logger.Infof("-> %s", f.GetPath())
					f.ReadContent().ReplaceVarsSetPlaceholder(variablesSet, "{{", "}}").Save()
				}
			}
		}
	}
}

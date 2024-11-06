package cicd_envsubst

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/file"
	"cicd_envsubst/utils/logger"
	pth "cicd_envsubst/utils/path"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
)

var _logger = logger.New()

// Uncomment the main function when you want to run the application
// func main(){
// 	Execute()
// }

var envsubstOpts struct {
	Prefix    string   `short:"p" long:"prefix" description:"Placeholder prefix" required:"false" default:"{{"`
	Suffix    string   `short:"x" long:"suffix" description:"Placeholder suffix" required:"false" default:"}}"`
	RegexMask string   `short:"m" long:"regex-mask" description:"Placeholder regex mask" required:"false" default:"[A-Z_0-9]+"`
	EnvFile   string   `short:"e" long:"env-file" description:"Source of environment variables" required:"false"`
	Sets      []string `short:"s" long:"set" description:"Variables set name" required:"false"`
}

func Execute() {
	paths := ReadCliOptionsEnvsubst()
	variablesSet := loadAllVariables()
	ProcessPaths(paths, variablesSet)
}

func ReadCliOptionsEnvsubst() (args []string) {
	args, err := flags.Parse(&envsubstOpts)

	if err != nil {
		_logger.Fatalf("Unable to read CLI options: %v", err)
	}

	return
}

func loadAllVariables() map[string]string {
	variables := make(map[string]string)

	// 1. Load variables from env file (lowest priority)
	if envsubstOpts.EnvFile != "" {
		_logger.Infof("Loading environment variables from file: %s", envsubstOpts.EnvFile)
		envFileVars := loadEnvFile(envsubstOpts.EnvFile)
		for k, v := range envFileVars {
			variables[k] = v
		}
	}

	// 2. Load variables from sets (middle priority)
	if len(envsubstOpts.Sets) > 0 {
		_logger.Infof("Loading variables from sets: %v", strings.Join(envsubstOpts.Sets, ", "))
		setVars := envsubstReadVariablesFromSets()
		for k, v := range setVars {
			variables[k] = v // Sets override env file variables
		}
	}

	// 3. Load existing environment variables (highest priority)
	_logger.Infof("Loading existing environment variables")
	envVars := loadEnvironmentVariables()
	for k, v := range envVars {
		variables[k] = v // Environment variables override sets and env file variables
	}

	return variables
}

func loadEnvFile(envFilePath string) map[string]string {
	variables := make(map[string]string)
	f := file.File{}
	f.SetPath(envFilePath)

	if f.IsExists() {
		content := f.ReadContent().GetContent()
		contentString := string(content)
		// Updated regex to ignore lines starting with '#' entirely
		envVarRegex := regexp.MustCompile(`^[ \t]*([a-zA-Z_][a-zA-Z0-9_]*)=(.*)$`)
		matches := envVarRegex.FindAllStringSubmatch(contentString, -1)
		for _, match := range matches {
			if len(match) == 3 {
				key := match[1]
				value := strings.Trim(strings.Trim(match[2], `"`), `'`)
				variables[key] = value
			}
		}
	} else {
		_logger.Warnf("Env file does not exist: %s", envFilePath)
	}

	return variables
}

func loadEnvironmentVariables() map[string]string {
	variables := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			variables[parts[0]] = parts[1]
		}
	}
	return variables
}

func envsubstReadVariablesFromSets() map[string]string {
	res := map[string]string{}

	for _, setName := range envsubstOpts.Sets {
		variablesSetYaml := ""

		ev := env_var.EnvVar{}
		if ev.
			SetName(setName).
			IsExist() {
			variablesSetYaml = ev.Value()
		} else {
			_logger.Fatalf("Unable to locate variables set: %s", setName)
			os.Exit(1)
		}

		tmpObj := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(variablesSetYaml), &tmpObj)
		if err != nil {
			_logger.Fatalf("Unable to unmarshal YAML content of variables set '%s': %v", setName, err)
		}

		for varName, varValue := range tmpObj {
			// Ensure the value is a string
			if strVal, ok := varValue.(string); ok {
				res[varName] = strVal
			} else {
				_logger.Warnf("Variable '%s' in set '%s' is not a string. Skipping.", varName, setName)
			}
		}
	}
	// let's print the names of variables
	// _logger.Infof("Variables from %v loaded successfully, extracted variables:", envsubstOpts.Sets)
	// for k, _ := range res {
	// 	_logger.Infof("-> %v", k)
	// }
	// _logger.Info("Variables set loaded successfully")
	// _logger.Infof("Variables set: %v", res)
	return res
}

func ProcessPaths(paths []string, variablesSet map[string]string) {
	_logger.Infof("Replacing variables in the following locations: { %v }", strings.Join(paths, ", "))
	_logger.Infof("Using prefix '%s' and suffix '%s'", envsubstOpts.Prefix, envsubstOpts.Suffix)
	_logger.Infof("PROCESSING")
	for _, path := range paths {
		processPath(path, variablesSet)
	}
	_logger.Infof("DONE")
}

func processPath(path string, variablesSet map[string]string) {
	replaceVariablesInPath(path, variablesSet)
}

func replaceVariablesInPath(path string, variablesSet map[string]string) {
	p := pth.New(path)

	if p.IsExistAndFile() {
		replaceEnvironmentVariablesInFile(path, variablesSet)
	} else if p.IsExistAndDirectory() {
		replaceEnvironmentVariablesInDirectory(path, variablesSet)
	} else {
		_logger.Warnf("Path does not exist or is not a file/directory: %s", path)
	}
}

func replaceEnvironmentVariablesInDirectory(path string, variablesSet map[string]string) {
	files := file.FindFilesRecursively(path)
	for _, file := range files {
		name := file.GetPath()
		replaceEnvironmentVariablesInFile(name, variablesSet)
	}
}

func replaceEnvironmentVariablesInFile(path string, variablesSet map[string]string) {
	if isFile(path) {
		var f = file.File{}
		f.SetPath(path)
		_logger.Infof("-> %s", f.GetPath())
		f.ReadContent().
			ReplaceVarsSetPlaceholder(variablesSet, envsubstOpts.Prefix, envsubstOpts.Suffix).
			// ReplaceEnvVarsPlaceholderByExplicitRegex(envsubstOpts.Prefix, envsubstOpts.Suffix, envsubstOpts.RegexMask).
			Save()
	}
}

func isFile(path string) bool {
	return pth.New(path).IsExistAndFile()
}

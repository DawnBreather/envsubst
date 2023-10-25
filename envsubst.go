package cicd_envsubst

import (
	"cicd_envsubst/utils/file"
	"cicd_envsubst/utils/logger"
	pth "cicd_envsubst/utils/path"
	"github.com/jessevdk/go-flags"
	"os"
	"regexp"
	"strings"
)

var _logger = logger.New()

//func main(){
//	Execute()
//}

var envsubstOpts struct {
	Prefix    string `short:"p" long:"prefix" description:"Placeholder prefix" required:"false" default:"{{"`
	Suffix    string `short:"s" long:"suffix" description:"Placeholder suffix" required:"false" default:"}}"`
	RegexMask string `short:"m" long:"regex-mask" description:"Placeholder regex mask" required:"false" default:"[A-Z_0-9]*"`
	EnvFile   string `short:"e" long:"env-file" description:"Source of environment variables" required:"false"`
}

func Execute() {
	paths := ReadCliOptionsEnvsubst()
	ProcessPaths(paths)
}

func setEnvironmentVariables() {
	f := file.File{}
	f.SetPath(envsubstOpts.EnvFile)

	if f.IsExists() {
		content := f.ReadContent().GetContent()
		contentString := string(content)
		envVarRegex := regexp.MustCompile(`[#]*[a-zA-Z_][a-zA-Z0-9_]*=.*`)
		matches := envVarRegex.FindAllString(contentString, -1)
		for _, match := range matches {
			envVarKeyValue := strings.SplitN(match, "=", 2)
			if !strings.Contains(envVarKeyValue[0], "#") { // if not commented
				os.Setenv(envVarKeyValue[0], strings.Trim(strings.Trim(envVarKeyValue[1], `"`), `'`))
			}
		}
	}
}

func ProcessPaths(paths []string) {
	_logger.Infof("Replacing environment variables in the following locations: { %v }", strings.Join(paths, ", "))
	_logger.Infof("PROCESSING")
	for _, path := range paths {
		processPath(path)
	}
	_logger.Infof("DONE")
}

func processPath(path string) {
	setEnvironmentVariables()
	replaceEnvironmentVariablesInFile(path)
	replaceEnvironmentVariablesInDirectory(path)

}

func isFile(path string) bool {
	return pth.New(path).IsExistAndFile()
}

func isDirectory(path string) bool {
	return pth.New(path).IsExistAndDirectory()
}

func replaceEnvironmentVariablesInDirectory(path string) {
	if isDirectory(path) {
		files := file.FindFilesRecursively(path)
		for _, file := range files {
			var name = file.GetPath()
			replaceEnvironmentVariablesInFile(name)
		}
	}
}

func replaceEnvironmentVariablesInFile(path string) {
	if isFile(path) {
		var f = file.File{}
		f.SetPath(path)
		_logger.Infof("-> %s", f.GetPath())
		f.ReadContent().
			ReplaceEnvVarsPlaceholderByExplicitRegex(envsubstOpts.Prefix, envsubstOpts.Suffix, envsubstOpts.RegexMask).
			Save()
	}
}

func ReadCliOptionsEnvsubst() (args []string) {
	args, err := flags.Parse(&envsubstOpts)

	if err != nil {
		_logger.Fatalf("Unable to read cli options: %v", err)
	}

	return
}

package cicd_envsubst

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/file"
	"cicd_envsubst/utils/path"
	"cicd_envsubst/utils/placeholder"
	"fmt"
	"github.com/thoas/go-funk"
	"os"
	"sort"
	"strings"
)

func ExecuteEnvMaker() {

	envDestination := os.Args[1]
	pathsForProcessing := os.Args[2:]

	var p = path.Path{}
	var f = file.File{}
	var ev = env_var.EnvVar{}

	var fDestination = file.File{}
	fDestination.SetPath(envDestination)
	var destinationFileContent string

	var ph = placeholder.New("{{", "}}", `[A-Z0-9_]*`)

	var allMatches []string

	_logger.Infof("Compiling env file { %s } based on the placeholders found in: { %v }", envDestination, strings.Join(pathsForProcessing, ", "))
	_logger.Infof("PROCESSING")

	for _, pfp := range pathsForProcessing {

		p.SetPath(pfp)
		if p.IsFile() {
			_logger.Infof("-> %s", p.GetPath())
			content := f.SetPath(p.GetPath()).ReadContent().GetContent()
			allMatches = append(allMatches, ph.FindAllStringMatchesWithoutSuffixesAndPrefixes(string(content))...)
		}

		if p.IsDirectory() {
			for _, fl := range file.FindFilesRecursively(p.GetPath()) {
				_logger.Infof("-> %s", fl.GetPath())
				content := fl.ReadContent().GetContent()
				allMatches = append(allMatches, ph.FindAllStringMatchesWithoutSuffixesAndPrefixes(string(content))...)
			}
		}
	}

	sort.Strings(allMatches)
	for _, envVarPlaceholderName := range funk.UniqString(allMatches) {
		ev.SetName(envVarPlaceholderName)
		var val string
		if ev.IsExist() {
			val = ev.Value()
		} else {
			val = "{{" + envVarPlaceholderName + "}}"
		}
		destinationFileContent += fmt.Sprintf("%s=%s\n", envVarPlaceholderName, val)
	}

	fDestination.
		SetContent([]byte(destinationFileContent)).
		Save()

	_logger.Infof("DONE")

}

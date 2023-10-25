package cicd_envsubst

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/file"
	path2 "cicd_envsubst/utils/path"
	"cicd_envsubst/utils/placeholder"
	"encoding/base64"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/jessevdk/go-flags"
	"github.com/thoas/go-funk"
	"os"
	"strings"
)

var set2SecretOpts = set2SecretOptsT{}

type set2SecretOptsT struct {
  Sets       []string `short:"s" long:"set" description:"variable set name" required:"true"`
  SecretName string   `short:"n" long:"secret-name" description:"k8s secret name" required:"true"`
  RefFiles   []string `short:"r" long:"ref" description:"reference files with placeholders" required:"false"`

  jointVariablesSet map[string]string
  targetPaths       []string
}

func ExecuteSet2SecretOpts() {

  //os.Setenv("SECRET", playground.TestEnvVar)

  set2SecretOpts.ReadCliOptions()

  _logger.Infof("Generating k8s secret/opaque based on variables sets { %v } to the following locations: { %v }", strings.Join(set2SecretOpts.Sets, ", "), strings.Join(set2SecretOpts.targetPaths, ", "))

  set2SecretOpts.
    ReadVariablesFromSets().
    CompareVariablesToRefs().
    GenerateK8sSecretManifest()

  _logger.Infof("DONE")

}

func (s *set2SecretOptsT) ReadCliOptions() *set2SecretOptsT {
  var err error
  s.targetPaths, err = flags.Parse(s)

  if err != nil {
    _logger.Fatalf("Unable to read cli options: %v", err)
  }

  return s
}

func (s *set2SecretOptsT) ReadVariablesFromSets() *set2SecretOptsT {

  if s.jointVariablesSet == nil {
    s.jointVariablesSet = map[string]string{}
  }

  for _, setName := range s.Sets {

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
      s.jointVariablesSet[varName] = fmt.Sprintf("%v", varValue)
    }
  }

  return s
}

func (s *set2SecretOptsT) GenerateK8sSecretManifest() {

  _logger.Infof("SAVING")

  secretName := s.SecretName

  manifest := map[string]interface{}{}
  manifest["apiVersion"] = "v1"
  manifest["kind"] = "Secret"
  manifest["metadata"] = map[string]string{"name": secretName}
  manifest["type"] = "Opaque"
  manifest["data"] = map[string]string{}

  for name, value := range s.jointVariablesSet {
    manifest["data"].(map[string]string)[name] = base64.StdEncoding.EncodeToString([]byte(value))
  }

  yamlBytes, err := yaml.Marshal(manifest)
  if err != nil {
    _logger.Fatalf("Unable to marshal k8s manifest object into YAML: %v", err)
  }

  for _, path := range s.targetPaths {
    f := file.File{}
    f.
      SetPath(path).
      SetContent(yamlBytes).
      Save()
    _logger.Infof("-> %s", path)
  }

}

func (s *set2SecretOptsT) CompareVariablesToRefs() *set2SecretOptsT {

  var jointVariableNamesFromRefFiles = map[string][]string{}

  for _, path := range s.RefFiles {
    p := path2.Path{}
    p.SetPath(path)
    if p.Exists() {
      if p.IsFile() {

        if jointVariableNamesFromRefFiles[path] == nil {
          jointVariableNamesFromRefFiles[path] = []string{}
        }

        ph := placeholder.New("{{", "}}", "[A-Z0-9_]*")
        f := file.File{}
        fContent := f.SetPath(path).ReadContent().GetContent()
        jointVariableNamesFromRefFiles[path] = append(jointVariableNamesFromRefFiles[path], ph.FindAllStringMatchesWithoutSuffixesAndPrefixes(string(fContent))...)
      }
    }
  }

  if len(jointVariableNamesFromRefFiles) > 0 {
    _logger.Infof("VERIFYING")
  }

  var notFoundError = false
  for path, vars := range jointVariableNamesFromRefFiles {
    _logger.Infof("-> %s", path)
    for _, name := range vars {
      if !funk.Contains(s.jointVariablesSet, name) {
        _logger.Warnf("{ %s } not found", name)
        notFoundError = true
      }
    }
  }

  if len(s.RefFiles) > 0 {
    _logger.Infof("-> %s", "sets for redundancy")
    for name, _ := range s.jointVariablesSet {
      if !funk.Contains(jointVariableNamesFromRefFiles, func(fileName string, varNames []string) bool {
        return funk.Contains(varNames, name) // or `name == "Florent"` for the value type
      }) {
        _logger.Warnf("{ %s } is redundant", name)
      }
    }
  }

  if notFoundError {
    os.Exit(1)
  }

  return s

}

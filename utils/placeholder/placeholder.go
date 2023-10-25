package placeholder

import (
	"cicd_envsubst/utils/env_var"
	"cicd_envsubst/utils/logger"
	"regexp"
	"strings"
)

var _logger = logger.New()

func New(prefix, suffix, core string) *Placeholder {
	return &Placeholder{
		prefix:        prefix,
		suffix:        suffix,
		coreRegexMask: core,
	}
}

type Placeholder struct {
	prefix        string
	suffix        string
	coreRegexMask string
	varsSet       map[string]string
}

func (p *Placeholder) SetPrefix(prefix string) *Placeholder {
	p.prefix = prefix
	return p
}
func (p *Placeholder) SetSuffix(suffix string) *Placeholder {
	p.suffix = suffix
	return p
}
func (p *Placeholder) SetCoreRegexpMask(coreRegexpMask string) *Placeholder {
	p.coreRegexMask = coreRegexpMask
	return p
}
func (p *Placeholder) SetVarsSet(varsSet map[string]string) *Placeholder {
	p.varsSet = varsSet
	return p
}
func (p *Placeholder) GetPrefix() string {
	return p.prefix
}
func (p *Placeholder) GetSuffix() string {
	return p.suffix
}
func (p *Placeholder) GetCoreRegexpMask() string {
	return p.coreRegexMask
}

func (p *Placeholder) getRegexPattern() *regexp.Regexp {
	regexMask := p.prefix + p.coreRegexMask + p.suffix
	pattern, err := regexp.Compile(regexMask)
	if err != nil {
		_logger.Errorf("Unable to compile regexp mask '%s': %v", regexMask, err)
	}

	return pattern
}

func (p *Placeholder) FindAllStringMatches(string string) []string {
	pattern := p.getRegexPattern()
	return pattern.FindAllString(string, -1)
}

func (p *Placeholder) FindAllStringMatchesWithoutSuffixesAndPrefixes(string string) (res []string) {
	pattern := p.getRegexPattern()
	matches := pattern.FindAllString(string, -1)
	for _, m := range matches {

		m1 := strings.TrimPrefix(m, p.prefix)
		m2 := strings.TrimSuffix(m1, p.suffix)

		res = append(res, m2)
	}

	return res
}

func (p *Placeholder) ReplacePlaceholdersWithEnvVars(baseString string) string {

	var res = baseString

	for _, match := range p.FindAllStringMatches(baseString) {
		envVarName := strings.ReplaceAll(match, p.prefix, "")
		envVarName = strings.ReplaceAll(envVarName, p.suffix, "")
		envVar := env_var.New(envVarName)
		envVarValue := envVar.Value()
		res = strings.ReplaceAll(res, match, envVarValue)
		if !envVar.IsExist() {
			_logger.Warnf("Environment variable not found: { %s }", envVarName)
		}
	}

	return res
}

func (p *Placeholder) ReplacePlaceholdersWithVarsSet(baseString string) string {

	var res = baseString

	for _, match := range p.FindAllStringMatches(baseString) {
		varName := strings.ReplaceAll(match, p.prefix, "")
		varName = strings.ReplaceAll(varName, p.suffix, "")

		if value, ok := p.varsSet[varName]; ok {
			res = strings.ReplaceAll(res, match, value)
		} else {
			_logger.Warnf("Variable not found: { %s }", varName)
		}
	}

	return res
}

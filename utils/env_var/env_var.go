package env_var

import (
	"os"
)

const(
	REGEX_MASK = "[A-Z_0-9]*"
)

func New(name string) *EnvVar{
	return &EnvVar{
		name: name,
	}
}

func (e *EnvVar) SetName(name string) *EnvVar{
	e.name = name
	return e
}

type EnvVar struct{
	name string
}

func (e EnvVar) IsExist() bool{
	_, present := os.LookupEnv(e.name)
	return present
}

func (e EnvVar) Value() string{
	val, _ := os.LookupEnv(e.name)
	return val
}
package gotdot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var Config GotDot

func init() {
	Config = NewGotDot()
}

type IConfig interface {
	Get(string) (*IConfigVar, error)
	Load()
}

type IConfigVar interface {
	Cast(string) interface{}
	Export()
}

type GotDot struct {
	Path    string
	mapping map[string]*ConfigVar
}

func NewGotDot() GotDot {
	var c GotDot
	c.Path = "./.env"
	c.mapping = make(map[string]*ConfigVar)
	return c
}

func (c *GotDot) Get(key string) (*ConfigVar, error) {
	value, exists := c.mapping[key]
	if !exists {
		err := fmt.Errorf("Env Variable with key '%s' doesn't exist", key)
		return nil, err
	}
	return value, nil
}

func (c *GotDot) Load() {
	dotfile, err := os.Open(c.Path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(dotfile)

	for scanner.Scan() {
		k, v := parseEnvFileLine(scanner.Text())
		if k == "" {
			continue
		}
		os.Setenv(k, v)
		configVar := NewConfigVar(k, v)
		c.mapping[k] = &configVar
	}
}

type ConfigVar struct {
	Key         string
	initValue   string
	castedValue interface{}
}

func NewConfigVar(key, value string) ConfigVar {
	return ConfigVar{
		Key:         key,
		initValue:   value,
		castedValue: nil,
	}
}

func (v *ConfigVar) Cast(t string) *ConfigVar {
	var err error
	switch t {
	case "string":
		v.castedValue = v.initValue
	case "int":
		v.castedValue, err = strconv.ParseInt(v.initValue, 10, 64)
	case "float":
		v.castedValue, err = strconv.ParseFloat(v.initValue, 64)
	case "bool":
		v.castedValue, err = strconv.ParseBool(v.initValue)
	default:
		err = fmt.Errorf("Unknown casting type - %s", t)
	}

	if err != nil {
		log.Fatalf("Error in casting ConfigVar - %s", err)
	}

	return v
}

func (v *ConfigVar) Import() interface{} {
	if v.castedValue != nil {
		return v.castedValue
	}
	return v.initValue
}

func parseEnvFileLine(line string) (k, v string) {
	line = trimCommentFromLine(line) // trim comment
	kv := strings.SplitN(line, "=", 2)
	if len(kv) != 2 {
		return "", ""
	}
	k, v = strings.ToUpper(kv[0]), kv[1]
	return strings.TrimSpace(k), strings.TrimSpace(v)
}

func trimCommentFromLine(line string) string {
	if idx := strings.Index(line, "#"); idx != -1 {
		return line[:idx]
	}
	return line
}

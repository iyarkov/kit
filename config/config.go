package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Password struct {
	value string
}

func NewPassword(value string) Password {
	return Password{
		value: value,
	}
}

func (p *Password) Value() string {
	return p.value
}

type DbConfig struct {
	Host     string
	Port     uint16
	User     string
	Password Password
	DbName   string
}

func Read(val any) error {
	if err := fromFile(val); err != nil {
		return err
	}
	if err := fromEnv(val); err != nil {
		return err
	}
	return fromCommandLine(val)
}

func fromFile(val any) error {
	fileName := getFlag("-f", os.Args)
	if fileName != "" {
		absFilePath, err := filepath.Abs(fileName)
		if err != nil {
			return fmt.Errorf("can not get absolute file path for file %s : %w", fileName, err)
		}
		fmt.Printf("Reading configuration from file %s\n", absFilePath)
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("can not open config file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Printf("unexpected error while closing config file %v", err)
			}
		}(file)
		decoder := json.NewDecoder(file)
		return decoder.Decode(val)
	}
	return nil
}

func fromEnv(val any) error {
	cfgMap := make(map[string]string, 0)
	for _, a := range os.Environ() {
		idx := strings.IndexRune(a, '=')
		if idx != -1 && idx != len(a)-1 {
			key := normalizeEnvKey(a[:idx])
			value := a[idx+1:]
			cfgMap[key] = value
		}
	}
	if len(cfgMap) > 0 {
		fmt.Println("Reading configuration from environment variables")
		return updateConfig(val, &cfgMap)
	}
	return nil
}

func normalizeEnvKey(key string) string {
	keyRunes := []rune(key)
	mustBeUpper := true
	for i, r := range keyRunes {
		if r == '_' {
			keyRunes[i] = '.'
			mustBeUpper = true
		} else if mustBeUpper {
			keyRunes[i] = unicode.ToUpper(r)
			mustBeUpper = false
		} else {
			keyRunes[i] = unicode.ToLower(r)
		}

	}
	return string(keyRunes)
}

func fromCommandLine(val any) error {
	cfgMap := make(map[string]string, 0)
	for _, a := range os.Args {
		idx := strings.IndexRune(a, '=')
		if idx != -1 && idx != len(a)-1 {
			key := a[:idx]
			value := a[idx+1:]
			cfgMap[key] = value
		}
	}
	if len(cfgMap) > 0 {
		fmt.Println("Reading configuration from command line")
		return updateConfig(val, &cfgMap)
	}
	return nil
}

func updateConfig(conf any, params *map[string]string) error {
	for key, value := range *params {
		if err := updateConfigField(reflect.ValueOf(conf), key, value); err != nil {
			return fmt.Errorf("cfg error: %s %w", key, err)
		}
	}
	return nil
}

func updateConfigField(conf reflect.Value, key, value string) error {
	var fieldName = key
	var subpath string
	if idx := strings.IndexRune(key, '.'); idx != -1 {
		fieldName = key[:idx]
		subpath = key[idx+1:]
	}
	if conf.Kind() == reflect.Pointer {
		conf = conf.Elem()
	}
	field := conf.FieldByName(fieldName)
	if !field.IsValid() {
		// Subfield does not exist
		return nil
	}
	if subpath == "" {
		// Set the value of the field
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err == nil {
				field.SetBool(boolVal)
			} else {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				field.SetInt(intVal)
			} else {
				return err
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := strconv.ParseUint(value, 10, 64)
			if err == nil {
				field.SetUint(intVal)
			} else {
				return err
			}
		case reflect.Float32, reflect.Float64:
			floatVal, err := strconv.ParseFloat(value, 64)
			if err == nil {
				field.SetFloat(floatVal)
			} else {
				return err
			}
		}
	}

	if field.Type() == reflect.TypeOf(Password{}) {
		field.Set(reflect.ValueOf(NewPassword(value)))
	} else if field.Kind() == reflect.Struct {
		if err := updateConfigField(field, subpath, value); err != nil {
			return err
		}
	}

	return nil
}

func getFlag(flag string, args []string) string {
	for _, a := range args {
		if strings.HasPrefix(a, flag) {
			return a[len(flag):]
		}
	}
	return ""
}

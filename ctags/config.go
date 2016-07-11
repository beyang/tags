package ctags

import (
	"os/exec"
	"path/filepath"
	"strings"
)

type Lang struct {
	// Name is the name of the language
	Name string

	// Files is the set of filenames mapped to this language by ctags
	Files []string

	// Exts is the set of file extensions mapped to this language by ctags
	Exts []string
}

type Config struct {
	Langs      []Lang
	extToLang  map[string]string
	fileToLang map[string]string
}

func getConfig() (*Config, error) {
	out, err := exec.Command("ctags", "--list-maps").Output()
	if err != nil {
		return nil, err
	}

	// Parse languages from ctags output
	var config Config
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 || line == "CTagsSelfTest" {
			continue
		}
		fields := strings.Fields(line)
		lang := Lang{Name: fields[0]}
		for _, field := range fields[1:] {
			if strings.HasPrefix(field, "*.") {
				lang.Exts = append(lang.Exts, field[1:])
			} else {
				lang.Files = append(lang.Files, field)
			}
		}
		config.Langs = append(config.Langs, lang)
	}

	// Create lookup tables
	config.extToLang = make(map[string]string)
	config.fileToLang = make(map[string]string)
	for _, lang := range config.Langs {
		for _, file := range lang.Files {
			config.fileToLang[file] = lang.Name
		}
		for _, ext := range lang.Exts {
			config.extToLang[ext] = lang.Name
		}
	}

	return &config, nil
}

func (c *Config) Lang(filename string) string {
	if l, exists := c.fileToLang[filename]; exists {
		return l
	}
	return c.extToLang[filepath.Ext(filename)]
}

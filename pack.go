package pack

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v3"
)

type PackConfig struct {
	Name     string            `yaml:"name"`
	Version  string            `yaml:"version"`
	Values   map[string]string `yaml:"values"`
	Template string            `yaml:"template"`
}

func LoadPack(path string) (*PackConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pack file: %w", err)
	}
	var cfg PackConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse pack file: %w", err)
	}
	return &cfg, nil
}

func Render(cfg *PackConfig, overrides map[string]string) ([]byte, error) {
	// Merge values: overrides take precedence
	merged := make(map[string]string)
	for k, v := range cfg.Values {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}

	tplData, err := os.ReadFile(cfg.Template)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	tmpl, err := template.New("stack").Parse(string(tplData))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, merged); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}
	return buf.Bytes(), nil
}

func Install(name string, rendered []byte) error {
	cmd := exec.Command("docker", "stack", "deploy", "-c", "-", "--with-registry-auth", name)
	cmd.Stdin = bytes.NewReader(rendered)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Uninstall(name string) error {
	cmd := exec.Command("docker", "stack", "rm", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package lab

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ProfilesDirName    = ".xoctx"
	ProfilesDirEnv     = "XOCTX_DIR"
	CurrentLabFileName = "current_lab"
)

type Profile struct {
	Name string
	Path string
	Vars map[string]string
}

func ProfilesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot determine home directory: %v", err))
	}

	if dir := strings.TrimSpace(os.Getenv(ProfilesDirEnv)); dir != "" {
		return dir
	}
	return filepath.Join(home, ProfilesDirName)
}

func CurrentLabFile() string {
	return filepath.Join(ProfilesDir(), CurrentLabFileName)
}

func ListLabs() ([]string, error) {
	dir := ProfilesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read profiles dir: %w", err)
	}

	var labs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".env") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".env")
		labs = append(labs, name)
	}
	return labs, nil
}

func LoadProfile(name string) (*Profile, error) {
	path := filepath.Join(ProfilesDir(), name+".env")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("lab %q does not exist", name)
	}

	p := &Profile{
		Name: name,
		Path: path,
		Vars: make(map[string]string),
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open profile: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		p.Vars[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read profile: %w", err)
	}

	return p, nil
}

func CurrentLab() (string, error) {
	path := CurrentLabFile()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read current lab: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func SetCurrentLab(name string) error {
	path := CurrentLabFile()
	if err := os.MkdirAll(ProfilesDir(), 0755); err != nil {
		return fmt.Errorf("create profiles dir: %w", err)
	}
	return os.WriteFile(path, []byte(name), 0644)
}

func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	return strings.Contains(upper, "TOKEN") || strings.Contains(upper, "PASSWORD")
}

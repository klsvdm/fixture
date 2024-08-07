package fixtures

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type Fixture struct {
	data map[string][]byte
}

func Load(path string) (Fixture, error) {
	data, err := traverseDir(path, "")
	if err != nil {
		return Fixture{}, fmt.Errorf("failed to load fixtures: %w", err)
	}

	return Fixture{data: data}, nil
}

func MustLoad(path string) Fixture {
	f, err := Load(path)
	if err != nil {
		panic(err)
	}

	return f
}

func traverseDir(path, prefix string) (map[string][]byte, error) {
	files, err := os.ReadDir(filepath.Join(path, prefix))
	if err != nil {
		return nil, fmt.Errorf("failed to read fixtures folder '%s': %w", path, err)
	}

	data := make(map[string][]byte, len(files))

	for _, file := range files {
		var (
			ext      = filepath.Ext(file.Name())
			filePath = filepath.Join(path, prefix, file.Name())
		)

		if file.IsDir() {
			dirData, err := traverseDir(path, filepath.Join(prefix, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read fixtures folder '%s': %w", path, err)
			}

			maps.Copy(data, dirData)

			continue
		}

		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		name := prefix + "/" + strings.TrimSuffix(file.Name(), ext)
		if prefix == "" {
			name = name[1:]
		}

		data[name] = content
	}

	return data, nil
}

func Get[T any](t *testing.T, f Fixture, name string, opts ...Option[T]) T {
	var data T

	if err := f.get(name, &data); err != nil {
		t.Fatalf(err.Error())
	}

	o := applyOptions(opts)

	if o.editor != nil {
		o.editor(&data)
	}

	return data
}

func GetList[T any](t *testing.T, f Fixture, name string, opts ...Option[T]) []T {
	data := make([]T, 0)

	if err := f.get(name, &data); err != nil {
		t.Fatalf(err.Error())
	}

	o := applyOptions(opts)

	if o.editor != nil {
		for i := range data {
			o.editor(&data[i])
		}
	}

	return data
}

func GetMap[T any](t *testing.T, f Fixture, name string) map[string]T {
	data := make(map[string]T)

	if err := f.get(name, &data); err != nil {
		t.Fatalf(err.Error())
	}

	return data
}

func (f *Fixture) get(name string, value any) error {
	content, ok := f.data[name]
	if !ok {
		return fmt.Errorf("fixture '%s' not found", name)
	}

	if err := yaml.Unmarshal(content, value); err != nil {
		return fmt.Errorf("failed to unmarshal fixture '%s': %s", name, err)
	}

	return nil
}

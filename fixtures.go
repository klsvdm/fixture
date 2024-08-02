package fixtures

import (
	"fmt"
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
	files, err := os.ReadDir(path)
	if err != nil {
		return Fixture{}, fmt.Errorf("failed to read fixtures folder '%s': %w", path, err)
	}

	data := make(map[string][]byte, len(files))

	for _, file := range files {
		var (
			ext      = filepath.Ext(file.Name())
			filePath = filepath.Join(path, file.Name())
		)

		if file.IsDir() || (ext != ".yaml" && ext != ".yml") {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ext)
		data[name] = content
	}

	return Fixture{data: data}, nil
}

func MustLoad(path string) Fixture {
	fixture, err := Load(path)
	if err != nil {
		panic(err)
	}

	return fixture
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

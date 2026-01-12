package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// mockConverter is a test implementation of Converter.
type mockConverter struct {
	name       string
	outputFile string
	convertErr error
	content    []byte
}

func (m *mockConverter) Name() string {
	return m.name
}

func (m *mockConverter) OutputFileName() string {
	return m.outputFile
}

func (m *mockConverter) Convert(ctx *Context) ([]byte, error) {
	if m.convertErr != nil {
		return nil, m.convertErr
	}
	return m.content, nil
}

func (m *mockConverter) WriteFile(ctx *Context, path string) error {
	data, err := m.Convert(ctx)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, DefaultFileMode)
}

func TestConverterRegistry(t *testing.T) {
	t.Run("Register and Get", func(t *testing.T) {
		registry := NewConverterRegistry()
		mock := &mockConverter{name: "test", outputFile: "TEST.md"}

		registry.Register(mock)

		got, ok := registry.Get("test")
		if !ok {
			t.Fatal("expected to find registered converter")
		}
		if got.Name() != "test" {
			t.Errorf("expected name 'test', got '%s'", got.Name())
		}
	})

	t.Run("Get nonexistent", func(t *testing.T) {
		registry := NewConverterRegistry()

		_, ok := registry.Get("nonexistent")
		if ok {
			t.Error("should not find nonexistent converter")
		}
	})

	t.Run("Names", func(t *testing.T) {
		registry := NewConverterRegistry()
		registry.Register(&mockConverter{name: "a", outputFile: "A.md"})
		registry.Register(&mockConverter{name: "b", outputFile: "B.md"})

		names := registry.Names()
		if len(names) != 2 {
			t.Fatalf("expected 2 names, got %d", len(names))
		}
	})
}

func TestConverterRegistryConvert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		registry := NewConverterRegistry()
		registry.Register(&mockConverter{
			name:       "test",
			outputFile: "TEST.md",
			content:    []byte("# Test"),
		})

		ctx := NewContext("test-project")
		data, err := registry.Convert(ctx, "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "# Test" {
			t.Errorf("expected '# Test', got '%s'", string(data))
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		registry := NewConverterRegistry()

		ctx := NewContext("test")
		_, err := registry.Convert(ctx, "unsupported")
		if err == nil {
			t.Fatal("expected error for unsupported format")
		}
		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Errorf("expected ConversionError, got %T", err)
		}
	})
}

func TestConverterRegistryWriteFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		registry := NewConverterRegistry()
		registry.Register(&mockConverter{
			name:       "test",
			outputFile: "TEST.md",
			content:    []byte("# Test Content"),
		})

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "output.md")

		ctx := NewContext("test-project")
		err := registry.WriteFile(ctx, "test", path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read written file: %v", err)
		}
		if string(data) != "# Test Content" {
			t.Errorf("expected '# Test Content', got '%s'", string(data))
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		registry := NewConverterRegistry()

		ctx := NewContext("test")
		err := registry.WriteFile(ctx, "unsupported", "/tmp/test.md")
		if err == nil {
			t.Fatal("expected error for unsupported format")
		}
	})
}

func TestConverterRegistryGenerateAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		registry := NewConverterRegistry()
		registry.Register(&mockConverter{
			name:       "test1",
			outputFile: "TEST1.md",
			content:    []byte("# Test 1"),
		})
		registry.Register(&mockConverter{
			name:       "test2",
			outputFile: "TEST2.md",
			content:    []byte("# Test 2"),
		})

		tmpDir := t.TempDir()
		ctx := NewContext("test-project")

		err := registry.GenerateAll(ctx, tmpDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check both files were created
		for _, file := range []string{"TEST1.md", "TEST2.md"} {
			path := filepath.Join(tmpDir, file)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("expected file %s to exist", path)
			}
		}
	})

	t.Run("with empty dir", func(t *testing.T) {
		registry := NewConverterRegistry()
		registry.Register(&mockConverter{
			name:       "test",
			outputFile: "TEST.md",
			content:    []byte("# Test"),
		})

		// Using empty dir should create file in current working dir
		tmpDir := t.TempDir()
		origDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := os.Chdir(origDir); err != nil {
				panic(err)
			}
		}()
		if err := os.Chdir(tmpDir); err != nil {
			panic(err)
		}

		ctx := NewContext("test-project")
		err = registry.GenerateAll(ctx, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat("TEST.md"); os.IsNotExist(err) {
			t.Error("expected TEST.md to exist")
		}
	})
}

func TestDefaultRegistry(t *testing.T) {
	// These use the global DefaultRegistry

	t.Run("RegisterConverter and GetConverter", func(t *testing.T) {
		mock := &mockConverter{name: "global-test", outputFile: "GLOBAL.md"}
		RegisterConverter(mock)

		got, ok := GetConverter("global-test")
		if !ok {
			t.Fatal("expected to find registered converter")
		}
		if got.Name() != "global-test" {
			t.Errorf("expected name 'global-test', got '%s'", got.Name())
		}
	})

	t.Run("ConvertTo", func(t *testing.T) {
		mock := &mockConverter{
			name:       "convert-test",
			outputFile: "CONVERT.md",
			content:    []byte("# Converted"),
		}
		RegisterConverter(mock)

		ctx := NewContext("test")
		data, err := ConvertTo(ctx, "convert-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "# Converted" {
			t.Errorf("expected '# Converted', got '%s'", string(data))
		}
	})
}

func TestBaseConverter(t *testing.T) {
	t.Run("Name and OutputFileName", func(t *testing.T) {
		bc := NewBaseConverter("test-conv", "TEST.md")

		if bc.Name() != "test-conv" {
			t.Errorf("expected name 'test-conv', got '%s'", bc.Name())
		}
		if bc.OutputFileName() != "TEST.md" {
			t.Errorf("expected output 'TEST.md', got '%s'", bc.OutputFileName())
		}
	})

	t.Run("WriteFileWithData", func(t *testing.T) {
		bc := NewBaseConverter("test", "TEST.md")
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "output.md")

		err := bc.WriteFileWithData([]byte("# Content"), path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != "# Content" {
			t.Errorf("expected '# Content', got '%s'", string(data))
		}
	})

	t.Run("WriteFileWithData error", func(t *testing.T) {
		bc := NewBaseConverter("test", "TEST.md")

		err := bc.WriteFileWithData([]byte("content"), "/nonexistent/dir/file.md")
		if err == nil {
			t.Fatal("expected error for invalid path")
		}
		var writeErr *WriteError
		if !errors.As(err, &writeErr) {
			t.Errorf("expected WriteError, got %T", err)
		}
	})

	t.Run("WriteFileWithDataAndMode", func(t *testing.T) {
		bc := NewBaseConverter("test", "TEST.md")
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "output.md")

		err := bc.WriteFileWithDataAndMode([]byte("# Content"), path, 0600)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("failed to stat file: %v", err)
		}
		// Check permissions (may vary by umask on some systems)
		if info.Mode().Perm() != 0600 {
			t.Logf("note: file permissions %v (expected 0600, may differ due to umask)", info.Mode().Perm())
		}
	})
}

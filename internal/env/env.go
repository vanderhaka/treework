package env

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CopyEnvFiles copies .env* files from src to dst, skipping files that already exist in dst.
func CopyEnvFiles(srcDir, dstDir string) ([]string, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}

	var copied []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, ".env") {
			continue
		}

		dstPath := filepath.Join(dstDir, name)
		if _, err := os.Stat(dstPath); err == nil {
			continue // already exists
		}

		srcPath := filepath.Join(srcDir, name)
		if err := copyFile(srcPath, dstPath); err != nil {
			return copied, err
		}
		copied = append(copied, name)
	}

	return copied, nil
}

func copyFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

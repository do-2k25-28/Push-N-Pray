package utils

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// ResolvePath resolves path relative to base and verifies the result stays
// within base. Absolute paths and traversals that escape base are rejected.
// Symlinks are resolved before the containment check to prevent symlink-based
// escapes within the workspace.
func ResolvePath(base, path string) (string, error) {
	var candidate string
	if filepath.IsAbs(path) {
		// Absolute paths would bypass the workspace root entirely.
		return "", fmt.Errorf("path %q must be relative to the workspace", path)
	}
	candidate = filepath.Join(base, path)

	// Resolve symlinks so a link pointing outside the workspace is caught.
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", fmt.Errorf("cannot resolve path %q: %w", path, err)
	}

	// Ensure the resolved path is inside base.
	cleanBase := filepath.Clean(base)
	if !strings.HasPrefix(resolved+string(filepath.Separator), cleanBase+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes workspace root", path)
	}

	return resolved, nil
}

func FindAvailablePort(defaultPort string) string {
	port, err := strconv.Atoi(defaultPort)
	if err != nil {
		port = 4000
	}

	for port <= 65535 {
		addr := fmt.Sprintf(":%d", port)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			_ = l.Close()
			return strconv.Itoa(port)
		}
		port++
	}

	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer func() {
			_ = l.Close()
		}()
		return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}

	return defaultPort
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds the global Pulse configuration.
type Config struct {
	PeerID       string `toml:"peer_id"`
	DefaultRelay string `toml:"default_relay"`
}

// BaseDir returns the root directory used by Pulse (~/.pulse).
func BaseDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".pulse")
	os.MkdirAll(dir, 0o700)
	return dir
}

// GroupsDir returns the directory where group configs are stored.
func GroupsDir() string {
	dir := filepath.Join(BaseDir(), "groups")
	os.MkdirAll(dir, 0o700)
	return dir
}

// IdentityKeyPath returns the path to the private key file.
func IdentityKeyPath() string {
	return filepath.Join(BaseDir(), "identity.key")
}

// ConfigPath returns the path to config.toml.
func ConfigPath() string {
	return filepath.Join(BaseDir(), "config.toml")
}

// PidDir returns the directory where listener PID files are stored.
func PidDir() string {
	dir := filepath.Join(BaseDir(), "pids")
	os.MkdirAll(dir, 0o700)
	return dir
}

// Load reads the config from disk. Returns zero-value Config if missing.
func Load() (Config, error) {
	var cfg Config
	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, fmt.Errorf("pulse not initialized: run 'pulse init' first")
	}
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}

// Save writes the config to disk.
func Save(cfg Config) error {
	f, err := os.OpenFile(ConfigPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// IsInitialized checks whether pulse init has been run.
func IsInitialized() bool {
	_, err := os.Stat(ConfigPath())
	return err == nil
}

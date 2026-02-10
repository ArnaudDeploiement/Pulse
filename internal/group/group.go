package group

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pulse/internal/config"

	"github.com/BurntSushi/toml"
)

// Group holds a group configuration.
type Group struct {
	Name     string   `toml:"name"`
	Protocol string   `toml:"protocol"`
	Relay    string   `toml:"relay"`
	Secret   string   `toml:"secret"`
	Members  []string `toml:"members"`
}

// Create creates a new group and writes it to disk.
func Create(name, relay string) (*Group, error) {
	if name == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}
	if relay == "" {
		return nil, fmt.Errorf("relay address cannot be empty")
	}

	path := groupPath(name)
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("group %q already exists", name)
	}

	protoID := make([]byte, 16)
	if _, err := rand.Read(protoID); err != nil {
		return nil, fmt.Errorf("generating protocol ID: %w", err)
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("generating secret: %w", err)
	}

	g := &Group{
		Name:     name,
		Protocol: "/pulse/" + base64.RawURLEncoding.EncodeToString(protoID) + "/2.0",
		Relay:    relay,
		Secret:   base64.RawURLEncoding.EncodeToString(secret),
		Members:  []string{},
	}

	if err := save(g); err != nil {
		return nil, err
	}
	return g, nil
}

// Load reads a group config from disk by name.
func Load(name string) (*Group, error) {
	path := groupPath(name)
	var g Group
	_, err := toml.DecodeFile(path, &g)
	if err != nil {
		return nil, fmt.Errorf("loading group %q: %w", name, err)
	}
	return &g, nil
}

// AddMember adds a peer ID to a group's member list.
func AddMember(name string, peerID string) error {
	g, err := Load(name)
	if err != nil {
		return err
	}

	for _, m := range g.Members {
		if m == peerID {
			return fmt.Errorf("peer %s is already a member of %q", peerID, name)
		}
	}

	g.Members = append(g.Members, peerID)
	return save(g)
}

// RemoveMember removes a peer ID from a group.
func RemoveMember(name string, peerID string) error {
	g, err := Load(name)
	if err != nil {
		return err
	}

	found := false
	members := make([]string, 0, len(g.Members))
	for _, m := range g.Members {
		if m == peerID {
			found = true
			continue
		}
		members = append(members, m)
	}
	if !found {
		return fmt.Errorf("peer %s is not a member of %q", peerID, name)
	}

	g.Members = members
	return save(g)
}

// List returns all group names.
func List() ([]Group, error) {
	dir := config.GroupsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var groups []Group
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".toml") {
			name := strings.TrimSuffix(e.Name(), ".toml")
			g, err := Load(name)
			if err != nil {
				continue
			}
			groups = append(groups, *g)
		}
	}
	return groups, nil
}

// Delete removes a group config file.
func Delete(name string) error {
	path := groupPath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("group %q does not exist", name)
	}
	return os.Remove(path)
}

// Exists checks if a group exists.
func Exists(name string) bool {
	_, err := os.Stat(groupPath(name))
	return err == nil
}

func groupPath(name string) string {
	safe := strings.ReplaceAll(name, string(filepath.Separator), "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	return filepath.Join(config.GroupsDir(), safe+".toml")
}

func save(g *Group) error {
	path := groupPath(g.Name)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("writing group file: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(g)
}

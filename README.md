# Pulse

P2P file sharing CLI. No servers. No cloud. Just peers.



## Installation

```bash
go build -o pulse .
```

## Quickstart

```bash
# 1. Initialize identity (once)
pulse init --relay "/ip4/<relay-ip>/tcp/4001/p2p/<relay-peerID>"

# 2. Create a group
pulse group create friends

# 3. Add members
pulse group add friends 12D3KooW...

# 4. Send a file
pulse send friends document.pdf

# 5. Receive files (on another machine)
pulse listen friends --dir ./downloads
```

## Commands

| Command | Description |
|---------|-------------|
| `pulse init` | Generate identity & config |
| `pulse whoami` | Display your PeerID |
| `pulse group create <name>` | Create a group |
| `pulse group add <group> <peerID...>` | Add members |
| `pulse group remove <group> <peerID>` | Remove a member |
| `pulse group list` | List all groups |
| `pulse group info <name>` | Show group details |
| `pulse group delete <name>` | Delete a group |
| `pulse send <group> <file>` | Send file to group members |
| `pulse listen <group>` | Listen for incoming files |
| `pulse status` | Show active listeners |
| `pulse stop <group>` | Stop a listener |

## Architecture

```
pulse/
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── config/             # TOML config, paths
│   ├── identity/           # Ed25519 key management
│   ├── group/              # Group CRUD
│   ├── transport/          # libp2p relay, streams, retry
│   ├── crypto/             # BLAKE3 integrity, NaCl encryption
│   └── ui/                 # Bubbletea models, Lipgloss styles
└── main.go
```

# Fincher Deployment Guide (NixOS + deploy-rs + agenix)

Reproducible, declarative deployment for **Fincher** to a cloud VPS (Google Compute Engine, Hetzner, etc.) using **NixOS**, **deploy-rs**, and **agenix**.

---

## Architecture

```
                      ┌────────────────────────────────────────────────────────┐
                      │ NixOS Host (e2-medium, 30GB Persistent Disk)           │
                      │                                                        │
┌───────────────┐     │  ┌───────────────┐     ┌────────────────────────────┐  │
│               │     │  │    Caddy      │────>│   Fincher Systemd Service  │  │
│    Browser    │────>│  │  (Port 80)    │     │   (127.0.0.1:8080)         │  │
│  (Port 80)    │     │  └───────────────┘     └───────┬────────────┬───────┘  │
│               │     │                                │            │          │
└───────────────┘     │          native TCP (9000)     │            │          │
                      │          ┌─────────────────────┘            │          │
                      │          │                                  │          │
                      │          ▼                                  ▼          │
                      │  ┌────────────────────────┐    ┌────────────────────┐  │
                      │  │   Native ClickHouse    │<───│  MCP ClickHouse    │  │
                      │  │   (/var/lib/clickhouse)│    │  (127.0.0.1:8000)  │  │
                      │  └────────────────────────┘    └────────────────────┘  │
                      │                                                        │
                      │  ┌──────────────────────────────────────────────────┐  │
                      │  │   Local SQLite with WAL Mode                     │  │
                      │  │   (/var/lib/fincher/fincher.db)                  │  │
                      │  └──────────────────────────────────────────────────┘  │
                      └────────────────────────────────────────────────────────┘
```

---

## Prerequisites

Enter the development shell containing all required deployment tools (`deploy-rs`, `agenix`, `age`, `gcloud`):

```bash
nix develop
```

---

## 1. Secrets Setup (`agenix`)

Secrets are encrypted with `age` using SSH public keys in `deploy/keys.nix`.

### Edit production secrets:

```bash
agenix -e deploy/secrets/fincher.env.age
```

Add your production environment variables (only Gemini API key is needed since SQLite and ClickHouse are hosted locally on disk):

```env
FINCHER_GEMINI_API_KEY=your_gemini_api_key
```

*(Optional overrides like `FINCHER_GEMINI_MODEL=gemini-2.5-flash` or `FINCHER_GEMINI_OPTIONS=location=asia-south1` can also be added if needed).*

Save and exit. The file is automatically re-encrypted.

---

## 2. Initial NixOS Bootstrap (`nixos-anywhere`)

If the target VM is newly provisioned or running standard Linux (Ubuntu/Debian), install NixOS over SSH in one command:

```bash
nix run github:nix-community/nixos-anywhere -- --flake .#fincher-prod root@<VM_IP_ADDRESS>
```

This will partition the disk via `disko`, install NixOS, install SSH keys, configure native ClickHouse, Caddy, Fincher, and reboot into NixOS.

---

## 3. Routine Deployments (`deploy-rs`)

For all subsequent code updates, schema migrations, and service changes:

```bash
deploy .#fincher-prod
```

`deploy-rs` will:
1. Build the Fincher binary (with embedded Preact UI) and NixOS closure locally.
2. Copy closures to the server over SSH.
3. Switch the live system to the new generation.
4. Automatically roll back if health verification fails.

---

## 4. Seeding & Database Migrations

To populate Turso Cloud and the local ClickHouse instance:

```bash
# Run seed on the remote server via SSH:
ssh deploy@<VM_IP_ADDRESS> "sudo -u fincher fincher-seed"

# Or reset and re-seed:
ssh deploy@<VM_IP_ADDRESS> "sudo -u fincher fincher-seed --reset"
```

---

## 5. Operations & Troubleshooting

### Check Service Status
```bash
ssh deploy@<VM_IP_ADDRESS> "systemctl status fincher"
ssh deploy@<VM_IP_ADDRESS> "systemctl status clickhouse"
ssh deploy@<VM_IP_ADDRESS> "systemctl status caddy"
```

### View Live Logs
```bash
ssh deploy@<VM_IP_ADDRESS> "journalctl -u fincher -f"
```

### Rollback to Previous Generation
```bash
# Instant rollback on the server:
ssh deploy@<VM_IP_ADDRESS> "sudo nix-env --rollback -p /nix/var/nix/profiles/system && sudo /nix/var/nix/profiles/system/bin/switch-to-configuration switch"
```

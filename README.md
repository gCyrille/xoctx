# XO Context Manager

A Go CLI tool to manage multiple Xen Orchestra contexts for testing and development.

## Installation

```bash
go build -o ~/bin/xoctx .
```

Or install directly:

```bash
go install github.com/gCyrille/xoctx@latest
```

## Profile Setup

Create `.env` files in `~/.xoctx/`:

```bash
mkdir -p ~/.xoctx
cat > ~/.xoctx/my-lab.env << 'EOF'
XOA_URL=https://xoa-my-lab.example.com
XOA_TOKEN=your-token-here
XOA_POOL=my-lab-pool
XOA_STORAGE=local-sr
EOF
```

## Usage

```
xoctx                     Reload last used context & show summary
xoctx use <name>          Switch to a context
xoctx env <name>          Print export statements for a context
xoctx init zsh            Print zsh integration (wrapper + xoctx_ps1)
xoctx list                List available contexts
xoctx current             Show active context name
xoctx which               Show current context summary (no reload)
xoctx show <name>         Display context variables (masked)
xoctx help                Show this help

Profiles: ~/.xoctx/
```

### List Available Contexts

```bash
xoctx list
```

### Switch to a Context

```bash
xoctx use my-lab
```

Or interactively (prompts for selection):

```bash
xoctx use
```

### Shell Integration (zsh)

Add this to your `~/.zshrc`:

```zsh
eval "$(xoctx init zsh)"
```

This defines a shell wrapper named `xoctx` that keeps the CLI behavior and also exports context environment variables after `xoctx use ...`.

**Note:** `xoctx` must be in your `PATH` for this to work.

```bash
xoctx use my-lab  # switch to my-lab and export env vars
xoctx use         # interactive selection + export selected context vars
xoctx list        # list contexts
xoctx show my-lab # show context config
xoctx current     # show current context
xoctx which       # show current context summary
```

### Prompt Integration

Display the current context in your shell prompt using `xoctx_ps1` (defined by `xoctx init zsh`):

```zsh
# ~/.zshrc
eval "$(xoctx init zsh)"
RPROMPT='$(xoctx_ps1)'
# or
PROMPT='$(xoctx_ps1) %n@%m %~ $ '
```

`xoctx_ps1` only displays the context name when `XOA_URL` is exported. If context environment variables are not loaded in the current shell, the prompt shows nothing.

Customize the prompt format:

```zsh
XO_PS1_PREFIX='⬢ '
XO_PS1_SUFFIX=''
```

Output: `⬢ my-lab` (default is `(|my-lab)`)

### Show Current Context

```bash
xoctx current
```

### Show Context Configuration

```bash
xoctx show my-lab
```

Sensitive values (tokens, passwords) are masked in output.

### Show Help

```bash
xoctx --help
```

### Print Shell Exports

```bash
xoctx env my-lab
```

This prints `export KEY=VALUE` lines suitable for `eval`.

## Environment Variables

Each lab profile exports these variables:

| Variable | Required | Purpose |
|----------|----------|---------|
| `XOA_URL` | Yes | XO API URL |
| `XOA_TOKEN` | Yes* | Authentication token |
| `XOA_USER` | Yes* | Username (if no token) |
| `XOA_PASSWORD` | Yes* | Password (if no token) |
| `XOA_POOL` | For tests | Pool name for integration tests |
| `XOA_STORAGE` | For tests | SR name for integration tests |
| `XOA_TEMPLATE` | Optional | VM template for tests |
| `XOA_INSECURE` | Optional | Skip TLS verification (true/false) |
| `XOA_DEVELOPMENT` | Optional | Enable debug logging (true/false) |
| `XOA_TEST_PREFIX` | Optional | Resource naming prefix |
| `XOA_TEST_PBD_ID` | Optional | PBD UUID for plug/unplug tests |
| `XOA_RETRY_MODE` | Optional | Retry mode (none/backoff) |
| `XOA_RETRY_MAX_TIME` | Optional | Max retry duration |

*Either `XOA_TOKEN` OR both `XOA_USER` and `XOA_PASSWORD` must be set.

## Tips

- Context profiles are stored as simple `.env` files in `~/.xoctx/`
- Edit them directly with a text editor — plain `KEY=VALUE`, no `export` needed
- Use `xoctx which` to view the active context summary without switching
- The active context persists across sessions (tracked in `~/.xoctx/current_lab`)
- Use descriptive context names (e.g., `my-lab`, `staging`, `devops-tools`)

## Troubleshooting

### The `use` command doesn't set environment variables

`xoctx` is a binary and cannot modify its parent shell's environment directly. Enable shell integration in your `~/.zshrc`:

```zsh
eval "$(xoctx init zsh)"
```

After that, `xoctx use <name>` both switches the active context and exports the corresponding variables in your current shell.

### Context configuration not persisting

The active context is saved to `~/.xoctx/current_lab`. If this file is deleted, no context will be active until you run `xoctx use` again.

## License

MIT

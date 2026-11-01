# Personal lazyjira build

This fork carries the customizations previously installed through Home Manager:

- Full, wrapped issue summaries above descriptions in the Details tab.
- Optional stacked panels at any terminal width with `gui.stackedLayout: true`.
- A taller issue list, up to 20 lines, with a stable list/detail split when focus
  moves between them and a compact three-line log panel.

The existing configuration in `~/.config/lazyjira/config.yml` remains in use:

```yaml
gui:
  stackedLayout: true
```

## Build and install

With Go installed, run from this checkout:

```bash
make check
mkdir -p "$HOME/bin"
go build -trimpath -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty)" -o lazyjira ./cmd/lazyjira
install -m 755 lazyjira "$HOME/bin/lazyjira"
```

Ensure `~/bin` is on `PATH`. The lazyjira Home Manager package should remain
disabled so it does not shadow this binary. `origin` points to the personal fork;
`upstream` points to `textfuel/lazyjira`.

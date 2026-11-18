# gocal

A colorful take on the classic Unix `cal` command, written in Go with no dependencies.

## Features

- A rainbow color for each month title
- Sundays and Saturdays are highlighted, and so are fixed-date holidays (New Year's Day, Valentine's Day, July 4th, Halloween, Christmas)
- Today stands out
- Single month, three-month (`-3`), or whole-year (`-y`) views
- Weeks can start on Monday (`-m`)
- Color turns off automatically when output isn't a terminal, when `NO_COLOR` is set, or when you pass `-no-color`
- Works in Windows consoles (VT escape processing is turned on automatically)

## Install

### Prebuilt binaries

Download the archive for your platform from the [Releases](https://github.com/ssta604/gocal/releases) page. Builds are provided for:

| OS      | Architectures  |
|---------|----------------|
| Linux   | amd64, arm64   |
| macOS   | amd64, arm64   |
| Windows | amd64, arm64   |

Extract it and put `gocal` (or `gocal.exe`) somewhere on your `PATH`.

> **macOS:** the binaries are unsigned. If Gatekeeper blocks one, run
> `xattr -d com.apple.quarantine ./gocal`.

### With Go

```sh
go install github.com/ssta604/gocal@latest
```

## Usage

```
gocal [-3] [-y] [-m] [-no-color] [[month] year]
```

| Flag        | Description                              |
|-------------|------------------------------------------|
| `-3`        | Show the previous, current and next month |
| `-y`        | Show the whole year                      |
| `-m`        | Start weeks on Monday                    |
| `-no-color` | Disable colors                           |

### Examples

```sh
gocal            # current month
gocal -3         # last, this and next month
gocal -y         # the whole current year
gocal 2027       # the whole year 2027
gocal 2 2028     # February 2028
gocal -m 11 2026 # November 2026, weeks starting on Monday
```

## Building from source

Requires Go 1.27 or newer.

```sh
git clone https://github.com/ssta604/gocal.git
cd gocal
go build
go test ./...
```

## Releasing

Pushing a tag that starts with `v` triggers the [release workflow](.github/workflows/release.yml). It cross-compiles for every platform above and publishes a GitHub Release with the archives and a `checksums.txt`:

```sh
git tag v0.1.0
git push origin v0.1.0
```

## License

[MIT](LICENSE.md)

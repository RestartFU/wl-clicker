# Moyai Clicker

A minimal hold-clicker for Linux and Windows.
Hold your trigger button/key to autoclick, and use a toggle keybind to pause/resume.

## Build

```bash
make build
```

This creates `./clicker`.

To build a Windows executable:

```bash
make build-windows
```

This creates `./clicker.exe`.
If you want to build for Windows while on Linux, install a MinGW-w64
cross-compiler first because Fyne's desktop backend depends on cgo.

Arch:

```bash
sudo pacman -S --needed mingw-w64-gcc
```

Fedora:

```bash
sudo dnf install -y mingw64-gcc mingw64-gcc-c++
```

Debian-based:

```bash
sudo apt install -y gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64
```

## Install

```bash
sudo make install
```

By default, this installs to `/usr/local/bin/clicker`.

UI reference: https://github.com/Dasciam/autoclicker-mcpe-go

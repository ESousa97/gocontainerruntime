<div align="center">
  <h1>gocontainerruntime</h1>
  <p>A minimal container runtime (study purposes) implemented in Go using Linux namespaces and cgroups.</p>

  <img src="assets/github-go.png" alt="gocontainerruntime Banner" width="600px">

  <br>

[![CI](https://github.com/ESousa97/gocontainerruntime/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/ESousa97/gocontainerruntime/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gocontainerruntime?=v2)](https://goreportcard.com/report/github.com/ESousa97/gocontainerruntime?=v9)
[![CodeFactor](https://www.codefactor.io/repository/github/esousa97/gocontainerruntime/badge)](https://www.codefactor.io/repository/github/esousa97/gocontainerruntime)
[![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gocontainerruntime.svg)](https://pkg.go.dev/github.com/ESousa97/gocontainerruntime)
[![License](https://img.shields.io/github/license/ESousa97/gocontainerruntime)](https://github.com/ESousa97/gocontainerruntime/blob/master/LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gocontainerruntime)](https://github.com/ESousa97/gocontainerruntime)
[![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gocontainerruntime)](https://github.com/ESousa97/gocontainerruntime/commits/master)

</div>

---

`gocontainerruntime` is a lightweight, educational container runtime that demonstrates how modern containers work under the hood. It implements process isolation using Linux Namespaces, resource control via Cgroups, and filesystem isolation using Chroot.

## Demonstration

```bash
# Pull Alpine rootfs
sudo ./gocontainer pull

# Run an isolated shell (Requires sudo for namespaces/cgroups)
sudo ./gocontainer run /bin/sh
```

## Technology Stack

| Technology     | Role                                                               |
| -------------- | ------------------------------------------------------------------ |
| Go 1.25.0      | Core language and logic                                            |
| Linux Syscalls | Namespaces (CLONE_NEWNS, CLONE_NEWUTS, CLONE_NEWPID, CLONE_NEWNET) |
| Cgroups v1     | Resource limits (100MB Memory, 512 CPU shares)                     |
| Cobra          | CLI Framework                                                      |
| Alpine Linux   | Lightweight rootfs for the container                               |

## Prerequisites

- Go >= 1.22
- Linux Kernel >= 4.x (with support for namespaces and cgroups v1)
- Root Privileges (required for namespace and network manipulation)

## Installation and Usage

### As a binary

```bash
go install github.com/ESousa97/gocontainerruntime@latest
```

### From source

```bash
git clone https://github.com/ESousa97/gocontainerruntime.git
cd gocontainerruntime
make build
# Optional: Pull default rootfs
make pull
# Run shell
make run
```

## Makefile Targets

| Target  | Description                                                    |
| ------- | -------------------------------------------------------------- |
| `build` | Compiles the `gocontainer` binary                              |
| `clean` | Removes binary and cache files                                 |
| `test`  | Executes the unit test suite                                   |
| `pull`  | Downloads and extracts the Alpine Linux minirootfs             |
| `run`   | Starts an interactive container with `/bin/sh` (requires sudo) |

## Architecture

The runtime operates in two main stages to ensure complete isolation:

<div align="center">

```mermaid
graph TD
    Start([User: gocontainer run]) --> Parent[Stage 1: Parent Process]

    Parent --> NS[Namespaces Isolation]
    NS --> CG[Cgroups Resource Control]
    CG --> Net[Network Setup]

    Net --> Child[Stage 2: Child Process]

    Child --> Host[Set Hostname]
    Host --> Chroot[Chroot Isolation]
    Chroot --> Mount[Mount /proc]

    Mount --> Final[/Exec: User Command/]

    style Final fill:#2da44e,stroke:#fff,stroke-width:1px,color:#fff
```

</div>

1. **Stage 1 (Parent)**: Creates new namespaces (UTS, PID, NS, NET), generates memory/CPU cgroups, and re-executes itself by calling the internal `child` command.
2. **Stage 2 (Child)**: Already inside the namespaces, sets the hostname (`gocontainer`), performs the `chroot` to the rootfs, mounts `/proc`, and executes the user's final command.

See more details in [docs/architecture.md](docs/architecture.md).

## API Reference

Detailed documentation for internal functions and packages is available at:
[pkg.go.dev/github.com/ESousa97/gocontainerruntime](https://pkg.go.dev/github.com/ESousa97/gocontainerruntime)

## Configuration

The CLI accepts flags for database configuration (legacy terminology from reference, but applicable to runtime setup).

| Variable    | Description                     | Type   | Default                  |
| ----------- | ------------------------------- | ------ | ------------------------ |
| `cacheDir`  | Directory for rootfs extraction | String | `./cache/alpine_rootfs`  |
| `alpineURL` | Alpine download URL             | String | Alpine 3.19.1 Minirootfs |

## Roadmap

- [x] Phase 1: Isolated Fork (Namespaces)
- [x] Phase 2: File Isolation (Chroot)
- [x] Phase 3: Resource Control (Cgroups)
- [x] Phase 4: Basic Networking (Netns)
- [x] Phase 5: Professional Interface and Images

## Contributing

Contributions are welcome! See the full guide at [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

<div align="center">

## Author

**Enoque Sousa**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?style=flat&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/enoque-sousa-bb89aa168/)
[![GitHub](https://img.shields.io/badge/GitHub-100000?style=flat&logo=github&logoColor=white)](https://github.com/ESousa97)
[![Portfolio](https://img.shields.io/badge/Portfolio-FF5722?style=flat&logo=target&logoColor=white)](https://enoquesousa.vercel.app)

**[⬆ Back to top](#gocontainerruntime)**

Made with ❤️ by [Enoque Sousa](https://github.com/ESousa97)

**Project Status:** Active — Study Project

</div>

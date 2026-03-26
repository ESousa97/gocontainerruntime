# Go Container Runtime

A minimalist implementation of a **Container Runtime** in Go, utilizing **Linux Namespaces** and **Cgroups** for process isolation and resource management.

## 🚀 How it Works

The runtime implements three layers of containerization:

### 1. Isolation (Namespaces)
- **UTS (`CLONE_NEWUTS`):** Isolates the hostname.
- **PID (`CLONE_NEWPID`):** Isolates the process ID tree. The container command becomes PID 1.
- **Mount (`CLONE_NEWNS`):** Isolates mount points. Used with `chroot` to provide a dedicated filesystem.
- **Network (`CLONE_NEWNET`):** Isolates the network stack. A `veth` pair bridges the host to the container.

### 2. Resource Control (Cgroups V1)
- **Memory:** Limits memory usage to 100MB (`memory.limit_in_bytes`).
- **CPU:** Sets CPU shares to 512 (`cpu.shares`), ensuring fair access to CPU cycles.
- **Cleanup:** Cgroup directories and veth interfaces are automatically deleted upon container exit.

### 3. Networking
- Creates a virtual ethernet pair (`veth-host` <-> `veth-child`).
- Host IP: `10.0.0.1`
- Container IP: `10.0.0.2`
- Enables communication via a virtual bridge.

## 🛠️ Requirements

- **Linux Kernel** (or WSL2).
- **Root Privileges** (`sudo`) for namespace and cgroup manipulation.

## 💻 Usage

### 1. Compile
```bash
go build -o gocontainer main.go
```

### 2. Pull a Basic Image
Downloads and extracts an Alpine Linux rootfs to `./cache`.
```bash
sudo ./gocontainer pull
```

### 3. Run a Container
If no rootfs is provided, it uses the cached Alpine image.
```bash
# Using cached image
sudo ./gocontainer run /bin/sh

# Using a custom rootfs
sudo ./gocontainer run ./my_rootfs /bin/bash
```

## 🔐 Security Note
This is an educational implementation. While it uses real Linux kernel primitives, it does not include advanced security features like **Seccomp**, **AppArmor**, or **User Namespaces** (`CLONE_NEWUSER`), which are essential for production-grade runtimes like `runc`.

---
*Developed for educational purposes to demonstrate the core pillars of Linux containers.*

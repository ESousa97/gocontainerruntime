# Architecture Details

This document explains the technical implementation of `gocontainerruntime`.

## Process Isolation: Namespaces

The core of container isolation relies on Linux Namespaces. When the parent process starts, it uses the `CLONE_NEWNS`, `CLONE_NEWUTS`, `CLONE_NEWPID`, and `CLONE_NEWNET` flags via `syscall.SysProcAttr`.

- **UTS (Unix Timesharing System)**: Allows the container to have its own hostname (`gocontainer`) separate from the host.
- **PID (Process ID)**: Ensures the container sees itself as PID 1, and cannot see processes running on the host or other containers.
- **Mount (NS)**: Provides an isolated mount table. This is crucial for `chroot` and mounting `/proc` safely.
- **Network (NET)**: Gives the container its own network stack, including its own interfaces, routing tables, and firewall rules.

## Resource Control: Cgroups v1

To prevent one container from exhausting host resources, we use Control Groups (Cgroups). The runtime currently manages two controllers:

### Memory Controller
The memory limit is set to **100MB** (`104857600` bytes) by writing to:
`/sys/fs/cgroup/memory/gocontainer/memory.limit_in_bytes`

### CPU Controller
CPU shares are set to **512** (relative weight) by writing to:
`/sys/fs/cgroup/cpu/gocontainer/cpu.shares`

The child process PID is written to `cgroup.procs` in both directories to start enforcement.

## Filesystem Isolation: Chroot

The runtime pulls an Alpine Linux minirootfs. The isolation is achieved in three steps:
1. `syscall.Chroot(rootfs)`: Changes the root directory context to the extracted Alpine folder.
2. `os.Chdir("/")`: Moves the current working directory to the new root.
3. `syscall.Mount("proc", "/proc", "proc", 0, "")`: Mounts the special `/proc` filesystem so that tools like `ps` and `top` work correctly within the container.

## Networking: Veth Pairs

The runtime sets up basic networking using a Virtual Ethernet (veth) pair:
1. `veth-host` remains in the host namespace.
2. `veth-child` is moved into the container's network namespace.
3. IP addresses `10.0.0.1/24` (host) and `10.0.0.2/24` (container) are assigned to enable connectivity.

---

> [!NOTE]
> This implementation is designed for educational purposes. For production use, consider OCI-compliant runtimes like `runc`.

# Go Container Runtime

A minimalist implementation of a **Container Runtime** in Go, utilizing **Linux Namespaces** for process isolation (PID), identity (UTS), and filesystem (Mount).

## 🚀 How it Works

The program uses the **re-execution** technique (`/proc/self/exe`) to properly apply isolation settings:

1.  **Stage 1 (run):** Prepares the namespaces (`CLONE_NEWUTS`, `CLONE_NEWPID`, `CLONE_NEWNS`) and re-executes itself within the new environment.
2.  **Stage 2 (child):** Inside the new namespaces, it sets the hostname (`container-runtime`) and mounts the `/proc` filesystem to ensure visible process isolation.
3.  **Final Stage:** Replaces the current process with the user's final command via `syscall.Exec`.

## 🛠️ Requirements

This project directly depends on **Linux Kernel** features.

- **Linux:** Works natively (requires `root` privileges to create namespaces).
- **Windows:** Use **WSL2** (Windows Subsystem for Linux).

## 💻 How to Run (Linux or WSL2)

### 1. Entering the Linux Environment
If you are on **Windows**, open your **Ubuntu** terminal (or whichever distro you installed on WSL).

Linux and Windows share files! To navigate to this project folder within Linux, use the `cd` (Change Directory) command:
```bash
# Replace with your actual path. Windows drives are mounted under /mnt/
cd /mnt/c/your/project/path/gocontainerruntime
```

### 2. Basic Navigation Commands:
- `ls`: List files in the current folder.
- `pwd`: Show the full path of your current location.
- `cd ..`: Go back one folder.

### 3. Compiling and Running the Container
Once you are in the correct folder inside the Linux terminal:

1.  **Compile the binary:**
    ```bash
    go build -o myruntime main.go
    ```

2.  **Run with superuser privileges (sudo):**
    ```bash
    sudo ./myruntime run /bin/bash
    ```

### Validating Isolation:

Inside the new shell:
- **Hostname:** Type `hostname`. It should return `container-runtime`.
- **Process Isolation:** Type `ps aux`. You should only see `/bin/bash` and `ps`, with the shell running as PID 1.

---
*This project is for educational purposes to demonstrate how containers work under the hood.*

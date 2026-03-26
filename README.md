# Go Container Runtime

Uma implementação minimalista de um **Container Runtime** em Go, utilizando **Linux Namespaces** para isolamento de processos (PID), rede/identidade (UTS) e sistema de arquivos (Mount).

## 🚀 Como funciona

O programa utiliza a técnica de **re-execução** (`/proc/self/exe`) para aplicar as configurações de isolamento corretamente:

1.  **Stage 1 (run):** Prepara os namespaces (`CLONE_NEWUTS`, `CLONE_NEWPID`, `CLONE_NEWNS`) e executa a si mesmo novamente dentro do novo ambiente.
2.  **Stage 2 (child):** Já dentro dos novos namespaces, define o hostname (`container-runtime`) e monta o sistema de arquivos `/proc` para garantir o isolamento visível de processos.
3.  **Final Stage:** Substitui o processo atual pelo comando final do usuário via `syscall.Exec`.

## 🛠️ Requisitos

Este projeto depende diretamente de funcionalidades do **Kernel Linux**.

- **Linux:** Funciona nativamente (requer privilégios de `root` para criar namespaces).
- **Windows:** Utilize o **WSL2** (Windows Subsystem for Linux).

## 💻 Como executar (Linux ou WSL2)

1.  Clone o repositório e acesse a pasta:
    ```bash
    git clone <URL_DO_REPOSITORIO>
    cd gocontainerruntime
    ```

2.  Compile o binário para Linux:
    ```bash
    go build -o myruntime main.go
    ```

3.  Execute o container runtime com privilégios de superusuário:
    ```bash
    sudo ./myruntime run /bin/bash
    ```

### Validando o isolamento:

Dentro do novo shell:
- **Hostname:** Digite `hostname`. Deve retornar `container-runtime`.
- **Isolamento de Processos:** Digite `ps aux`. Você verá apenas o `/bin/bash` e o `ps`, com o PID 1.

---
*Este projeto é para fins educacionais de como containers funcionam por baixo do capô.*

# gocontainerruntime

> A minimal container runtime (study purposes) implemented in Go using Linux namespaces and cgroups

![CI](https://github.com/ESousa97/gocontainerruntime/actions/workflows/ci.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/ESousa97/gocontainerruntime)
![Go Reference](https://pkg.go.dev/badge/github.com/ESousa97/gocontainerruntime.svg)
![License](https://img.shields.io/github/license/ESousa97/gocontainerruntime)
![Go Version](https://img.shields.io/github/go-mod/go-version/ESousa97/gocontainerruntime)
![Last Commit](https://img.shields.io/github/last-commit/ESousa97/gocontainerruntime)

---

`gocontainerruntime` is a lightweight, educational container runtime that demonstrates how modern containers work under the hood. It implements process isolation using Linux Namespaces, resource control via Cgroups, and filesystem isolation using Chroot.

## Demonstração

```bash
# Baixar rootfs do Alpine
sudo ./gocontainer pull

# Rodar um shell isolado (Requires sudo for namespaces/cgroups)
sudo ./gocontainer run /bin/sh
```

## Stack Tecnológico

| Tecnologia | Papel |
|---|---|
| Go 1.25.0 | Linguagem principal e chamadas de sistema |
| Linux Syscalls | Namespaces (CLONE_NEWNS, CLONE_NEWUTS, CLONE_NEWPID, CLONE_NEWNET) |
| Cgroups v1 | Limite de memória (100MB) e CPU (512 shares) |
| Cobra | CLI Framework |
| Alpine Linux | Rootfs leve para o container |

## Pré-requisitos

- Go >= 1.22
- Linux Kernel >= 4.x (com suporte a namespaces e cgroups v1)
- Privilégios de Root (necessário para manipulação de namespaces e redes)

## Instalação e Uso

### Como binário

```bash
go install github.com/ESousa97/gocontainerruntime@latest
```

### A partir do source

```bash
git clone https://github.com/ESousa97/gocontainerruntime.git
cd gocontainerruntime
make build
# Opcional: Baixar rootfs padrão
make pull
# Rodar shell
make run
```

## Makefile Targets

| Target | Descrição |
|---|---|
| `build` | Compila o binário `gocontainer` |
| `clean` | Remove binário e arquivos de cache |
| `test` | Executa a suíte de testes unitários |
| `pull` | Faz o download e extração do minirootfs do Alpine Linux |
| `run` | Inicia um container interativo com `/bin/sh` (requer sudo) |

## Arquitetura

O runtime opera em dois estágios principais para garantir a isolação completa:

1. **Stage 1 (Parent)**: Cria novos namespaces (UTS, PID, NS, NET), gera os cgroups de memória/CPU e re-executa o próprio binário chamando o comando interno `child`.
2. **Stage 2 (Child)**: Já dentro dos namespaces, define o hostname (`gocontainer`), realiza o `chroot` para o rootfs, monta `/proc` e executa o comando final do usuário.

Veja mais detalhes em [docs/architecture.md](docs/architecture.md).

## API Reference

A documentação detalhada das funções internas e pacotes está disponível em:
[pkg.go.dev/github.com/ESousa97/gocontainerruntime](https://pkg.go.dev/github.com/ESousa97/gocontainerruntime)

## Configuração

| Variável | Descrição | Tipo | Padrão |
|---|---|---|---|
| `cacheDir` | Diretório para extração do rootfs | String | `./cache/alpine_rootfs` |
| `alpineURL` | URL de download do Alpine | String | Alpine 3.19.1 Minirootfs |

## Roadmap de Estudo (Finalizado)

- [x] **Fase 1: O Fork Isolado (Namespaces)**: Criação de namespaces PID, UTS e Mount com re-exec.
- [x] **Fase 2: Isolação de Arquivos (Chroot)**: Isolação da raiz do sistema e montagem do `/proc`.
- [x] **Fase 3: Controle de Recursos (Cgroups)**: Limitação de memória (100MB) e CPU (shares).
- [x] **Fase 4: Networking Básico (Netns)**: Configuração de veth pairs e IPs estáticos.
- [x] **Fase 5: Interface Profissional e Imagens**: CLI completa com Cobra e download de rootfs Alpine.

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md) para saber como participar do projeto.

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](LICENSE) para mais informações.

## Autor

**Enoque Sousa**
- [Portfolio](https://enoquesousa.vercel.app)
- [GitHub](https://github.com/ESousa97)

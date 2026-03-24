# Ketsu Notes

> Um editor de notas CLI minimalista para o samurai moderno.

## Quick Start

```bash
# Clone e compile
git clone https://github.com/seu-user/ketsu.git
cd ketsu
go mod download
go build -o ketsu ./cmd/ketsu

# Execute
./ketsu
```

## Documentação

| Arquivo | Descrição |
|---------|-----------|
| [SETUP.md](SETUP.md) | Guia de configuração e instalação |
| [PROJECT.md](PROJECT.md) | Detalhes técnicos do projeto |
| [PROGRESS.md](PROGRESS.md) | Progresso e status do desenvolvimento |
| [COMMANDS.md](COMMANDS.md) | Referência rápida de comandos |

## Funcionalidades

- **Editor Vim-like** - Modos normal e insert
- **File Explorer** - Árvore de notas lateral
- **Preview Markdown** - Renderização em tempo real
- **Scripting Lua** - Automação com hot reload
- **Plugins** - Extensibilidade via GitHub
- **LLM** - Integração com Ollama

## Stack

- **Go 1.21+** - Linguagem principal
- **Bubble Tea** - Framework TUI
- **Gopher-Lua** - Interpreter Lua
- **Goldmark** - Parser Markdown

## Comandos

```bash
ketsu              # Abrir interface
ketsu --help      # Ver ajuda
ketsu --version   # Ver versão
```

---

*Minimal. Powerful. Ketsu.*
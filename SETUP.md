# Setup do Ketsu Notes

## Pré-requisitos

- **Go 1.21+** - https://go.dev/dl/
- **Git** - https://git-scm.com/
- **Terminal moderno** - Windows Terminal (recomendado para Windows)

## Instalação

### 1. Clone o projeto

```bash
git clone https://github.com/seu-user/ketsu.git
cd ketsu
```

### 2. Instale dependências

```bash
go mod download
```

### 3. Compile

```bash
go build -o ketsu ./cmd/ketsu
```

### 4. Execute

```bash
./ketsu
# ou no Windows
ketsu.exe
```

## Configuração

O Ketsu cria um arquivo de configuração em:
- Linux/macOS: `~/.config/ketsu/config.toml`
- Windows: `%APPDATA%\ketsu\config.toml`

### Configuração Padrão

```toml
# Configuração básica
vault = "~/notes"          # Pasta principal de notas
theme = "default"           # Tema de cores
font = "default"            # Fonte do terminal

# Editor
indent_size = 2             # Tamanho da indentação
tab_size = 2                # Tamanho do tab
auto_save = true            # Auto salvar
auto_save_delay = 1000      # Delay em ms

# Lua scripting
lua_scripts = "~/ketsu/scripts"  # Pasta de scripts
hot_reload = true                # Recarregar scripts automaticamente

# Plugins
plugins_dir = "~/ketsu/plugins"  # Pasta de plugins
auto_update_plugins = true        # Atualizar plugins automaticamente

# Ollama (LLM)
ollama_host = "http://localhost:11434"
ollama_model = "llama2"
```

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `KETSU_VAULT` | Pasta de notas | `~/notes` |
| `KETSU_CONFIG` | Arquivo de config | Padrão do sistema |
| `KETSU_DEBUG` | Ativar debug | `false` |

## Estrutura de Arquivos

```
~/notes/                    # Seu vault
├── .ketsu/
│   └── config.toml         # Config específica do vault
├── notas/
│   └── *.md                # Suas notas
└── scripts/                 # Scripts Lua (opcional)
    └── *.lua
```

## Comandos Úteis

```bash
ketsu --help               # Ver ajuda
ketsu --version            # Ver versão
ketsu --vault /path        # Abrir vault específico
```

## Troubleshooting

### Terminal não suporta cores
Configure a variável `COLORTERM`:
```bash
export COLORTERM=truecolor
```

### Problemas no Windows
Use o **Windows Terminal** em vez do CMD padrão:
https://github.com/microsoft/terminal

### Imagens não aparecem
O suporte a imagens requer terminal compatível:
- Windows Terminal (com SIXEL ou Kitty protocol)
- iTerm2 (macOS)
- Foot (Linux)

---

Para mais informações, veja [PROJECT.md](PROJECT.md)
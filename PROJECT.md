# Ketsu Notes

> Um editor de notas CLI minimalista para o samurai moderno. Como uma katana afiada nas mãos certas - simples, preciso, poderoso.

## Filosofia

- **Minimalista**: Interface limpa, sem distrações
- **Poderoso**: Funcionalidades avançadas de forma abstrata
- **Extensível**: Scripting Lua e plugins via GitHub
- **Cross-platform**: Funciona no Windows (PowerShell/Terminal), Linux e macOS
- **Híbrido**: Suporta tanto notas markdown quanto dados estruturados

## Stack Técnica

| Componente | Biblioteca | Referência |
|------------|------------|------------|
| **TUI Framework** | Bubble Tea | charmbracelet/bubbletea |
| **Styling** | Lip Gloss | charmbracelet/lipgloss |
| **Scripting** | Gopher-Lua | yuin/gopher-lua |
| **Terminal** | tcell v3 | gdamore/tcell |
| **Markdown** | goldmark | yuin/goldmark |
| **YAML Parser** | go-yaml | go-yaml/yaml |
| **LLM/Ollama** | go-ollama | ollama/ollama |
| **Fuzzy Search** | go-fuzzyfinder |/skratchdot/go-fuzzyfinder |

## Conceitos Principais

### Space (Espaço de Trabalho)

Equivalent ao "vault" do Obsidian, mas chamado de **Space**.
- Qualquer pasta pode ser um Space
- Um Space é reconhecido pela presença de `.ketsu/` (pasta de configuração)
- Suporta múltiplos Spaces simultâneos

### Tipos de Arquivos

```
.md    → Nota markdown padrão (texto livre)
.db    → Dados estruturados em YAML (frontmatter)
```

### Arquivo .db (Dados Estruturados)

Arquivo YAML com frontmatter que funciona como um "banco de dados" minimalista:

```yaml
---
name: project_alpha
created: 2024-01-15
tags: [work, important]
fields:
  title: "Project Alpha"
  status: "in-progress"
  priority: 3
  assigned_to: "john"
  due_date: "2024-02-01"
---

# Entry points (opcional)
# markdown tradicional como conteúdo
- entry: started
  timestamp: 2024-01-15T10:00:00Z
```

#### Campos especiais

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `name` | string | Identificador único |
| `tags` | array | Tags para organização |
| `fields` | object | Dados estruturados |
| `created` | timestamp | Data de criação |
| `modified` | timestamp | Última modificação |

#### Exemplo de uso

```yaml
# tasks.db
---
name: tasks
fields:
  todo: []
  doing: []
  done: []
---
- task: "Implementar busca"
  status: doing
  priority: high
- task: "Testar no Windows"
  status: todo
  priority: medium
```

### Ketsu Shell (Console Interativo)

Um REPL (Read-Eval-Print Loop) integrado no app - como um Redis ouShell minimalista:

```
> Abstração
Simplificar

> Abstração = O ato de remover detalhes desnecessários
(ok)

> project
{'name': 'alpha', 'status': 'in-progress'}

> project.status = 'done'
(ok)

> keys *
abstraction, project, todo_list
```

#### Sintaxe do Shell (modo interativo)

Digite `\` para entrar no shell:

```
> key              → Buscar valor
> key = value     → Definir/atualizar valor
> keys *          → Listar todas as chaves
> key++           → Incrementar número
> key--           → Decrementar número
> clear           → Limpar display
```

#### Comandos do Editor (modo Vim)

Digite `:` para entrar em modo comando:

```
:key = value      → Salvar no KV diretamente
:Wq               → Salvar e sair
:w                → Salvar arquivo
:q                → Sair
:q!               → Forçar saída
:e arquivo.md     → Abrir arquivo
:new              → Novo arquivo
:split            → Dividir janela
```

**Exemplos:**
```
:Test = Testando    → Salva "Testando" na chave "Test"
:Wq                 → Salva e sai
:Nota = Minha nota → Salva no KV
```

#### Armazenamento

Armazenamento em memória (RAM) com persistência opcional:
- **Ephemeral**: Só em memória (mais rápido)
- **Persistente**: Salvo em `.ketsu/data.kv` (formato KV minimalista)

Formato do arquivo `.ketsu/data.kv`:
```
# comment
abstração=Simplificar
project={"name":"alpha","status":"in-progress"}
counter=42
```

#### Busca em Tempo Real

O modo busca (diferente do shell) funciona como fzf:
- **Busca por tags**: `#tag` filtra notas com tag específica
- **Busca por texto**: `palavra` filtra pelo conteúdo
- **Busca por arquivo**: `path/to/file` encontra arquivos específicos
- **Busca por campo**: `field:value` filtra arquivos .db por campos

O componente de busca atualiza em tempo real conforme o usuário digita.

### Auto-detecção de Space

```
~/notes/
├── .ketsu/              # Pasta de config do Ketsu
│   └── config.toml
├── notas/
│   ├── meeting.md
│   └── project.db
└── scripts/
    └── auto-save.lua
```

## Funcionalidades

### Fase 1 - Core (MVP)
- [x] Editor de texto Vim-like (modos normal/insert)
- [x] File explorer (árvore de arquivos)
- [ ] Preview Markdown renderizado
- [x] Sistema de Space (pasta de notas com .ketsu/)
- [x] Criar/editar notas .md
- [x] **Ketsu Shell** (console KV interativo)
- [x] Sistema KV em memória (ultra-rápido)
- [x] Persistência opcional em .ketsu/data.kv

### Fase 2 - Busca Avançada
- [ ] Fuzzy finder integrado (estilo fzf)
- [ ] Busca em tempo real com preview
- [ ] Busca por tags (#tag)
- [ ] Busca em arquivos .db por campos
- [ ] Histórico de buscas

### Fase 3 - Scripting
- [ ] Interpreter Lua integrado (Gopher-Lua)
- [ ] Hot reload de scripts
- [ ] API Lua para ações do editor
- [ ] Hooks: on_open, on_save, on_change
- [ ] Scripts para automatizar dados .db

### Fase 4 - Plugins
- [ ] Sistema de plugins
- [ ] Plugin manager com GitHub API
- [ ] Auto-update de plugins
- [ ] Repositório oficial de plugins

### Fase 5 - Avançado
- [ ] Integração Ollama (LLM local)
- [ ] LSP básico para Markdown
- [ ] Suporte a imagens no terminal (Kitty/Sixel)
- [ ] Queries para arquivos .db (estilo DataView)

## Estrutura do Projeto

```
ketsu/
├── cmd/ketsu/main.go       # Entry point
├── internal/
│   ├── core/               # Lógica principal
│   │   └── core.go         # Core (Space, KV, comandos)
│   ├── editor/             # Editor de texto Vim-like
│   │   └── editor.go       # Editor com Visual Mode, Undo/Redo
│   ├── ui/                 # Componentes TUI
│   │   ├── main.go         # App principal (Bubble Tea)
│   │   └── theme.go        # Sistema de temas
│   ├── kv/                 # Sistema Key-Value
│   │   └── store.go        # KV store em memória
│   ├── markdown/           # ✓ Parser/Renderer Markdown
│   │   └── markdown.go     # Goldmark + syntax highlighting
│   ├── lua/                # Interpreter scripting (TODO)
│   └── plugins/            # Sistema de plugins (TODO)
├── scripts/                # Scripts Lua do usuário
├── plugins/                # Plugins instalados
├── go.mod                  # Dependências Go
├── go.sum                  # Checksums
├── README.md               # Documentação principal
├── PROJECT.md              # Detalhes técnicos
├── SETUP.md                # Guia de instalação
├── COMMANDS.md             # Referência de comandos
└── PROGRESS.md             # Progresso do projeto
```

## Configuração de um Space

Arquivo: `.ketsu/config.toml`

```toml
[space]
name = "My Notes"
description = "Notas pessoais"

[editor]
indent_size = 2
tab_size = 2
auto_save = true
auto_save_delay = 1000
line_numbers = true

[search]
fuzzy_threshold = 0.6
max_results = 100

[lua]
enabled = true
scripts_dir = "scripts"

[plugins]
enabled = true
plugins_dir = "plugins"
auto_update = true

[ollama]
enabled = false
host = "http://localhost:11434"
model = "llama2"
```

## Referências

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [Gopher-Lua](https://github.com/yuin/gopher-lua)
- [go-fuzzyfinder](https://github.com/skratchdot/go-fuzzyfinder)
- [Goldmark](https://github.com/go-yaml/goldmark)
- [go-yaml](https://github.com/go-yaml/yaml)

---

*Minimal. Powerful. Ketsu.*
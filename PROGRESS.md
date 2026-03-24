# Ketsu Notes - Progresso do Projeto

> Última atualização: 24/03/2026 - 15:30

## Resumo Geral

| Fase | Status | Progresso |
|------|--------|-----------|
| Fase 1 - Core (MVP) | 🟢 Quase completo | ~85% |
| Fase 2 - Busca Avançada | 🟡 Em andamento | 30% |
| Fase 3 - Scripting | 🔴 Não iniciado | 0% |
| Fase 4 - Plugins | 🔴 Não iniciado | 0% |
| Fase 5 - Avançado | 🔴 Não iniciado | 5% |

---

## ✅ Funcionalidades Implementadas

### Core (Fase 1)

- [x] Editor de texto Vim-like (modos Normal/Insert/Command/Visual)
- [x] **Visual Mode** - Seleção de texto (v, V)
- [x] Navegação Vim: h/j/k/l, w, b, 0, $, gg, G, ^, e
- [x] Inserção de texto em modo Insert
- [x] **Undo/Redo** (u, Ctrl+R) - Histórico de até 100 estados
- [x] **Comandos Vim avançados:**
  - [x] dd - Deletar linha
  - [x] cc - Mudar linha
  - [x] yy - Yank (copiar) linha
  - [x] p/P - Colar antes/depois
  - [x] J - Juntar linhas
  - [x] ~/ Toggle case
  - [x] >> / << - Indentar
  - [x] o/O - Nova linha abaixo/acima
- [x] Comandos Vim: :w, :q, :wq, :e, :new, :split
- [x] **Preview Markdown** - Renderização em tempo real
- [x] **Syntax highlighting Markdown** no editor
- [x] File explorer (árvore de arquivos lateral)
- [x] Ícones por extensão de arquivo
- [x] Sistema de Space (pasta .ketsu/ detectada)
- [x] Criar/editar notas .md
- [x] Salvar e carregar arquivos
- [x] Sistema KV em memória (ultra-rápido)
- [x] Persistência em .ketsu/data.kv
- [x] Ketsu Shell (console KV interativo)
- [x] Comandos KV: key=value, keys *, key++, key--
- [x] **Layout flexível** - Explorer + Editor + Preview

### Interface (Fase 1)

- [x] Interface TUI com Bubble Tea
- [x] Header com nome do arquivo e modo
- [x] Footer com status e ajuda
- [x] Sistema de temas (Catppuccin, Tokyo Night, Nord)
- [x] Alternar tema com Ctrl+T ou :theme
- [x] Toggle explorer com Tab ou Ctrl+H

### UI/UX

- [x] Linhas numeradas no editor
- [x] Destaque da linha atual
- [x] Cursor visual no modo Insert
- [x] Barra de ferramentas com botões
- [x] Indicador de arquivo modificado
- [x] Scroll automático do editor

---

## ❌ Funcionalidades NÃO Implementadas

### Módulos Vazios

- [ ] `internal/lua/` - Interpreter Lua
- [ ] `internal/plugins/` - Sistema de plugins

### ✅ Recém Implementados

- [x] `internal/markdown/` - Parser/Renderer Markdown **✓ IMPLEMENTADO**
- [x] Preview Markdown renderizado **✓ IMPLEMENTADO**
- [x] Visual Mode (v, V) **✓ IMPLEMENTADO**
- [x] Undo/Redo (u, Ctrl+R) **✓ IMPLEMENTADO**
- [x] Comandos: dd, cc, yy, p **✓ IMPLEMENTADO**

### Fase 2 - Busca Avançada (em andamento)

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

### Editor - Recursos Vim Faltando

- [ ] Modo Visual Block (Ctrl+V)
- [ ] Modo Replace (R)
- [ ] Substituição (:s)
- [ ] Text objects (diw, ciw, etc.)
- [ ] Marcas (m, ')
- [ ] Macros (q, @)

### Testes e Qualidade

- [ ] Testes unitários
- [ ] Testes de integração
- [ ] CI/CD pipeline

---

## 📊 Estrutura do Código

### Arquivos Implementados

```
cmd/ketsu/main.go          ✅ Entry point funcional
internal/core/core.go      ✅ Lógica principal (152 linhas)
internal/editor/editor.go  ✅ Editor Vim-like (266 linhas)
internal/ui/main.go        ✅ Interface TUI (685 linhas)
internal/ui/theme.go       ✅ Sistema de temas (310 linhas)
internal/kv/store.go       ✅ Sistema KV (233 linhas)
```

### Arquivos Vazios/Não existentes

```
internal/lua/              ❌ Vazio (0 arquivos)
internal/markdown/         ❌ Vazio (0 arquivos)
internal/plugins/          ❌ Vazio (0 arquivos)
```

---

## 🔧 Dependências Atuais (go.mod)

```go
require (
    github.com/charmbracelet/bubbletea v0.25.0  // TUI Framework
    github.com/charmbracelet/lipgloss v0.11.0   // Styling
)
```

### Dependências Faltando (conforme PROJECT.md)

- [ ] `github.com/yuin/gopher-lua` - Lua interpreter
- [ ] `github.com/go-yaml/goldmark` - Markdown parser
- [ ] `github.com/go-yaml/yaml` - YAML parser
- [ ] `github.com/ollama/ollama` - Ollama client
- [ ] `github.com/skratchdot/go-fuzzyfinder` - Fuzzy finder

---

## 📝 Últimas Mudanças (24/03/2026)

### Commit: `b74b7fc`

**Novos arquivos criados:**
- `internal/markdown/markdown.go` - Parser Markdown com Goldmark

**Arquivos atualizados:**
- `internal/editor/editor.go` - Reescrito com Visual Mode, Undo/Redo, mais comandos
- `internal/ui/main.go` - Reescrito com Preview, layout flexível, novos modos
- `go.mod` - Adicionadas dependências goldmark e fuzzy
- `COMMANDS.md` - Atualizado com novos comandos
- `PROGRESS.md` - Este arquivo
- `PROJECT.md` - Atualizado com estrutura atual

---

## 📝 Próximos Passos Recomendados

### Prioridade Alta (MVP)

1. **Implementar preview Markdown** - Criar `internal/markdown/`
2. **Implementar fuzzy finder** - Adicionar dependência e UI
3. **Corrigir persistência de config** - Config.toml não está sendo salva
4. **Adicionar mais comandos Vim** - dd, dw, Visual mode

### Prioridade Média

5. **Implementar Lua scripting** - Criar `internal/lua/`
6. **Implementar sistema de plugins** - Criar `internal/plugins/`
7. **Adicionar suporte a .db files** - Dados estruturados YAML
8. **Implementar Undo/Redo** - Essential para editor

### Prioridade Baixa

9. **Adicionar testes unitários** - Cobertura mínima
10. **Integração Ollama** - LLM features
11. **LSP para Markdown** - Autocompletar

---

## 🐛 Issues Conhecidos

1. Go não está instalado no ambiente de desenvolvedo (mas exe funciona)
2. Não há validação de salvar antes de sair (perde dados)
3. File explorer não mostra subdiretórios recursivamente
4. KV save pode falhar silenciosamente

---

## 📚 Notas de Desenvolvimento

- O projeto compila e gera `ketsu.exe` funcional
- Interface TUI está completa e visualmente bonita
- Sistema KV está funcional mas limitado a strings
- Faltam módulos essenciais: Lua, Markdown, Plugins
- Documentação (README, PROJECT, SETUP) está completa

---

*Minimal. Powerful. Ketsu.*

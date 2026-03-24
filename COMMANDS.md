# Referência Rápida - Ketsu Notes

> Atualizado: 24/03/2026

## Modos de Operação

| Tecla | Ação |
|-------|------|
| `\` | Entrar no **Shell** (REPL) |
| `:` | Entrar em **Comando** (Vim-style) |
| `/` | Buscar no editor |
| `?` | Modo **Busca** global (fzf-style) |
| `Esc` | Voltar ao modo Normal |
| `Ctrl+S` | Salvar tudo |
| `Ctrl+C` | Sair |

---

## Layout

| Tecla | Ação |
|-------|------|
| `Tab` | Toggle Explorer |
| `Shift+Tab` | Toggle Preview |
| `Ctrl+H` | Toggle Explorer |
| `Ctrl+P` | Toggle Preview |
| `Ctrl+T` | Mudar tema |

---

## Shell (REPL) - Digite `\`

```
> Abstração
Simplificar

> Abstração = O ato de simplificar
(ok)

> keys *
abstração, project

> project = {"name": "alpha", "status": "done"}
(ok)
```

### Comandos do Shell

| Comando | Ação |
|---------|------|
| `key` | Buscar valor |
| `key = value` | Definir/atualizar valor |
| `keys *` | Listar todas as chaves |
| `key++` | Incrementar número |
| `key--` | Decrementar número |
| `clear` | Limpar tela |

---

## Editor (Vim-like) - Modos

### Modo Normal (padrão)

#### Movimento Básico

| Tecla | Ação |
|-------|------|
| `h` / `←` | Mover esquerda |
| `j` / `↓` | Mover baixo |
| `k` / `↑` | Mover cima |
| `l` / `→` | Mover direita |
| `w` | Próxima palavra |
| `b` | Palavra anterior |
| `e` | Final da palavra |
| `0` | Início da linha |
| `^` | Primeiro caractere não-espaço |
| `$` | Final da linha |
| `gg` | Primeira linha |
| `G` | Última linha |

#### Edição

| Tecla | Ação |
|-------|------|
| `i` | Entrar em modo Insert |
| `a` | Entrar em Insert após cursor |
| `A` | Entrar em Insert no final da linha |
| `I` | Entrar em Insert no início da linha |
| `o` | Nova linha abaixo + Insert |
| `O` | Nova linha acima + Insert |
| `x` / `Delete` | Deletar caractere à direita |
| `X` | Deletar caractere à esquerda |
| `~` | Toggle maiúscula/minúscula |
| `J` | Juntar com linha seguinte |
| `>>` | Indentar linha |
| `<<` | Remover indentação |

#### Operações (digite duas vezes)

| Tecla | Ação |
|-------|------|
| `dd` | Deletar linha |
| `cc` | Mudar linha (deletar + Insert) |
| `yy` | Copiar linha |
| `dw` | Deletar palavra |
| `cw` | Mudar palavra |

#### Colar

| Tecla | Ação |
|-------|------|
| `p` | Colar após cursor |
| `P` | Colar antes do cursor |

#### Undo/Redo

| Tecla | Ação |
|-------|------|
| `u` | Desfazer |
| `Ctrl+R` | Refazer |

#### Visual Mode

| Tecla | Ação |
|-------|------|
| `v` | Entrar Visual Mode (caracteres) |
| `V` | Entrar Visual Line Mode |
| `d` / `x` | Deletar seleção |
| `y` / `Y` | Copiar seleção |
| `c` | Mudar seleção |

#### Busca

| Tecla | Ação |
|-------|------|
| `/` | Buscar para frente |
| `n` | Próxima ocorrência |
| `N` | Ocorrência anterior |
| `*` | Buscar palavra sob cursor |

### Modo Insert

| Tecla | Ação |
|-------|------|
| `Esc` | Voltar ao modo Normal |
| `Ctrl+C` | Voltar ao modo Normal |

---

## Comandos (digite `:`)

```
:w                   → Salvar arquivo atual
:wq                  → Salvar e sair
:q                   → Sair (se não modificado)
:q!                  → Forçar sair
:e notas.md          → Abrir arquivo
:new                 → Novo arquivo
:split               → Dividir janela
:preview             → Toggle preview markdown
:theme catppuccin    → Mudar tema
:key = value         → Salvar no KV
```

### Temas Disponíveis

- `catppuccin` (padrão)
- `tokyo-night`
- `nord`

---

## Preview Markdown

O preview é exibido automaticamente para arquivos `.md` quando ativado.

| Tecla | Ação |
|-------|------|
| `Shift+Tab` | Toggle preview |
| `Ctrl+P` | Toggle preview |
| `:preview` | Toggle preview via comando |

---

## Atalhos Globais

| Atalho | Ação |
|--------|------|
| `\` | Abrir Shell |
| `:` | Abrir linha de comando |
| `/` | Buscar no editor |
| `?` | Busca global |
| `Ctrl+S` | Salvar |
| `Ctrl+Q` | Sair |
| `Tab` | Toggle Explorer |
| `Shift+Tab` | Toggle Preview |
| `Ctrl+T` | Próximo tema |

---

## Estrutura de Arquivos

```
.space/
├── .ketsu/
│   ├── config.toml      # Configuração
│   └── data.kv          # Dados KV
├── notas.md             # Notas markdown
├── tarefas.db           # Dados YAML
└── scripts/             # Scripts Lua
```

---

*Minimal. Powerful. Ketsu.*

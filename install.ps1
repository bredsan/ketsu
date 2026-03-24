# Ketsu Notes - Instalador Automático

Write-Host "Ketsu Notes Installer" -ForegroundColor Cyan
Write-Host "=====================" -ForegroundColor Cyan

$ErrorActionPreference = "Stop"

# Cores
$Green = [ConsoleColor]::Green
$Yellow = [ConsoleColor]::Yellow
$Cyan = [ConsoleColor]::Cyan

# Detectar local de instalação
$InstallDir = "$env:USERPROFILE\.ketsu"
$NotesDir = "$env:USERPROFILE\notes"
$ExePath = "$InstallDir\ketsu.exe"

# Baixar última versão ou usar local
Write-Host "`n1. Preparando instalação..." -ForegroundColor $Cyan

if (Test-Path "$PSScriptRoot\ketsu.exe") {
    Write-Host "   Usando executável local..." -ForegroundColor $Green
    Copy-Item "$PSScriptRoot\ketsu.exe" $ExePath -Force
} else {
    Write-Host "   Executável não encontrado no diretório!" -ForegroundColor $Yellow
    Write-Host "   Coloque ketsu.exe na mesma pasta deste script." -ForegroundColor $Yellow
    exit 1
}

# Criar estrutura de diretórios
Write-Host "`n2. Criando estrutura..." -ForegroundColor $Cyan

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path $NotesDir | Out-Null

# Criar arquivo de configuração
$ConfigContent = @"
[space]
name = "My Notes"
description = "Notas pessoais"

[editor]
indent_size = 2
tab_size = 2
auto_save = true
auto_save_delay = 1000
line_numbers = true
"@

Set-Content -Path "$NotesDir\.ketsu\config.toml" -Value $ConfigContent -Force
Write-Host "   Configuração criada em $NotesDir\.ketsu\" -ForegroundColor $Green

# Criar atalho no Desktop
Write-Host "`n3. Criando atalho..." -ForegroundColor $Cyan

$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\Ketsu Notes.lnk")
$Shortcut.TargetPath = $ExePath
$Shortcut.WorkingDirectory = $NotesDir
$Shortcut.Description = "Ketsu Notes - Editor CLI"
$Shortcut.Save()

Write-Host "   Atalho criado na área de trabalho!" -ForegroundColor $Green

# Adicionar ao PATH do usuário
$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
    Write-Host "   Adicionado ao PATH do usuário" -ForegroundColor $Green
}

# Mostrar resumo
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Instalação concluída!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Executável: $ExePath" -ForegroundColor White
Write-Host "  Notas em:   $NotesDir" -ForegroundColor White
Write-Host ""
Write-Host "  Para usar:" -ForegroundColor Yellow
Write-Host "    - Clique no atalho na área de trabalho" -ForegroundColor White
Write-Host "    - Ou execute: ketsu" -ForegroundColor White
Write-Host ""
Write-Host "  Comandos:" -ForegroundColor Yellow
Write-Host "    \    - Modo Shell (digitar key = value)" -ForegroundColor White
Write-Host "    :    - Modo Comando Vim" -ForegroundColor White
Write-Host "    ?    - Busca" -ForegroundColor White
Write-Host "    Tab  - Alternar explorer" -ForegroundColor White
Write-Host "    Ctrl+G - Voltar ao explorer" -ForegroundColor White
Write-Host "    Ctrl+S - Salvar" -ForegroundColor White
Write-Host ""

# Abrir pasta de notas
Invoke-Item $NotesDir
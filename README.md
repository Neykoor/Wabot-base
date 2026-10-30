# WaBot - Base de bot de WhatsApp en Go

Bot de WhatsApp en Go con [whatsmeow](https://github.com/tulir/whatsmeow), vinculado por código de emparejamiento (sin QR).

## Estructura

| Archivo | Función |
|---|---|
| `main.go` | Conexión y vinculación por código |
| `bot.go` | Eventos, despacho de comandos, `Reply` y `Send` |
| `commands.go` | Comandos de ejemplo y registro |
| `config.go` | Configuración por variables de entorno |
| `go.mod` | Módulo de Go |

## Instalación en Termux

### 1. Preparar Termux

```
pkg update && pkg upgrade -y
pkg install -y golang git clang make tmux
```

### 2. Descomprimir el proyecto

Copia `wabot-base.zip` a tu almacenamiento y descomprímelo:

```
termux-setup-storage
pkg install -y unzip
cp ~/storage/downloads/wabot-base.zip ~
unzip wabot-base.zip
cd wabot
```

### 3. Descargar dependencias

```
export CGO_ENABLED=1
export CC=clang
go mod tidy
```

La primera compilación de `go-sqlite3` tarda unos minutos. Es normal.

### 4. Compilar

```
go build -o wabot .
```

### 5. Ejecutar y vincular

```
./wabot
```

La primera vez te pide el número del bot con código de país, sin `+` ni espacios (ejemplo: `521234567890`). Luego muestra un código de 8 caracteres. En el teléfono del bot:

1. WhatsApp > Dispositivos vinculados > Vincular dispositivo
2. Elige **Vincular con el número de teléfono**
3. Escribe el código mostrado en Termux

Cuando aparezca `Vinculado correctamente` ya queda lista la sesión en `session.db`.

### 6. Mantenerlo activo en segundo plano

```
termux-wake-lock
tmux new -s bot
./wabot
```

Para salir sin cerrar el bot: `Ctrl+B` y luego `D`. Para volver: `tmux attach -t bot`.

## Variables de entorno opcionales

| Variable | Por defecto | Descripción |
|---|---|---|
| `BOT_PHONE` | (pregunta) | Número para vincular sin prompt |
| `BOT_PREFIX` | `!` | Prefijo de comandos |
| `BOT_NAME` | `WaBot` | Nombre en el menú |
| `BOT_SESSION_DB` | `session.db` | Archivo de sesión |
| `BOT_LOG_LEVEL` | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |

Ejemplo:

```
BOT_PREFIX="." BOT_NAME="MiBot" ./wabot
```

## Comandos incluidos

`!ping`, `!menu`, `!echo <texto>`, `!id`, `!grupo`

## Agregar un comando

En `commands.go`, dentro de `registerCommands`:

```go
b.Register(&Command{
	Name:        "hola",
	Description: "Saluda",
	Run: func(c *Ctx) {
		c.Reply("Hola!")
	},
})
```

## Problemas comunes

- **Error de CGO / gcc not found**: instala `clang` y exporta `CGO_ENABLED=1` y `CC=clang`.
- **Sesión cerrada (`LoggedOut`)**: borra `session.db` y vuelve a vincular.
- **El código no aparece o falla**: verifica el número con código de país y que tengas internet estable.
- **Quitar el bot**: elimina el dispositivo desde WhatsApp > Dispositivos vinculados.

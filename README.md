<h1 align="center">WaBot Base</h1>

<p align="center">
  <b>Base de bot de WhatsApp desarrollada en Go</b>
  <br>
  Ligera, sencilla y personalizable.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/Platform-Termux-000000?style=for-the-badge&logo=gnometerminal" />
  <img src="https://img.shields.io/badge/Library-whatsmeow-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
</p>

---

## ✧ Sobre el proyecto

WaBot Base es una plantilla de bot de WhatsApp creada en Go con [whatsmeow](https://github.com/tulir/whatsmeow).

Su objetivo es ofrecer una base sencilla para desarrollar bots de WhatsApp: agregar comandos, personalizar funciones y ampliar el proyecto sin estructuras complejas.

La vinculación se hace mediante código de emparejamiento, sin escanear QR.

## ✧ Características

- Desarrollado completamente en Go.
- Conexión mediante whatsmeow.
- Vinculación por código de emparejamiento, con regeneración automática si la conexión se corta o el código expira.
- Sistema de comandos con alias y restricción a grupos.
- Respuestas con `Reply` y `Send`.
- Configuración por variables de entorno.
- Sesión almacenada localmente en SQLite.
- Compatible con Termux.

## ✧ Estructura del proyecto

| Archivo | Descripción |
|---|---|
| `main.go` | Inicio, conexión y vinculación |
| `bot.go` | Eventos y despacho de comandos |
| `commands.go` | Registro y comandos de ejemplo |
| `config.go` | Configuración del bot |
| `go.mod` | Módulo y dependencias |
| `go.sum` | Se genera al ejecutar `go mod tidy` |

## ✧ Requisitos

- Android con Termux, o cualquier entorno compatible con Go.
- Conexión a Internet.
- Go, Git, Clang y Make (CGO es necesario para SQLite).

## ✧ Instalación en Termux

**1. Actualizar paquetes e instalar herramientas**

```bash
pkg update && pkg upgrade -y
pkg install -y golang git clang make tmux
```

**2. Descargar el repositorio**

```bash
cd ~
git clone https://github.com/Neykoor/Wabot-base.git wabot
cd ~/wabot
```

**3. Configurar el entorno, descargar dependencias y compilar**

```bash
export CGO_ENABLED=1
export CC=clang
go mod tidy
go build -o wabot .
```

La primera compilación de `go-sqlite3` tarda unos minutos. Es normal.

**4. Ejecutar**

```bash
termux-wake-lock
./wabot
```

## ✧ Vinculación de WhatsApp

Al iniciar por primera vez, el bot pide un número con código de país, por ejemplo `521234567890`. Después muestra un código de emparejamiento.

1. Abre WhatsApp.
2. Entra en **Dispositivos vinculados**.
3. Toca **Vincular dispositivo**.
4. Elige **Vincular con el número de teléfono**.
5. Introduce el código que muestra WaBot.

Al terminar, la sesión queda guardada en `session.db`.

> **Importante:** no compartas `session.db`, contiene datos sensibles de tu conexión.

### Consejos para vincular sin errores

- Ejecuta `termux-wake-lock` antes de iniciar. Si Android pone Termux en segundo plano, corta la red y el código deja de servir.
- Desactiva la optimización de batería para Termux (Ajustes, Apps, Termux, Batería, Sin restricciones).
- Usa pantalla dividida o ventana flotante para ver Termux mientras ingresas el código.
- Si el código expira o se corta la conexión, WaBot genera uno nuevo automáticamente.
- No lo ejecutes una y otra vez seguidas: WhatsApp limita los códigos (ver el error 429 más abajo).

## ✧ Configuración

| Variable | Valor predeterminado | Descripción |
|---|---|---|
| `BOT_PHONE` | vacío | Número para vincular sin preguntarlo |
| `BOT_PREFIX` | `!` | Prefijo de comandos |
| `BOT_NAME` | `WaBot` | Nombre del bot |
| `BOT_SESSION_DB` | `session.db` | Archivo de sesión |
| `BOT_LOG_LEVEL` | `INFO` | Nivel de registros |

Ejemplo:

```bash
BOT_PREFIX="." BOT_NAME="MiBot" ./wabot
```

## ✧ Comandos disponibles

| Comando | Alias | Descripción |
|---|---|---|
| `!ping` | `!p` | Mide la latencia del bot |
| `!menu` | `!help`, `!ayuda` | Muestra los comandos disponibles |
| `!echo <texto>` | | Repite el texto que escribas |
| `!id` | | Muestra el JID del chat y del remitente |
| `!grupo` | | Información básica del grupo (solo en grupos) |

## ✧ Crear nuevos comandos

Los comandos se registran en `commands.go`:

```go
b.Register(&Command{
    Name:        "hola",
    Aliases:     []string{"hi"},
    Description: "Saluda",
    Run: func(c *Ctx) {
        c.Reply("¡Hola!")
    },
})
```

Después de modificar el código, compila de nuevo:

```bash
go build -o wabot .
```

## ✧ Actualizar el bot

```bash
cd ~/wabot
git pull
go mod tidy
go build -o wabot .
./wabot
```

Detén el bot antes de actualizar y no borres `session.db`.

> `git pull` solo funciona si instalaste con `git clone`. Si descomprimiste un zip, la carpeta no es un repositorio: reemplaza los archivos a mano.

## ✧ Ejecutar en segundo plano

```bash
termux-wake-lock
tmux new -s bot
cd ~/wabot
./wabot
```

- Salir sin detener el proceso: `CTRL + B` y después `D`.
- Volver a la sesión: `tmux attach -t bot`.

## ✧ Solución de problemas

**Error de CGO o del compilador**

```bash
export CGO_ENABLED=1
export CC=clang
```

**Error al descargar dependencias**

```bash
go mod tidy
```

**`./wabot: No such file or directory`**

El ejecutable no existe todavía. Compila con `go build -o wabot .`

**`fatal: not a git repository`**

La carpeta se creó desde un zip. Clona el repositorio con `git clone` o actualiza los archivos manualmente.

**`status 429: rate-overlimit`**

WhatsApp bloqueó temporalmente la generación de códigos por demasiados intentos. Cierra el bot con `CTRL + C`, espera al menos 1 hora sin intentar de nuevo y vuelve a ejecutar `./wabot`. Cada intento durante el bloqueo lo prolonga.

**`software caused connection abort`**

Android cortó la red de Termux, normalmente por pasar a otra app. Usa `termux-wake-lock`, quita la optimización de batería y ver Termux en pantalla dividida mientras vinculas.

**Sesión cerrada desde el teléfono**

Borra la base de sesión y vincula otra vez:

```bash
rm -f session.db
./wabot
```

## ✧ Contribuciones

Las contribuciones, sugerencias y reportes de errores son bienvenidos. Puedes crear un Fork y enviar un Pull Request.

## ✧ Autor

**Neykoor**

- GitHub: [@Neykoor](https://github.com/Neykoor)
- Repositorio: [Neykoor/Wabot-base](https://github.com/Neykoor/Wabot-base)

---

<p align="center">
  Hecho con Go y dedicación.
  <br>
  <b>WaBot Base</b>
</p>

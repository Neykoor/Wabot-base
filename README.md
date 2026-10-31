WaBot - Base de bot de WhatsApp en Go

Bot de WhatsApp desarrollado en Go utilizando "whatsmeow" (https://github.com/tulir/whatsmeow). Incluye vinculación mediante código de emparejamiento, sin necesidad de QR.

Estructura del proyecto

Archivo| Función
"main.go"| Conexión y vinculación por código
"bot.go"| Eventos, comandos, Reply y Send
"commands.go"| Comandos de ejemplo y registro
"config.go"| Configuración mediante variables de entorno
"go.mod"| Dependencias del proyecto
"go.sum"| Versiones verificadas de dependencias

Instalación en Termux

1. Preparar Termux

Actualizar los paquetes:

pkg update

pkg upgrade -y

Instalar las herramientas necesarias:

pkg install golang

pkg install git

pkg install clang

pkg install make

pkg install tmux

2. Descargar el proyecto desde GitHub

Ir al directorio principal:

cd ~

Clonar el repositorio:

git clone https://github.com/Neykoor/Wabot-base.git wabot

Entrar en la carpeta:

cd ~/wabot

3. Preparar las dependencias

Habilitar CGO:

export CGO_ENABLED=1

Configurar Clang:

export CC=clang

Descargar y organizar las dependencias:

go mod tidy

4. Compilar el bot

go build -o wabot .

5. Iniciar el bot

./wabot

La primera vez solicitará el número de teléfono con código de país, sin espacios ni signo "+".

Ejemplo:

521234567890

Después aparecerá un código de emparejamiento de 8 caracteres.

En WhatsApp:

1. Abrir Dispositivos vinculados.
2. Seleccionar Vincular dispositivo.
3. Elegir Vincular con número de teléfono.
4. Introducir el código mostrado en Termux.

Al finalizar aparecerá el mensaje de vinculación correcta y se guardará la sesión en "session.db".

Mantener el bot activo

Activar el bloqueo de suspensión:

termux-wake-lock

Crear una sesión de TMUX:

tmux new -s bot

Iniciar el bot:

cd ~/wabot

./wabot

Para salir sin cerrar el bot:

"CTRL + B" y después "D".

Para volver a la sesión:

tmux attach -t bot

Actualizar el bot desde GitHub

Entrar en la carpeta:

cd ~/wabot

Descargar los últimos cambios:

git pull

Actualizar dependencias:

go mod tidy

Compilar nuevamente:

go build -o wabot .

Iniciar el bot:

./wabot

Importante: si el bot está ejecutándose, primero debes detenerlo antes de actualizarlo. No elimines "session.db", ya que contiene la sesión de WhatsApp.

Variables de entorno

Variable| Valor predeterminado| Descripción
"BOT_PHONE"| Pregunta al iniciar| Número de vinculación
"BOT_PREFIX"| "!"| Prefijo de comandos
"BOT_NAME"| "WaBot"| Nombre del bot
"BOT_SESSION_DB"| "session.db"| Archivo de sesión
"BOT_LOG_LEVEL"| "INFO"| Nivel de registros

Ejemplo de configuración:

BOT_PREFIX="." BOT_NAME="MiBot" ./wabot

Comandos incluidos

Comando| Función
"!ping"| Comprobar respuesta
"!menu"| Mostrar menú
"!echo <texto>"| Repetir mensaje
"!id"| Mostrar identificador
"!grupo"| Información del grupo

Agregar nuevos comandos

Abrir "commands.go" y agregar el siguiente código dentro de "registerCommands":

b.Register(&Command{
    Name:        "hola",
    Description: "Saluda",
    Run: func(c *Ctx) {
        c.Reply("¡Hola!")
    },
})

Después de modificar el código, volver a compilar:

go build -o wabot .

Problemas comunes

Error de CGO o GCC no encontrado

pkg install clang

export CGO_ENABLED=1

export CC=clang

Sesión cerrada ("LoggedOut")

Será necesario volver a vincular el dispositivo. No elimines la sesión sin comprobar primero el motivo del error.

El código de vinculación no aparece

- Verificar el número con código de país.
- Comprobar la conexión a Internet.
- Revisar que WhatsApp permita vincular otro dispositivo.

Eliminar el bot

Puedes detenerlo y eliminar la carpeta del proyecto cuando ya no necesites sus archivos. Ten presente que perderás la sesión local si eliminas "session.db".

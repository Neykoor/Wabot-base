WaBot Base

<p align="center">
  **Base de bot de WhatsApp desarrollada en Go**
  <br>
  Ligera, sencilla y personalizable.
</p><p align="center">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/Platform-Termux-000000?style=for-the-badge&logo=terminal" />
  <img src="https://img.shields.io/badge/Library-whatsmeow-blue?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
</p>---

✧ Sobre el proyecto

WaBot Base es una plantilla de bot de WhatsApp creada en Go utilizando "whatsmeow" (https://github.com/tulir/whatsmeow).

Su objetivo es proporcionar una base sencilla para desarrollar bots de WhatsApp, permitiendo agregar comandos, personalizar funciones y ampliar el proyecto sin depender de estructuras demasiado complejas.

La vinculación se realiza mediante código de emparejamiento, sin necesidad de escanear un código QR.

✧ Características

- Desarrollado completamente en Go.
- Conexión mediante whatsmeow.
- Vinculación mediante código de emparejamiento.
- Sistema de comandos personalizable.
- Respuestas mediante "Reply" y "Send".
- Configuración mediante variables de entorno.
- Almacenamiento local de sesión.
- Compatible con Termux.
- Código organizado y fácil de modificar.

✧ Estructura del proyecto

Archivo| Descripción
"main.go"| Inicio, conexión y vinculación
"bot.go"| Eventos y administración de comandos
"commands.go"| Registro y comandos de ejemplo
"config.go"| Configuración del bot
"go.mod"| Dependencias
"go.sum"| Verificación de dependencias

✧ Requisitos

- Android con Termux o un entorno compatible con Go.
- Conexión a Internet.
- Go.
- Git.
- Clang y Make para compilaciones que necesiten CGO.

✧ Instalación en Termux

1. Actualizar paquetes

pkg update

pkg upgrade -y

2. Instalar herramientas

pkg install golang

pkg install git

pkg install clang

pkg install make

pkg install tmux

3. Descargar el repositorio

cd ~

git clone https://github.com/Neykoor/Wabot-base.git wabot

4. Entrar en el proyecto

cd ~/wabot

5. Configurar el entorno

export CGO_ENABLED=1

export CC=clang

6. Instalar y organizar dependencias

go mod tidy

7. Compilar

go build -o wabot .

8. Ejecutar

./wabot

✧ Vinculación de WhatsApp

Al iniciar por primera vez, el bot solicitará un número telefónico con código de país.

Ejemplo:

521234567890

Después de introducir el número, aparecerá un código de emparejamiento.

Para vincular:

1. Abre WhatsApp.
2. Entra en Dispositivos vinculados.
3. Selecciona Vincular dispositivo.
4. Elige Vincular con número de teléfono.
5. Introduce el código proporcionado por WaBot.

Una vez completado el proceso, la sesión quedará almacenada en "session.db".

Importante: no compartas tu archivo de sesión, ya que contiene información sensible de tu conexión.

✧ Configuración

WaBot permite modificar su comportamiento mediante variables de entorno.

Variable| Valor predeterminado| Descripción
"BOT_PHONE"| Vacío| Número para vincular
"BOT_PREFIX"| "!"| Prefijo de comandos
"BOT_NAME"| "WaBot"| Nombre del bot
"BOT_SESSION_DB"| "session.db"| Archivo de sesión
"BOT_LOG_LEVEL"| "INFO"| Nivel de registros

Ejemplo de configuración:

BOT_PREFIX="." BOT_NAME="MiBot" ./wabot

✧ Comandos disponibles

Comando| Descripción
"!ping"| Comprueba la respuesta del bot
"!menu"| Muestra el menú
"!echo <texto>"| Repite un mensaje
"!id"| Muestra el identificador
"!grupo"| Muestra información del grupo

✧ Crear nuevos comandos

Los comandos se registran en "commands.go".

Ejemplo:

b.Register(&Command{
    Name:        "hola",
    Description: "Saluda",
    Run: func(c *Ctx) {
        c.Reply("¡Hola!")
    },
})

Después de modificar el código, vuelve a compilar:

go build -o wabot .

✧ Actualizar el bot

Para obtener los últimos cambios del repositorio:

1. Entrar en el proyecto

cd ~/wabot

2. Descargar actualizaciones

git pull

3. Actualizar dependencias

go mod tidy

4. Compilar nuevamente

go build -o wabot .

5. Iniciar

./wabot

Si el bot está ejecutándose, detenlo antes de actualizarlo.

No elimines "session.db" durante las actualizaciones.

✧ Ejecutar en segundo plano

Puedes utilizar TMUX para mantener el proceso activo mientras utilizas otras terminales.

Activar bloqueo de suspensión:

termux-wake-lock

Crear una sesión:

tmux new -s bot

Entrar al proyecto:

cd ~/wabot

Ejecutar:

./wabot

Para salir sin detener el proceso:

"CTRL + B" y después "D".

Para regresar:

tmux attach -t bot

✧ Solución de problemas

Error de CGO o compilador

Comprueba que Clang esté instalado y configura:

export CGO_ENABLED=1

export CC=clang

Error al descargar dependencias

go mod tidy

El ejecutable no existe

Vuelve a compilar:

go build -o wabot .

Sesión desconectada

Comprueba el estado de la conexión y el registro de errores. Si WhatsApp cerró la sesión, será necesario realizar una nueva vinculación.

---

✧ Contribuciones

Las contribuciones, sugerencias y reportes de errores son bienvenidos.

Si deseas mejorar el proyecto, puedes crear un Fork y enviar un Pull Request.

✧ Autor

Neykoor

GitHub: "@Neykoor" (https://github.com/Neykoor)

Repositorio: "Neykoor/Wabot-base" (https://github.com/Neykoor/Wabot-base)

---

<p align="center">
  Hecho con Go y dedicación.
  <br>
  **WaBot Base**
</p>

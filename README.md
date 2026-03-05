# Gestor de Procesos y Concurrencia

# Dairon Stivens Layos Rico

El proyecto es un gestor de proceso y concurrencias que simulará procesos, utilizara concurrencia real con goroutines, implementara sincronización y demostrara interacción con el SO real

## Qué funciona (E1)
- Se crea 3 proceso con burst aleatorios
- Muestra en consola los procesos creados con su información.
- Crea metricas y las muestra en consola.

## Estructura
- src/procesos/proceso.go
- docs/propuesta.pdf
- Simulacion de procesos.drawio

## Requisitos
- Go 1.26 o superior
- Sistema operativo: Windows, Linux o macOS

## Instalación y Ejecución

```bash
# Clonar o crear el proyecto
mkdir gestion-procesos
cd gestion-procesos

# Ejecutar
go run main.go O go run .

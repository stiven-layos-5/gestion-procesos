# Gestor de Procesos y Concurrencia

# Dairon Stivens Layos Rico

El proyecto es un gestor de proceso y concurrencias que simulará procesos, utilizara concurrencia real con goroutines, implementara sincronización y demostrara interacción con el SO real

## Qué funciona (E1)

- Se crea 3 proceso con burst aleatorios
- Muestra en consola los procesos creados con su información.
- Crea metricas y las muestra en consola.

## Estructura

```
gestion-procesos/
├── src
│   ├── services/
│       ├── scheduler/
│       │   └── scheduler.go         ← Módulo 1: Simulación académica (Round-Robin + Prioridad)
│       ├── process/
│       │   └── manager.go           ← Módulo 2: Procesos reales del SO (ps, start, kill)
│       ├── concurrency/
│       │   └── worker_pool.go       ← Módulo 3: Worker Pool con goroutines reales
│       ├── sync/
│       │   └── sincronizacion.go    ← Módulo 4: Race condition + Mutex + Semáforo
│       └── monitor/
│           └── monitor.go           ← Módulo 5: CPU, Memoria y Goroutines en tiempo real
│   ├── models/
│       └── process.go               ← Modelos compartidos (ProcesoSimulado, ProcesoReal, etc.)
├── go.mod
├── main.go
└── README.md
```

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
```

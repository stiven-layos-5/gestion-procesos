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
#

# Ejecutar
go run main.go
go run .
go run -race .

```

## Conceptos claves

### Módulo 1 — Planificación

- **Estado de proceso:** NUEVO → LISTO → EJECUTANDO → (ESPERANDO) → TERMINADO
- **Round-Robin:** Quantum fijo, equitativo, apropiativo por tiempo
- **Prioridad:** Menor número = mayor prioridad; no apropiativo en esta impl.
- **Métricas:** Tiempo de espera = T.Fin - T.Llegada - BurstTime

### Módulo 2 — Procesos reales

- **API del SO:** `exec.Cmd` (fork+exec), `os.FindProcess`, `Process.Kill`
- **`ps aux`:** Información de USER, PID, STAT, COMMAND
- **Señal de terminación:** `Kill()` envía SIGKILL al proceso destino

### Módulo 3 — Concurrencia

- **Goroutine:** Hilo ligero administrado por el runtime de Go (M:N threading)
- **Channel:** Pipe tipado para comunicación entre goroutines (CSP model)
- **Worker Pool:** Patrón que limita goroutines activas para controlar recursos

### Módulo 4 — Sincronización

- **Race condition:** Dos goroutines acceden al mismo dato sin coordinación
- **Sección crítica:** El bloque `leer → incrementar → escribir` del contador
- **Mutex:** `Lock()` garantiza que solo una goroutine esté en la sección crítica
- **Semáforo:** Channel buffereado de capacidad N, `send` = adquirir, `recv` = liberar

### Módulo 5 — Monitoreo

- **CPU:** Diferencia de ticks idle/total entre dos lecturas de `/proc/stat`
- **Memoria:** `runtime.MemStats.Alloc` (heap en uso), `Sys` (total del sistema)
- **Goroutines:** `runtime.NumGoroutine()` — valor real del runtime de Go

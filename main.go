package main

import (
	"fmt"
	"time"

	"gestion-procesos/src/internal/concurrency"
	"gestion-procesos/src/internal/monitor"
	"gestion-procesos/src/internal/process"
	"gestion-procesos/src/internal/scheduler"
	sincronizacion "gestion-procesos/src/internal/sync"
	"gestion-procesos/src/models"
)

func separador(titulo string) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Printf( "║  %-62s║\n", titulo)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║        GESTOR DE PROCESOS Y CONCURRENCIA — Entrega 3        ║")
	fmt.Println("║              Modalidad 1: Proyecto Base (Go)                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	// ═══════════════════════════════════════════════════════
	// MÓDULO 1: Simulación académica
	// ═══════════════════════════════════════════════════════
	separador("MÓDULO 1 — Simulación académica de planificación de CPU")

	fmt.Println("\n▶ Algoritmo: Round-Robin (quantum = 3)")
	planRR := scheduler.NuevoPlanificador(scheduler.RoundRobin, 3)
	procesosRR := []*models.ProcesoSimulado{
		{ID: 1, Nombre: "Editor",    BurstTime: 8},
		{ID: 2, Nombre: "Compilador",BurstTime: 4},
		{ID: 3, Nombre: "Linker",    BurstTime: 5},
		{ID: 4, Nombre: "Debugger",  BurstTime: 3},
	}
	for _, p := range procesosRR {
		planRR.AgregarProceso(p)
	}
	planRR.Ejecutar()

	fmt.Println("\n▶ Algoritmo: Prioridad (menor número = mayor prioridad)")
	planPrio := scheduler.NuevoPlanificador(scheduler.Prioridad, 0)
	procesosPrio := []*models.ProcesoSimulado{
		{ID: 1, Nombre: "Kernel",    BurstTime: 6, Prioridad: 1},
		{ID: 2, Nombre: "Servicio",  BurstTime: 4, Prioridad: 3},
		{ID: 3, Nombre: "App",       BurstTime: 5, Prioridad: 5},
		{ID: 4, Nombre: "Daemon",    BurstTime: 3, Prioridad: 2},
	}
	for _, p := range procesosPrio {
		planPrio.AgregarProceso(p)
	}
	planPrio.Ejecutar()

	// ═══════════════════════════════════════════════════════
	// MÓDULO 2: Procesos reales
	// ═══════════════════════════════════════════════════════
	separador("MÓDULO 2 — Gestión de procesos reales del SO")

	mgr := process.NuevoManager()

	fmt.Println("\n▶ Listando procesos reales del sistema operativo...")
	procesos, err := mgr.Listar()
	if err != nil {
		fmt.Printf("  Error al listar: %v\n", err)
	} else {
		process.ImprimirProcesos(procesos)
	}

	fmt.Println("\n▶ Iniciando proceso real: 'sleep 2'")
	pid, err := mgr.Iniciar("sleep", "2")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Proceso 'sleep 2' corriendo con PID: %d\n", pid)
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("\n▶ Finalizando proceso PID %d...\n", pid)
		if err := mgr.Finalizar(pid); err != nil {
			fmt.Printf("  Error al finalizar: %v\n", err)
		}
	}

	// ═══════════════════════════════════════════════════════
	// MÓDULO 3: Concurrencia
	// ═══════════════════════════════════════════════════════
	separador("MÓDULO 3 — Concurrencia real con Worker Pool (goroutines)")

	fmt.Println("\n▶ Creando Worker Pool con 3 workers y 10 tareas...")
	pool := concurrency.NuevoWorkerPool(3)
	tareas := concurrency.GenerarTareas(10)
	pool.Iniciar()
	pool.EnviarTareas(tareas)
	resultados := pool.RecolectarResultados()

	fmt.Printf("\n  Resumen: %d tareas completadas por %d workers concurrentes.\n",
		len(resultados), 3)

	// ═══════════════════════════════════════════════════════
	// MÓDULO 4: Sincronización
	// ═══════════════════════════════════════════════════════
	separador("MÓDULO 4 — Sincronización: race condition y mecanismos de control")

	fmt.Println("\n▶ PASO 1: Condición de carrera SIN protección")
	fmt.Println("  (El resultado incorrecto demuestra el problema)")
	sincronizacion.DemostrarCondicionCarrera(10, 100)

	fmt.Println("▶ PASO 2: Sección crítica protegida con sync.Mutex")
	sincronizacion.DemostrarMutex(10, 100)

	fmt.Println("▶ PASO 3: Control de acceso con Semáforo (channel buffereado)")
	sincronizacion.DemostrarSemaforo(3, 8)

	// ═══════════════════════════════════════════════════════
	// MÓDULO 5: Monitoreo de recursos
	// ═══════════════════════════════════════════════════════
	separador("MÓDULO 5 — Monitoreo de recursos reales del sistema")

	fmt.Println("\n▶ Tomando 5 snapshots de CPU, Memoria y Goroutines...")
	mon := monitor.NuevoMonitor(500 * time.Millisecond)
	mon.Monitorear(5)

	fmt.Println("\n▶ Resumen del historial de monitoreo:")
	mon.ResumenHistorial()

	// ═══════════════════════════════════════════════════════
	// FIN
	// ═══════════════════════════════════════════════════════
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                 Ejecución completada exitosamente           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}

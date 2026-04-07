package main

import (
	"fmt"
	"gestion-procesos/src/logger"
	proceso "gestion-procesos/src/procesos"
	simulator "gestion-procesos/src/simulador"
	"math/rand"
	"time"
)

// func main() {

// 	archivoLog, logger := setupLogger()

// 	defer archivoLog.Close()

// 	logger.Printf("=== INICIO DE SIMULACIÓN DE CONCURRENCIA ===")
// 	fmt.Println("\n╔══════════════════════════════════════════╗")
// 	fmt.Println("║   SIMULACIÓN DE CONCURRENCIA EN GO       ║")
// 	fmt.Println("╚══════════════════════════════════════════╝\n")

// 	fmt.Println("━━━ Procesos concurrentes con sincronización ━━━")
// 	logger.Printf("--- Goroutines con WaitGroup ---")

// 	almacen := concurrencia.NuevoAlmacen(logger)

// 	procesos := []concurrencia.Proceso{
// 		{ID: 1, Nombre: "CPU-Scheduler", Duracion: 300 * time.Millisecond},
// 		{ID: 2, Nombre: "IO-Manager", Duracion: 500 * time.Millisecond},
// 		{ID: 3, Nombre: "Memory-Mgr", Duracion: 200 * time.Millisecond},
// 		{ID: 4, Nombre: "Network-Svc", Duracion: 400 * time.Millisecond},
// 	}

// 	duracion := concurrencia.EjecutarProcesos(procesos, almacen, logger)

// 	logger.Printf("Todos los procesos terminaron en: %v", duracion)
// 	fmt.Printf("\n  Tiempo total (paralelo): %v\n", duracion)
// 	fmt.Printf("   (Sin concurrencia sería 1400ms)\n\n")

// 	fmt.Println(" Estado del almacén compartido:")
// 	for _, p := range procesos {
// 		if val, ok := almacen.Leer(p.Nombre); ok {
// 			fmt.Printf("   %s → %d\n", p.Nombre, val)
// 		}
// 	}

// 	fmt.Println("\n━━━ PARTE 2: Productor-Consumidor con canales ━━━")
// 	logger.Printf("--- Parte 2: Canales productor-consumidor ---")

// 	tareas := []string{"Tarea-A", "Tarea-B", "Tarea-C", "Tarea-D"}
// 	concurrencia.EjecutarProductorConsumidor(tareas, 2, logger)

// 	// --- Fin ---
// 	logger.Printf("=== SIMULACIÓN COMPLETADA ===")
// 	fmt.Println("\n Simulación completada.")
// 	fmt.Println(" Log guardado en: simulacion.txt")

// }

func main() {

	log, err := logger.NewLogger("logs/simulacion.log")

	if err != nil {
		panic(err)
	}

	defer log.Close()

	const quantum = 3

	procesos := []proceso.Proceso{}

	rand.Seed(time.Now().UnixNano())

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================")

	fmt.Println("[COMPONENTE 1] Simulación de Procesos")
	fmt.Println("----------------------------------------")

	for i := 0; i < 5; i++ {

		min := 1
		max := 11

		duracion := rand.Intn(max-min+1) + min
		nombre := i

		procesos = append(procesos, proceso.NuevoProceso(nombre, duracion, i))
	}

	simulator.RoundRobin(procesos, quantum)
}

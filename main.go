package main

import (
	"fmt"
	"gestion-procesos/src/logger"
	proceso "gestion-procesos/src/procesos"
	"gestion-procesos/src/scheduler"
	"math/rand"
	"time"
)

func main() {

	log, err := logger.NewLogger("logs/simulacion.log")

	if err != nil {
		panic(err)
	}

	defer log.Close()

	const quantum = 3

	sched := scheduler.NewScheduler(quantum, 16, log)

	rand.Seed(time.Now().UnixNano())

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================")

	for i := 0; i < 5; i++ {

		min := 1
		max := 11

		duracion := rand.Intn(max-min+1) + min
		tamañoMemoria := rand.Intn(max-min+1) + min

		proc := &proceso.Proceso{
			PID:        i,
			Llegada:    i,
			Rafaga:     duracion,
			Restante:   duracion,
			TamMemoria: tamañoMemoria,
		}

		sched.AddProcess(proc)
	}

	fmt.Println("=== SIMULADOR CONCURRENTE ROUND ROBIN CON MEMORIA ===")
	fmt.Println("Revisar simulacion.log para eventos detallados.\n")
	sched.Run()

	sched.MostrarEstadisticas()

}

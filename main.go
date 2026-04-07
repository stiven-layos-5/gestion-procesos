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

	const quantum = 2
	const totalPaginas = 4

	sched := scheduler.NewScheduler(quantum, totalPaginas, log)

	rand.Seed(time.Now().UnixNano())

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================")

	// procs := []struct {
	// 	pid     int
	// 	llegada int
	// 	rafaga  int
	// 	tamMem  int
	// }{
	// 	{1, 0, 5, 8192},
	// 	{2, 1, 3, 12288},
	// 	{3, 2, 8, 4096},
	// 	{4, 3, 6, 20480},
	// }

	// for _, p := range procs {
	// 	proc := &proceso.Proceso{
	// 		PID:        p.pid,
	// 		Llegada:    p.llegada,
	// 		Rafaga:     p.rafaga,
	// 		Restante:   p.rafaga,
	// 		TamMemoria: p.tamMem,
	// 	}
	// 	sched.AddProcess(proc)
	// }

	for i := 0; i < 10; i++ {

		rafagaMin := 1
		rafagaMax := 11

		memoryMin := 1
		memoryMax := 5

		duracion := rand.Intn(rafagaMax-rafagaMin+1) + rafagaMin
		tamañoMemoria := rand.Intn(memoryMax-memoryMin+1) + memoryMin

		proc := &proceso.Proceso{
			PID:        i + 1,
			Llegada:    i,
			Rafaga:     duracion,
			Restante:   duracion,
			TamMemoria: tamañoMemoria * 4096,
		}

		sched.AddProcess(proc)
	}

	fmt.Println("=== SIMULADOR CONCURRENTE ROUND ROBIN CON MEMORIA ===")
	fmt.Println("Revisar simulacion.log para eventos detallados.\n")
	sched.Run()

	sched.MostrarEstadisticas()

}

package main

import (
	"fmt"
	proceso "gestion-procesos/src/procesos"
	"math/rand"
	"time"
)

func main() {

	const quantum = 3

	procesos := []proceso.Proceso{}

	rand.Seed(time.Now().UnixNano())

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================")

	fmt.Println("[COMPONENTE 1] Simulación de Procesos")
	fmt.Println("----------------------------------------")

	for i := 0; i < 3; i++ {

		min := 1
		max := 11

		duracion := rand.Intn(max-min+1) + min
		nombre := fmt.Sprintf("P %d", i)

		procesos = append(procesos, proceso.NuevoProceso(nombre, duracion, i))
	}

	proceso.ImprimirProcesos(procesos)

	completed := proceso.EjecutarRoundRobin(procesos, quantum)

	metricas, promedioEspera, promedioRetorno := proceso.CalcularMetricas(completed)

	fmt.Println("\n╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║              METRICAS POR PROCESO               	║")
	fmt.Println("╠═══════╦════════════════════╦══════════════════════════╣")
	fmt.Println("║  PID  ║  Tiempo de espera  ║     Tiempo de Retorno    ║")
	fmt.Println("╠═══════╬════════════════════╬══════════════════════════╣")
	for _, m := range metricas {
		fmt.Printf("║  %-4s ║ %-10d	     ║ %-14d           ║\n",
			m.ID, m.TiempoEspera, m.TiempoRetorno)
	}
	fmt.Println("╚═══════╩════════════════════╩══════════════════════════╝")
	fmt.Printf("\n  Average waiting time:    %.2f ms\n", promedioEspera)
	fmt.Printf("  Average turnaround time: %.2f ms\n", promedioRetorno)

}

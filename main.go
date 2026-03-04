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
		nombre := fmt.Sprintf("Proceso %d", i)

		procesos = append(procesos, proceso.NuevoProceso(nombre, duracion, i))
	}

	proceso.ImprimirProcesos(procesos)

	proceso.EjecutarRoundRobin(procesos, quantum)

	fmt.Println()
}

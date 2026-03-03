package main

import (
	"fmt"
	"gestion-procesos/src/procesos"
	"math/rand"
	"time"
)

func main() {

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================")

	fmt.Println("[COMPONENTE 1] Simulación de Procesos")
	fmt.Println("----------------------------------------")
	simulador := procesos.NuevoSimuladorProcesos()

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 3; i++ {

		min := 1
		max := 11

		duracion := rand.Intn(max-min+1) + min
		nombre := fmt.Sprintf("Proceso %d", i)

		simulador.CrearProcesoSimulado(nombre, duracion)

	}

	//simulador.CrearProcesoSimulado("Proceso A", 5)
	//simulador.CrearProcesoSimulado("Proceso B", 3)
	//simulador.CrearProcesoSimulado("Proceso C", 7)

	simulador.EjecutarSimulacion()
	fmt.Println()
}

package main

import (
	"fmt"
	"gestion-procesos/src/procesos"
)

func main() {

	fmt.Println("========================================")
	fmt.Println("  GESTOR DE PROCESOS Y CONCURRENCIA")
	fmt.Println("========================================\n")

	fmt.Println("[COMPONENTE 1] Simulación de Procesos")
	fmt.Println("----------------------------------------")
	simulador := procesos.NuevoSimuladorProcesos()

	simulador.CrearProcesoSimulado("Proceso A", 5)
	simulador.CrearProcesoSimulado("Proceso B", 3)
	simulador.CrearProcesoSimulado("Proceso C", 7)

	simulador.EjecutarSimulacion()
	fmt.Println()
}

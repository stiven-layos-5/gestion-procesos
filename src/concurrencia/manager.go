package concurrencia

import (
	"fmt"
	proceso "gestion-procesos/src/procesos"
)

func StartConcurrentSimulation(procesos []proceso.Proceso, quantum int) {
	fmt.Println("\n=== SIMULACIÓN CONCURRENTE (Entrega 2) ===")
	fmt.Println("Próximamente: cada proceso será una goroutine,")
	fmt.Println("el planificador usará canales y sincronización con mutex.")

}

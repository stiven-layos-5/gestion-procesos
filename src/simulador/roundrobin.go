package simulator

import (
	"container/list"
	"fmt"
	proceso "gestion-procesos/src/procesos"
)

func RoundRobin(procesos []proceso.Proceso, quantum int) {
	tiempoActual := 0
	indice := 0
	completados := 0
	ordenEjecucion := []int{}
	n := len(procesos)
	cola := list.New()

	procs := make([]proceso.Proceso, n)
	copy(procs, procesos)

	for i := 0; i < n; i++ {
		procs[i].Restante = procs[i].Rafaga
		procs[i].PrimeraEjec = true
	}

	for completados < n {
		// Agregar procesos que han llegado
		for indice < n && procs[indice].Llegada <= tiempoActual {
			cola.PushBack(indice)
			indice++
		}

		if cola.Len() == 0 {
			if indice < n {
				tiempoActual = procs[indice].Llegada
				continue
			} else {
				break
			}
		}

		front := cola.Front()
		cola.Remove(front)
		iProc := front.Value.(int)

		if procs[iProc].PrimeraEjec {
			procs[iProc].PrimeraEjec = false
			procs[iProc].InicioPrimero = tiempoActual
		}

		tiempoEjec := quantum
		if procs[iProc].Restante < quantum {
			tiempoEjec = procs[iProc].Restante
		}

		ordenEjecucion = append(ordenEjecucion, procs[iProc].PID)
		procs[iProc].Restante -= tiempoEjec
		tiempoActual += tiempoEjec

		if procs[iProc].Restante == 0 {
			procs[iProc].Completado = true
			procs[iProc].Retorno = tiempoActual
			completados++
			procs[iProc].Espera = procs[iProc].Retorno - procs[iProc].Llegada - procs[iProc].Rafaga
		} else {
			// Reencolar después de agregar nuevos llegados
			for indice < n && procs[indice].Llegada <= tiempoActual {
				cola.PushBack(indice)
				indice++
			}
			cola.PushBack(iProc)
		}
	}

	fmt.Println("\n=== RESULTADOS ROUND ROBIN (secuencial) ===")
	fmt.Printf("Quantum: %d\n\n", quantum)
	fmt.Println("PID\tLlegada\tRáfaga\tEspera\tRetorno")
	fmt.Println("------------------------------------------------")

	sumaEspera := 0
	sumaRetorno := 0

	for _, p := range procs {
		fmt.Printf("%d\t%d\t%d\t%d\t%d\n", p.PID, p.Llegada, p.Rafaga, p.Espera, p.Retorno)
		sumaEspera += p.Espera
		sumaRetorno += p.Retorno
	}

	fmt.Println("------------------------------------------------")
	fmt.Printf("Tiempo promedio de espera: %.2f\n", float64(sumaEspera)/float64(n))
	fmt.Printf("Tiempo promedio de retorno: %.2f\n", float64(sumaRetorno)/float64(n))

	fmt.Println("\n--- Diagrama de Gantt (orden de ejecución) ---")

	for _, pid := range ordenEjecucion {
		fmt.Printf("| P%d ", pid)
	}

	fmt.Println("|")
}

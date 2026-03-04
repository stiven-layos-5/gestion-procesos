package proceso

import "fmt"

type Proceso struct {
	ID             string
	TiempoTotal    int
	TiempoRestante int
	TiempoLlegada  int
	TiempoInicio   int
	TiempoFin      int
	Completado     bool
}

type ResultadoEjecucion struct {
	ProcesoID string
	Inicio    int
	Fin       int
}

func NuevoProceso(id string, burst, arrival int) Proceso {

	return Proceso {
		ID: id,
		TiempoTotal: burst,
		TiempoRestante: burst,
		TiempoLlegada: arrival,
		TiempoInicio: -1,
	}
}

func EjecutarRoundRobin(procesos []Proceso, quantum int) ([]ResultadoEjecucion, []Proceso) {

	historial := []ResultadoEjecucion{}
	tiempoActual := 0
	completados := 0
	n := len(procesos)
	cola := []int{}

	for i := range procesos {
		if procesos[i].TiempoLlegada <= tiempoActual {
			cola = append(cola, i)
		}
	}

	for completados < n {
		
		if len(cola) == 0 {
			tiempoActual++
			for i := range procesos {
				if !procesos[i].Completado && procesos[i].TiempoLlegada == tiempoActual {
					cola = append(cola, i)
				}
			}
			continue
		}

		idx := cola[0]
		cola = cola[1:]
		p := &procesos[idx]

		if p.TiempoInicio == -1 {
			p.TiempoInicio = tiempoActual
		}

		ejecutar := quantum
		if p.TiempoRestante < quantum {
			ejecutar = p.TiempoRestante
		}

		inicio := tiempoActual
		tiempoActual += ejecutar
		p.TiempoRestante -= ejecutar

		historial = append(historial, ResultadoEjecucion{
			ProcesoID: p.ID,
			Inicio:    inicio,
			Fin:       tiempoActual,
		})

		for i := range procesos {
			if !procesos[i].Completado && i != idx &&
				procesos[i].TiempoLlegada > inicio &&
				procesos[i].TiempoLlegada <= tiempoActual {
				if !estaEnCola(cola, i) {
					cola = append(cola, i)
				}
			}
		}

		if p.TiempoRestante == 0 {
			p.Completado = true
			p.TiempoFin = tiempoActual
			completados++
		} else {
			cola = append(cola, idx)
		}
	}

	return historial, procesos
}

func ImprimirProcesos(procesos []Proceso) {
	fmt.Println("\n  Procesos cargados:")
	fmt.Println("  ──────────────────────────────────────")
	for _, p := range procesos {
		fmt.Printf("  • %-4s  Ráfaga: %3d ms  |  Llegada: t=%d\n",
			p.ID, p.TiempoTotal, p.TiempoLlegada)
	}
}

func estaEnCola(cola []int, idx int) bool {
	for _, c := range cola {
		if c == idx {
			return true
		}
	}
	return false
}
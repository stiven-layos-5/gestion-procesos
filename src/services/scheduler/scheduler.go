package scheduler

import (
	"fmt"
	"gestion-procesos/src/models"
)

type Algoritmo string

const (
	RoundRobin Algoritmo = "RoundRobin"
	Prioridad  Algoritmo = "Prioridad"
)

type Planificador struct {
	Algoritmo  Algoritmo
	Quantum    int // Solo para Round-Robin
	ColaListos []*models.ProcesoSimulado
	Terminados []*models.ProcesoSimulado
	tickActual int
}

func NuevoPlanificador(alg Algoritmo, quantum int) *Planificador {
	return &Planificador{
		Algoritmo: alg,
		Quantum:   quantum,
	}
}

func (p *Planificador) AgregarProceso(proc *models.ProcesoSimulado) {
	proc.Estado = models.EstadoNuevo
	fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s\n", p.tickActual, proc.ID, proc.Nombre, proc.Estado)
	proc.Estado = models.EstadoListo
	proc.TiempoLlegada = p.tickActual
	proc.TiempoRestante = proc.BurstTime
	fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s\n", p.tickActual, proc.ID, proc.Nombre, proc.Estado)
	p.ColaListos = append(p.ColaListos, proc)
}

func (p *Planificador) Ejecutar() {
	fmt.Printf("\n=== Iniciando planificador [%s] ===\n\n", p.Algoritmo)
	switch p.Algoritmo {
	case RoundRobin:
		p.ejecutarRoundRobin()
	case Prioridad:
		p.ejecutarPrioridad()
	}
	p.imprimirEstadisticas()
}

func (p *Planificador) ejecutarRoundRobin() {
	for len(p.ColaListos) > 0 {
		proc := p.ColaListos[0]
		p.ColaListos = p.ColaListos[1:]

		if proc.TiempoInicio == 0 && p.tickActual > 0 {
			proc.TiempoInicio = p.tickActual
		} else if p.tickActual == 0 {
			proc.TiempoInicio = 0
		}

		proc.Estado = models.EstadoEjecutando
		fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s (restante: %d)\n",
			p.tickActual, proc.ID, proc.Nombre, proc.Estado, proc.TiempoRestante)

		ejecutado := min(p.Quantum, proc.TiempoRestante)
		proc.TiempoRestante -= ejecutado
		p.tickActual += ejecutado

		if proc.TiempoRestante > 0 {
			proc.Estado = models.EstadoEsperando
			fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s (rebota a cola)\n",
				p.tickActual, proc.ID, proc.Nombre, proc.Estado)
			proc.Estado = models.EstadoListo
			p.ColaListos = append(p.ColaListos, proc)
		} else {
			proc.Estado = models.EstadoTerminado
			proc.TiempoFin = p.tickActual
			fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s\n",
				p.tickActual, proc.ID, proc.Nombre, proc.Estado)
			p.Terminados = append(p.Terminados, proc)
		}
	}
}

func (p *Planificador) ejecutarPrioridad() {
	for len(p.ColaListos) > 0 {
		idx := 0
		for i, proc := range p.ColaListos {
			if proc.Prioridad < p.ColaListos[idx].Prioridad {
				idx = i
			}
		}
		proc := p.ColaListos[idx]
		p.ColaListos = append(p.ColaListos[:idx], p.ColaListos[idx+1:]...)

		proc.TiempoInicio = p.tickActual
		proc.Estado = models.EstadoEjecutando
		fmt.Printf("  [Tick %02d] Proceso %d (%s) prioridad=%d → %s\n",
			p.tickActual, proc.ID, proc.Nombre, proc.Prioridad, proc.Estado)

		p.tickActual += proc.BurstTime
		proc.TiempoRestante = 0
		proc.Estado = models.EstadoTerminado
		proc.TiempoFin = p.tickActual
		fmt.Printf("  [Tick %02d] Proceso %d (%s) → %s\n",
			p.tickActual, proc.ID, proc.Nombre, proc.Estado)
		p.Terminados = append(p.Terminados, proc)
	}
}

func (p *Planificador) imprimirEstadisticas() {
	fmt.Println("\n--- Estadísticas de planificación ---")
	fmt.Printf("%-5s %-12s %-10s %-12s %-12s %-12s\n",
		"PID", "Nombre", "BurstTime", "T.Llegada", "T.Fin", "T.Espera")
	totalEspera := 0
	for _, proc := range p.Terminados {
		espera := proc.TiempoFin - proc.TiempoLlegada - proc.BurstTime
		if espera < 0 {
			espera = 0
		}
		totalEspera += espera
		fmt.Printf("%-5d %-12s %-10d %-12d %-12d %-12d\n",
			proc.ID, proc.Nombre, proc.BurstTime, proc.TiempoLlegada, proc.TiempoFin, espera)
	}
	if len(p.Terminados) > 0 {
		fmt.Printf("\nTiempo de espera promedio: %.2f ticks\n",
			float64(totalEspera)/float64(len(p.Terminados)))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

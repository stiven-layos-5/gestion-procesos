package models

import "time"

type Estado string

const (
	EstadoNuevo      Estado = "NUEVO"
	EstadoListo      Estado = "LISTO"
	EstadoEjecutando Estado = "EJECUTANDO"
	EstadoEsperando  Estado = "ESPERANDO"
	EstadoTerminado  Estado = "TERMINADO"
)

type ProcesoSimulado struct {
	ID             int
	Nombre         string
	Estado         Estado
	Prioridad      int
	BurstTime      int
	TiempoRestante int
	TiempoLlegada  int
	TiempoInicio   int
	TiempoFin      int
	Creado         time.Time
}

type ProcesoReal struct {
	PID     int
	Nombre  string
	Estado  string
	Memoria uint64
}

type RecursosSnapshot struct {
	Timestamp     time.Time
	CPUPorcentaje float64
	MemTotal      uint64
	MemUsada      uint64
	NumGoroutines int
}

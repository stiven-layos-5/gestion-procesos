package procesos

import (
	"fmt"
	"sync"
	"time"
)

type ProcesoSimulado struct {
	ID           int
	Nombre       string
	Duracion     int
	Estado       string
	TiempoInicio time.Time
	TiempoFin    time.Time
	mutex        sync.Mutex
}

type SimuladorProcesos struct {
	Procesos   []*ProcesoSimulado
	contadorID int
	mutex      sync.Mutex
}

func NuevoSimuladorProcesos() *SimuladorProcesos {
	return &SimuladorProcesos{
		Procesos:   make([]*ProcesoSimulado, 0),
		contadorID: 0,
	}
}

func (s *SimuladorProcesos) CrearProcesoSimulado(nombre string, duracion int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.contadorID++
	proceso := &ProcesoSimulado{
		ID:       s.contadorID,
		Nombre:   nombre,
		Duracion: duracion,
		Estado:   "Nuevo",
	}

	s.Procesos = append(s.Procesos, proceso)
	fmt.Printf("Proceso simulado creado: ID=%d, Nombre=%s, Duración=%ds\n",
		proceso.ID, proceso.Nombre, proceso.Duracion)
}

func (s *SimuladorProcesos) EjecutarSimulacion() {
	if len(s.Procesos) == 0 {
		fmt.Println("No hay procesos para simular")
		return
	}

	fmt.Println("\n--- SIMULACIÓN DE PROCESOS ---")

	var wg sync.WaitGroup

	for _, proceso := range s.Procesos {
		wg.Add(1)
		go s.ejecutarProcesoSimulado(proceso, &wg)
	}

	wg.Wait()
	fmt.Println("--- SIMULACIÓN COMPLETADA ---")
}

func (s *SimuladorProcesos) ejecutarProcesoSimulado(p *ProcesoSimulado, wg *sync.WaitGroup) {

	defer wg.Done()

	p.mutex.Lock()
	p.Estado = "Listo"
	p.mutex.Unlock()

	fmt.Printf("Proceso %d (%s) está LISTO\n", p.ID, p.Nombre)
	time.Sleep(500 * time.Millisecond)

	p.mutex.Lock()
	p.Estado = "Ejecutando"
	p.TiempoInicio = time.Now()
	p.mutex.Unlock()

	fmt.Printf("Proceso %d (%s) está EJECUTANDO (duración: %ds)\n", p.ID, p.Nombre, p.Duracion)

	for i := 0; i < p.Duracion; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("Proceso %d (%s) ejecutándose... %d/%d\n",
			p.ID, p.Nombre, i+1, p.Duracion)

		if i == p.Duracion/2 {
			p.mutex.Lock()
			p.Estado = "Bloqueado"
			p.mutex.Unlock()

			fmt.Printf("Proceso %d (%s) se ha BLOQUEADO\n", p.ID, p.Nombre)
			time.Sleep(1 * time.Second)

			p.mutex.Lock()
			p.Estado = "Listo"
			p.mutex.Unlock()

			fmt.Printf("Proceso %d (%s) vuelve a ESTAR LISTO\n", p.ID, p.Nombre)
			time.Sleep(500 * time.Millisecond)

			p.mutex.Lock()
			p.Estado = "Ejecutando"
			p.mutex.Unlock()
		}
	}

	p.mutex.Lock()
	p.Estado = "Terminado"
	p.TiempoFin = time.Now()
	duracionReal := p.TiempoFin.Sub(p.TiempoInicio)
	p.mutex.Unlock()

	fmt.Printf("Proceso %d (%s) TERMINADO (duración real: %v)\n", p.ID, p.Nombre, duracionReal)
}

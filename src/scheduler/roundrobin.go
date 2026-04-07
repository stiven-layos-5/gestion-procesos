package scheduler

import (
	"fmt"
	"gestion-procesos/src/logger"
	"gestion-procesos/src/memory"
	proceso "gestion-procesos/src/procesos"
	"sync"
	"time"
)

type ProcessTask struct {
	Proc      *proceso.Proceso
	Ejecutar  chan int
	Respuesta chan bool
}

type Scheduler struct {
	quantum        int
	reloj          int
	mutex          sync.Mutex
	colaListos     []*ProcessTask
	colaNuevos     chan *ProcessTask
	colaBloqueados []*ProcessTask
	memManager     *memory.MemoryManager
	log            *logger.Logger

	pendMutex          sync.Mutex
	procesosPendientes int
	done               chan struct{}
	estadisticas       []*proceso.Proceso
}

func NewScheduler(quantum, totalPaginas int, log *logger.Logger) *Scheduler {
	return &Scheduler{
		quantum:            quantum,
		reloj:              0,
		colaListos:         []*ProcessTask{},
		colaNuevos:         make(chan *ProcessTask, 100),
		colaBloqueados:     []*ProcessTask{},
		memManager:         memory.NewMemoryManager(totalPaginas),
		log:                log,
		procesosPendientes: 0,
		estadisticas:       []*proceso.Proceso{},
	}
}

func (s *Scheduler) AddProcess(proc *proceso.Proceso) {
	s.pendMutex.Lock()
	s.procesosPendientes++
	s.pendMutex.Unlock()

	task := &ProcessTask{
		Proc:      proc,
		Ejecutar:  make(chan int),
		Respuesta: make(chan bool),
	}
	s.estadisticas = append(s.estadisticas, proc)

	// Lanzar la goroutine del proceso
	go s.runProcess(task)

	// Simular tiempo de llegada (escala: 1 unidad = 100 ms para visibilidad)
	go func(llegada int) {
		time.Sleep(time.Duration(llegada) * 100 * time.Millisecond)
		s.colaNuevos <- task
	}(proc.Llegada)
}

func (s *Scheduler) runProcess(task *ProcessTask) {
	proc := task.Proc
	for {
		quantum := <-task.Ejecutar
		if quantum == -1 { // señal de terminación forzada (no se usa)
			break
		}
		// Ejecutar durante el quantum (o menos si la ráfaga restante es menor)
		ejecucion := quantum
		if proc.Restante < quantum {
			ejecucion = proc.Restante
		}
		proc.Restante -= ejecucion
		terminado := proc.Restante == 0
		task.Respuesta <- terminado
		if terminado {
			break
		}
	}
	// Notificar al planificador que este proceso terminó
	s.pendMutex.Lock()
	s.procesosPendientes--
	if s.procesosPendientes == 0 && s.done != nil {
		close(s.done)
	}
	s.pendMutex.Unlock()
}

func (s *Scheduler) Run() {
	s.done = make(chan struct{})
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			s.log.Log("=== FIN DE LA SIMULACIÓN ===")
			return
		default:
			// 1. Manejar llegada de nuevos procesos
			select {
			case newTask := <-s.colaNuevos:
				paginasReq := (newTask.Proc.TamMemoria + 4095) / 4096 // páginas de 4 KB
				paginas := s.memManager.AsignarPaginas(paginasReq)
				if paginas == nil {
					s.log.Log(fmt.Sprintf("T=%d: Proceso %d SIN MEMORIA (necesita %d páginas) -> swap out",
						s.reloj, newTask.Proc.PID, paginasReq))
					s.colaBloqueados = append(s.colaBloqueados, newTask)
				} else {
					newTask.Proc.Paginas = paginas
					s.log.Log(fmt.Sprintf("T=%d: Proceso %d obtiene %d páginas (memoria asignada)",
						s.reloj, newTask.Proc.PID, len(paginas)))
					s.colaListos = append(s.colaListos, newTask)
				}
			default:
			}

			// 2. Si no hay procesos listos, esperar un poco
			if len(s.colaListos) == 0 {
				<-ticker.C
				continue
			}

			// 3. Round Robin: tomar el primero de la cola
			task := s.colaListos[0]
			s.colaListos = s.colaListos[1:]

			// 4. Enviar quantum al proceso y esperar respuesta
			task.Ejecutar <- s.quantum
			terminado := <-task.Respuesta

			// 5. Avanzar el reloj lógico
			s.mutex.Lock()
			s.reloj += s.quantum
			tiempoActual := s.reloj
			s.mutex.Unlock()

			if terminado {
				// Proceso terminado: calcular estadísticas y liberar memoria
				task.Proc.Retorno = tiempoActual
				task.Proc.Espera = tiempoActual - task.Proc.Llegada - task.Proc.Rafaga
				s.log.Log(fmt.Sprintf("T=%d: Proceso %d TERMINADO (ráfaga=%d, espera=%d, retorno=%d)",
					tiempoActual, task.Proc.PID, task.Proc.Rafaga, task.Proc.Espera, task.Proc.Retorno))
				s.memManager.LiberarPaginas(task.Proc.Paginas)
				// Revisar procesos bloqueados que puedan ahora obtener memoria
				s.revisarBloqueados()
			} else {
				// El proceso continúa: se reencola al final
				s.colaListos = append(s.colaListos, task)
				s.log.Log(fmt.Sprintf("T=%d: Proceso %d ejecutó %d unidades (restante=%d)",
					tiempoActual, task.Proc.PID, s.quantum, task.Proc.Restante))
			}
		}
	}
}

func (s *Scheduler) revisarBloqueados() {
	nuevaCola := []*ProcessTask{}
	for _, task := range s.colaBloqueados {
		paginasReq := (task.Proc.TamMemoria + 4095) / 4096
		paginas := s.memManager.AsignarPaginas(paginasReq)
		if paginas != nil {
			task.Proc.Paginas = paginas
			s.log.Log(fmt.Sprintf("T=%d: Proceso %d DESBLOQUEADO (memoria disponible) -> pasa a listo",
				s.reloj, task.Proc.PID))
			s.colaListos = append(s.colaListos, task)
		} else {
			nuevaCola = append(nuevaCola, task)
		}
	}
	s.colaBloqueados = nuevaCola
}

func (s *Scheduler) MostrarEstadisticas() {
	fmt.Println("\n=== ESTADÍSTICAS FINALES ===")
	fmt.Println("PID\tLlegada\tRáfaga\tEspera\tRetorno")
	fmt.Println("----------------------------------------")
	sumaEspera := 0
	sumaRetorno := 0
	for _, p := range s.estadisticas {
		fmt.Printf("%d\t%d\t%d\t%d\t%d\n", p.PID, p.Llegada, p.Rafaga, p.Espera, p.Retorno)
		sumaEspera += p.Espera
		sumaRetorno += p.Retorno
	}
	n := len(s.estadisticas)
	fmt.Printf("\nTiempo promedio de espera: %.2f\n", float64(sumaEspera)/float64(n))
	fmt.Printf("Tiempo promedio de retorno: %.2f\n", float64(sumaRetorno)/float64(n))
}

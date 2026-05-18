package concurrency

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Tarea struct {
	ID         int
	Nombre     string
	DuracionMs int
}

type Resultado struct {
	TareaID     int
	WorkerID    int
	GoroutineID string
	Duracion    time.Duration
	Ok          bool
}

type WorkerPool struct {
	numWorkers int
	tareas     chan Tarea
	resultados chan Resultado
	wg         sync.WaitGroup
}

func NuevoWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		tareas:     make(chan Tarea, 20),
		resultados: make(chan Resultado, 20),
	}
}

func goroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := string(buf[:n])
	start := len("goroutine ")
	end := start
	for end < len(stack) && stack[end] != ' ' {
		end++
	}
	if end > start {
		return "goroutine-" + stack[start:end]
	}
	return "goroutine-?"
}

func (wp *WorkerPool) Iniciar() {
	fmt.Printf("  Iniciando pool de %d workers concurrentes...\n", wp.numWorkers)
	for i := 1; i <= wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	go func() {
		wp.wg.Wait()
		close(wp.resultados)
	}()
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	gid := goroutineID()

	for tarea := range wp.tareas {
		inicio := time.Now()
		fmt.Printf("    [Worker %d | %s] ← procesando Tarea %d (%s)\n",
			id, gid, tarea.ID, tarea.Nombre)

		// Simular trabajo real con tiempo variable
		time.Sleep(time.Duration(tarea.DuracionMs) * time.Millisecond)

		duracion := time.Since(inicio)
		fmt.Printf("    [Worker %d | %s] ✓ Tarea %d completada en %v\n",
			id, gid, tarea.ID, duracion.Round(time.Millisecond))

		wp.resultados <- Resultado{
			TareaID:     tarea.ID,
			WorkerID:    id,
			GoroutineID: gid,
			Duracion:    duracion,
			Ok:          true,
		}
	}
}

func (wp *WorkerPool) EnviarTareas(tareas []Tarea) {
	go func() {
		for _, t := range tareas {
			wp.tareas <- t
		}
		close(wp.tareas)
	}()
}

func (wp *WorkerPool) RecolectarResultados() []Resultado {
	var resultados []Resultado
	for r := range wp.resultados {
		resultados = append(resultados, r)
	}
	return resultados
}

func GenerarTareas(n int) []Tarea {
	rand.Seed(time.Now().UnixNano())
	tareas := make([]Tarea, n)
	nombres := []string{"CompilarModulo", "ProcesarImagen", "ConsultarDB",
		"EnviarEmail", "GenerarReporte", "ValidarToken", "SincronizarCache",
		"AnalizarLog", "RespaldarDatos", "ActualizarIndice"}
	for i := 0; i < n; i++ {
		tareas[i] = Tarea{
			ID:         i + 1,
			Nombre:     nombres[i%len(nombres)],
			DuracionMs: 50 + rand.Intn(200), // 50–250ms
		}
	}
	return tareas
}

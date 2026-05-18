package sync

import (
	"fmt"
	gosync "sync"
	"time"
)

type ContadorInseguro struct {
	valor int
}

func (c *ContadorInseguro) Incrementar() {
	v := c.valor
	time.Sleep(time.Microsecond)
	c.valor = v + 1
}

func (c *ContadorInseguro) Valor() int {
	return c.valor
}

func DemostrarCondicionCarrera(numGoroutines, incrementosPorGoroutine int) {
	contador := &ContadorInseguro{}
	var wg gosync.WaitGroup
	esperado := numGoroutines * incrementosPorGoroutine

	fmt.Printf("  Lanzando %d goroutines, cada una incrementa %d veces.\n",
		numGoroutines, incrementosPorGoroutine)
	fmt.Printf("  Resultado esperado: %d\n", esperado)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementosPorGoroutine; j++ {
				contador.Incrementar()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("  Resultado REAL (sin mutex): %d  ← ¡INCORRECTO por race condition!\n\n",
		contador.Valor())
}

type ContadorSeguro struct {
	mu    gosync.Mutex
	valor int
}

func (c *ContadorSeguro) Incrementar() {
	
	c.mu.Lock()   
	defer c.mu.Unlock() 

	v := c.valor
	time.Sleep(time.Microsecond)
	c.valor = v + 1
}

func (c *ContadorSeguro) Valor() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.valor
}

func DemostrarMutex(numGoroutines, incrementosPorGoroutine int) {
	contador := &ContadorSeguro{}
	var wg gosync.WaitGroup
	esperado := numGoroutines * incrementosPorGoroutine

	fmt.Printf("  Lanzando %d goroutines con sync.Mutex.\n", numGoroutines)
	fmt.Printf("  Resultado esperado: %d\n", esperado)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementosPorGoroutine; j++ {
				contador.Incrementar()
			}
		}()
	}
	wg.Wait()

	resultado := contador.Valor()
	estado := "✓ CORRECTO"
	if resultado != esperado {
		estado = "✗ INCORRECTO"
	}
	fmt.Printf("  Resultado CON mutex: %d  ← %s\n\n", resultado, estado)
}

type Semaforo struct {
	ch chan struct{}
}

func NuevoSemaforo(n int) *Semaforo {
	return &Semaforo{ch: make(chan struct{}, n)}
}

func (s *Semaforo) Adquirir() {
	s.ch <- struct{}{}
}

func (s *Semaforo) Liberar() {
	<-s.ch
}

func DemostrarSemaforo(capacidad, numGoroutines int) {
	sem := NuevoSemaforo(capacidad)
	var wg gosync.WaitGroup
	var mu gosync.Mutex
	dentroActual := 0
	maxDentro := 0

	fmt.Printf("  Semáforo con capacidad %d, %d goroutines compitiendo.\n",
		capacidad, numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Adquirir()
			defer sem.Liberar()

			mu.Lock()
			dentroActual++
			if dentroActual > maxDentro {
				maxDentro = dentroActual
			}
			fmt.Printf("    Goroutine %2d ENTRÓ  | activos: %d\n", id, dentroActual)
			mu.Unlock()

			time.Sleep(50 * time.Millisecond) // trabajo en sección crítica

			mu.Lock()
			dentroActual--
			fmt.Printf("    Goroutine %2d SALIÓ  | activos: %d\n", id, dentroActual)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("  Máximo simultáneo observado: %d (límite: %d)\n\n", maxDentro, capacidad)
}

package monitor

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gestion-procesos/src/models"
)

type Monitor struct {
	Intervalo time.Duration
	Historial []models.RecursosSnapshot
}

func NuevoMonitor(intervalo time.Duration) *Monitor {
	return &Monitor{Intervalo: intervalo}
}

func (m *Monitor) TomarSnapshot() models.RecursosSnapshot {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	snap := models.RecursosSnapshot{
		Timestamp:     time.Now(),
		CPUPorcentaje: leerCPULinux(),
		MemTotal:      mem.Sys / 1024,
		MemUsada:      mem.Alloc / 1024,
		NumGoroutines: runtime.NumGoroutine(),
	}
	m.Historial = append(m.Historial, snap)
	return snap
}

func leerCPULinux() float64 {
	leer := func() (idle, total uint64) {
		f, err := os.Open("/proc/stat")
		if err != nil {
			return 0, 0
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			linea := scanner.Text()
			if !strings.HasPrefix(linea, "cpu ") {
				continue
			}
			campos := strings.Fields(linea)[1:] // saltar "cpu"
			var vals []uint64
			for _, c := range campos {
				v, _ := strconv.ParseUint(c, 10, 64)
				vals = append(vals, v)
			}
			if len(vals) < 4 {
				return 0, 0
			}
			// user, nice, system, idle, iowait, irq, softirq, steal
			idle = vals[3]
			for _, v := range vals {
				total += v
			}
			return
		}
		return 0, 0
	}

	idle1, total1 := leer()
	time.Sleep(200 * time.Millisecond)
	idle2, total2 := leer()

	deltaTotal := total2 - total1
	deltaIdle := idle2 - idle1

	if deltaTotal == 0 {
		return 0
	}
	return (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100.0
}

func (m *Monitor) Monitorear(n int) {
	fmt.Printf("  %-26s %-10s %-14s %-14s %-12s\n",
		"Timestamp", "CPU%", "MemTotal(KB)", "MemUsada(KB)", "Goroutines")
	fmt.Println("  " + strings.Repeat("-", 78))

	for i := 0; i < n; i++ {
		snap := m.TomarSnapshot()
		fmt.Printf("  %-26s %-10.2f %-14d %-14d %-12d\n",
			snap.Timestamp.Format("15:04:05.000"),
			snap.CPUPorcentaje,
			snap.MemTotal,
			snap.MemUsada,
			snap.NumGoroutines,
		)
		if i < n-1 {
			time.Sleep(m.Intervalo)
		}
	}
}

func (m *Monitor) ResumenHistorial() {
	if len(m.Historial) == 0 {
		fmt.Println("  Sin datos en el historial.")
		return
	}
	var sumCPU float64
	var maxMem uint64
	for _, s := range m.Historial {
		sumCPU += s.CPUPorcentaje
		if s.MemUsada > maxMem {
			maxMem = s.MemUsada
		}
	}
	fmt.Printf("  CPU promedio: %.2f%% | Memoria pico: %d KB | Muestras: %d\n",
		sumCPU/float64(len(m.Historial)), maxMem, len(m.Historial))
}

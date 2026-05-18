package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"gestion-procesos/src/models"
)

type Manager struct {
	procesosIniciados []*exec.Cmd
}

func NuevoManager() *Manager {
	return &Manager{}
}

func (m *Manager) Listar() ([]models.ProcesoReal, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist", "/fo", "csv", "/nh")
	} else {
		cmd = exec.Command("ps", "aux")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar ps: %w", err)
	}

	var procesos []models.ProcesoReal
	lineas := strings.Split(string(out), "\n")

	if len(lineas) > 1 {
		lineas = lineas[1:]
	}

	for _, linea := range lineas {
		campos := strings.Fields(linea)
		if len(campos) < 11 {
			continue
		}
		pid, err := strconv.Atoi(campos[1])
		if err != nil {
			continue
		}
		nombre := campos[10]
		if len(nombre) > 25 {
			nombre = nombre[:25]
		}
		procesos = append(procesos, models.ProcesoReal{
			PID:    pid,
			Nombre: nombre,
			Estado: campos[7],
		})
	}
	return procesos, nil
}

func (m *Manager) Iniciar(comando string, args ...string) (int, error) {
	cmd := exec.Command(comando, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("error al iniciar proceso '%s': %w", comando, err)
	}

	m.procesosIniciados = append(m.procesosIniciados, cmd)
	fmt.Printf("  ✓ Proceso iniciado: %s (PID: %d)\n", comando, cmd.Process.Pid)
	return cmd.Process.Pid, nil
}

func (m *Manager) Finalizar(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("proceso PID %d no encontrado: %w", pid, err)
	}
	if err := proc.Kill(); err != nil {
		return fmt.Errorf("error al finalizar PID %d: %w", pid, err)
	}
	fmt.Printf("  ✓ Proceso PID %d finalizado.\n", pid)
	return nil
}

func ImprimirProcesos(procesos []models.ProcesoReal) {
	fmt.Printf("\n%-8s %-28s %-10s\n", "PID", "Nombre", "Estado")
	fmt.Println(strings.Repeat("-", 50))
	limite := 15
	if len(procesos) < limite {
		limite = len(procesos)
	}
	for _, p := range procesos[:limite] {
		fmt.Printf("%-8d %-28s %-10s\n", p.PID, p.Nombre, p.Estado)
	}
	fmt.Printf("... (%d procesos en total)\n", len(procesos))
}

package memory

import "sync"

type MemoryManager struct {
	totalPaginas  int
	paginasLibres []bool
	mutex         sync.Mutex
}

func NewMemoryManager(totalPaginas int) *MemoryManager {
	return &MemoryManager{
		totalPaginas:  totalPaginas,
		paginasLibres: make([]bool, totalPaginas),
	}
}

// Asignar páginas a un proceso (modelo de asignación simple: primeras libres)
// Retorna slice de índices o nil si no hay suficientes páginas libres
func (mm *MemoryManager) AsignarPaginas(numPaginas int) []int {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	paginas := []int{}
	for i := 0; i < mm.totalPaginas && len(paginas) < numPaginas; i++ {
		if !mm.paginasLibres[i] {
			paginas = append(paginas, i)
			mm.paginasLibres[i] = true
		}
	}

	if len(paginas) == numPaginas {
		return paginas
	}

	// No hay suficientes, liberar las que se asignaron parcialmente
	for _, p := range paginas {
		mm.paginasLibres[p] = false
	}
	return nil
}

// Liberar páginas previamente asignadas
func (mm *MemoryManager) LiberarPaginas(paginas []int) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	for _, p := range paginas {
		mm.paginasLibres[p] = false
	}
}

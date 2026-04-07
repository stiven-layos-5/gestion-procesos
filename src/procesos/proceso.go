package proceso

type Proceso struct {
	PID        int
	Llegada    int
	Rafaga     int
	Restante   int
	Espera     int
	Retorno    int
	TamMemoria int
	Paginas    []int
}

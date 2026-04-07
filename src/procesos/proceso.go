package proceso

type Proceso struct {
	PID           int
	Llegada       int
	Rafaga        int
	Restante      int
	Espera        int
	Retorno       int
	Completado    bool
	PrimeraEjec   bool
	InicioPrimero int
}

func NuevoProceso(id int, arrival int, burst int) Proceso {

	return Proceso{
		PID:     id,
		Llegada: arrival,
		Rafaga:  burst,
	}
}

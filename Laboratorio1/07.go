package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ej_07() {
	// Slice compartido para almacenar los eventos del log
	var log []string

	// Mutex para proteger el acceso concurrente al log
	var mutex sync.Mutex

	// Crear un WaitGroup para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup

	// Número de nodos (goroutines)
	const numNodos = 10

	// Inicializar el generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Lista de eventos críticos posibles
	eventosCriticos := []string{
		"temperatura alta",
		"pérdida de conexión",
		"uso elevado de CPU",
		"falla en el disco",
		"error de memoria",
	}

	// Crear las goroutines
	for i := 1; i <= numNodos; i++ {
		wg.Add(1) // Incrementar el contador del WaitGroup
		go func(id int) {
			defer wg.Done() // Indicar que esta goroutine ha terminado

			for j := 0; j < 5; j++ { // Cada nodo registra 5 eventos
				// Seleccionar un evento crítico aleatorio
				evento := eventosCriticos[rand.Intn(len(eventosCriticos))]

				// Formatear el mensaje del evento
				mensaje := fmt.Sprintf("nodo-%d: %s", id, evento)

				// Proteger el acceso al log con el mutex
				mutex.Lock()
				log = append(log, mensaje)
				mutex.Unlock()

				// Esperar 0.5 segundos antes de registrar el siguiente evento
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Mostrar el contenido del log
	fmt.Println("Log de eventos:")
	for _, evento := range log {
		fmt.Println(evento)
	}
}

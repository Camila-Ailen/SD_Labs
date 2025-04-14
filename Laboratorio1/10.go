package main

import (
	"fmt"
	"sync"
)

// Variable global
var x int

// Mutex para proteger el acceso a x
var mutex sync.Mutex

// Función para incrementar x
func incrementar() {
	mutex.Lock()         // Bloquear el acceso a x
	defer mutex.Unlock() // Liberar el bloqueo
	x += 5               // Incrementar x
}

func ej_10() {
	// Crear un WaitGroup para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup

	// Lanzar 100 goroutines
	for i := 0; i < 100; i++ {
		wg.Add(1) // Incrementar el contador del WaitGroup
		go func() {
			defer wg.Done() // Decrementar el contador al finalizar
			incrementar()
		}()
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	// Imprimir el valor final de x
	fmt.Printf("El valor final de x es: %v\n", x)
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ej_08() {
	// Slice compartido para almacenar los nodos con menor latencia
	var resultados []string

	// Mutex para proteger el acceso concurrente al slice
	var mutex sync.Mutex

	// Lista de nodos
	nodos := []string{"nodo-1", "nodo-2", "nodo-3"}

	// Inicializar el generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Crear un WaitGroup para sincronizar la goroutine
	var wg sync.WaitGroup
	wg.Add(1) // Incrementar el contador del WaitGroup

	// Función para simular el ping a los nodos
	ping := func(nodo string) int {
		latencia := rand.Intn(401) + 100 // Latencia aleatoria entre 100 y 500 ms
		time.Sleep(time.Duration(latencia) * time.Millisecond)
		return latencia
	}

	// Goroutine para realizar el monitoreo
	go func() {
		defer wg.Done() // Indicar que la goroutine ha terminado

		for ronda := 1; ronda <= 10; ronda++ {
			fmt.Printf("Ronda %d:\n", ronda)

			// Variables para rastrear el nodo con menor latencia
			menorLatencia := 501 // Un valor mayor al máximo posible (500 ms)
			nodoMenorLatencia := ""

			// Realizar ping a cada nodo
			for _, nodo := range nodos {
				latencia := ping(nodo)
				fmt.Printf("%s respondió en %d ms\n", nodo, latencia)

				// Actualizar el nodo con menor latencia
				if latencia < menorLatencia {
					menorLatencia = latencia
					nodoMenorLatencia = nodo
				}
			}

			// Proteger el acceso al slice con el mutex
			mutex.Lock()
			resultados = append(resultados, nodoMenorLatencia)
			mutex.Unlock()

			// Esperar 2 segundos antes de la siguiente ronda
			time.Sleep(2 * time.Second)
		}
	}()

	// Esperar a que la goroutine termine
	wg.Wait()

	// Mostrar los resultados
	fmt.Println("\nResultados:")
	for i, nodo := range resultados {
		fmt.Printf("Ronda %d: %s\n", i+1, nodo)
	}
}

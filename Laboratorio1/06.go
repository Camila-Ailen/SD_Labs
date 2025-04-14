package main

import (
	"fmt"
	"time"
)

func ej_06() {
	// Crear canales para conectar los nodos
	const numNodos = 5
	canales := make([]chan string, numNodos)
	for i := 0; i < numNodos; i++ {
		canales[i] = make(chan string)
	}

	// Crear los nodos como goroutines
	for i := 0; i < numNodos; i++ {
		entrada := canales[i]
		salida := canales[(i+1)%numNodos] // El último nodo se conecta al primero
		go nodo(i+1, entrada, salida)
	}

	// Iniciar el sistema enviando el primer mensaje al nodo 1
	go func() {
		canales[0] <- "Inicio del sistema"
	}()

	// Ejecutar el sistema durante 1 minuto
	time.Sleep(1 * time.Minute)

	// Cerrar todos los canales para finalizar las goroutines
	for _, canal := range canales {
		close(canal)
	}

	fmt.Println("El sistema en anillo ha terminado.")
}

func nodo(id int, entrada <-chan string, salida chan<- string) {
	for mensaje := range entrada {
		// Recibir el mensaje del nodo anterior
		fmt.Printf("Nodo %v recibió: %v\n", id, mensaje)

		// Esperar 1 segundo antes de enviar el mensaje al siguiente nodo
		time.Sleep(1 * time.Second)

		// Enviar el mensaje al siguiente nodo
		salida <- fmt.Sprintf("Heartbeat desde nodo %v", id)
	}
}

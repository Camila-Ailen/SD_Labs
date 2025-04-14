package main

import (
	"fmt"
	"time"
)

func ej_09() {
	// Crear canales para los suscriptores
	const numSuscriptores = 3
	suscriptores := make([]chan string, numSuscriptores)
	for i := 0; i < numSuscriptores; i++ {
		suscriptores[i] = make(chan string)
	}

	// Canal de control para detener las goroutines
	done := make(chan struct{})

	// Función para el publicador
	publicador := func() {
		eventoID := 1
		for {
			select {
			case <-done:
				// Detener el publicador si recibe una señal en el canal `done`
				fmt.Println("Publicador detenido.")
				return
			default:
				// Enviar un mensaje a todos los suscriptores
				mensaje := fmt.Sprintf("evento-%d", eventoID)
				fmt.Printf("Publicador envió: %s\n", mensaje)
				for _, sub := range suscriptores {
					sub <- mensaje
				}
				eventoID++
				time.Sleep(1 * time.Second) // Esperar 1 segundo antes de enviar el siguiente mensaje
			}
		}
	}

	// Función para los suscriptores
	suscriptor := func(id int, canal <-chan string) {
		for {
			select {
			case mensaje, ok := <-canal:
				if !ok {
					// Detener el suscriptor si el canal está cerrado
					fmt.Printf("Suscriptor %d detenido.\n", id)
					return
				}
				fmt.Printf("Suscriptor %d recibió: %s\n", id, mensaje)
			case <-done:
				// Detener el suscriptor si recibe una señal en el canal `done`
				fmt.Printf("Suscriptor %d detenido.\n", id)
				return
			}
		}
	}

	// Iniciar las goroutines para los suscriptores
	for i := 0; i < numSuscriptores; i++ {
		go suscriptor(i+1, suscriptores[i])
	}

	// Iniciar la goroutine para el publicador
	go publicador()

	// Ejecutar el sistema durante 10 segundos
	time.Sleep(10 * time.Second)

	// Notificar a las goroutines que deben detenerse
	close(done)

	// Cerrar los canales para finalizar las goroutines de los suscriptores
	for _, sub := range suscriptores {
		close(sub)
	}

	fmt.Println("El ejercicio 9 ha terminado.")
}

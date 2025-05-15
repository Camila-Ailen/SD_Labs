package main

import (
	"context"
	"grpc-pg-1/proto"
	"log"
	"time"
	"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var nombres = []string{"Juan", "Maria", "Pedro", "Ana", "Luis", "Claudio", "Sofia", "Diego", "Valeria", "Javier"}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()
	c := proto.NewServicioClient(conn)
	
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(nombre string) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.Hola(ctx, &proto.Requerimiento{Nombre: nombre})
			if err != nil {
				log.Fatalf("Error al llamar al servidro: %v", err)
			}
			log.Printf("Respuesta: %s", r.Personas[0])
		}(nombres[i%len(nombres)])
	}
	wg.Wait()
	log.Println("Todas las goroutines han terminado.")

}

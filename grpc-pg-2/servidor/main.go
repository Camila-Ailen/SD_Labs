package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"grpc-pg-2/proto"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedMonitorServer
	mu          sync.Mutex
	ultimaVista map[string]time.Time
}

func (s *servidor) EnviarHeartbeat(stream proto.Monitor_EnviarHeartbeatServer) error {
	var nodoId string
	for {
		hb, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Ack{Mensaje: "Stream cerrado"})
		}

		if err != nil {
			// SE QUITA LOS LOGS DE STREAM
			// log.Printf("Error en stream: %v", err)
			return err
		}
		nodoId = hb.NodoId
		s.mu.Lock()
		s.ultimaVista[nodoId] = time.Unix(hb.MarcaTiempo, 0)
		s.mu.Unlock()
		// SE QUITA LOS LOGS DE LOS HEARTBEATS
		// log.Printf("[HEARTBEAT] %v %v", nodoId, hb.MarcaTiempo)

	}
}

func (s *servidor) detectorFallas(intervalo time.Duration) {
	for {
		time.Sleep(intervalo)
		s.mu.Lock()
		ahora := time.Now()

		var conectados []string

		for nodo, ultimo := range s.ultimaVista {
			if ahora.Sub(ultimo) > 3*intervalo {
				log.Printf("Fallo en Nodo %v inactivo desde	hace %.0fs", nodo, ahora.Sub(ultimo).Seconds())
				
				for i, v := range conectados {
					if v == nodo {
						conectados = append(conectados[:i], conectados[i+1:]...)
						break
					}
				}
			}else{
				conectados = append(conectados, nodo)
			}
		}

		if len(conectados) > 0 {
			log.Printf("Nodos conectados: %v", conectados)
		} else {
			log.Println("No hay nodos conectados actualmente")
		}

		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	s := grpc.NewServer()
	servidor := &servidor{ultimaVista: make(map[string]time.Time)}
	proto.RegisterMonitorServer(s, servidor)

	go servidor.detectorFallas(5 * time.Second)
	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}

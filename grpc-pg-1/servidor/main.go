package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type servidor struct {
	proto.UnimplementedServicioServer
}

func (s *servidor) Hola(ctx context.Context, req *proto.Requerimiento) (*proto.ListadoPersonas, error) {
	log.Printf("Recibido: %s", req.Nombre)
	// return &proto.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
	personas := []*proto.Persona{
        {Nombre: req.Nombre}, // Usar el nombre recibido en la solicitud
    }

	return &proto.ListadoPersonas{Personas: personas}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterServicioServer(s, &servidor{})
	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}

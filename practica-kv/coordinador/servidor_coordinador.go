package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync/atomic"

	pb "github.com/Camila-Ailen/SD_Labs/practica-kv/proto" // Importa el paquete generado por protoc
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServidorCoordinador implementa pb.CoordinadorServer.
type ServidorCoordinador struct {
	pb.UnimplementedCoordinadorServer
	listaReplicas []string
	indiceRR      uint64 // índice para Round Robin
}

// NewServidorCoordinador crea un Coordinador con direcciones de réplica.
func NewServidorCoordinador(replicas []string) *ServidorCoordinador {
	return &ServidorCoordinador{
		listaReplicas: replicas,
		indiceRR:      0,
	}
}

// elegirReplicaParaEscritura: round-robin simple (ignora la clave).
func (c *ServidorCoordinador) elegirReplicaParaEscritura(clave string) string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// elegirReplicaParaLectura: también round-robin.
func (c *ServidorCoordinador) elegirReplicaParaLectura() string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// Obtener: redirige petición de lectura a una réplica.
func (c *ServidorCoordinador) Obtener(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
	replicaAddr := c.elegirReplicaParaLectura()
	log.Printf("Coordinador: Redirigiendo Obtener clave '%s' a réplica %s", req.Clave, replicaAddr)

	conn, err := grpc.NewClient(replicaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Coordinador: Error al conectar con réplica %s: %v", replicaAddr, err)
		return nil, err
	}
	defer conn.Close()

	clienteReplica := pb.NewReplicaClient(conn)
	return clienteReplica.ObtenerLocal(ctx, req)
}

// Guardar: redirige petición de escritura a una réplica elegida.
func (c *ServidorCoordinador) Guardar(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	replicaAddr := c.elegirReplicaParaEscritura(req.Clave)
	log.Printf("Coordinador: Redirigiendo Guardar clave '%s' a réplica %s", req.Clave, replicaAddr)

	conn, err := grpc.NewClient(replicaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Coordinador: Error al conectar con réplica %s: %v", replicaAddr, err)
		return nil, err
	}
	defer conn.Close()

	clienteReplica := pb.NewReplicaClient(conn)
	return clienteReplica.GuardarLocal(ctx, req)
}

// Eliminar: redirige petición de eliminación a una réplica elegida.
func (c *ServidorCoordinador) Eliminar(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
	replicaAddr := c.elegirReplicaParaEscritura(req.Clave) // Usa la misma lógica que para escritura
	log.Printf("Coordinador: Redirigiendo Eliminar clave '%s' a réplica %s", req.Clave, replicaAddr)

	conn, err := grpc.NewClient(replicaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Coordinador: Error al conectar con réplica %s: %v", replicaAddr, err)
		return nil, err
	}
	defer conn.Close()

	clienteReplica := pb.NewReplicaClient(conn)
	return clienteReplica.EliminarLocal(ctx, req)
}

func main() {
	// Definir bandera para la dirección de escucha del Coordinador.
	listenAddr := flag.String("listen", ":6000", "dirección para que escuche el Coordinador (p.ej., :6000)")
	flag.Parse()

	replicaAddrs := flag.Args()
	if len(replicaAddrs) < 3 {
		log.Fatalf("Debe proveer al menos 3 direcciones de réplicas, p.ej.: go run servidor_coordinador.go -listen :6000 :50051 :50052 :50053")
	}

	log.Printf("Coordinador escuchando en %s", *listenAddr)
	log.Printf("Réplicas configuradas: %v", replicaAddrs)

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Error al escuchar en %s: %v", *listenAddr, err)
	}

	grpcServer := grpc.NewServer()
	servidorCoordinador := NewServidorCoordinador(replicaAddrs)
	pb.RegisterCoordinadorServer(grpcServer, servidorCoordinador)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar servidor gRPC del Coordinador: %v", err)
	}
}

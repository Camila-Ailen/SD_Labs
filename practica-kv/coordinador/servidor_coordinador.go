package main

import (
	"context"
	"flag"
	"log"
	"sync/atomic"
)

// ServidorCoordinador implementa pb.CoordinadorServer.
type ServidorCoordinador struct {
	pb.UnimplementedCoordinadorServer
	listaReplicas []string
	indiceRR uint64 // índice para Round Robin
}

// NewServidorCoordinador crea un Coordinador con direcciones de réplica.
func NewServidorCoordinador(replicas []string) *ServidorCoordinador {
	return &ServidorCoordinador{
		listaReplicas: replicas,
		indiceRR:     0,
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
}
// Guardar: redirige petición de escritura a una réplica elegida.
func (c *ServidorCoordinador) Guardar(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
}
// Eliminar: redirige petición de eliminación a una réplica elegida.
func (c *ServidorCoordinador) Eliminar(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
}

func main() {
 // Definir bandera para la dirección de escucha del Coordinador.
	listen := flag.String("listen", ":6000", "dirección para que escuche el Coordinador (p.ej., :6000)")
	flag.Parse()
	replicas := flag.Args()
	if len(replicas) < 3 {
	log.Fatalf("Debe proveer al menos 3 direcciones de réplicas, p.ej.: go run servidor_coordinador.go -listen :6000 :50051 :50052 :50053")
	}
}
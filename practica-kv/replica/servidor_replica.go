package main

import (
	"context"
	"encoding/binary"
	"sync"
)

// ValorConVersion guarda el valor y su reloj vectorial asociado.
type ValorConVersion struct {
	Valor       []byte
	RelojVector VectorReloj
}

// ServidorReplica implementa pb.ReplicaServer
type ServidorReplica struct {
	pb.UnimplementedReplicaServer
	mu           sync.Mutex
	almacen      map[string]ValorConVersion // map[clave]ValorConVersion
	relojVector  VectorReloj
	idReplica    int                // 0, 1 o 2
	clientesPeer []pb.ReplicaClient // stubs gRPC a las otras réplicas
}

// VectorReloj representa un reloj vectorial de longitud 3 (tres réplicas).
type VectorReloj [3]uint64

// Incrementar aumenta en 1 el componente correspondiente a la réplica que llama.
func (vr *VectorReloj) Incrementar(idReplica int) {

}

// Fusionar toma el máximo elemento a elemento entre dos vectores.
func (vr *VectorReloj) Fusionar(otro VectorReloj) {
}

// AntesDe devuelve true si vr < otro en el sentido estricto (strictlyless).
func (vr VectorReloj) AntesDe(otro VectorReloj) bool {
	menor := false
	return menor
}

// encodeVector serializa el VectorReloj a []byte para enviarlo por gRPC.
func encodeVector(vr VectorReloj) []byte {
	buf := make([]byte, 8*3)
	for i := 0; i < 3; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], vr[i])
	}
	return buf
}

// decodeVector convierte []byte a VectorReloj.
func decodeVector(b []byte) VectorReloj {
	var vr VectorReloj
	for i := 0; i < 3; i++ {
		vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return vr
}

// NewServidorReplica crea una instancia de ServidorReplica
// idReplica: 0, 1 o 2
// peerAddrs: direcciones gRPC de los otros dos peers (ej.:
// []string{":50052", ":50053"})
func NewServidorReplica(idReplica int, peerAddrs []string) *ServidorReplica{

}


// GuardarLocal recibe la petición del Coordinador para almacenar clave/valor.
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	r.mu.Lock()
 // 1. Incrementar nuestro componente del reloj vectorial
 // 2. Guardar en el mapa local
 // 3. Construir mutación para replicar a peers
 // 4. Replicar asíncronamente a cada peer
 // 5. Responder al Coordinador con el nuevo reloj vectorial
}


// EliminarLocal recibe la petición del Coordinador para borrar una clave.
func (r *ServidorReplica) EliminarLocal(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
	r.mu.Lock()
 // 1. Incrementar nuestro componente del reloj vectorial
 // 2. Borrar del mapa local (si existe)
 // 3. Construir mutación de eliminación
 // 4. Replicar a peers
 // 5. Responder al Coordinador
}


// ObtenerLocal retorna el valor y reloj vectorial de una clave en esta réplica.
func (r *ServidorReplica) ObtenerLocal(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {

}


// ReplicarMutacion recibe una mutación de otra réplica y la aplica localmente.
func (r *ServidorReplica) ReplicarMutacion(ctx context.Context, m *pb.Mutacion) (*pb.Reconocimiento, error) {
	r.mu.Lock()
 // 1. Decodificar el reloj vectorial de la mutación
// Considere que si no existía, o la mutación es “más nueva”, sobrescribir
// Si existe y nuestra versión local está “por delante”, ignoramos (conflicto resuelto a favor local).
 // 2. Fusionar nuestro reloj vectorial con el remoto
 // 3. Responder con nuestro reloj actualizado
}



func main() {
 // Uso: go run servidor_replica.go <idReplica> <direccionEscucha> <peer1> <peer2>
 // Ejemplo: go run servidor_replica.go 0 :50051 :50052 :50053
	if len(os.Args) != 5 {
		log.Fatalf("Uso: %s <idReplica> <direccionEscucha> <peer1> <peer2>",os.Args[0])
	}
 // 1. Inicializar servidor gRPC
 // 2. Crear instancia de ServidorReplica
 // 3. Iniciar servidor
}
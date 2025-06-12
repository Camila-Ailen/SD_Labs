package main

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	pb "github.com/Camila-Ailen/SD_Labs/practica-kv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	vr[idReplica]++
}

// Fusionar toma el máximo elemento a elemento entre dos vectores.
func (vr *VectorReloj) Fusionar(otro VectorReloj) {
	for i := 0; i < 3; i++ {
		if otro[i] > vr[i] {
			vr[i] = otro[i]
		}
	}
}

// AntesDe devuelve true si vr < otro en el sentido estricto (strictlyless).
// vr es estrictamente anterior a otro si vr[i] <= otro[i] para todo i,
// y existe al menos un j tal que vr[j] < otro[j].
func (vr VectorReloj) AntesDe(otro VectorReloj) bool {
	alMenosUnoMenor := false
	for i := 0; i < 3; i++ {
		if vr[i] > otro[i] {
			return false // Si algún componente de vr es mayor, no es AntesDe
		}
		if vr[i] < otro[i] {
			alMenosUnoMenor = true
		}
	}
	return alMenosUnoMenor // Devuelve true si todos son <= y al menos uno es <
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
func NewServidorReplica(idReplica int, peerAddrs []string) *ServidorReplica {
	s := &ServidorReplica{
		almacen:      make(map[string]ValorConVersion),
		relojVector:  VectorReloj{}, // Se inicializará a [0,0,0]
		idReplica:    idReplica,
		clientesPeer: make([]pb.ReplicaClient, 2), // Asumiendo 2 peers
	}

	// Conectar a los otros peers
	peerIndex := 0
	for i := 0; i < 3; i++ {
		if i == idReplica {
			continue // No conectar a sí mismo
		}
		if peerIndex < len(peerAddrs) {
			conn, err := grpc.NewClient(peerAddrs[peerIndex], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("Error al conectar con peer %s: %v", peerAddrs[peerIndex], err)
				// Podríamos decidir si fallar o continuar sin este peer
				s.clientesPeer[peerIndex] = nil // Marcar como no disponible
			} else {
				s.clientesPeer[peerIndex] = pb.NewReplicaClient(conn)
			}
			peerIndex++
		} else {
			log.Printf("Advertencia: No se proporcionó la dirección para el peer esperado.")
		}
	}
	return s
}

// GuardarLocal recibe la petición del Coordinador para almacenar clave/valor.
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Incrementar nuestro componente del reloj vectorial
	r.relojVector.Incrementar(r.idReplica)

	// 2. Guardar en el mapa local
	valorActual, existe := r.almacen[req.Clave]
	nuevoValorConVersion := ValorConVersion{
		Valor:       req.Valor,
		RelojVector: r.relojVector, // Usar el reloj actualizado de la réplica
	}

	// Lógica de resolución de conflictos: si ya existe y la versión actual no es anterior a la nueva
	if existe && !valorActual.RelojVector.AntesDe(nuevoValorConVersion.RelojVector) && !nuevoValorConVersion.RelojVector.AntesDe(valorActual.RelojVector) {
		// Conflicto concurrente, podríamos tener una política (ej. timestamp, id de réplica)
		// Por ahora, simplemente sobrescribimos si no es estrictamente anterior.
		// O podríamos decidir no sobrescribir si la entrante no es estrictamente posterior.
		// Para este ejemplo, si no hay orden causal claro, la nueva escritura gana si no es "más vieja".
		// Si la versión actual es posterior o concurrente, y la nueva no es estrictamente posterior, no hacer nada o registrar.
		// Aquí, para simplificar, si no es anterior, se considera que la nueva puede aplicar o es concurrente.
	}

	r.almacen[req.Clave] = nuevoValorConVersion
	log.Printf("Replica %d: Guardado local clave '%s', valor '%s', reloj %v", r.idReplica, req.Clave, req.Valor, r.relojVector)

	// 3. Construir mutación para replicar a peers
	mutacion := &pb.Mutacion{
		Clave:       req.Clave,
		Valor:       req.Valor,
		RelojVector: encodeVector(r.relojVector),
		Tipo:        pb.Mutacion_GUARDAR, 
	}

	// 4. Replicar asíncronamente a cada peer
	for i, cliente := range r.clientesPeer {
		if cliente == nil {
			log.Printf("Replica %d: Peer %d no disponible para replicar Guardar", r.idReplica, i)
			continue
		}
		go func(c pb.ReplicaClient, peerIdx int) {
			_, err := c.ReplicarMutacion(context.Background(), mutacion)
			if err != nil {
				log.Printf("Replica %d: Error al replicar Guardar a peer %d: %v", r.idReplica, peerIdx, err)
			} else {
				log.Printf("Replica %d: Guardar replicado a peer %d para clave '%s'", r.idReplica, peerIdx, req.Clave)
			}
		}(cliente, i)
	}

	// 5. Responder al Coordinador con el nuevo reloj vectorial
	return &pb.RespuestaGuardar{Exito: true, NuevoRelojVector: encodeVector(r.relojVector)}, nil
}

// EliminarLocal recibe la petición del Coordinador para borrar una clave.
func (r *ServidorReplica) EliminarLocal(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Incrementar nuestro componente del reloj vectorial
	r.relojVector.Incrementar(r.idReplica)
	log.Printf("Replica %d: Eliminando clave '%s', reloj antes de eliminar %v", r.idReplica, req.Clave, r.almacen[req.Clave].RelojVector)

	// 2. Borrar del mapa local (si existe)
	_, existe := r.almacen[req.Clave]
	if existe {
		delete(r.almacen, req.Clave)
		log.Printf("Replica %d: Clave '%s' eliminada localmente. Reloj actualizado: %v", r.idReplica, req.Clave, r.relojVector)
	} else {
		log.Printf("Replica %d: Clave '%s' no encontrada para eliminar. Reloj actualizado: %v", r.idReplica, req.Clave, r.relojVector)
		// Aunque no exista, se propaga la eliminación para asegurar consistencia (marcar como borrado con reloj actual)
	}

	// 3. Construir mutación de eliminación
	mutacion := &pb.Mutacion{
		Clave:       req.Clave,
		RelojVector: encodeVector(r.relojVector),
		Tipo:        pb.Mutacion_ELIMINAR, // Corregido de Operacion a Tipo
	}

	// 4. Replicar a peers
	for i, cliente := range r.clientesPeer {
		if cliente == nil {
			log.Printf("Replica %d: Peer %d no disponible para replicar Eliminar", r.idReplica, i)
			continue
		}
		go func(c pb.ReplicaClient, peerIdx int) {
			_, err := c.ReplicarMutacion(context.Background(), mutacion)
			if err != nil {
				log.Printf("Replica %d: Error al replicar Eliminar a peer %d: %v", r.idReplica, peerIdx, err)
			} else {
				log.Printf("Replica %d: Eliminar replicado a peer %d para clave '%s'", r.idReplica, peerIdx, req.Clave)
			}
		}(cliente, i)
	}

	// 5. Responder al Coordinador
	return &pb.RespuestaEliminar{Exito: true, NuevoRelojVector: encodeVector(r.relojVector)}, nil
}

// ObtenerLocal retorna el valor y reloj vectorial de una clave en esta réplica.
func (r *ServidorReplica) ObtenerLocal(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	vc, existe := r.almacen[req.Clave]
	if !existe {
		// Si no existe, devolvemos un valor vacío y un reloj vectorial nulo o el actual de la réplica.
		// Devolver el reloj actual de la réplica puede ser más informativo para el coordinador.
		log.Printf("Replica %d: Clave '%s' no encontrada localmente. Reloj de réplica: %v", r.idReplica, req.Clave, r.relojVector)
		return &pb.RespuestaObtener{Valor: nil, RelojVector: encodeVector(r.relojVector), Existe: false}, nil
	}

	log.Printf("Replica %d: Obteniendo clave '%s', valor '%s', reloj %v", r.idReplica, req.Clave, vc.Valor, vc.RelojVector)
	return &pb.RespuestaObtener{Valor: vc.Valor, RelojVector: encodeVector(vc.RelojVector), Existe: true}, nil
}

// ReplicarMutacion recibe una mutación de otra réplica y la aplica localmente.
func (r *ServidorReplica) ReplicarMutacion(ctx context.Context, m *pb.Mutacion) (*pb.Reconocimiento, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	relojMutacion := decodeVector(m.RelojVector)
	log.Printf("Replica %d: Recibida mutación para clave '%s', op: %s, reloj mutación: %v. Mi reloj: %v", r.idReplica, m.Clave, m.Tipo, relojMutacion, r.relojVector) // Corregido de m.Operacion a m.Tipo

	valorLocal, existeLocal := r.almacen[m.Clave]

	// 1. Decodificar el reloj vectorial de la mutación
	// Considere que si no existía, o la mutación es “más nueva”, sobrescribir
	// Si existe y nuestra versión local está “por delante”, ignoramos (conflicto resuelto a favor local).

	aplicarMutacion := false
	if !existeLocal {
		aplicarMutacion = true // No existe localmente, aplicar siempre
		log.Printf("Replica %d: Clave '%s' no existe localmente, aplicando mutación.", r.idReplica, m.Clave)
	} else {
		// La clave existe localmente, comparar relojes
		if relojMutacion.AntesDe(valorLocal.RelojVector) {
			// La versión local es estrictamente posterior, ignorar mutación
			log.Printf("Replica %d: Mutación para clave '%s' es más antigua (reloj mut: %v, reloj local: %v). Ignorando.", r.idReplica, m.Clave, relojMutacion, valorLocal.RelojVector)
			aplicarMutacion = false
		} else if valorLocal.RelojVector.AntesDe(relojMutacion) {
			// La mutación es estrictamente posterior, aplicar
			aplicarMutacion = true
			log.Printf("Replica %d: Mutación para clave '%s' es más nueva (reloj mut: %v, reloj local: %v). Aplicando.", r.idReplica, m.Clave, relojMutacion, valorLocal.RelojVector)
		} else {
			// Relojes concurrentes. Aquí se podría aplicar una política de resolución de conflictos.
			// Por ejemplo, si es GUARDAR vs GUARDAR, podría ganar el ID de réplica mayor/menor, o el timestamp si se tuviera.
			// Para ELIMINAR, generalmente la eliminación "gana" sobre un GUARDAR concurrente si se quiere consistencia eventual hacia la eliminación.
			// Para este ejemplo, si son concurrentes, la mutación entrante gana para asegurar la propagación.
			log.Printf("Replica %d: Mutación para clave '%s' es concurrente (reloj mut: %v, reloj local: %v). Aplicando por defecto.", r.idReplica, m.Clave, relojMutacion, valorLocal.RelojVector)
			aplicarMutacion = true
		}
	}

	if aplicarMutacion {
		if m.Tipo == pb.Mutacion_GUARDAR { // Corregido de m.Operacion a m.Tipo
			r.almacen[m.Clave] = ValorConVersion{
				Valor:       m.Valor,
				RelojVector: relojMutacion,
			}
			log.Printf("Replica %d: Mutación GUARDAR aplicada para clave '%s'. Nuevo reloj local para clave: %v", r.idReplica, m.Clave, relojMutacion)
		} else if m.Tipo == pb.Mutacion_ELIMINAR { // Corregido de m.Operacion a m.Tipo
			// Si la clave existe, se elimina. Si no, se considera "eliminada" con el reloj de la mutación.
			// Esto es para manejar el caso de que una réplica reciba una eliminación para una clave que nunca vio o que ya eliminó con un reloj anterior.
			// Almacenar una "tumba" con el reloj de la eliminación podría ser una opción, pero aquí simplemente la borramos.
			// El reloj de la réplica se fusionará de todas formas.
			delete(r.almacen, m.Clave)
			log.Printf("Replica %d: Mutación ELIMINAR aplicada para clave '%s'.", r.idReplica, m.Clave)
			// Podríamos guardar una "tumba" con relojMutacion si quisiéramos evitar que una escritura posterior "más vieja" la resucite.
			// r.almacen[m.Clave] = ValorConVersion{Valor: nil, RelojVector: relojMutacion, EsTumba: true} // Ejemplo conceptual
		}
	}

	// 2. Fusionar nuestro reloj vectorial con el remoto (independientemente de si se aplicó la mutación o no, para propagar el conocimiento causal)
	r.relojVector.Fusionar(relojMutacion)
	log.Printf("Replica %d: Reloj fusionado con el de la mutación. Mi reloj ahora: %v", r.idReplica, r.relojVector)

	// 3. Responder con nuestro reloj actualizado
	return &pb.Reconocimiento{Ok: true, RelojVectorAck: encodeVector(r.relojVector)}, nil
}

func main() {
	// Uso: go run servidor_replica.go <idReplica> <direccionEscucha> <peer1> <peer2>
	// Ejemplo: go run servidor_replica.go 0 :50051 :50052 :50053
	if len(os.Args) != 5 {
		log.Fatalf("Uso: %s <idReplica> <direccionEscucha> <peer1> <peer2>", os.Args[0])
	}

	idReplicaStr := os.Args[1]
	direccionEscucha := os.Args[2]
	peer1Addr := os.Args[3]
	peer2Addr := os.Args[4]

	idReplica, err := strconv.Atoi(idReplicaStr)
	if err != nil {
		log.Fatalf("ID de réplica inválido: %s", idReplicaStr)
	}
	if idReplica < 0 || idReplica > 2 {
		log.Fatalf("ID de réplica debe ser 0, 1 o 2. Recibido: %d", idReplica)
	}

	peerAddrs := []string{}
	// El orden de peerAddrs en NewServidorReplica debe ser consistente
	// para que el índice del clientePeer corresponda al ID de la réplica esperada.
	// Aquí asumimos que peer1 y peer2 son las otras dos réplicas.
	// Si idReplica es 0, peers son 1 (:50052) y 2 (:50053)
	// Si idReplica es 1, peers son 0 (:50051) y 2 (:50053)
	// Si idReplica es 2, peers son 0 (:50051) y 1 (:50052)
	// La implementación de NewServidorReplica debe manejar esto correctamente o
	// los argumentos deben pasarse de forma que se sepa a qué ID de réplica corresponde cada dirección.

	// Para simplificar, NewServidorReplica espera las direcciones de los *otros* peers.
	// El main debe pasar las direcciones correctas.
	// Si tenemos 3 réplicas en total (0, 1, 2) y sus direcciones son addr0, addr1, addr2.
	// Si esta es la réplica 0, peerAddrs para NewServidorReplica será [addr1, addr2]
	// Si esta es la réplica 1, peerAddrs para NewServidorReplica será [addr0, addr2]
	// Si esta es la réplica 2, peerAddrs para NewServidorReplica será [addr0, addr1]

	// Los argumentos peer1 y peer2 son las direcciones de las *otras* réplicas.
	// No necesariamente en orden de ID.
	// NewServidorReplica los tomará en el orden que se le pasen.
	// Es importante que el coordinador sepa qué dirección corresponde a qué ID de réplica.
	// Y que cada réplica sepa su propio ID.
	peerAddrs = append(peerAddrs, peer1Addr)
	peerAddrs = append(peerAddrs, peer2Addr)

	// 1. Inicializar servidor gRPC
	lis, err := net.Listen("tcp", direccionEscucha)
	if err != nil {
		log.Fatalf("Error al escuchar en %s: %v", direccionEscucha, err)
	}
	log.Printf("Replica %d escuchando en %s", idReplica, direccionEscucha)

	grpcServer := grpc.NewServer()

	// 2. Crear instancia de ServidorReplica
	servidorReplica := NewServidorReplica(idReplica, peerAddrs)
	pb.RegisterReplicaServer(grpcServer, servidorReplica)

	log.Printf("Replica %d: Conectándose a peers: %v", idReplica, peerAddrs)

	// 3. Iniciar servidor
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar servidor gRPC: %v", err)
	}
}

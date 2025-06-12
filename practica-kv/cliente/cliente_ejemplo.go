package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/Camila-Ailen/SD_Labs/practica-kv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Dirección del coordinador desde los argumentos o usa un valor por defecto.
	coordinadorAddr := ":6000" // Valor por defecto
	if len(os.Args) > 1 {
		coordinadorAddr = os.Args[1]
	}

	log.Printf("Conectando al coordinador en %s", coordinadorAddr)

	// Establecer conexión con el coordinador
	conn, err := grpc.NewClient(coordinadorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al coordinador: %v", err)
	}
	defer conn.Close()

	clienteCoordinador := pb.NewCoordinadorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clave := "usuario123"
	valor := []byte("datosImportantes")

	// 1. Guardar la clave "usuario123" con valor "datosImportantes".
	log.Printf("Intentando Guardar: Clave='%s', Valor='%s'", clave, string(valor))
	respGuardar, err := clienteCoordinador.Guardar(ctx, &pb.SolicitudGuardar{
		Clave: clave,
		Valor: valor,
	})
	if err != nil {
		log.Fatalf("Error al Guardar: %v", err)
	}
	if !respGuardar.Exito {
		log.Fatalf("Fallo al Guardar la clave '%s'. Reloj: %x", clave, respGuardar.NuevoRelojVector)
	}
	log.Printf("Guardar exitoso para clave '%s'. Nuevo Reloj Vectorial: %x", clave, respGuardar.NuevoRelojVector)
	relojActual := respGuardar.NuevoRelojVector // Guardamos el reloj para la siguiente operación

	// 2. Obtener la clave "usuario123" e imprime valor + reloj vectorial.
	log.Printf("Intentando Obtener: Clave='%s'", clave)
	respObtener, err := clienteCoordinador.Obtener(ctx, &pb.SolicitudObtener{Clave: clave})
	if err != nil {
		log.Fatalf("Error al Obtener: %v", err)
	}
	if !respObtener.Existe {
		log.Printf("Obtener: La clave '%s' no existe. Reloj: %x", clave, respObtener.RelojVector)
	} else {
		log.Printf("Obtener exitoso: Clave='%s', Valor='%s', Reloj Vectorial: %x", clave, string(respObtener.Valor), respObtener.RelojVector)
		relojActual = respObtener.RelojVector // Actualizamos el reloj con el obtenido
	}

	// 3. Eliminar la misma clave, enviando el reloj vectorial que recibió.
	log.Printf("Intentando Eliminar: Clave='%s', Reloj Vectorial: %x", clave, relojActual)
	respEliminar, err := clienteCoordinador.Eliminar(ctx, &pb.SolicitudEliminar{
		Clave:       clave,
		RelojVector: relojActual, // Usamos el reloj de la operación anterior (Guardar u Obtener)
	})
	if err != nil {
		log.Fatalf("Error al Eliminar: %v", err)
	}
	if !respEliminar.Exito {
		log.Fatalf("Fallo al Eliminar la clave '%s'. Reloj: %x", clave, respEliminar.NuevoRelojVector)
	}
	log.Printf("Eliminar exitoso para clave '%s'. Nuevo Reloj Vectorial: %x", clave, respEliminar.NuevoRelojVector)

	// 4. Obtener nuevamente para verificar que ya no existe.
	log.Printf("Intentando Obtener (después de eliminar): Clave='%s'", clave)
	respObtenerDespuesEliminar, err := clienteCoordinador.Obtener(ctx, &pb.SolicitudObtener{Clave: clave})
	if err != nil {
		log.Fatalf("Error al Obtener (después de eliminar): %v", err)
	}
	if respObtenerDespuesEliminar.Existe {
		log.Printf("ERROR: La clave '%s' todavía existe después de eliminar. Valor: '%s', Reloj: %x", clave, string(respObtenerDespuesEliminar.Valor), respObtenerDespuesEliminar.RelojVector)
	} else {
		log.Printf("Verificación exitosa: La clave '%s' no existe después de eliminar. Reloj devuelto: %x", clave, respObtenerDespuesEliminar.RelojVector)
	}

	log.Println("Cliente de ejemplo finalizado.")
}

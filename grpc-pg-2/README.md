# Practica Guiada 2
## Grupo 2
### Integrantes 
- Gomez Camila
- Henker Franco
- Mendez Camila

### Instrucción para compilar el archivo .proto 
Dentro de la raíz del proyecto ejecutar:

`protoc --go_out=. --go-grpc_out=. proto/servicio.proto`

### Para ejecutar el servidor y el cliente:
- Para correr el servidor ejecutar en una terminal `go run servidor/main.go`
- Para correr el cliente ejecutar en otra terminal `go run cliente/main.go {id}` (id es un argumento que representa el nombre del nodo que se conecta al servidor)

### Frecuencia de los heartbeats

- ### Aumentar la frecuencia de los heartbeats:
  - Pro: Implica que los nodos se monitorean más frecuentemente, lo que permite detectar más rápido si un nodo ha dejado de responder o está caído. Esto reduce el tiempo de detección de fallos y puede acelerar la conmutación por error o la recuperación. 
  - Contra: Incrementa el tráfico de red y el consumo de recursos, lo que puede afectar la eficiencia del sistema si la frecuencia es demasiado alta.

- ### Disminuir la frecuencia de los heartbeats:
  - Pro: Reduce el tráfico de monitoreo y el uso de recursos.
  - Contra: Aumenta el tiempo que tarda el sistema en detectar que un nodo está inactivo o ha fallado. Esto puede provocar demoras en la recuperación o en la toma de acciones correctivas, afectando la disponibilidad y la respuesta del sistema ante fallos.
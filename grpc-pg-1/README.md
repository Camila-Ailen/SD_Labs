# Practica Guiada 1
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
- Para correr el cliente ejecutar en otra terminal `go run cliente/main.go`

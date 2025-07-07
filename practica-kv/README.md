# Practica Guiada 3
## Grupo 2
### Integrantes 
- Gomez Camila
- Henker Franco
- Mendez Camila
---

## Instrucción para compilar el archivo .proto 
Dentro de la raíz del proyecto ejecutar:

```bash
protoc --go_out=. --go-grpc_out=. proto/kv.proto
```

## Para ejecutar el coordinador, servidor y el cliente:
- Para correr el coordinador ejecutar en una terminal 
  - El coordinador escucha por el puerto 6000, los puertos de escucha definidos debajo seran los mismos de los servidores
```bash
go run coordinador/servidor_coordinador.go :<port1> :<port2> :<port3>
```

- Para correr los servidores ejecutar en diferentes terminales 
  - Indicar id del servidor, el puerto por el que escucha y luego los puertos de los demas servidores
```bash
go run servidor/servidor_replica.go <id> :<self_port> :<port> :<port>
```

- Para correr el cliente ejecutar en otra terminal 
  - El cliente enviara peticiones para escribir, leer y borrar un registro determinado
```bash
go run cliente/cliente_ejemplo.go
``` 


## Compilar archivos
Se puede compilar cada archivo .go desde su directorio ejecutando la instruccion
```bash
go build <archivo>.go
```


## Ejemplos
### Levantar servidores y coordinador
- Coordinador
```bash
go run coordinador/servidor_coordinador.go :6001 :6002 :6003
```
- Servidores
```bash
go run replica/servidor_replica.go 1 :6001 :6002 :6003
```
```bash
go run replica/servidor_replica.go 2 :6002 :6001 :6003
```
```bash
go run replica/servidor_replica.go 0 :6003 :6002 :6001
```

```bash
go run cliente/cliente_ejemplo.go
```

### Salida esperada 
- Coordinador
```bash
2025/07/04 13:25:15 Coordinador escuchando en :6000
2025/07/04 13:25:15 Réplicas configuradas: [:6001 :6002 :6003]
2025/07/04 13:25:30 Coordinador: Redirigiendo Guardar clave 'usuario123' a réplica :6002
2025/07/04 13:25:30 Coordinador: Redirigiendo Obtener clave 'usuario123' a réplica :6003
2025/07/04 13:25:30 Coordinador: Redirigiendo Eliminar clave 'usuario123' a réplica :6001
2025/07/04 13:25:30 Coordinador: Redirigiendo Obtener clave 'usuario123' a réplica :6002
```

- Servidores
```bash
2025/07/04 13:25:20 Replica 1 escuchando en :6001
2025/07/04 13:25:20 Replica 1: Conectándose a peers: [:6002 :6003]
2025/07/04 13:25:30 Replica 1: Recibida mutación para clave 'usuario123', op: GUARDAR, reloj mutación: [0 0 1]. Mi reloj: [0 0 0]
2025/07/04 13:25:30 Replica 1: Clave 'usuario123' no existe localmente, aplicando mutación.
2025/07/04 13:25:30 Replica 1: Mutación GUARDAR aplicada para clave 'usuario123'. Nuevo reloj local para clave: [0 0 1]
2025/07/04 13:25:30 Replica 1: Reloj fusionado con el de la mutación. Mi reloj ahora: [0 0 1]
2025/07/04 13:25:30 Replica 1: Eliminando clave 'usuario123', reloj antes de eliminar [0 0 1]
2025/07/04 13:25:30 Replica 1: Clave 'usuario123' eliminada localmente. Reloj actualizado: [0 1 1]
2025/07/04 13:25:30 Replica 1: Eliminar replicado a peer 0 para clave 'usuario123'
2025/07/04 13:25:30 Replica 1: Eliminar replicado a peer 1 para clave 'usuario123'



2025/07/04 13:25:21 Replica 2 escuchando en :6002
2025/07/04 13:25:21 Replica 2: Conectándose a peers: [:6003 :6001]
2025/07/04 13:25:30 Replica 2: Guardado local clave 'usuario123', valor 'datosImportantes', reloj [0 0 1]
2025/07/04 13:25:30 Replica 2: Guardar replicado a peer 1 para clave 'usuario123'
2025/07/04 13:25:30 Replica 2: Guardar replicado a peer 0 para clave 'usuario123'
2025/07/04 13:25:30 Replica 2: Recibida mutación para clave 'usuario123', op: ELIMINAR, reloj mutación: [0 1 1]. Mi reloj: [0 0 1]
2025/07/04 13:25:30 Replica 2: Mutación para clave 'usuario123' es más nueva (reloj mut: [0 1 1], reloj local: [0 0 1]). Aplicando.
2025/07/04 13:25:30 Replica 2: Mutación ELIMINAR aplicada para clave 'usuario123'.
2025/07/04 13:25:30 Replica 2: Reloj fusionado con el de la mutación. Mi reloj ahora: [0 1 1]
2025/07/04 13:25:30 Replica 2: Clave 'usuario123' no encontrada localmente. Reloj de réplica: [0 1 1]



2025/07/04 13:25:22 Replica 0 escuchando en :6003
2025/07/04 13:25:22 Replica 0: Conectándose a peers: [:6002 :6001]
2025/07/04 13:25:30 Replica 0: Clave 'usuario123' no encontrada localmente. Reloj de réplica: [0 0 0]
2025/07/04 13:25:30 Replica 0: Recibida mutación para clave 'usuario123', op: GUARDAR, reloj mutación: [0 0 1]. Mi reloj: [0 0 0]
2025/07/04 13:25:30 Replica 0: Clave 'usuario123' no existe localmente, aplicando mutación.
2025/07/04 13:25:30 Replica 0: Mutación GUARDAR aplicada para clave 'usuario123'. Nuevo reloj local para clave: [0 0 1]
2025/07/04 13:25:30 Replica 0: Reloj fusionado con el de la mutación. Mi reloj ahora: [0 0 1]
2025/07/04 13:25:30 Replica 0: Recibida mutación para clave 'usuario123', op: ELIMINAR, reloj mutación: [0 1 1]. Mi reloj: [0 0 1]
2025/07/04 13:25:30 Replica 0: Mutación para clave 'usuario123' es más nueva (reloj mut: [0 1 1], reloj local: [0 0 1]). Aplicando.
2025/07/04 13:25:30 Replica 0: Mutación ELIMINAR aplicada para clave 'usuario123'.
2025/07/04 13:25:30 Replica 0: Reloj fusionado con el de la mutación. Mi reloj ahora: [0 1 1]
```


- Cliente
```bash
2025/07/04 13:25:30 Conectando al coordinador en :6000
2025/07/04 13:25:30 Intentando Guardar: Clave='usuario123', Valor='datosImportantes'
2025/07/04 13:25:30 Guardar exitoso para clave 'usuario123'. Nuevo Reloj Vectorial: 000000000000000000000000000000000000000000000001
2025/07/04 13:25:30 Intentando Obtener: Clave='usuario123'
2025/07/04 13:25:30 Obtener: La clave 'usuario123' no existe. Reloj: 000000000000000000000000000000000000000000000000
2025/07/04 13:25:30 Intentando Eliminar: Clave='usuario123', Reloj Vectorial: 000000000000000000000000000000000000000000000001
2025/07/04 13:25:30 Eliminar exitoso para clave 'usuario123'. Nuevo Reloj Vectorial: 000000000000000000000000000000010000000000000001
2025/07/04 13:25:30 Intentando Obtener (después de eliminar): Clave='usuario123'
2025/07/04 13:25:30 Verificación exitosa: La clave 'usuario123' no existe después de eliminar. Reloj devuelto: 000000000000000000000000000000010000000000000001
2025/07/04 13:25:30 Cliente de ejemplo finalizado.
```


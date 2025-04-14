package main

import (
	"bufio"
	"fmt"
	"os"
)

// Función para mostrar el contenido de un archivo
func mostrarArchivo(nombreArchivo string) {
	// Intentar abrir el archivo
	archivo, err := os.Open(nombreArchivo)
	if err != nil {
		// Si ocurre un error (el archivo no existe), mostrar un mensaje y terminar
		fmt.Printf("No se pudo abrir el archivo: %v\n", err)
		return
	}
	defer archivo.Close() // Asegurarse de cerrar el archivo al finalizar

	// Crear un lector para leer el contenido del archivo línea por línea
	scanner := bufio.NewScanner(archivo)
	fmt.Println("Contenido del archivo:")
	for scanner.Scan() {
		// Imprimir cada línea del archivo
		fmt.Println(scanner.Text())
	}

	// Verificar si ocurrió algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error al leer el archivo: %v\n", err)
	}
}

// Para ejecutar desde el programa principal:
func ej_05() {

	// Para ejecutar desde la línea de comandos (hacer los cambios necesarios en main para llamar a esta función):
	// go run 05.go ejemplo.txt
	// func main() {

	// Verificar si se pasó un argumento desde la línea de comandos
	if len(os.Args) > 1 {
		// Si hay un argumento, usarlo como nombre del archivo
		nombreArchivo := os.Args[1]
		mostrarArchivo(nombreArchivo)
	} else {
		// Si no hay argumentos, llamar a la función con un nombre de archivo predeterminado
		fmt.Println("No se proporcionó un archivo como argumento.")
		fmt.Println("Llamando a la función con un archivo predeterminado...")
		mostrarArchivo("archivo.txt")
	}
}

package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("")
		fmt.Println("Seleccione un ejercicio para ejecutar:")
		fmt.Println("1. Ejercicio 1 (Sumar números pares)")
		fmt.Println("0. Salir")
		fmt.Print("Ingrese su opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			ej_01()
		case 0:
			fmt.Println("Saliendo del programa...")
			return
		default:
			fmt.Println("Opción no válida, intente nuevamente.")
		}
	}
}

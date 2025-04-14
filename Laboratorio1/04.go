package main

import (
	"fmt"
)

func ej_04() {
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Convertir de °C a °F")
		fmt.Println("2. Convertir de °F a °C")
		fmt.Println("0. Salir")
		fmt.Print("Ingrese su opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			convertirCaF()
		case 2:
			convertirFaC()
		case 0:
			fmt.Println("Saliendo del programa...")
			return
		default:
			fmt.Println("Opción no válida, intente nuevamente.")
		}
	}
}

func convertirCaF() {
	fmt.Print("Ingrese la temperatura en °C: ")
	var celsius float64
	fmt.Scanln(&celsius)
	fahrenheit := (celsius * 9 / 5) + 32
	fmt.Printf("%.4v °C equivalen a %.4v °F\n", celsius, fahrenheit)
}

func convertirFaC() {
	fmt.Print("Ingrese la temperatura en °F: ")
	var fahrenheit float64
	fmt.Scanln(&fahrenheit)
	celsius := (fahrenheit - 32) * 5 / 9
	fmt.Printf("%.4v °F equivalen a %.4v °C\n", fahrenheit, celsius)
}

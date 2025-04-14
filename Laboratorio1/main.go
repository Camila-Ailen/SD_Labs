package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("")
		fmt.Println("Seleccione un ejercicio para ejecutar:")
		fmt.Println("1. Ejercicio 1 (Sumar números pares)")
		fmt.Println("2. Ejercicio 2 (Contar palabras en una frase)")
		fmt.Println("3. Ejercicio 3 (Calcular promedio de alumnos)")
		fmt.Println("4. Ejercicio 4 (Convertir temperaturas)")
		fmt.Println("5. Ejercicio 5 (Mostrar contenido de un archivo)")
		fmt.Println("6. Ejercicio 6 (Sistema en anillo)")
		fmt.Println("7. Ejercicio 7 (Log de eventos)")
		fmt.Println("8. Ejercicio 8 (Monitoreo de gorutina)")
		fmt.Println("9. Ejercicio 9 (Publicador y suscriptores)")
		fmt.Println("0. Salir")
		fmt.Print("Ingrese su opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			ej_01()
		case 2:
			ej_02()
		case 3:
			ej_03()
		case 4:
			ej_04()
		case 5:
			ej_05()
		case 6:
			ej_06()
		case 7:
			ej_07()
		case 8:
			ej_08()
		case 9:
			ej_09()
		case 0:
			fmt.Println("Saliendo del programa...")
			return
		default:
			fmt.Println("Opción no válida, intente nuevamente.")
		}
	}
}

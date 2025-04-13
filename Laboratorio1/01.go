package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ej_01() {
	// Crear un slice vacío para almacenar los números
	var numeros []int

	fmt.Println("Ingrese números enteros separados por espacio (presione Enter para finalizar):")

	// Usar bufio.Scanner para leer toda la línea de entrada
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()

		// Dividir la entrada en partes y convertirlas a enteros
		for _, val := range strings.Fields(input) {
			num, err := strconv.Atoi(val)
			if err == nil {
				numeros = append(numeros, num)
			}
		}
	}

	// Llamar a la función SumarPares y mostrar el resultado
	resultado := SumarPares(numeros)
	fmt.Printf("La suma de los números pares es: %v\n", resultado)
}

func SumarPares(slice []int) int {
	suma := 0
	for _, num := range slice {
		if num%2 == 0 {
			suma += num
		}
	}
	return suma
}

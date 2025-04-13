package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ej_02() {
	// Solicitar al usuario que ingrese una frase
	fmt.Println("Ingrese una frase:")
	var frase string
	// Leer la frase desde la entrada estándar
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		frase = scanner.Text()
	}

	// Contar palabras en la frase
	cantidad := ContarPalabras(frase)
	// Mostrar resultado
	fmt.Println("La cantidad de palabras es:", cantidad)
}

func ContarPalabras(frase string) int {
	// Dividir la frase en palabras y contar
	palabras := strings.Fields(frase)
	return len(palabras)
}

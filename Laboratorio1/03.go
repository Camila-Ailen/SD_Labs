package main

import (
	"fmt"
)

// Definimos la estructura Alumno
type Alumno struct {
	Nombre string
	Notas  []float64
}

// Método Promedio para calcular el promedio de las notas del alumno
func (a Alumno) Promedio() float64 {
	suma := 0.0
	for _, nota := range a.Notas {
		suma += nota
	}
	return suma / float64(len(a.Notas))
}

func ej_03() {
	// Crear una lista de alumnos con sus notas
	alumnos := []Alumno{
		{"Camila", []float64{8.5, 9.0, 7.5}},
		{"Franco", []float64{6.0, 7.5, 8.0}},
		{"Ailen", []float64{9.0, 8.5, 9.5}},
	}

	// Calcular y mostrar el promedio de cada alumno
	for _, alumno := range alumnos {
		fmt.Printf("El promedio de %v es: %.3v\n", alumno.Nombre, alumno.Promedio())
	}
}

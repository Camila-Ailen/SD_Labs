package main

import (
	"fmt"
	"sync"
)

// Mutex global
var mx sync.Mutex
var mx2 sync.Mutex

// Función b()
func b() {
	mx2.Lock()         // Bloquear el mutex
	defer mx2.Unlock() // Desbloquear el mutex al finalizar
	fmt.Println("Hola mundo")
}

// Función a()
func a() {
	mx.Lock()         // Bloquear el mutex
	defer mx.Unlock() // Desbloquear el mutex al finalizar
	b()               // Invocar a la función b()
}

func ej_11() {
	// Invocar a la función a() desde la gorutina principal
	a()
}

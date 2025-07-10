package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var shardNodes = map[int][]string{
	0: {"http://127.0.0.1:11000", "http://127.0.0.1:11001", "http://127.0.0.1:11002"},
	1: {"http://127.0.0.1:11100", "http://127.0.0.1:11101", "http://127.0.0.1:11102"},
}

func getLeaderURL(shard int) (string, error) {
	for _, nodeURL := range shardNodes[shard] {
		resp, err := http.Get(nodeURL + "/status")
		if err != nil {
			continue // siguiente nodo
		}
		defer resp.Body.Close()

		var status struct {
			Me     struct{ ID string } `json:"me"`
			Leader struct{ ID string } `json:"leader"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			continue
		}

		if status.Me.ID == status.Leader.ID {
			return nodeURL, nil // este es el líder
		}
	}
	return "", fmt.Errorf("líder no encontrado para shard %d", shard)
}

func hashModulo(clave string) int {
	h := sha1.New()
	h.Write([]byte(clave))
	sum := h.Sum(nil)
	fmt.Printf("DEBUG clave=%s hash[0]=%d → shard=%d\n", clave, sum[0], int(sum[0])%2)
	return int(sum[0]) % 2
}

type SetRequest struct {
	Clave string `json:"clave"`
	Valor string `json:"valor"`
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Clave string `json:"clave"`
		Valor string `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	shard := hashModulo(req.Clave)
	leaderURL, err := getLeaderURL(shard)
	if err != nil {
		log.Printf("[SET] No se pudo detectar líder del shard %d: %v", shard, err)
		http.Error(w, "Líder no disponible", http.StatusServiceUnavailable)
		return
	}

	payload, _ := json.Marshal(map[string]string{req.Clave: req.Valor})
	resp, err := http.Post(leaderURL+"/key", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[SET] Error al POST a %s: %v", leaderURL, err)
		http.Error(w, "Error al contactar al líder", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	clave := strings.TrimPrefix(r.URL.Path, "/get/")
	if clave == "" {
		http.Error(w, "Clave faltante", http.StatusBadRequest)
		return
	}

	shard := hashModulo(clave)
	leaderURL, err := getLeaderURL(shard)
	if err != nil {
		log.Printf("[GET] No se pudo detectar líder del shard %d: %v", shard, err)
		http.Error(w, "Líder no disponible", http.StatusServiceUnavailable)
		return
	}

	resp, err := http.Get(leaderURL + "/key/" + clave)
	if err != nil {
		log.Printf("[GET] Error al GET a %s: %v", leaderURL, err)
		http.Error(w, "Error al contactar al líder", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get/", getHandler)

	fmt.Println("Distribuidor escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	handlers "cryptotrade/internal"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/prices", handlers.PricesHandler)

	// Статичные файлы (включая index.html)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Println("Server running on http://localhost:88")
	log.Fatal(http.ListenAndServe(":88", nil))
}

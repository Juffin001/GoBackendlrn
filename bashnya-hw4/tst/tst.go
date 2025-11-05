package main

import (
	"flag"
	"fmt"
)

// Объявляем обычные переменные (не указатели)
var (
	host  string
	port  int
	debug bool
)

func main() {
	// flag.StringVar() связывает флаг с существующей переменной
	// Параметры: (&переменная, "имя флага", значение по умолчанию, описание)
	flag.StringVar(&host, "host", "localhost", "Server hostname")

	// Аналогично для других типов
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")

	// Парсим флаги
	flag.Parse()

	// Теперь используем переменные напрямую (без разыменования)
	fmt.Printf("Host: %s, Port: %d, Debug: %t\n", host, port, debug)
}

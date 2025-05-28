package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Fuente struct {
	id     int
	nombre string
}

type Evento struct {
	id        int
	fuente_id int
	timestamp int
	valor     float64
}

var mapaArchivos map[string]*os.File = make(map[string]*os.File)

func main() {
	numFuentes := flag.Int("numFuentes", 1, "Numero de fuentes que se van a generar")
	numArchivos := flag.Int("numArchivos", 1, "Numero de archivos que se generan")
	numLineas := flag.Int("numLineas", 100, "Numero de lienas generadas por archivo")
	outputDir := flag.String("o", "./outputs", "Directorio de salida de los archivos")
	prefixFile := flag.String("pre", "", "Perfijo del archivo antes del numero")

	flag.Parse()

	listaFuentes := make([]Fuente, *numFuentes)
	for i := 0; i < *numFuentes; i++ {
		listaFuentes[i] = Fuente{id: i, nombre: fmt.Sprintf("Fuente-%d", i)}
	}

	for i := 0; i < *numArchivos; i++ {
		nombre := fmt.Sprintf("%s%d", *prefixFile, i)
		path := fmt.Sprintf("%s/%s.events", *outputDir, nombre)
		file, err := os.Create(path)
		mapaArchivos[path] = file
		check(err)
		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
	}

	lineas := uint(*numLineas)
	canalFin := make(chan int, *numFuentes)
	for _, file := range mapaArchivos {
		go generateLine(lineas, &listaFuentes, file, canalFin)
	}

	completos := 0
	for range *numArchivos {
		println("esperando...")
		c := <-canalFin
		completos += c
		println("Completados:", completos)
	}
	// close(canalFin)
}

func generateLine(cantidad uint, fuentes *[]Fuente, file *os.File, canal chan int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bufferLineas := bufio.NewWriter(file)
	for range cantidad {
		fuente := r.Intn(len(*fuentes) - 1)
		timestamp := r.Intn(int(time.Now().Unix()))
		valor := r.Float64()
		linea := fmt.Sprintf("%d;%d;%f\n", timestamp, (*fuentes)[fuente].id, valor)
		if _, err := bufferLineas.WriteString(linea); err != nil {
			fmt.Println("Error escribiendo el archivo" + file.Name())
		}
	}
	bufferLineas.Flush()

	canal <- 1
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

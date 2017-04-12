/*
 * base-stats : Calculate Hub base container image distribution
 */

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("docker-hub-stats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	in := csv.NewReader(file)
	distroidx := -1
	distro := make(map[string]int)
	distrocnt := 0

	for {
		record, err := in.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if distroidx < 0 {
			for i, val := range record {
				if val == "distro" {
					distroidx = i
					break
				}
			}
		} else {
			d := record[distroidx]
			distro[d]++
			distrocnt++
		}
	}

	fmt.Println("Distro,%Images,Images")
	for k, v := range distro {
		pdistro := float64(v) / float64(distrocnt) * 100.0
		data := [...]string{
			k,
			strconv.FormatFloat(pdistro, 'f', 8, 64) + "%",
			strconv.Itoa(v),
		}
		fmt.Println(strings.Join(data[:], ","))
	}
}

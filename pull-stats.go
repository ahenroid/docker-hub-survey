/*
 * pull-stats : Calculate Hub repository pull count distribution
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
	pullbkt := [...]int{50000000, 10000000, 5000000, 1000000, 500000,
		100000, 50000, 10000, 5000, 1000, 0}

	file, err := os.Open("docker-hub-repos.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	in := csv.NewReader(file)
	pullidx := -1
	var img [len(pullbkt)]int
	imgcnt := 0
	var pull [len(pullbkt)]int
	pullcnt := 0
	for {
		record, err := in.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if pullidx < 0 {
			for i, val := range record {
				if val == "pull_count" {
					pullidx = i
					break
				}
			}
		} else {
			p, _ := strconv.Atoi(record[pullidx])
			for i, val := range pullbkt {
				if p >= val {
					img[i]++
					imgcnt++
					pull[i] += p
					pullcnt += p
					break
				}
			}
		}
	}

	fmt.Println("Pull Count,%Images,%Pulls,Images,Pulls")
	for i, val := range pullbkt {
		pimg := float64(img[i]) / float64(imgcnt) * 100.0
		ppull := float64(pull[i]) / float64(pullcnt) * 100.0

		data := [...]string{
			strconv.Itoa(val),
			strconv.FormatFloat(pimg, 'f', 8, 64) + "%",
			strconv.FormatFloat(ppull, 'f', 8, 64) + "%",
			strconv.Itoa(img[i]),
			strconv.Itoa(pull[i]),
		}
		fmt.Println(strings.Join(data[:], ","))
	}
}

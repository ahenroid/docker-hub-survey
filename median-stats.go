/*
 * median-stats : Calculate Hub median statistics
 */

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	distroheavy := map[string]int{
		"debian": 1,
		"ubuntu": 1,
		"centos": 1,
		"fedora": 1}
	
	file, err := os.Open("docker-hub-stats.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	in := csv.NewReader(file)
	idx := make(map[string]int)
	idx["distro"] = -1
	heavy := make(map[string][]int)
	light := make(map[string][]int)
	
	for {
		record, err := in.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if idx["distro"] < 0 {
			for i, val := range record {
				switch val {
				case "distro", "image_size", "packages", "files", "package_files":
					idx[val] = i
				}
			}
		} else {
			distro := record[idx["distro"]]
			sz, _ := strconv.Atoi(record[idx["image_size"]])
			pkgs, _ := strconv.Atoi(record[idx["packages"]])
			files, _ := strconv.Atoi(record[idx["files"]])
			pfiles, _ := strconv.Atoi(record[idx["package_files"]])
			if _, ok := distroheavy[distro]; ok {
				heavy["image_size"] =
					append(heavy["image_size"], sz)
				heavy["packages"] =
					append(heavy["packages"], pkgs)
				heavy["files"] =
					append(heavy["files"], files)
				heavy["ufiles"] =
					append(heavy["ufiles"], files - pfiles)
			} else {
				light["image_size"] =
					append(light["image_size"], sz)
				light["packages"] =
					append(light["packages"], pkgs)
				light["files"] =
					append(light["files"], files)
				light["ufiles"] =
					append(light["ufiles"], files - pfiles)
			}
		}
	}

	var median = func(data []int) int{
		sort.Ints(data)
		return data[len(data) / 2]
	}
	
	fmt.Println("Type,Image Size,Packages,Files,Unmanaged Files")

	datah := [...]string{
		"Heavyweight",
		strconv.Itoa(median(heavy["image_size"][:])),
		strconv.Itoa(median(heavy["packages"][:])),
		strconv.Itoa(median(heavy["files"][:])),
		strconv.Itoa(median(heavy["ufiles"][:])),
	}
	fmt.Println(strings.Join(datah[:], ","))

	datal := [...]string{
		"Lightweight",
		strconv.Itoa(median(light["image_size"][:])),
		strconv.Itoa(median(light["packages"][:])),
		strconv.Itoa(median(light["files"][:])),
		strconv.Itoa(median(light["ufiles"][:])),
	}
	fmt.Println(strings.Join(datal[:], ","))
}

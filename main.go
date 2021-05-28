package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const modu = "/go.mod"

type Rename struct {
	name                      string
	dirFiles                  []string
	modName, newName, dirRoot string
}

func (r *Rename) filePath() {

	err := filepath.Walk(r.dirRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dirs := strings.Split(path, "/n")

		for i := 0; i < len(dirs); i++ {
			if strings.Contains(dirs[i], ".go") {
				r.dirFiles = append(r.dirFiles, dirs[i])
			}
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}
	r.refacFile()

}

func (r Rename) refacFile() {

	for i := 0; i < len(r.dirFiles); i++ {

		file, err := os.ReadFile(r.dirFiles[i])

		if err != nil {
			return
		}

		err = r.doRefac(file, r.dirFiles[i])

		if err != nil {
			log.Println(err)
		}

	}
	fmt.Println("Succesfuly")
}

func (r Rename) doRefac(file []byte, dir string) error {
	newconten := strings.Replace(string(file), r.modName, r.newName, 1)
	err := os.WriteFile(dir, []byte(newconten), os.FileMode(0755))

	if err != nil {
		return err
	}
	return nil
}

func (r *Rename) renameMod() {
	dirModu := r.dirRoot + modu
	file, err := os.ReadFile(dirModu)

	if err != nil {
		log.Println(err)
		return
	}

	lines := strings.Split(string(file), "\n")
	r.modName = (lines[0])[len("module")+1:]

	err = r.doRefac(file, dirModu)

	if err != nil {
		log.Println(err)
		return
	}
	r.filePath()
}

func main() {

	runApp()

}

func runApp() {
	var rename = new(Rename)

	root, _ := os.Getwd()
	rename.dirRoot = root

	fmt.Print("\nInto new mod name: ")
	fmt.Scanln(&rename.newName)

	rename.renameMod()
}

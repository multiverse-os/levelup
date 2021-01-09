package main

import (
	"fmt"

	archive "github.com/multiverse-os/levelup/backend/archive"
)

func main() {
	fmt.Println("tar example")
	fmt.Println("=============================================================")

	archive := archive.Tar("./test-data")

	fmt.Println("archive:", archive)

}

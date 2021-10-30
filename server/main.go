package main

import (
	"fmt"
	"log"

	"github.com/arxeiss/sample-terraform-provider/server/database"
)

func main() {
	db, err := database.Open("superdupercloud.db", "superdupercloud.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	vm, _ := db.VirtualMachines.FindByID(1)
	fmt.Printf("%d - %s - %q\n", vm.ID, vm.Name, vm.DisplayName)
	vm, _ = db.VirtualMachines.FindByID(2)
	fmt.Printf("%d - %s - %q\n", vm.ID, vm.Name, vm.DisplayName)

	fmt.Println("starting server")
}

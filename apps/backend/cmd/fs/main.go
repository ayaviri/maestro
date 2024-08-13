package main

import "maestro/internal"

func initialiseServer() {

}

func main() {
	internal.WithTimer("running file server", initialiseServer)
}

package main

func main() {
	server := &Server{}
	server.initialize()
	server.run("8083")
}

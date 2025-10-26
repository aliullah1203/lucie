package main

func main() {
	configs.LoadEnv()
	configs.ConnectDB()
	api.RunServer()
}

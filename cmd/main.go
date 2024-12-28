package main

import (
	"fmt"
	"internship_backend_2022/internal/api"
	"internship_backend_2022/internal/app"
	"internship_backend_2022/internal/repository"
	"internship_backend_2022/internal/service"
	"log"
	"net/http"
)



func main() {
	
	cfg := app.NewConfig()

	db, err := repository.InitDB(cfg.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	Repository := repository.NewRepository(db)
	Service := service.NewService(Repository)
	Handler := api.NewHandler(Service)

	router := api.SetupRouter(Handler)

	fmt.Println("Server is running on http://localhost:8080")
	

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
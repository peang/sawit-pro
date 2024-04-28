package main

import (
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	// e.Logger.Fatal(e.Start(":1323"))
	e.Logger.Fatal(e.Start(":8080"))
}

func newServer() *handler.Server {
	// dbDsn := os.Getenv("DATABASE_URL")
	dbDsn := "postgres://docker:docker@localhost/sawit?sslmode=disable"

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

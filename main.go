package main

import (
	"log"
	"os"
	"strconv"

	"github.com/bashmohandes/go-askme/web/askme"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/question/inmemory"
	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/gobuffalo/packr"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(newConfig)
	container.Provide(newFileProvider)
	container.Provide(framework.NewApp)
	container.Provide(framework.NewRouter)
	container.Provide(framework.NewRenderer)
	container.Provide(question.NewRepository)
	container.Provide(answer.NewRepository)
	container.Provide(user.NewAsksUsecase)
	container.Provide(user.NewAnswersUsecase)
	container.Provide(controllers.NewHomeController)
	container.Provide(controllers.NewProfileController)
	container.Provide(askme.NewApp)
	err := container.Invoke(func(app *askme.App) {
		if e := app.Start(); e != nil {
			log.Fatalln(e)
		}
	})

	if err != nil {
		panic(err)
	}
}

func newFileProvider(config *framework.Config) shared.FileProvider {
	log.Printf("Initializing box to path %s\n", config.PublicFolder)
	wd, _ := os.Getwd()
	log.Printf("Working directory %s\n", wd)
	return packr.NewBox(config.PublicFolder)
}

func newConfig() *framework.Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Incorrect format: %v\n", err)
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		log.Fatalf("Incorrect format: %v\n", err)
	}
	public := os.Getenv("PUBLIC_FOLDER")
	return &framework.Config{
		Debug:        debug,
		PublicFolder: public,
		Port:         port,
	}
}
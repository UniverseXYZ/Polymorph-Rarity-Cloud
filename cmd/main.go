package main

import (
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/vikinatora/rarity-cloud-function/handlers"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		godotenv.Load(args[0])
	} else {
		godotenv.Load()
	}

	polymorphDBName := os.Getenv("POLYMORPH_DB")
	rarityCollectionName := os.Getenv("RARITY_COLLECTION")

	funcframework.RegisterHTTPFunction("/", handlers.GetPolymorphs(polymorphDBName, rarityCollectionName))

	if err := funcframework.Start("8001"); err != nil {
		log.Errorln("funcframework.Start: %v\n", err)
	}
}

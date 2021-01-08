package main

import (
	"log"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

func addMongoDbTags(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			if model.Name == "Trade" && field.Name == "closePrice" {
				field.Tag += ` bson:"close_price"`
			}
		}
	}
	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		log.Fatal(err)
	}

	p := modelgen.Plugin{
		MutateHook: addMongoDbTags,
	}

	err = api.Generate(cfg,
		api.NoPlugins(),
		api.AddPlugin(&p),
	)

	if err != nil {
		log.Fatal(err)
	}
}

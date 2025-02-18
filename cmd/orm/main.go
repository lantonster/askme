package main

import (
	"github.com/lantonster/askme/internal/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./pkg/orm",
		Mode:    gen.WithDefaultQuery, // generate mode
	})

	g.ApplyBasic(
		model.Config{},
		model.SiteInfo{},
		model.User{},
	)

	g.Execute()
}

package main

import (
	"context"

	"github.com/lantonster/askme/cmd/wire"
)

func main() {
	wire.Init().Run(context.Background())
}

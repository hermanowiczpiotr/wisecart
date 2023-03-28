package cli

import (
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/repository"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/infrastructure/shopify"
	"github.com/urfave/cli/v2"
	"os"
)

type ShopifyCli struct {
	Client                 shopify.Client
	StoreProfileRepository repository.ProfileRepository
}

func (cs *ShopifyCli) Run() {
	app := cli.NewApp()
	app.Name = "shopify"
	app.Usage = "A simple CLI application"
	app.Commands = cli.Commands{
		{
			Name:  "fetch-products",
			Usage: "Fetch products by profile id",
			Action: func(c *cli.Context) error {
				storeProfile, _ := cs.StoreProfileRepository.GetByUserId(c.Args().Get(0))
				_, err := cs.Client.FetchProducts(storeProfile)
				if err != nil {
					return err
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

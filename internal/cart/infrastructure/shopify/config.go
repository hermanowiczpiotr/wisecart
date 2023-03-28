package shopify

import (
	goshopify "github.com/bold-commerce/go-shopify"
	"os"
)

var (
	app = goshopify.App{
		ApiKey:      os.Getenv("SHOPIFY_API_KEY"),
		ApiSecret:   os.Getenv("SHOPIFY_API_SECRET"),
		RedirectUrl: "http://localhost:8080/shopify/callback",
		Scope:       "read_products",
	}
)

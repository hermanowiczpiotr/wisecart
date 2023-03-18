package dto

type ProductDto struct {
	ShopifyId   int64
	Title       string
	Description string
}

type ProductDtoList struct {
	Products []ProductDto
}

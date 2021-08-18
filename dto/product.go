package dto

type Product struct {
	ProductID int64           `json:"product_id"`
	Name      UppercaseString `json:"name"`
	Price     int64           `json:"price"`
	CreatedBy string          `json:"created_by"`
	CreatedAt int64           `json:"crated_at"`
	Image     string          `json:"image"`
}

type ProductReq struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

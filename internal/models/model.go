package models


import (
	"time"
)

type OrderItem struct {
	SKUID       int     `bson:"sku_id" json:"sku_id"`
	Quantity    int     `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID            int         `bson:"_id,omitempty" json:"id"`
	SellerID      int         `bson:"seller_id" json:"seller_id"`
	TenantID      int         `bson:"tenant_id" json:"tenant_id"`
	HubID         int         `bson:"hub_id" json:"hub_id"`
	CreatedAt     time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `bson:"updated_at" json:"updated_at"`
	OrderStatus   string      `bson:"order_status" json:"order_status"`
	ModeOfPayment string      `bson:"mode_of_payment" json:"mode_of_payment"`
	OrderItems    []OrderItem `bson:"order_items" json:"order_items"`
}
type BulkOrderQueueMessage struct {
	CustomerID string `json:"customer_id"`
	FilePath   string `json:"file_path"`
}

type BulkOrderRequest struct {
	SellerID int    `json:"sellerID"`
	FilePath string `json:"filePath"`
}

type SQSMessage struct {
	SellerID int    `json:"sellerID"`
	FilePath string `json:"filePath"`
}
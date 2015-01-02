package exercise2

import (
	"github.com/jzipfler/htw-ava/server"
)

const (
	THRESHOLD_BUY_WHEN_ADVERTISEMENT_HEARD  = 5
	THRESHOLD_BUY_WHEN_FRIED_TELLS_HE_BUYED = 8
)

// This node represents a customer for exercise 2. The customer gets informed
// from companies about new products and buys one, if he hears the advertisement
// from a company. He also buys the product if enough friends told him, that
// they buyed the product already.
type CustomerNode struct {
	server.NetworkServer
	customerId          int
	friends             map[int]server.NetworkServer //Also neighbors.
	heardAdvertiements  map[int]int                  //An entry for each company
	heardFriedBuyedItem map[int]int                  //An entry for each friend
}

func NewCustomerNode() CustomerNode {
	return CustomerNode{server.New(), 0, nil, nil, nil}
}

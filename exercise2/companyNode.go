package exercise2

import (
	"github.com/jzipfler/htw-ava/server"
)

type CompanyNode struct {
	server.NetworkServer
	companyId         int
	product           string
	advertisingBudger int
	regularCustomers  map[int]CustomerNode
}

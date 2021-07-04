package main

import (
	"log"
	"net/http"

	"wallet_service/src/api/wallet"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/wallet/new", wallet.WalletNewHandler)
	r.POST("/wallet/topup", wallet.WalletTopupHandler)
	r.POST("/wallet/transfer", wallet.WalletTransferHandler)
	r.GET("/wallet/report", wallet.WalletReportHandler)

	r.Run()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package wallet

import (
	"net/http"
	"wallet_service/src/wallet"

	"github.com/gin-gonic/gin"
)

// DefaultUserID: we haven't registration and auth for now, just using default UserID
const DefaultUserID = 1

type WalletNewInput struct {
	Name string `json:"name" form:"name" binding:"required"`
}

func WalletNewHandler(c *gin.Context) {
	var input WalletNewInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := wallet.WalletAdd(DefaultUserID, input.Name)
	if err != nil && err != wallet.WalletExistsError {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err == wallet.WalletExistsError {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

type WallerTopupInput struct {
	WalletID            uint32 `json:"wallet_id" form:"wallet_id" binding:"required"`
	Amount              uint32 `json:"amount" form:"amount" binding:"required"`
	ClientOperationHash string `json:"client_operation_hash" form:"client_operation_hash" binding:"required"`
}

func WalletTopupHandler(c *gin.Context) {
	var input WallerTopupInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := wallet.WalletTopup(input.WalletID, input.Amount, input.ClientOperationHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

type WalletTransferInput struct {
	WalletIDFrom        uint32 `json:"wallet_id_from" form:"wallet_id_from" binding:"required"`
	WalletIDTo          uint32 `json:"wallet_id_to" form:"wallet_id_to" binding:"required"`
	Amount              uint32 `json:"amount" form:"amount" binding:"required"`
	ClientOperationHash string `json:"client_operation_hash" form:"client_operation_hash" binding:"required"`
}

func WalletTransferHandler(c *gin.Context) {
	var input WalletTransferInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := wallet.WalletTransfer(input.WalletIDFrom, input.WalletIDTo, input.Amount, input.ClientOperationHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

type WalletReportInput struct {
	WalletID uint32 `json:"wallet_id" form:"wallet_id" binding:"required"`
	DateFrom int64  `json:"date_from" form:"date_from" binding:"required"`
	DateTo   int64  `json:"date_to" form:"date_to" binding:"required"`
	TypeOp   int    `json:"type" form:"type"`
}

func WalletReportHandler(c *gin.Context) {
	var input WalletReportInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := wallet.WalletReport(input.WalletID, input.DateFrom, input.DateTo, input.TypeOp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": result})
}

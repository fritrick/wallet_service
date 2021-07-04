package wallet

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"
)

const (
	DBUser     = "postgres"
	DBPassword = "mysecretpassword"
	DBName     = "wallet"

	OpTypeAny      = 0
	OpTypeDeposit  = 1
	OpTypeWithdraw = 2

	ReportDefaultAllocSize = 100

	DateFormatTemplate = "2006-01-02"
)

var (
	dbCon *sql.DB

	WalletExistsError     = fmt.Errorf("wallet with this name already exists")
	WalletNotExists       = fmt.Errorf("wallet not exists")
	WalletNotEnougBalance = fmt.Errorf("wallet not enough balance")
)

func getDBConnection() (*sql.DB, error) {
	if dbCon != nil {
		return dbCon, nil
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DBUser, DBPassword, DBName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	dbCon = db
	return db, err
}

func WalletAdd(userID uint32, name string) error {
	exists, err := walletExistsByName(name)
	if err != nil {
		return err
	}
	if exists {
		return WalletExistsError
	}

	dbCon, err := getDBConnection()
	if err != nil {
		return err
	}

	_, err = dbCon.Exec("INSERT INTO wallet (name, user_id) VALUES ($1,$2)", name, userID)
	if err != nil {
		return err
	}
	return nil
}

func walletExistsByName(name string) (bool, error) {
	dbCon, err := getDBConnection()
	if err != nil {
		return false, err
	}

	var res int
	err = dbCon.QueryRow("SELECT 1 FROM wallet WHERE name = $1", name).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func walletExists(walletID uint32) (bool, error) {
	dbCon, err := getDBConnection()
	if err != nil {
		return false, err
	}

	var res int
	err = dbCon.QueryRow("SELECT 1 FROM wallet WHERE wid = $1", walletID).Scan(&res)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func WalletTopup(walletID uint32, amount uint32, clientOperationHash string) error {
	dbCon, err := getDBConnection()
	if err != nil {
		return err
	}

	var walletBalance uint32
	err = dbCon.QueryRow("SELECT balance FROM wallet WHERE wid = $1", walletID).Scan(&walletBalance)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := dbCon.Begin()
	if err != nil {
		return err
	}

	_, err = dbCon.Exec("INSERT INTO transactions (wid, amount, client_operation_hash) VALUES ($1,$2,$3)",
		walletID, amount, clientOperationHash)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = dbCon.Exec("UPDATE wallet SET balance = $1 WHERE wid = $2",
		walletBalance+amount, walletID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func WalletTransfer(walletIDFrom uint32, walletIDTo uint32, amount uint32, clientOperationHash string) error {
	dbCon, err := getDBConnection()
	if err != nil {
		return err
	}

	walletToExists, err := walletExists(walletIDTo)
	if err != nil {
		return err
	}
	if !walletToExists {
		return WalletNotExists
	}

	walletFromExists, err := walletExists(walletIDFrom)
	if err != nil {
		return err
	}
	if !walletFromExists {
		return WalletNotExists
	}

	var walletBalance uint32
	err = dbCon.QueryRow("SELECT balance FROM wallet WHERE wid = $1", walletIDFrom).Scan(&walletBalance)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if walletBalance < amount {
		return WalletNotEnougBalance
	}

	var walletBalanceTo uint32
	err = dbCon.QueryRow("SELECT balance FROM wallet WHERE wid = $1", walletIDTo).Scan(&walletBalanceTo)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := dbCon.Begin()
	if err != nil {
		return err
	}

	_, err = dbCon.Exec("INSERT INTO transactions (wid, amount, client_operation_hash) VALUES ($1,$2,$3)",
		walletIDFrom, -1*int32(amount), clientOperationHash)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = dbCon.Exec("UPDATE wallet SET balance = $1 WHERE wid = $2",
		walletBalance-amount, walletIDFrom)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	ksID := ksuid.New()
	_, err = dbCon.Exec("INSERT INTO transactions (wid, amount, client_operation_hash) VALUES ($1,$2,$3)",
		walletIDTo, amount, ksID.String())

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = dbCon.Exec("UPDATE wallet SET balance = $1 WHERE wid = $2",
		walletBalanceTo+amount, walletIDTo)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

type ReportLine struct {
	Date   time.Time
	Amount int32
}

type Report []ReportLine

func WalletReport(walletID uint32, dateFrom int64, dateTo int64, opType int) (Report, error) {
	dbCon, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	dateFromStr := time.Unix(dateFrom, 0).Format(DateFormatTemplate)
	dateToStr := time.Unix(dateTo, 0).Format(DateFormatTemplate)

	query := "SELECT amount, create_date FROM transactions WHERE wid = $1 AND create_date > $2 AND create_date < $3 %s"
	switch opType {
	case OpTypeDeposit:
		query = fmt.Sprintf(query, " AND amount > 0")
	case OpTypeWithdraw:
		query = fmt.Sprintf(query, " AND amount < 0")
	default:
		query = fmt.Sprintf(query, "")
	}
	rows, err := dbCon.Query(query, walletID, dateFromStr, dateToStr)
	if err != nil {
		return nil, err
	}

	result := make(Report, 0, ReportDefaultAllocSize)

	for rows.Next() {
		var line ReportLine
		err = rows.Scan(&line.Amount, &line.Date)
		if err != nil {
			return nil, err
		}
		result = append(result, line)
	}

	return result, nil
}

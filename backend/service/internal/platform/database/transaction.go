package database

import (
	"context"
	"database/sql"
	"errors"

	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
)

type contextKey string

const (
	txKey contextKey = "db_transaction"
)

var CAN_NOT_START_TRANSACTION = customerror.NewInternalError(errors.New("Transaction Error"))

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// Transaction'ı context'e paketler
func (tm *TransactionManager) StartTransaction(ctx context.Context) (context.Context, error) {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, CAN_NOT_START_TRANSACTION
	}

	ctx = context.WithValue(ctx, txKey, tx)
	return ctx, nil
}

func (tm *TransactionManager) StartReadOnlyTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		//BUrada tüm akışı kesmek istiyorum
		return nil, CAN_NOT_START_TRANSACTION
	}
	ctx = context.WithValue(ctx, txKey, tx)
	return tx, nil
}

// Context içinden transaction'ı çıkarır
func GetTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

func DeferRollback(ctx context.Context, err *error) {
	tx, ok := GetTxFromContext(ctx)
	if !ok {
		return // Context içinde aktif bir transaction yoksa işlem yapma
	}

	// Recover ile bir panik olup olmadığını kontrol et
	if p := recover(); p != nil {
		tx.Rollback()
		panic(p) // Paniği durdurma, rollback sonrası tekrar fırlat (re-panic)
	}

	// Eğer işaret edilen hata nil değilse (bir hata oluşmuşsa) rollback yap
	if err != nil && *err != nil {
		tx.Rollback()
	}
}

func Commit(ctx context.Context) error {
	tx, ok := GetTxFromContext(ctx)
	if !ok {
		// Eğer bir transaction yoksa commit edilecek bir şey de yoktur.
		return nil
	}

	return tx.Commit()
}

// TransactionManager içine bu metodu ekle
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) (err error) {
	// 1. KONTROL: Eğer zaten bir transaction varsa, mevcudu kullan ve ÇIK
	if _, ok := GetTxFromContext(ctx); ok {
		// Zaten bir tx içindeyiz. Yeni bir tx başlatmıyoruz.
		// Bu yüzden defer DeferRollback veya Commit de çağırmıyoruz.
		// Sorumluluk bu transaction'ı başlatan en dıştaki fonksiyondur.
		return fn(ctx)
	}

	// 2. Eğer tx yoksa (ilk defa çağrılıyorsa) başlat
	txCtx, err := tm.StartTransaction(ctx)
	if err != nil {
		return err
	}

	// Bu kısımdan sonrası sadece "Transaction Sahibi" (en dıştaki metod) için çalışır
	defer DeferRollback(txCtx, &err)

	err = fn(txCtx)
	if err != nil {
		return err
	}

	return Commit(txCtx)
}

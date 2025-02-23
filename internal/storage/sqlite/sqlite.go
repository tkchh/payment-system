package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"infotecsTest/internal/models"
	"infotecsTest/internal/storage"
	"log"
	"strings"
	"time"
)

// Storage представляет SQLite-хранилище с подготовленными запросами.
// Содержит подключение к БД и подготовленные SQL-выражения.
type Storage struct {
	db                     *sql.DB
	stmtSelectWallet       *sql.Stmt
	stmtInsertTransaction  *sql.Stmt
	stmtSelectTransactions *sql.Stmt
}

// New инициализирует новое подключение к SQLite.
// Создает таблицы и индексы при первом запуске.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Лучше использовать INTEGER с преобразованием в коде
	//REAL не подходит для операций с деньгами

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS wallets(
		id INTEGER PRIMARY KEY,
		address TEXT NOT NULL UNIQUE,
		balance REAL DEFAULT 0.0 
	);
	CREATE INDEX IF NOT EXISTS idx_address ON wallets(address);

	CREATE TABLE IF NOT EXISTS transactions(
	    id INTEGER PRIMARY KEY,
	    from_address TEXT NOT NULL,
	    to_address TEXT NOT NULL,
	    amount REAL NOT NULL,
	    timestamp DATE DEFAULT CURRENT_DATE
	);
	CREATE INDEX IF NOT EXISTS idx_transactions ON transactions(from_address,to_address);
	`); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmtSelectWallet, err := db.Prepare("SELECT balance FROM wallets WHERE address = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmtInsertTransaction, err := db.Prepare("INSERT INTO transactions(from_address, to_address, amount, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmtSelectTransactions, err := db.Prepare(`
		SELECT from_address, to_address, amount, timestamp
		FROM transactions
		ORDER BY timestamp DESC 
		LIMIT ?
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = checkDB(db); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:                     db,
		stmtSelectWallet:       stmtSelectWallet,
		stmtInsertTransaction:  stmtInsertTransaction,
		stmtSelectTransactions: stmtSelectTransactions,
	}, nil
}

// checkD автоматически создает 10 кошельки при первом запуске.
func checkDB(db *sql.DB) error {
	const op = "storage.sqlite.checkDB"

	var count int
	if err := db.QueryRow("SELECT count(*) FROM wallets").Scan(&count); err != nil {
		return fmt.Errorf("%s: check number of wallets: %w", op, err)
	}

	if count == 0 {
		for count < 10 {
			address := uuid.NewString()
			balance := 100.0
			if _, err := db.Exec("INSERT INTO wallets(address, balance) VALUES (?,?)", address, balance); err != nil {
				if errors.Is(err, sqlite3.ErrConstraintUnique) {
					log.Println("collision")
					continue
				}
				return fmt.Errorf("%s: cant insert wallet: %w", op, err)
			}
			count++
		}
	}
	return nil
}

// GetWalletBalance возвращает баланс кошелька по адресу.
// Возвращает ErrWalletNotFound если кошелек не существует.
func (s *Storage) GetWalletBalance(address string) (models.Wallet, error) {
	const op = "storage.sqlite.GetWalletBalance"

	var wallet models.Wallet

	err := s.stmtSelectWallet.QueryRow(address).Scan(&wallet.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet, storage.ErrWalletNotFound
		}

		return wallet, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	wallet.Address = address
	return wallet, nil
}

// AddTransaction выполняет перевод между кошельками.
// Проверяет: сумму перевода, разные адреса, достаточный баланс.
func (s *Storage) AddTransaction(from, to string, amount float64) error {
	const op = "storage.sqlite.AddTransaction"

	if amount <= 0 {
		return storage.ErrIncorrectAmount
	}
	if from == to {
		return storage.ErrAddressesEqual
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	var fromBalance, toBalance float64
	err = tx.Stmt(s.stmtSelectWallet).QueryRow(from).Scan(&fromBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrWalletNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	err = tx.Stmt(s.stmtSelectWallet).QueryRow(to).Scan(&toBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrWalletNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if fromBalance < amount {
		return storage.ErrInsufficientFunds
	}

	err = s.updateBalance(tx, from, -amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.updateBalance(tx, to, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Stmt(s.stmtInsertTransaction).Exec(from, to, amount, time.Now())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// updateBalance изменяет баланс кошелька атомарно в транзакции.
// Используется внутри метода AddTransaction
func (s *Storage) updateBalance(tx *sql.Tx, address string, delta float64) error {
	const op = "storage.sqlite.updateBalance"

	_, err := tx.Exec(`
		UPDATE wallets 
		SET balance = balance + ? 
		WHERE address = ?
		`, delta, address)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetNTransactions возвращает N последних транзакций.
// Транзакции сортируются от новых к старым.
// При N <= 0 возвращает ErrInvalidRequest.
func (s *Storage) GetNTransactions(N int) ([]models.Transaction, error) {
	const op = "storage.sqlite.GetNTransactions"

	if N <= 0 {
		return nil, storage.ErrInvalidRequest
	}

	var txs []models.Transaction

	rows, err := s.stmtSelectTransactions.Query(N)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var tx models.Transaction
		var date time.Time
		if err = rows.Scan(&tx.From, &tx.To, &tx.Amount, &date); err != nil {
			return txs, fmt.Errorf("%s: %w", op, err)
		}
		tx.Time = date.Format("15:04:05 02-01-2006")
		txs = append(txs, tx)
	}

	return txs, nil
}

// Close Закрывает все подготовленные выражения и соединение.
// Возвращает объединенные ошибки при их наличии.
func (s *Storage) Close() error {
	const op = "storage.sqlite.Close"

	errs := make([]string, 0, 4)
	if err := s.stmtSelectWallet.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := s.stmtInsertTransaction.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := s.stmtSelectTransactions.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if err := s.db.Close(); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) != 0 {
		return fmt.Errorf("%s: %s", op, strings.Join(errs, ", "))
	}
	return nil
}

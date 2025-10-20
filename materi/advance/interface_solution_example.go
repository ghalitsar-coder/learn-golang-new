package main

import (
	"fmt"
	"time"
)

// PaymentRequest adalah struct yang berisi semua informasi yang mungkin dibutuhkan oleh berbagai metode pembayaran
type PaymentRequest struct {
	Amount      float64
	Currency    string
	Destination string
	Source      string
	Reference   string
	// Tambahkan field lain yang mungkin dibutuhkan oleh implementasi tertentu
}

// PaymentMethod adalah interface untuk metode pembayaran
type PaymentMethod interface {
	ProcessPayment(req PaymentRequest) (string, error)
	GetPaymentInfo() string
}

// CreditCard adalah struct untuk pembayaran dengan kartu kredit
type CreditCard struct {
	CardNumber string
	ExpiryDate string
	CVV        string
}

// ProcessPayment untuk CreditCard hanya menggunakan beberapa field dari PaymentRequest
func (cc CreditCard) ProcessPayment(req PaymentRequest) (string, error) {
	// Simulasi pemrosesan pembayaran kartu kredit hanya dengan amount
	transactionID := fmt.Sprintf("CC-%d", time.Now().Unix())
	fmt.Printf("Memroses pembayaran kartu kredit sebesar $%.2f\n", req.Amount)
	return transactionID, nil
}

func (cc CreditCard) GetPaymentInfo() string {
	return fmt.Sprintf("Kartu Kredit dengan nomor ending in %s", cc.CardNumber[len(cc.CardNumber)-4:])
}

// PayPal adalah struct untuk pembayaran dengan PayPal
type PayPal struct {
	Email    string
	Password string
}

// ProcessPayment untuk PayPal bisa menggunakan lebih banyak field dari PaymentRequest
func (pp PayPal) ProcessPayment(req PaymentRequest) (string, error) {
	// Simulasi pemrosesan pembayaran PayPal dengan lebih banyak informasi
	transactionID := fmt.Sprintf("PP-%d", time.Now().Unix())
	fmt.Printf("Memroses pembayaran PayPal sebesar $%.2f dari %s ke %s dalam mata uang %s\n", 
		req.Amount, pp.Email, req.Destination, req.Currency)
	return transactionID, nil
}

func (pp PayPal) GetPaymentInfo() string {
	return fmt.Sprintf("Akun PayPal %s", pp.Email)
}

// ProcessOrder adalah fungsi yang menggunakan interface PaymentMethod
func ProcessOrder(paymentMethod PaymentMethod, req PaymentRequest) {
	fmt.Println("Menggunakan metode pembayaran:", paymentMethod.GetPaymentInfo())
	
	transactionID, err := paymentMethod.ProcessPayment(req)
	if err != nil {
		fmt.Printf("Gagal memproses pembayaran: %v\n", err)
		return
	}
	
	fmt.Printf("Pembayaran berhasil! ID Transaksi: %s\n", transactionID)
	fmt.Println("---")
}

func main() {
	creditCard := CreditCard{
		CardNumber: "1234567890123456",
		ExpiryDate: "12/25",
		CVV:        "123",
	}
	
	paypal := PayPal{
		Email:    "user@example.com",
		Password: "password", // dalam dunia nyata, jangan simpan password di sini!
	}
	
	// Request pembayaran untuk CreditCard - hanya menggunakan amount
	creditCardReq := PaymentRequest{
		Amount: 99.99,
	}
	
	// Request pembayaran untuk PayPal - menggunakan lebih banyak field
	paypalReq := PaymentRequest{
		Amount:      49.99,
		Currency:    "USD",
		Destination: "merchant@example.com",
	}
	
	ProcessOrder(creditCard, creditCardReq)
	ProcessOrder(paypal, paypalReq)
}
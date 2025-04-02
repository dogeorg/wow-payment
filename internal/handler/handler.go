package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dogeorg/wow-payment/internal/config"
	"github.com/dogeorg/wow-payment/internal/database"
	"github.com/dogeorg/wow-payment/internal/models"
)

type EmailRequest struct {
	ReplyToEmail string `json:"reply_to_email"`
	ReplyToName  string `json:"reply_to_name"`
	ToEmail      string `json:"to_email"`
	ToName       string `json:"to_name"`
	Subject      string `json:"subject"`
	HTML         string `json:"html"`
}

type GigaWalletAccount struct {
	ForeignID       string `json:"foreign_id"`
	PayoutAddress   string `json:"payout_address,omitempty"`
	PayoutThreshold string `json:"payout_threshold,omitempty"`
	PayoutFrequency string `json:"payout_frequency,omitempty"`
}

type GigaWalletInvoiceRequest struct {
	RequiredConfirmations int `json:"required_confirmations"`
	Items                 []struct {
		Type     string  `json:"type"`
		Name     string  `json:"name"`
		Sku      string  `json:"sku"`
		Value    float64 `json:"value"`
		Quantity int     `json:"quantity"`
	} `json:"items"`
}

type GigaWalletInvoiceResponse struct {
	ID        string  `json:"id"`
	ForeignID string  `json:"foreign_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
}

func RegisterHandler(cfg config.Config) http.HandlerFunc {
	db, err := database.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Much Sad! DB init failed: %v", err)
	}

	client := &http.Client{}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Such Sad! Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RegistrationRequest // Use models.RegistrationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Wow! Invalid JSON", http.StatusBadRequest)
			return
		}

		// Create GigaWallet Account (using admin endpoint)
		gigaAdminURL := fmt.Sprintf("%s:%d/account/%s", cfg.GigaWallet.Host, cfg.GigaWallet.AdminPort, req.DogeAddress)
		accountReq := GigaWalletAccount{
			ForeignID:       req.DogeAddress,
			PayoutAddress:   req.DogeAddress,
			PayoutThreshold: "0",
			PayoutFrequency: "0",
		}
		accountJSON, _ := json.Marshal(accountReq)
		gigaReq, _ := http.NewRequest("POST", gigaAdminURL, bytes.NewBuffer(accountJSON))
		gigaReq.Header.Set("Content-Type", "application/json")
		gigaReq.Header.Set("Authorization", "Bearer "+cfg.GigaWallet.AdminBearerToken)
		_, err = client.Do(gigaReq)
		if err != nil {
			log.Printf("Wow! GigaWallet account failed: %v", err)
			http.Error(w, "Such Sad! GigaWallet error", http.StatusInternalServerError)
			return
		}

		// Create GigaWallet Invoice (using admin endpoint)
		invoiceReq := GigaWalletInvoiceRequest{
			RequiredConfirmations: 1,
			Items: []struct {
				Type     string  `json:"type"`
				Name     string  `json:"name"`
				Sku      string  `json:"sku"`
				Value    float64 `json:"value"`
				Quantity int     `json:"quantity"`
			}{
				{Type: "item", Name: req.Sku, Sku: req.Sku, Value: req.Amount, Quantity: 1},
			},
		}
		invoiceJSON, _ := json.Marshal(invoiceReq)
		invoiceURL := fmt.Sprintf("%s:%d/account/%s/invoice/", cfg.GigaWallet.Host, cfg.GigaWallet.AdminPort, req.DogeAddress)
		invoiceHTTPReq, _ := http.NewRequest("POST", invoiceURL, bytes.NewBuffer(invoiceJSON))
		invoiceHTTPReq.Header.Set("Content-Type", "application/json")
		invoiceHTTPReq.Header.Set("Authorization", "Bearer "+cfg.GigaWallet.AdminBearerToken)
		resp, err := client.Do(invoiceHTTPReq)
		if err != nil {
			log.Printf("Wow! Invoice creation failed: %v", err)
			http.Error(w, "Such Sad! Invoice error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var invoice GigaWalletInvoiceResponse
		err = json.NewDecoder(resp.Body).Decode(&invoice)
		if err != nil {
			log.Printf("Wow! Invoice decode failed: %v", err)
			http.Error(w, "Such Sad! Invoice parse error", http.StatusInternalServerError)
			return
		}
		req.PaytoDogeAddress = invoice.ID

		// Store in database
		_, err = database.InsertShibe(db, req)
		if err != nil {
			log.Printf("Much Sad! DB error: %v", err)
			http.Error(w, "Such Sad! Failed to save", http.StatusInternalServerError)
			return
		}

		// Send email via much-sender
		emailReq := EmailRequest{
			ReplyToEmail: cfg.EmailService.ReplyToEmail,
			ReplyToName:  cfg.EmailService.ReplyToName,
			ToEmail:      req.Email,
			ToName:       req.Name,
			Subject:      cfg.EmailService.Subject,
			HTML:         fmt.Sprintf("<h1>Much Wow!</h1><p>Welcome %s! Pay %.2f DOGE to %s for %s</p>", req.Name, req.Amount, req.PaytoDogeAddress, req.Sku),
		}
		emailJSON, _ := json.Marshal(emailReq)
		emailURL := fmt.Sprintf("%s:%d/send-email", cfg.EmailService.Host, cfg.EmailService.Port)
		emailReqHTTP, _ := http.NewRequest("POST", emailURL, bytes.NewBuffer(emailJSON))
		emailReqHTTP.Header.Set("Authorization", "Bearer "+cfg.EmailService.BearerToken)
		emailReqHTTP.Header.Set("Content-Type", "application/json")
		emailResp, err := client.Do(emailReqHTTP)
		if err != nil || emailResp.StatusCode != http.StatusOK {
			log.Printf("Wow! Email send failed: %v", err)
		}
		if emailResp != nil {
			emailResp.Body.Close()
		}

		// Response with invoice details
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":          "Such Success!",
			"paytoDogeAddress": req.PaytoDogeAddress,
			"invoice":          invoice,
		})
	}
}

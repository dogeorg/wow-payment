# Wow-Payment

Wow Payment is a Golang payment registration service that stores user details in SQLite, integrates with GigaWallet for Dogecoin payments, and sends email notifications via much-sender.

## Features
- Accepts JSON POST requests for user registration
- Stores user data in SQLite database
- Creates GigaWallet accounts and invoices for Dogecoin payments
- Returns GigaWallet invoice details in JSON response
- Sends email notifications through much-sender
- Configurable via TOML file
- Much Doge-inspired goodness

## Installation
```bash
git clone https://github.com/dogeorg/wow-payment.git
cd wow-payment
go mod init github.com/dogeorg/wow-payment
go get
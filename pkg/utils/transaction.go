package utils

type PayStackPayment struct {
	Key      string
	Email    string
	Amount   int64
	Ref      string
	Currency string
	Channels []string
	Label    string
}

type InitializePayStackPaymentRequest struct {
	Email       string `json:"email"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Reference   string `json:"reference"`
	CallbackURL string `json:"callback_url"`
	Channels    string `json:"channels"`
	Bearer      string `json:"bearer"`
}

// https://paystack.com/docs/payments/accept-payments/
func GeneratePayStackTransaction() {

}

package indto

type SendMailParams struct {
	Recipient string
	Subject   string
	Message   string
}

type UserForgotPasswordHTML struct {
	Username string
	Email    string
	OTPCode  string
}

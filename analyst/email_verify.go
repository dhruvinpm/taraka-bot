package analyst

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

type EmailVerifier struct {
	verifier *emailverifier.Verifier
}

func NewEmailVerifier() *EmailVerifier {
	v := emailverifier.NewVerifier().EnableSMTPCheck()
	return &EmailVerifier{verifier: v}
}

func (e *EmailVerifier) Verify(email string) (bool, error) {
	ret, err := e.verifier.Verify(email)
	if err != nil {
		return false, err
	}
	return ret.Syntax.Valid, nil
}

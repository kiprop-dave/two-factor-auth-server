package twofa

import (
	"fmt"
	"time"

	"github.com/xlzd/gotp"
)

type TwoFaRegistration struct {
	Secret string
	Uri    string
}

type TwoFaService interface {
	VerifyCode(secret string, code string) (bool, error)
	GenerateTwoFa(email string) (*TwoFaRegistration, error)
}

type Authy struct{}

func (a *Authy) GenerateTwoFa(email string) (*TwoFaRegistration, error) {
	secret := gotp.RandomSecret(32)
	if len(secret) == 0 {
		return nil, fmt.Errorf("failed to generate secret")
	}
	totp := gotp.NewDefaultTOTP(secret)
	uri := totp.ProvisioningUri(email, "Secure2Fa")
	return &TwoFaRegistration{Secret: secret, Uri: uri}, nil
}

func (a *Authy) VerifyCode(secret string, code string) (bool, error) {
	totp := gotp.NewDefaultTOTP(secret)
	valid := totp.Verify(code, time.Now().Unix())
	return valid, nil
}

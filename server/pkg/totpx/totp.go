package totpx

import "github.com/pquerna/otp/totp"

func Generate(accountName string) (string, string, error) {
	// 生成一个随机的密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Apollo11",
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func Validate(verifiCode, secretKey string) bool {
	return totp.Validate(verifiCode, secretKey)
}

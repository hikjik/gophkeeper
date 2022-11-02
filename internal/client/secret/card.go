package secret

import (
	"bytes"
	"text/template"

	"github.com/rs/zerolog/log"
)

var _ Secret = (*Card)(nil)

// Card данные банковской карты
type Card struct {
	Number       string
	ExpiryDate   string
	SecurityCode string
	Holder       string
}

// Type возвращает тип хранимой информации
func (c Card) Type() string {
	return secretTypeCard
}

// String функция отображения приватной информации
func (c Card) String() string {
	tmpl := "Number: {{.Number}}, ExpiryDate: {{.ExpiryDate}}, SecurityCode: {{.SecurityCode}}, Holder: {{.Holder}}"

	t := template.Must(template.New("secret").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "secret", c); err != nil {
		log.Warn().Err(err)
	}
	return buf.String()
}

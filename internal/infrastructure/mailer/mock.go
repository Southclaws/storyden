package mailer

import "fmt"

type Mock struct{}

func (m *Mock) Mail(toname, toaddr, subj, rich, text string) error {
	fmt.Printf(`Mock email sender: %s <%s>\nFROM %s\nRICH:\n%s\nTEXT:\n%s\m`, toname, toaddr, subj, rich, text)
	return nil
}

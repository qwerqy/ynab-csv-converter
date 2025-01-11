package bank

type Bank struct {
	HSBC *HSBCBank
}

func NewBank() Bank {
	return Bank{
		HSBC: &HSBCBank{
			Credit: &HSBCCredit{},
			Debit:  &HSBCDebit{},
		},
	}
}

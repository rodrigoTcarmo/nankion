package ofx

type OFXOperationType string

// OFX Transaction Types (TRNTYPE)
const (
	OFXOperationTypeCredit      OFXOperationType = "CREDIT"      // Generic credit (money in)
	OFXOperationTypeDebit       OFXOperationType = "DEBIT"       // Generic debit (money out)
	OFXOperationTypeInt         OFXOperationType = "INT"         // Interest earned
	OFXOperationTypeDiv         OFXOperationType = "DIV"         // Dividend
	OFXOperationTypeFee         OFXOperationType = "FEE"         // FI fee
	OFXOperationTypeSrvChg      OFXOperationType = "SRVCHG"      // Service charge
	OFXOperationTypeDep         OFXOperationType = "DEP"         // Deposit
	OFXOperationTypeATM         OFXOperationType = "ATM"         // ATM debit or credit
	OFXOperationTypePOS         OFXOperationType = "POS"         // Point of sale debit or credit
	OFXOperationTypeXfer        OFXOperationType = "XFER"        // Transfer
	OFXOperationTypeCheck       OFXOperationType = "CHECK"       // Check
	OFXOperationTypePayment     OFXOperationType = "PAYMENT"     // Electronic payment
	OFXOperationTypeCash        OFXOperationType = "CASH"        // Cash withdrawal
	OFXOperationTypeDirectDep   OFXOperationType = "DIRECTDEP"   // Direct deposit
	OFXOperationTypeDirectDebit OFXOperationType = "DIRECTDEBIT" // Merchant initiated debit
	OFXOperationTypeRepeatPmt   OFXOperationType = "REPEATPMT"   // Repeating payment
	OFXOperationTypeOther       OFXOperationType = "OTHER"       // Other
)

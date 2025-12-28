package properties

import (
	"github.com/dstotijn/go-notion"
	"github.com/rodrigoTcarmo/nankion/pkg/models/ofx"
)

var RequiredProperties = map[string]*notion.DatabaseProperty{
	"Transaction Date": {
		ID:   "transaction_date",
		Name: "Transaction Date",
		Type: notion.DBPropTypeDate,
		Date: &notion.EmptyMetadata{},
	},
	"Operation": {
		ID:   "operation",
		Name: "Operation",
		Type: notion.DBPropTypeSelect,
		Select: &notion.SelectMetadata{Options: []notion.SelectOptions{
			{Name: string(ofx.OFXOperationTypeCredit)},
			{Name: string(ofx.OFXOperationTypeDebit)},
			{Name: string(ofx.OFXOperationTypeInt)},
			{Name: string(ofx.OFXOperationTypeDiv)},
			{Name: string(ofx.OFXOperationTypeFee)},
			{Name: string(ofx.OFXOperationTypeSrvChg)},
			{Name: string(ofx.OFXOperationTypeDep)},
			{Name: string(ofx.OFXOperationTypeATM)},
			{Name: string(ofx.OFXOperationTypePOS)},
			{Name: string(ofx.OFXOperationTypeXfer)},
			{Name: string(ofx.OFXOperationTypeCheck)},
			{Name: string(ofx.OFXOperationTypePayment)},
			{Name: string(ofx.OFXOperationTypeCash)},
			{Name: string(ofx.OFXOperationTypeDirectDep)},
			{Name: string(ofx.OFXOperationTypeDirectDebit)},
			{Name: string(ofx.OFXOperationTypeRepeatPmt)},
			{Name: string(ofx.OFXOperationTypeOther)},
		}},
	},
	"Destination": {
		ID:       "destination",
		Name:     "Destination",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Amount": {
		ID:     "amount",
		Name:   "Amount",
		Type:   notion.DBPropTypeNumber,
		Number: &notion.NumberMetadata{Format: notion.NumberFormatReal},
	},
	"Memo": {
		ID:       "memo",
		Name:     "Memo",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Transaction ID": {
		ID:       "transaction_id",
		Name:     "Transaction ID",
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
}

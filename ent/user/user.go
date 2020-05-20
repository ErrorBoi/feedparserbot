// Code generated by entc, DO NOT EDIT.

package user

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"    // FieldEmail holds the string denoting the email vertex property in the database.
	FieldEmail       = "email" // FieldTgID holds the string denoting the tg_id vertex property in the database.
	FieldTgID        = "tg_id" // FieldPaymentInfo holds the string denoting the payment_info vertex property in the database.
	FieldPaymentInfo = "payment_info"

	// Table holds the table name of the user in the database.
	Table = "users"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldEmail,
	FieldTgID,
	FieldPaymentInfo,
}
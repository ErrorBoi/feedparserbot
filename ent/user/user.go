// Code generated by entc, DO NOT EDIT.

package user

import (
	"fmt"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"           // FieldEmail holds the string denoting the email vertex property in the database.
	FieldEmail       = "email"        // FieldTgID holds the string denoting the tg_id vertex property in the database.
	FieldTgID        = "tg_id"        // FieldPaymentInfo holds the string denoting the payment_info vertex property in the database.
	FieldPaymentInfo = "payment_info" // FieldRole holds the string denoting the role vertex property in the database.
	FieldRole        = "role"

	// EdgeSettings holds the string denoting the settings edge name in mutations.
	EdgeSettings = "settings"
	// EdgeSources holds the string denoting the sources edge name in mutations.
	EdgeSources = "sources"

	// Table holds the table name of the user in the database.
	Table = "users"
	// SettingsTable is the table the holds the settings relation/edge.
	SettingsTable = "user_settings"
	// SettingsInverseTable is the table name for the UserSettings entity.
	// It exists in this package in order to avoid circular dependency with the "usersettings" package.
	SettingsInverseTable = "user_settings"
	// SettingsColumn is the table column denoting the settings relation/edge.
	SettingsColumn = "user_settings"
	// SourcesTable is the table the holds the sources relation/edge. The primary key declared below.
	SourcesTable = "user_sources"
	// SourcesInverseTable is the table name for the Source entity.
	// It exists in this package in order to avoid circular dependency with the "source" package.
	SourcesInverseTable = "sources"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldEmail,
	FieldTgID,
	FieldPaymentInfo,
	FieldRole,
}

var (
	// SourcesPrimaryKey and SourcesColumn2 are the table columns denoting the
	// primary key for the sources relation (M2M).
	SourcesPrimaryKey = []string{"user_id", "source_id"}
)

var ()

// Role defines the type for the role enum field.
type Role string

// RoleUser is the default Role.
const DefaultRole = RoleUser

// Role values.
const (
	RoleUser   Role = "user"
	RoleEditor Role = "editor"
	RoleAdmin  Role = "admin"
)

func (s Role) String() string {
	return string(s)
}

// RoleValidator is a validator for the "r" field enum values. It is called by the builders before save.
func RoleValidator(r Role) error {
	switch r {
	case RoleUser, RoleEditor, RoleAdmin:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for role field: %q", r)
	}
}

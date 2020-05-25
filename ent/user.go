// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
	"github.com/facebookincubator/ent/dialect/sql"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Email holds the value of the "email" field.
	Email *string `json:"email,omitempty"`
	// TgID holds the value of the "tg_id" field.
	TgID int `json:"tg_id,omitempty"`
	// PaymentInfo holds the value of the "payment_info" field.
	PaymentInfo *string `json:"payment_info,omitempty"`
	// Role holds the value of the "role" field.
	Role user.Role `json:"role,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserQuery when eager-loading is set.
	Edges UserEdges `json:"edges"`
}

// UserEdges holds the relations/edges for other nodes in the graph.
type UserEdges struct {
	// Settings holds the value of the settings edge.
	Settings *UserSettings
	// Sources holds the value of the sources edge.
	Sources []*Source
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// SettingsOrErr returns the Settings value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserEdges) SettingsOrErr() (*UserSettings, error) {
	if e.loadedTypes[0] {
		if e.Settings == nil {
			// The edge settings was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: usersettings.Label}
		}
		return e.Settings, nil
	}
	return nil, &NotLoadedError{edge: "settings"}
}

// SourcesOrErr returns the Sources value or an error if the edge
// was not loaded in eager-loading.
func (e UserEdges) SourcesOrErr() ([]*Source, error) {
	if e.loadedTypes[1] {
		return e.Sources, nil
	}
	return nil, &NotLoadedError{edge: "sources"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // email
		&sql.NullInt64{},  // tg_id
		&sql.NullString{}, // payment_info
		&sql.NullString{}, // role
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(values ...interface{}) error {
	if m, n := len(values), len(user.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	u.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field email", values[0])
	} else if value.Valid {
		u.Email = new(string)
		*u.Email = value.String
	}
	if value, ok := values[1].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field tg_id", values[1])
	} else if value.Valid {
		u.TgID = int(value.Int64)
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field payment_info", values[2])
	} else if value.Valid {
		u.PaymentInfo = new(string)
		*u.PaymentInfo = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field role", values[3])
	} else if value.Valid {
		u.Role = user.Role(value.String)
	}
	return nil
}

// QuerySettings queries the settings edge of the User.
func (u *User) QuerySettings() *UserSettingsQuery {
	return (&UserClient{config: u.config}).QuerySettings(u)
}

// QuerySources queries the sources edge of the User.
func (u *User) QuerySources() *SourceQuery {
	return (&UserClient{config: u.config}).QuerySources(u)
}

// Update returns a builder for updating this User.
// Note that, you need to call User.Unwrap() before calling this method, if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return (&UserClient{config: u.config}).UpdateOne(u)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v", u.ID))
	if v := u.Email; v != nil {
		builder.WriteString(", email=")
		builder.WriteString(*v)
	}
	builder.WriteString(", tg_id=")
	builder.WriteString(fmt.Sprintf("%v", u.TgID))
	if v := u.PaymentInfo; v != nil {
		builder.WriteString(", payment_info=")
		builder.WriteString(*v)
	}
	builder.WriteString(", role=")
	builder.WriteString(fmt.Sprintf("%v", u.Role))
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}

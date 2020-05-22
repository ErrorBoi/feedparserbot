// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
	"github.com/facebookincubator/ent/dialect/sql"
)

// UserSettings is the model entity for the UserSettings schema.
type UserSettings struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UrgentWords holds the value of the "urgent_words" field.
	UrgentWords []string `json:"urgent_words,omitempty"`
	// BannedWords holds the value of the "banned_words" field.
	BannedWords []string `json:"banned_words,omitempty"`
	// Language holds the value of the "language" field.
	Language usersettings.Language `json:"language,omitempty"`
	// SendingFrequency holds the value of the "sending_frequency" field.
	SendingFrequency usersettings.SendingFrequency `json:"sending_frequency,omitempty"`
	// LastSending holds the value of the "last_sending" field.
	LastSending time.Time `json:"last_sending,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserSettingsQuery when eager-loading is set.
	Edges         UserSettingsEdges `json:"edges"`
	user_settings *int
}

// UserSettingsEdges holds the relations/edges for other nodes in the graph.
type UserSettingsEdges struct {
	// User holds the value of the user edge.
	User *User
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserSettingsEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// The edge user was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserSettings) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&[]byte{},         // urgent_words
		&[]byte{},         // banned_words
		&sql.NullString{}, // language
		&sql.NullString{}, // sending_frequency
		&sql.NullTime{},   // last_sending
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*UserSettings) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // user_settings
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserSettings fields.
func (us *UserSettings) assignValues(values ...interface{}) error {
	if m, n := len(values), len(usersettings.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	us.ID = int(value.Int64)
	values = values[1:]

	if value, ok := values[0].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field urgent_words", values[0])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &us.UrgentWords); err != nil {
			return fmt.Errorf("unmarshal field urgent_words: %v", err)
		}
	}

	if value, ok := values[1].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field banned_words", values[1])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &us.BannedWords); err != nil {
			return fmt.Errorf("unmarshal field banned_words: %v", err)
		}
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field language", values[2])
	} else if value.Valid {
		us.Language = usersettings.Language(value.String)
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field sending_frequency", values[3])
	} else if value.Valid {
		us.SendingFrequency = usersettings.SendingFrequency(value.String)
	}
	if value, ok := values[4].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field last_sending", values[4])
	} else if value.Valid {
		us.LastSending = value.Time
	}
	values = values[5:]
	if len(values) == len(usersettings.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field user_settings", value)
		} else if value.Valid {
			us.user_settings = new(int)
			*us.user_settings = int(value.Int64)
		}
	}
	return nil
}

// QueryUser queries the user edge of the UserSettings.
func (us *UserSettings) QueryUser() *UserQuery {
	return (&UserSettingsClient{config: us.config}).QueryUser(us)
}

// Update returns a builder for updating this UserSettings.
// Note that, you need to call UserSettings.Unwrap() before calling this method, if this UserSettings
// was returned from a transaction, and the transaction was committed or rolled back.
func (us *UserSettings) Update() *UserSettingsUpdateOne {
	return (&UserSettingsClient{config: us.config}).UpdateOne(us)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (us *UserSettings) Unwrap() *UserSettings {
	tx, ok := us.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserSettings is not a transactional entity")
	}
	us.config.driver = tx.drv
	return us
}

// String implements the fmt.Stringer.
func (us *UserSettings) String() string {
	var builder strings.Builder
	builder.WriteString("UserSettings(")
	builder.WriteString(fmt.Sprintf("id=%v", us.ID))
	builder.WriteString(", urgent_words=")
	builder.WriteString(fmt.Sprintf("%v", us.UrgentWords))
	builder.WriteString(", banned_words=")
	builder.WriteString(fmt.Sprintf("%v", us.BannedWords))
	builder.WriteString(", language=")
	builder.WriteString(fmt.Sprintf("%v", us.Language))
	builder.WriteString(", sending_frequency=")
	builder.WriteString(fmt.Sprintf("%v", us.SendingFrequency))
	builder.WriteString(", last_sending=")
	builder.WriteString(us.LastSending.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// UserSettingsSlice is a parsable slice of UserSettings.
type UserSettingsSlice []*UserSettings

func (us UserSettingsSlice) config(cfg config) {
	for _i := range us {
		us[_i].config = cfg
	}
}
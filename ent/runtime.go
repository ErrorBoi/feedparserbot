// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/ErrorBoi/feedparserbot/ent/post"
	"github.com/ErrorBoi/feedparserbot/ent/schema"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	postFields := schema.Post{}.Fields()
	_ = postFields
	// postDescCreatedAt is the schema descriptor for created_at field.
	postDescCreatedAt := postFields[10].Descriptor()
	// post.DefaultCreatedAt holds the default value on creation for the created_at field.
	post.DefaultCreatedAt = postDescCreatedAt.Default.(func() time.Time)
	// postDescUpdatedAt is the schema descriptor for updated_at field.
	postDescUpdatedAt := postFields[11].Descriptor()
	// post.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	post.DefaultUpdatedAt = postDescUpdatedAt.Default.(func() time.Time)
	// post.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	post.UpdateDefaultUpdatedAt = postDescUpdatedAt.UpdateDefault.(func() time.Time)
	sourceFields := schema.Source{}.Fields()
	_ = sourceFields
	userFields := schema.User{}.Fields()
	_ = userFields
	usersettingsFields := schema.UserSettings{}.Fields()
	_ = usersettingsFields
	// usersettingsDescLastSending is the schema descriptor for last_sending field.
	usersettingsDescLastSending := usersettingsFields[4].Descriptor()
	// usersettings.DefaultLastSending holds the default value on creation for the last_sending field.
	usersettings.DefaultLastSending = usersettingsDescLastSending.Default.(func() time.Time)
}

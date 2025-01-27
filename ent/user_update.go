// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/ErrorBoi/feedparserbot/ent/predicate"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks      []Hook
	mutation   *UserMutation
	predicates []predicate.User
}

// Where adds a new predicate for the builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.predicates = append(uu.predicates, ps...)
	return uu
}

// SetEmail sets the email field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmail(s *string) *UserUpdate {
	if s != nil {
		uu.SetEmail(*s)
	}
	return uu
}

// ClearEmail clears the value of email.
func (uu *UserUpdate) ClearEmail() *UserUpdate {
	uu.mutation.ClearEmail()
	return uu
}

// SetTgID sets the tg_id field.
func (uu *UserUpdate) SetTgID(i int) *UserUpdate {
	uu.mutation.ResetTgID()
	uu.mutation.SetTgID(i)
	return uu
}

// AddTgID adds i to tg_id.
func (uu *UserUpdate) AddTgID(i int) *UserUpdate {
	uu.mutation.AddTgID(i)
	return uu
}

// SetPaymentInfo sets the payment_info field.
func (uu *UserUpdate) SetPaymentInfo(s string) *UserUpdate {
	uu.mutation.SetPaymentInfo(s)
	return uu
}

// SetNillablePaymentInfo sets the payment_info field if the given value is not nil.
func (uu *UserUpdate) SetNillablePaymentInfo(s *string) *UserUpdate {
	if s != nil {
		uu.SetPaymentInfo(*s)
	}
	return uu
}

// ClearPaymentInfo clears the value of payment_info.
func (uu *UserUpdate) ClearPaymentInfo() *UserUpdate {
	uu.mutation.ClearPaymentInfo()
	return uu
}

// SetRole sets the role field.
func (uu *UserUpdate) SetRole(u user.Role) *UserUpdate {
	uu.mutation.SetRole(u)
	return uu
}

// SetNillableRole sets the role field if the given value is not nil.
func (uu *UserUpdate) SetNillableRole(u *user.Role) *UserUpdate {
	if u != nil {
		uu.SetRole(*u)
	}
	return uu
}

// SetSettingsID sets the settings edge to UserSettings by id.
func (uu *UserUpdate) SetSettingsID(id int) *UserUpdate {
	uu.mutation.SetSettingsID(id)
	return uu
}

// SetNillableSettingsID sets the settings edge to UserSettings by id if the given value is not nil.
func (uu *UserUpdate) SetNillableSettingsID(id *int) *UserUpdate {
	if id != nil {
		uu = uu.SetSettingsID(*id)
	}
	return uu
}

// SetSettings sets the settings edge to UserSettings.
func (uu *UserUpdate) SetSettings(u *UserSettings) *UserUpdate {
	return uu.SetSettingsID(u.ID)
}

// AddSourceIDs adds the sources edge to Source by ids.
func (uu *UserUpdate) AddSourceIDs(ids ...int) *UserUpdate {
	uu.mutation.AddSourceIDs(ids...)
	return uu
}

// AddSources adds the sources edges to Source.
func (uu *UserUpdate) AddSources(s ...*Source) *UserUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return uu.AddSourceIDs(ids...)
}

// ClearSettings clears the settings edge to UserSettings.
func (uu *UserUpdate) ClearSettings() *UserUpdate {
	uu.mutation.ClearSettings()
	return uu
}

// RemoveSourceIDs removes the sources edge to Source by ids.
func (uu *UserUpdate) RemoveSourceIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveSourceIDs(ids...)
	return uu
}

// RemoveSources removes sources edges to Source.
func (uu *UserUpdate) RemoveSources(s ...*Source) *UserUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return uu.RemoveSourceIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	if v, ok := uu.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}

	var (
		err      error
		affected int
	)
	if len(uu.hooks) == 0 {
		affected, err = uu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uu.mutation = mutation
			affected, err = uu.sqlSave(ctx)
			return affected, err
		})
		for i := len(uu.hooks) - 1; i >= 0; i-- {
			mut = uu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	if ps := uu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if uu.mutation.EmailCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uu.mutation.TgID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: user.FieldTgID,
		})
	}
	if value, ok := uu.mutation.AddedTgID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: user.FieldTgID,
		})
	}
	if value, ok := uu.mutation.PaymentInfo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPaymentInfo,
		})
	}
	if uu.mutation.PaymentInfoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldPaymentInfo,
		})
	}
	if value, ok := uu.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if uu.mutation.SettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SettingsTable,
			Columns: []string{user.SettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersettings.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.SettingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SettingsTable,
			Columns: []string{user.SettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersettings.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uu.mutation.RemovedSourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.SourcesTable,
			Columns: user.SourcesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: source.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.SourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.SourcesTable,
			Columns: user.SourcesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: source.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// SetEmail sets the email field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmail(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetEmail(*s)
	}
	return uuo
}

// ClearEmail clears the value of email.
func (uuo *UserUpdateOne) ClearEmail() *UserUpdateOne {
	uuo.mutation.ClearEmail()
	return uuo
}

// SetTgID sets the tg_id field.
func (uuo *UserUpdateOne) SetTgID(i int) *UserUpdateOne {
	uuo.mutation.ResetTgID()
	uuo.mutation.SetTgID(i)
	return uuo
}

// AddTgID adds i to tg_id.
func (uuo *UserUpdateOne) AddTgID(i int) *UserUpdateOne {
	uuo.mutation.AddTgID(i)
	return uuo
}

// SetPaymentInfo sets the payment_info field.
func (uuo *UserUpdateOne) SetPaymentInfo(s string) *UserUpdateOne {
	uuo.mutation.SetPaymentInfo(s)
	return uuo
}

// SetNillablePaymentInfo sets the payment_info field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillablePaymentInfo(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetPaymentInfo(*s)
	}
	return uuo
}

// ClearPaymentInfo clears the value of payment_info.
func (uuo *UserUpdateOne) ClearPaymentInfo() *UserUpdateOne {
	uuo.mutation.ClearPaymentInfo()
	return uuo
}

// SetRole sets the role field.
func (uuo *UserUpdateOne) SetRole(u user.Role) *UserUpdateOne {
	uuo.mutation.SetRole(u)
	return uuo
}

// SetNillableRole sets the role field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableRole(u *user.Role) *UserUpdateOne {
	if u != nil {
		uuo.SetRole(*u)
	}
	return uuo
}

// SetSettingsID sets the settings edge to UserSettings by id.
func (uuo *UserUpdateOne) SetSettingsID(id int) *UserUpdateOne {
	uuo.mutation.SetSettingsID(id)
	return uuo
}

// SetNillableSettingsID sets the settings edge to UserSettings by id if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableSettingsID(id *int) *UserUpdateOne {
	if id != nil {
		uuo = uuo.SetSettingsID(*id)
	}
	return uuo
}

// SetSettings sets the settings edge to UserSettings.
func (uuo *UserUpdateOne) SetSettings(u *UserSettings) *UserUpdateOne {
	return uuo.SetSettingsID(u.ID)
}

// AddSourceIDs adds the sources edge to Source by ids.
func (uuo *UserUpdateOne) AddSourceIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddSourceIDs(ids...)
	return uuo
}

// AddSources adds the sources edges to Source.
func (uuo *UserUpdateOne) AddSources(s ...*Source) *UserUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return uuo.AddSourceIDs(ids...)
}

// ClearSettings clears the settings edge to UserSettings.
func (uuo *UserUpdateOne) ClearSettings() *UserUpdateOne {
	uuo.mutation.ClearSettings()
	return uuo
}

// RemoveSourceIDs removes the sources edge to Source by ids.
func (uuo *UserUpdateOne) RemoveSourceIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveSourceIDs(ids...)
	return uuo
}

// RemoveSources removes sources edges to Source.
func (uuo *UserUpdateOne) RemoveSources(s ...*Source) *UserUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return uuo.RemoveSourceIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	if v, ok := uuo.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}

	var (
		err  error
		node *User
	)
	if len(uuo.hooks) == 0 {
		node, err = uuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uuo.mutation = mutation
			node, err = uuo.sqlSave(ctx)
			return node, err
		})
		for i := len(uuo.hooks) - 1; i >= 0; i-- {
			mut = uuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	u, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return u
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (u *User, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		},
	}
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing User.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := uuo.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if uuo.mutation.EmailCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uuo.mutation.TgID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: user.FieldTgID,
		})
	}
	if value, ok := uuo.mutation.AddedTgID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: user.FieldTgID,
		})
	}
	if value, ok := uuo.mutation.PaymentInfo(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPaymentInfo,
		})
	}
	if uuo.mutation.PaymentInfoCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldPaymentInfo,
		})
	}
	if value, ok := uuo.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if uuo.mutation.SettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SettingsTable,
			Columns: []string{user.SettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersettings.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.SettingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SettingsTable,
			Columns: []string{user.SettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usersettings.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := uuo.mutation.RemovedSourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.SourcesTable,
			Columns: user.SourcesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: source.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.SourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   user.SourcesTable,
			Columns: user.SourcesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: source.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	u = &User{config: uuo.config}
	_spec.Assign = u.assignValues
	_spec.ScanValues = u.scanValues()
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return u, nil
}

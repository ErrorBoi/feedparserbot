// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
)

// UserCreate is the builder for creating a User entity.
type UserCreate struct {
	config
	mutation *UserMutation
	hooks    []Hook
}

// SetEmail sets the email field.
func (uc *UserCreate) SetEmail(s string) *UserCreate {
	uc.mutation.SetEmail(s)
	return uc
}

// SetNillableEmail sets the email field if the given value is not nil.
func (uc *UserCreate) SetNillableEmail(s *string) *UserCreate {
	if s != nil {
		uc.SetEmail(*s)
	}
	return uc
}

// SetTgID sets the tg_id field.
func (uc *UserCreate) SetTgID(i int) *UserCreate {
	uc.mutation.SetTgID(i)
	return uc
}

// SetPaymentInfo sets the payment_info field.
func (uc *UserCreate) SetPaymentInfo(s string) *UserCreate {
	uc.mutation.SetPaymentInfo(s)
	return uc
}

// SetNillablePaymentInfo sets the payment_info field if the given value is not nil.
func (uc *UserCreate) SetNillablePaymentInfo(s *string) *UserCreate {
	if s != nil {
		uc.SetPaymentInfo(*s)
	}
	return uc
}

// SetRole sets the role field.
func (uc *UserCreate) SetRole(u user.Role) *UserCreate {
	uc.mutation.SetRole(u)
	return uc
}

// SetNillableRole sets the role field if the given value is not nil.
func (uc *UserCreate) SetNillableRole(u *user.Role) *UserCreate {
	if u != nil {
		uc.SetRole(*u)
	}
	return uc
}

// SetSettingsID sets the settings edge to UserSettings by id.
func (uc *UserCreate) SetSettingsID(id int) *UserCreate {
	uc.mutation.SetSettingsID(id)
	return uc
}

// SetNillableSettingsID sets the settings edge to UserSettings by id if the given value is not nil.
func (uc *UserCreate) SetNillableSettingsID(id *int) *UserCreate {
	if id != nil {
		uc = uc.SetSettingsID(*id)
	}
	return uc
}

// SetSettings sets the settings edge to UserSettings.
func (uc *UserCreate) SetSettings(u *UserSettings) *UserCreate {
	return uc.SetSettingsID(u.ID)
}

// AddSourceIDs adds the sources edge to Source by ids.
func (uc *UserCreate) AddSourceIDs(ids ...int) *UserCreate {
	uc.mutation.AddSourceIDs(ids...)
	return uc
}

// AddSources adds the sources edges to Source.
func (uc *UserCreate) AddSources(s ...*Source) *UserCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return uc.AddSourceIDs(ids...)
}

// Save creates the User in the database.
func (uc *UserCreate) Save(ctx context.Context) (*User, error) {
	if _, ok := uc.mutation.TgID(); !ok {
		return nil, errors.New("ent: missing required field \"tg_id\"")
	}
	if _, ok := uc.mutation.Role(); !ok {
		v := user.DefaultRole
		uc.mutation.SetRole(v)
	}
	if v, ok := uc.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"role\": %v", err)
		}
	}
	var (
		err  error
		node *User
	)
	if len(uc.hooks) == 0 {
		node, err = uc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			uc.mutation = mutation
			node, err = uc.sqlSave(ctx)
			return node, err
		})
		for i := len(uc.hooks) - 1; i >= 0; i-- {
			mut = uc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UserCreate) SaveX(ctx context.Context) *User {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (uc *UserCreate) sqlSave(ctx context.Context) (*User, error) {
	var (
		u     = &User{config: uc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: user.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: user.FieldID,
			},
		}
	)
	if value, ok := uc.mutation.Email(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
		u.Email = &value
	}
	if value, ok := uc.mutation.TgID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: user.FieldTgID,
		})
		u.TgID = value
	}
	if value, ok := uc.mutation.PaymentInfo(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPaymentInfo,
		})
		u.PaymentInfo = &value
	}
	if value, ok := uc.mutation.Role(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
		u.Role = value
	}
	if nodes := uc.mutation.SettingsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := uc.mutation.SourcesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	u.ID = int(id)
	return u, nil
}

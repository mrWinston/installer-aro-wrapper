// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// VolumesClones volumes clones
//
// swagger:model VolumesClones
type VolumesClones struct {

	// list of volumes-clone requests
	VolumesClone []*VolumesClone `json:"volumesClone"`
}

// Validate validates this volumes clones
func (m *VolumesClones) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateVolumesClone(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VolumesClones) validateVolumesClone(formats strfmt.Registry) error {
	if swag.IsZero(m.VolumesClone) { // not required
		return nil
	}

	for i := 0; i < len(m.VolumesClone); i++ {
		if swag.IsZero(m.VolumesClone[i]) { // not required
			continue
		}

		if m.VolumesClone[i] != nil {
			if err := m.VolumesClone[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("volumesClone" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("volumesClone" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this volumes clones based on the context it is used
func (m *VolumesClones) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateVolumesClone(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VolumesClones) contextValidateVolumesClone(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.VolumesClone); i++ {

		if m.VolumesClone[i] != nil {

			if swag.IsZero(m.VolumesClone[i]) { // not required
				return nil
			}

			if err := m.VolumesClone[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("volumesClone" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("volumesClone" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *VolumesClones) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VolumesClones) UnmarshalBinary(b []byte) error {
	var res VolumesClones
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MultiHomingConfiguration Multi-homing configuration
//
// swagger:model MultiHomingConfiguration
type MultiHomingConfiguration struct {

	// List of devices used in multi-homing mode
	Devices []string `json:"devices"`
}

// Validate validates this multi homing configuration
func (m *MultiHomingConfiguration) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MultiHomingConfiguration) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MultiHomingConfiguration) UnmarshalBinary(b []byte) error {
	var res MultiHomingConfiguration
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
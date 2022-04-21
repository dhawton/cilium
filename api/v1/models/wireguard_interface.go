// Code generated by go-swagger; DO NOT EDIT.

// Copyright Authors of Cilium
// SPDX-License-Identifier: Apache-2.0

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WireguardInterface Status of a Wireguard interface
//
// +k8s:deepcopy-gen=true
//
// swagger:model WireguardInterface
type WireguardInterface struct {

	// Port on which the Wireguard endpoint is exposed
	ListenPort int64 `json:"listen-port,omitempty"`

	// Name of the interface
	Name string `json:"name,omitempty"`

	// Number of peers configured on this interface
	PeerCount int64 `json:"peer-count,omitempty"`

	// Optional list of wireguard peers
	Peers []*WireguardPeer `json:"peers"`

	// Public key of this interface
	PublicKey string `json:"public-key,omitempty"`
}

// Validate validates this wireguard interface
func (m *WireguardInterface) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePeers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WireguardInterface) validatePeers(formats strfmt.Registry) error {

	if swag.IsZero(m.Peers) { // not required
		return nil
	}

	for i := 0; i < len(m.Peers); i++ {
		if swag.IsZero(m.Peers[i]) { // not required
			continue
		}

		if m.Peers[i] != nil {
			if err := m.Peers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("peers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *WireguardInterface) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WireguardInterface) UnmarshalBinary(b []byte) error {
	var res WireguardInterface
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
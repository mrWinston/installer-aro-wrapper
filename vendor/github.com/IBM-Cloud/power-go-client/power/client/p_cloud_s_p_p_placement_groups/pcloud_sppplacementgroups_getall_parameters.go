// Code generated by go-swagger; DO NOT EDIT.

package p_cloud_s_p_p_placement_groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewPcloudSppplacementgroupsGetallParams creates a new PcloudSppplacementgroupsGetallParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPcloudSppplacementgroupsGetallParams() *PcloudSppplacementgroupsGetallParams {
	return &PcloudSppplacementgroupsGetallParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPcloudSppplacementgroupsGetallParamsWithTimeout creates a new PcloudSppplacementgroupsGetallParams object
// with the ability to set a timeout on a request.
func NewPcloudSppplacementgroupsGetallParamsWithTimeout(timeout time.Duration) *PcloudSppplacementgroupsGetallParams {
	return &PcloudSppplacementgroupsGetallParams{
		timeout: timeout,
	}
}

// NewPcloudSppplacementgroupsGetallParamsWithContext creates a new PcloudSppplacementgroupsGetallParams object
// with the ability to set a context for a request.
func NewPcloudSppplacementgroupsGetallParamsWithContext(ctx context.Context) *PcloudSppplacementgroupsGetallParams {
	return &PcloudSppplacementgroupsGetallParams{
		Context: ctx,
	}
}

// NewPcloudSppplacementgroupsGetallParamsWithHTTPClient creates a new PcloudSppplacementgroupsGetallParams object
// with the ability to set a custom HTTPClient for a request.
func NewPcloudSppplacementgroupsGetallParamsWithHTTPClient(client *http.Client) *PcloudSppplacementgroupsGetallParams {
	return &PcloudSppplacementgroupsGetallParams{
		HTTPClient: client,
	}
}

/*
PcloudSppplacementgroupsGetallParams contains all the parameters to send to the API endpoint

	for the pcloud sppplacementgroups getall operation.

	Typically these are written to a http.Request.
*/
type PcloudSppplacementgroupsGetallParams struct {

	/* CloudInstanceID.

	   Cloud Instance ID of a PCloud Instance
	*/
	CloudInstanceID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the pcloud sppplacementgroups getall params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PcloudSppplacementgroupsGetallParams) WithDefaults() *PcloudSppplacementgroupsGetallParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the pcloud sppplacementgroups getall params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PcloudSppplacementgroupsGetallParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) WithTimeout(timeout time.Duration) *PcloudSppplacementgroupsGetallParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) WithContext(ctx context.Context) *PcloudSppplacementgroupsGetallParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) WithHTTPClient(client *http.Client) *PcloudSppplacementgroupsGetallParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCloudInstanceID adds the cloudInstanceID to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) WithCloudInstanceID(cloudInstanceID string) *PcloudSppplacementgroupsGetallParams {
	o.SetCloudInstanceID(cloudInstanceID)
	return o
}

// SetCloudInstanceID adds the cloudInstanceId to the pcloud sppplacementgroups getall params
func (o *PcloudSppplacementgroupsGetallParams) SetCloudInstanceID(cloudInstanceID string) {
	o.CloudInstanceID = cloudInstanceID
}

// WriteToRequest writes these params to a swagger request
func (o *PcloudSppplacementgroupsGetallParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cloud_instance_id
	if err := r.SetPathParam("cloud_instance_id", o.CloudInstanceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

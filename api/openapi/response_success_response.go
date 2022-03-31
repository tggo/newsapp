/*
 * Boosters news server API
 *
 * Implement API test server
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type SuccessResponse struct {

	// Request status.
	Status string `json:"status"`

	// information description
	Message string `json:"message,omitempty"`

	// any data response, if need
	Data map[string]interface{} `json:"data,omitempty"`
}

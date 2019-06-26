package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// JSONResponse is the Accept header value to generate a JSON response
	JSONResponse string = "json"

	// XMLResponse is the Accept header value to generate a XML response
	XMLResponse string = "xml"

	// YAMLResponse is the Accept header value to geneate a YAML response
	YAMLResponse string = "yaml"
)

// IFormatResponse defines the methods to verifies response format that
// should be generated
type IFormatResponse interface {
	GetNext() IFormatResponse
	GetType() string
	Render()
}

// FormatResponse is the base struct to create the chain of responsabilities
// according to the format to render
type FormatResponse struct {
	code int
	next IFormatResponse

	c   *gin.Context
	obj interface{}
}

// GetNext return the next element in the chain
func (fr *FormatResponse) GetNext() IFormatResponse {
	return fr.next
}

// GetType determines the type of format to reponse
func (fr *FormatResponse) GetType() string {
	typ := fr.c.GetHeader("Accept")
	return typ
}

// JSON returns a json response
type JSON struct {
	*FormatResponse
}

// Render generates a reponse in json format
func (j *JSON) Render() {

	if j.GetType() == JSONResponse {
		j.c.JSON(j.code, j.obj)
		return
	}

	if j.next != nil {
		j.next.Render()
	}

}

// XML returns a xml response
type XML struct {
	*FormatResponse
}

// Render generates a reponse in XML format
func (x *XML) Render() {

	if x.GetType() == XMLResponse {
		x.c.XML(x.code, x.obj)
		return
	}

	if x.next != nil {
		x.next.Render()
	}
}

// YAML return a yaml response
type YAML struct {
	*FormatResponse
}

// Render generates a reponse in YAML format
func (y *YAML) Render() {

	if y.GetType() == YAMLResponse {
		y.c.YAML(y.code, y.obj)
		return
	}

	if y.next != nil {
		y.next.Render()
	}
}

// NotSupportFormat ...
type NotSupportFormat struct {
	*FormatResponse
}

// Render generates an error because is not exist or is not specified the
// Accept heder
func (n *NotSupportFormat) Render() {
	n.c.String(http.StatusBadRequest, "Accept header is not supported")
}

// NewBuilder returns the chain builder
func NewBuilder(c *gin.Context, code int, obj interface{}) IFormatResponse {
	notSupport := &NotSupportFormat{&FormatResponse{c: c, code: code, obj: obj}}
	json := &JSON{&FormatResponse{c: c, code: code, next: notSupport, obj: obj}}
	xml := &XML{&FormatResponse{c: c, next: json, code: code, obj: obj}}
	yml := &YAML{&FormatResponse{c: c, next: xml, code: code, obj: obj}}

	return yml
}

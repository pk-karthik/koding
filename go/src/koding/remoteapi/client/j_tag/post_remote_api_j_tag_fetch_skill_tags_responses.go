package j_tag

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"koding/remoteapi/models"
)

// PostRemoteAPIJTagFetchSkillTagsReader is a Reader for the PostRemoteAPIJTagFetchSkillTags structure.
type PostRemoteAPIJTagFetchSkillTagsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostRemoteAPIJTagFetchSkillTagsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPostRemoteAPIJTagFetchSkillTagsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewPostRemoteAPIJTagFetchSkillTagsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostRemoteAPIJTagFetchSkillTagsOK creates a PostRemoteAPIJTagFetchSkillTagsOK with default headers values
func NewPostRemoteAPIJTagFetchSkillTagsOK() *PostRemoteAPIJTagFetchSkillTagsOK {
	return &PostRemoteAPIJTagFetchSkillTagsOK{}
}

/*PostRemoteAPIJTagFetchSkillTagsOK handles this case with default header values.

Request processed successfully
*/
type PostRemoteAPIJTagFetchSkillTagsOK struct {
	Payload *models.DefaultResponse
}

func (o *PostRemoteAPIJTagFetchSkillTagsOK) Error() string {
	return fmt.Sprintf("[POST /remote.api/JTag.fetchSkillTags][%d] postRemoteApiJTagFetchSkillTagsOK  %+v", 200, o.Payload)
}

func (o *PostRemoteAPIJTagFetchSkillTagsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DefaultResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostRemoteAPIJTagFetchSkillTagsUnauthorized creates a PostRemoteAPIJTagFetchSkillTagsUnauthorized with default headers values
func NewPostRemoteAPIJTagFetchSkillTagsUnauthorized() *PostRemoteAPIJTagFetchSkillTagsUnauthorized {
	return &PostRemoteAPIJTagFetchSkillTagsUnauthorized{}
}

/*PostRemoteAPIJTagFetchSkillTagsUnauthorized handles this case with default header values.

Unauthorized request
*/
type PostRemoteAPIJTagFetchSkillTagsUnauthorized struct {
	Payload *models.UnauthorizedRequest
}

func (o *PostRemoteAPIJTagFetchSkillTagsUnauthorized) Error() string {
	return fmt.Sprintf("[POST /remote.api/JTag.fetchSkillTags][%d] postRemoteApiJTagFetchSkillTagsUnauthorized  %+v", 401, o.Payload)
}

func (o *PostRemoteAPIJTagFetchSkillTagsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.UnauthorizedRequest)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

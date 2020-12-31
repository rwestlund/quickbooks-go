package quickbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
)

type ContentType string

const (
	AI   ContentType = "application/postscript"
	CSV  ContentType = "text/csv"
	DOC  ContentType = "application/msword"
	DOCX ContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	EPS  ContentType = "application/postscript"
	GIF  ContentType = "image/gif"
	JPEG ContentType = "image/jpeg"
	JPG  ContentType = "image/jpg"
	ODS  ContentType = "application/vnd.oasis.opendocument.spreadsheet"
	PDF  ContentType = "application/pdf"
	PNG  ContentType = "image/png"
	RTF  ContentType = "text/rtf"
	TIF  ContentType = "image/tif"
	TXT  ContentType = "text/plain"
	XLS  ContentType = "application/vnd/ms-excel"
	XLSX ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	XML  ContentType = "text/xml"
)

type Attachable struct {
	ID                       string          `json:"Id,omitempty"`
	SyncToken                string          `json:",omitempty"`
	FileName                 string          `json:",omitempty"`
	Note                     string          `json:",omitempty"`
	Category                 string          `json:",omitempty"`
	ContentType              ContentType     `json:",omitempty"`
	PlaceName                string          `json:",omitempty"`
	AttachableRef            []AttachableRef `json:",omitempty"`
	Long                     string          `json:",omitempty"`
	Tag                      string          `json:",omitempty"`
	Lat                      string          `json:",omitempty"`
	MetaData                 MetaData        `json:",omitempty"`
	FileAccessUri            string          `json:",omitempty"`
	Size                     json.Number     `json:",omitempty"`
	ThumbnailFileAccessUri   string          `json:",omitempty"`
	TempDownloadUri          string          `json:",omitempty"`
	ThumbnailTempDownloadUri string          `json:",omitempty"`
}

type AttachableRef struct {
	IncludeOnSend bool   `json:",omitempty"`
	LineInfo      string `json:",omitempty"`
	NoRefOnly     bool   `json:",omitempty"`
	// CustomField[0..n]
	Inactive  bool          `json:",omitempty"`
	EntityRef ReferenceType `json:",omitempty"`
}

// CreateAttachable creates the attachable
func (c *Client) CreateAttachable(attachable *Attachable) (*Attachable, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/attachable"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(attachable)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Attachable Attachable
		Time       Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Attachable, err
}

// DeleteAttachable deletes the attachable
func (c *Client) DeleteAttachable(attachable *Attachable) error {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return err
	}
	u.Path = "/v3/company/" + c.RealmID + "/attachable"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	v.Add("operation", "delete")
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(attachable)
	if err != nil {
		return err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return parseFailure(res)
	}

	return nil
}

// DownloadAttachable downloads the attachable
func (c *Client) DownloadAttachable(attachableId string) (string, error) {

	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return "", err
	}
	u.Path = "/v3/company/" + c.RealmID + "/download/" + attachableId
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", parseFailure(res)
	}
	url, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(url), err
}

// GetAttachables gets the attachables
func (c *Client) GetAttachables(startpos int) ([]Attachable, error) {

	var r struct {
		QueryResponse struct {
			Attachable    []Attachable
			StartPosition int
			MaxResults    int
		}
	}
	q := "SELECT * FROM Attachable ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(queryPageSize)
	err := c.query(q, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Attachable == nil {
		r.QueryResponse.Attachable = make([]Attachable, 0)
	}
	return r.QueryResponse.Attachable, nil
}

// GetAttachable gets the attachable
func (c *Client) GetAttachable(attachableId string) (*Attachable, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/attachable/" + attachableId
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Attachable Attachable
		Time       Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Attachable, err
}

// UpdateAttachable updates the attachable
func (c *Client) UpdateAttachable(attachable *Attachable) (*Attachable, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/attachable"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var d = struct {
		*Attachable
		Sparse bool `json:"sparse"`
	}{
		Attachable: attachable,
		Sparse:     true,
	}
	var j []byte
	j, err = json.Marshal(d)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Attachable Attachable
		Time       Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Attachable, err
}

// UploadAttachable uploads the attachable
func (c *Client) UploadAttachable(attachable *Attachable, data io.Reader) (*Attachable, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/upload"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()

	var buffer bytes.Buffer
	mWriter := multipart.NewWriter(&buffer)

	// Add file metadata
	metadataHeader := make(textproto.MIMEHeader)
	metadataHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file_metadata_01", "attachment.json"))
	metadataHeader.Set("Content-Type", "application/json")
	var metadataContent io.Writer
	if metadataContent, err = mWriter.CreatePart(metadataHeader); err != nil {
		return nil, err
	}
	var j []byte
	j, err = json.Marshal(attachable)
	if err != nil {
		return nil, err
	}
	if _, err = metadataContent.Write(j); err != nil {
		return nil, err
	}

	// Add file content
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file_content_01", attachable.FileName))
	fileHeader.Set("Content-Type", string(attachable.ContentType))
	var fileContent io.Writer
	if fileContent, err = mWriter.CreatePart(fileHeader); err != nil {
		return nil, err
	}
	if _, err = io.Copy(fileContent, data); err != nil {
		return nil, err
	}

	mWriter.Close()

	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), &buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mWriter.FormDataContentType())
	req.Header.Add("Accept", "application/json")

	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		AttachableResponse []struct {
			Attachable Attachable
		}
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.AttachableResponse[0].Attachable, err
}

packahttpge pagetoken


import (
   "encoding/base64"
   "encoding/json"
)


// PageToken contains an information to paginate a list request
type PageToken struct {
   Offset   int `json:"offset"`
   PageSize int `json:"page_size"`
   // TODO add also [Filter[, Sort ?]]   map[string][]string `json:"filter"`
}


func NewPageToken(pageSize int) *PageToken {
   return &PageToken{
       PageSize: pageSize,
   }
}


// ParseToken decodes a page_token
func ParsePageToken(encodedToken string) (*PageToken, error) {
   if encodedToken == "" {
       return &PageToken{}, nil
   }


   var pageTokenBytes []byte
   var err error
   if pageTokenBytes, err = base64.RawURLEncoding.DecodeString(encodedToken); err != nil {
       return nil, err
   }


   var token PageToken
   if err = json.Unmarshal(pageTokenBytes, &token); err != nil {
       return nil, err
   }


   return &token, nil
}


// Encode creates a Base64 encoded string to use as next_page_token
func (pt *PageToken) Encode() (string, error) {
   var pageTokenBytes []byte
   var err error
   if pageTokenBytes, err = json.Marshal(pt); err != nil {
       return "", err
   }


   return base64.RawURLEncoding.EncodeToString(pageTokenBytes), nil
}


// GetNextPageTokenValue creates next_page_token and returns an encoded view of it
func (pt *PageToken) GetNextPageToken(itemsGotInThisPage int) *PageToken {
   if pt.PageSize == 0 || itemsGotInThisPage < pt.PageSize {
       return nil
   }


   return &PageToken{
       PageSize: pt.PageSize,
       Offset:   pt.Offset + itemsGotInThisPage,
   }
}

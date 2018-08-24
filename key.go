package rundeck

import "encoding/json"

// ListKeysResponse ...
type ListKeysResponse struct {
	Resources []*KeyResource `json:"resources"`
	sharedKeyMeta
}

// KeyResource ...
type KeyResource struct {
	Meta KeyMetadata `json:"meta"`
	Name string      `json:"name"`
	sharedKeyMeta
}

// KeyMetadata ...
type KeyMetadata struct {
	KeyType     string `json:"Rundeck-key-type"`
	ContentMask string `json:"Rundeck-content-mask"`
	ContentSize int64  `json:"Rundeck-content-size"`
	ContentType string `json:"Rundeck-content-type"`
}

type sharedKeyMeta struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Path string `json:"path"`
}

// KeyStore interacts with the storage facility with regards to keys
type KeyStore struct {
	c *Client
}

// KeyStorage returns a keystore for use
func (c *Client) KeyStorage() *KeyStore {
	return &KeyStore{c: c}
}

// List lists resources at the specified path
func (k *KeyStore) List(path string) (*ListKeysResponse, error) {
	rawURL := k.c.RundeckAddr + "/storage/keys/" + path + "/"

	res, err := k.c.checkResponseOK(k.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var keys ListKeysResponse
	return &keys, json.NewDecoder(res.Body).Decode(&keys)
}

// KeyMetadata returns the metadata about the stored key file
func (k *KeyStore) KeyMetadata(path, file string) (*KeyMetadata, error) {
	rawURL := k.c.RundeckAddr + "/storage/keys/" + path + "/" + file

	res, err := k.c.checkResponseOK(k.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var metadata KeyMetadata
	return &metadata, json.NewDecoder(res.Body).Decode(&metadata)
}

// Delete deletes the file if it exists
func (k *KeyStore) Delete(path, file string) error {
	rawURL := k.c.RundeckAddr + "/storage/keys/" + path + "/" + file

	_, err := k.c.checkResponseNoContent(k.c.delete(rawURL, nil))
	if err != nil {
		return err
	}
	return nil
}

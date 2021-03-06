package ams

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"
)

const (
	assetsEndpoint = "Assets"
)

const (
	OptionStorageEncrypted = 1 << iota
	OptionCommonEncryptionProtected
	OptionEnvelopeEncryptionProtected
	OptionNone = 0
)

const (
	StateInitialized = iota
	StatePublished   // The 'Publish' action has been deprecated. Remove the code that checks whether an asset is in the 'Published' state.
	StateDeleted
)

const (
	FormatOptionNoFormat          = 0
	FormatOptionAdaptiveStreaming = 1
)

type Asset struct {
	ID                 string `json:"Id"`
	State              int    `json:"State"`
	Created            string `json:"Created"`
	LastModified       string `json:"LastModified"`
	Name               string `json:"Name"`
	Options            int    `json:"Options"`
	FormatOption       int    `json:"FormatOption"`
	URI                string `json:"Uri"`
	StorageAccountName string `json:"StorageAccountName"`
}

func (c *Client) GetAsset(ctx context.Context, assetID string) (*Asset, error) {
	c.logger.Printf("[INFO] get asset #%s ...", assetID)

	endpoint := toAssetResource(assetID)
	var out Asset
	if err := c.get(ctx, endpoint, &out); err != nil {
		return nil, err
	}

	c.logger.Printf("[INFO] completed")
	return &out, nil
}

func (c *Client) GetAssets(ctx context.Context) ([]Asset, error) {
	c.logger.Printf("[INFO] get assets ...")

	var out struct {
		Assets []Asset `json:"value"`
	}
	if err := c.get(ctx, assetsEndpoint, &out); err != nil {
		return nil, err
	}

	c.logger.Printf("[INFO] completed")
	return out.Assets, nil
}

func (c *Client) CreateAsset(ctx context.Context, name string) (*Asset, error) {
	c.logger.Printf("[INFO] create asset [name=%#v] ...", name)

	params := map[string]interface{}{
		"Name": name,
	}
	var out Asset
	if err := c.post(ctx, assetsEndpoint, params, &out); err != nil {
		return nil, err
	}

	c.logger.Printf("[INFO] completed, new asset[#%s]", out.ID)
	return &out, nil
}

func (c *Client) GetAssetFiles(ctx context.Context, assetID string) ([]AssetFile, error) {
	c.logger.Printf("[INFO] get asset[#%s] files ...", assetID)

	endpoint := path.Join(toAssetResource(assetID), filesEndpoint)
	var out struct {
		AssetFiles []AssetFile `json:"value"`
	}
	if err := c.get(ctx, endpoint, &out); err != nil {
		return nil, err
	}

	c.logger.Printf("[INFO] completed")
	return out.AssetFiles, nil
}

func (c *Client) DeleteAsset(ctx context.Context, assetID string) error {
	endpoint := toAssetResource(assetID)
	req, err := c.newRequest(ctx, http.MethodDelete, endpoint)
	if err != nil {
		return errors.Wrap(err, "failed to construct request")
	}

	c.logger.Printf("[INFO] delete asset[#%s] ...")
	if err := c.do(req, http.StatusNoContent, nil); err != nil {
		return errors.Wrap(err, "request failed")
	}
	c.logger.Printf("[INFO] completed")
	return nil
}

func toAssetResource(assetID string) string {
	return toResource(assetsEndpoint, assetID)
}

func (c *Client) buildAssetURI(assetID string) string {
	return c.buildURI(toAssetResource(assetID))
}

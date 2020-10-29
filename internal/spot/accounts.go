package spot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
)

type accounts struct {
	client *client.Client
}

func (x *accounts) ListAccounts(ctx context.Context) ([]*Account, error) {
	log.V(1).Info("Listing all accounts")

	req := client.NewRequest(http.MethodGet, "/setup/account")
	resp, err := client.RequireOK(x.client.Do(ctx, req))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	accounts, err := accountsFromHttpResponse(resp)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func accountFromJSON(in []byte) (*Account, error) {
	c := new(Account)
	if err := json.Unmarshal(in, c); err != nil {
		return nil, err
	}
	return c, nil
}

func accountsFromJSON(in []byte) ([]*Account, error) {
	var rw client.Response
	if err := json.Unmarshal(in, &rw); err != nil {
		return nil, err
	}
	out := make([]*Account, len(rw.Response.Items))
	if len(out) == 0 {
		return out, nil
	}
	for i, rb := range rw.Response.Items {
		b, err := accountFromJSON(rb)
		if err != nil {
			return nil, err
		}
		out[i] = b
	}
	return out, nil
}

func accountsFromHttpResponse(resp *http.Response) ([]*Account, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return accountsFromJSON(body)
}

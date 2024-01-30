package restcli

import "github.com/go-resty/resty/v2"

type Params struct {
	Header  map[string]string
	Queries map[string]string
	Body    any

	Formdata map[string]string
}

type Options struct {
	BaseUrl string
	//TODO: add more options
}

type Client struct {
	client *resty.Client
}

func NewClient(opt ...Options) *Client {
	client := resty.New()
	if len(opt) > 0 {
		if opt[0].BaseUrl != "" {
			client.SetBaseURL(opt[0].BaseUrl)
		}
	}

	return &Client{
		client,
	}
}

func (thiz *Client) Get(url string, params *Params, result ...any) (*resty.Response, error) {
	req := genReq(thiz.client, params, result)
	return req.Get(url)
}

func (thiz *Client) Post(url string, params *Params, result ...any) (*resty.Response, error) {
	req := genReq(thiz.client, params, result)
	return req.Post(url)
}

func (thiz *Client) Put(url string, params *Params, result ...any) (*resty.Response, error) {
	req := genReq(thiz.client, params, result)
	return req.Put(url)
}

func (thiz *Client) Patch(url string, params *Params, result ...any) (*resty.Response, error) {
	req := genReq(thiz.client, params, result)
	return req.Patch(url)
}

func (thiz *Client) Delete(url string, params *Params, result ...any) (*resty.Response, error) {
	req := genReq(thiz.client, params, result)
	return req.Delete(url)
}

func genReq(cli *resty.Client, params *Params, result []any) *resty.Request {
	req := cli.R()

	if len(result) != 0 {
		req = req.SetResult(result[0])
	}

	if params == nil {
		return req
	}

	if params.Header != nil {
		req = req.SetHeaders(params.Header)
		if params.Header["Accept"] == "" {
			req = req.SetHeader("Accept", "application/json")
		}
	}

	if params.Queries != nil {
		req = req.SetQueryParams(params.Queries)
	}

	if params.Body != nil {
		req = req.SetBody(params.Body)
		if params.Header["Content-Type"] == "" {
			req = req.SetHeader("Content-Type", "application/json")
		}
	}

	if params.Formdata != nil {
		req = req.SetFormData(params.Formdata)
		if params.Header["Content-Type"] == "" {
			req = req.SetHeader("Content-Type", "multipart/form-data")
		}
	}
	return req
}

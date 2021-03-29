// http post eof error problem..

package ethclient

func (c *CallClient) Wallet(method string, req, rep interface{}) error {
	client := &http.Client{}
	data, _ :=  EncodeClientRequest(method, req)
	reqest, _ := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(data))
	resp, err := client.Do(reqest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
	return DecodeClientResponse(resp.Body, rep)
}

//with EncodeClientRquest && DecodeClientResponse

// EncodeClientRequest encodes parameters for a JSON-RPC client request.

func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Version: "2.0",
		Method: method,
		Params: [1]interface{}{args},
		Id:     uint64(rand.Int63()),
	}

	return json.Marshal(c)
}
// DecodeClientResponse decodes the response body of a client request into // the interface reply.

func DecodeClientResponse(r io.Reader, reply interface{}) error {
	var c clientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	if c.Error != nil {
		return fmt.Errorf("%v", c.Error)
	}
	if c.Result == nil {
		return errors.New("result is null")
	}
	return json.Unmarshal(*c.Result, reply)
}
//And I got error EOF.


//solve EOF version. this style..

package ethclient
//This line:
//
//io.Copy(ioutil.Discard, resp.Body)
//reads the whole resp.Body, leaving the reader with no more bytes to be read. Therefore any successive calls to resp.Body.Read will return EOF and the json.Decoder.Decode method does use the io.Reader.Read method when decoding the given reader's content, so...
//
//And since resp.Body is an io.ReadCloser, which is an interface that does not support "rewinding", and you want to read the body content more than once (ioutil.Discard and json.Decode), you'll have to read the body into a variable that you can re-read afterwards. It's up to you how you do that, a slice of bytes, or bytes.Reader, or something else.
//
//Example using bytes.Reader:
func (c *CallClient) Wallet(method string, req, rep interface{}) error {
	client := &http.Client{}
	data, err := EncodeClientRequest(method, req) //
	if err != nil {
		return err
	}
	reqest, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(data)) //
	if err != nil {
		return err
	}
	resp, err := client.Do(reqest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()  //여기까지 같고.

	// get a reader that can be "rewound"
	// begin - 차이..
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, resp.Body); err != nil {  //io.Copy(ioutil.Discard, resp.Body) -> this
		return err
	}
	br := bytes.NewReader(buf.Bytes())  // <--  추가

	if _, err := io.Copy(ioutil.Discard, br); err != nil {  // 추가
		return err
	}

	// rewind
	if _, err := br.Seek(0, 0); err != nil { //추가
		return err
	}
	return DecodeClientResponse(br, rep)
	// end
}

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



package cachestorage

type chache struct {
	storageUrl map[string]string
}

func New() *chache {
	return &chache{storageUrl: make(map[string]string)}
}

func (c *chache) Get(inputURL string) string {
	if foundurl, ok := c.storageUrl[inputURL]; ok {
		return foundurl
	}
	return ""
}

func (c *chache) Save(url1, url2 string) error {
	c.storageUrl[url1] = url2
	c.storageUrl[url2] = url1
	return nil
}

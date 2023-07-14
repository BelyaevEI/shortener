package cachestorage

type chache struct {
	storageURL map[string]string
}

func New() *chache {
	return &chache{storageURL: make(map[string]string)}
}

func (c *chache) Get(inputURL string) string {
	if foundurl, ok := c.storageURL[inputURL]; ok {
		return foundurl
	}
	return ""
}

func (c *chache) Save(url1, url2 string) error {
	c.storageURL[url1] = url2
	c.storageURL[url2] = url1
	return nil
}

func (c *chache) Ping() error {
	panic("No implemention")
}

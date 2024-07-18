package client

type CLI struct {
	Args []string
}

type Client struct {
	*CLI
}

func NewClient(args []string) *Client {
	return &Client{&CLI{Args: args}}
}

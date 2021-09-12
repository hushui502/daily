## Function Options
type dialOptions struct {
	CreateTime int
	ModifyTime int
	// when you want to add or change some parameter, you can just only add/delete/change there
}

type DialOptionFunc func(*dialOptions)

func WithCreateTime(time int) DialOptionFunc {
	return func(o *dialOptions) {
		o.CreateTime = time
	}
}

func dial(network, address string, options ...DialOptionFunc) (Conn, error) {
	do := dialOptions {
		dial: net.Dial,
	}

	for _, option := range options {
		option(&do)
	}
}

## Config
type Config struct {
	*pool.Config
	Addr string
	Auth string
	DialTimeout time.Duration
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

func NewConn(c *Config) (cn Conn, err error)


"json/yaml 配置加载走Config"
“不依赖配置走DialOption”

https://github.com/gopherchina/conference/blob/master/2020/1.6%20Functional%20options%20and%20config%20for%20APIs.pdf

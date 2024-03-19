package model

type Proxy struct {
	*CommonFunc
	Instance string
}

func NewProxy() *Proxy {
	return &Proxy{CommonFunc: NewCommonFunc()}
}

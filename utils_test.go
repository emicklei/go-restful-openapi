package restfulspec

import restful "github.com/emicklei/go-restful"

func dummy(i *restful.Request, o *restful.Response) {}

type Sample struct {
	ID    string `swagger:"required"`
	Root  Item   `json:"root" description:"root desc"`
	Items []Item
}

type Item struct {
	ItemName string `json:"name"`
}

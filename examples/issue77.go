package main

import (
	"log"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

type ItemTestResource struct {
	// normally one would use DAO (data access object)
	users map[string]ItemTest
}

func (u ItemTestResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	tags := []string{"users"}

	ws.Route(ws.GET("/").To(u.findAllItemTests).
		// docs
		Doc("get all users").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]ItemTest{}).
		Returns(200, "OK", []ItemTest{}))

	ws.Route(ws.GET("/{user-id}").To(u.findItemTest).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(ItemTest{}). // on the response
		Returns(200, "OK", ItemTest{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(u.updateItemTest).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(ItemTest{})) // from the request

	ws.Route(ws.PUT("").To(u.createItemTest).
		// docs
		Doc("create a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(ItemTest{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeItemTest).
		// docs
		Doc("delete a user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}

// GET http://localhost:8080/users
//
func (u ItemTestResource) findAllItemTests(request *restful.Request, response *restful.Response) {
	list := []ItemTest{}
	for _, each := range u.users {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
//
func (u ItemTestResource) findItemTest(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.ID) == 0 {
		response.WriteErrorString(http.StatusNotFound, "ItemTest could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <ItemTest><Id>1</Id><Name>Melissa Raspberry</Name></ItemTest>
//
func (u *ItemTestResource) updateItemTest(request *restful.Request, response *restful.Response) {
	usr := new(ItemTest)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/1
// <ItemTest><Id>1</Id><Name>Melissa</Name></ItemTest>
//
func (u *ItemTestResource) createItemTest(request *restful.Request, response *restful.Response) {
	usr := ItemTest{ID: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.users[usr.ID] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *ItemTestResource) removeItemTest(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	u := ItemTestResource{map[string]ItemTest{}}
	restful.DefaultContainer.Add(u.WebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/ItemTests/emicklei/Projects/swagger-ui/dist"))))

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	log.Printf("Get the API using http://localhost:8080/apidocs.json")
	log.Printf("Open Swagger UI using http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "ItemTestService",
			Description: "Resource for managing ItemTests",
			Version:     "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        "users",
		Description: "Managing users"}}}
}

type ItemTest struct {
	ID               string
	Type             *int64                `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	NumVal           *string               `protobuf:"bytes,2,opt,name=numVal" json:"numVal,omitempty"`
	BoolVal          *bool                 `protobuf:"varint,3,opt,name=boolVal" json:"boolVal,omitempty"`
	StrVal           *string               `protobuf:"bytes,4,opt,name=strVal" json:"strVal,omitempty"`
	MapVal           map[string]*ItemValue `protobuf:"bytes,5,rep,name=mapVal" json:"mapVal,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ListVal          []*ItemValue          `protobuf:"bytes,6,rep,name=listVal" json:"listVal,omitempty"`
	Cells            [][]string
	XXX_unrecognized []byte `json:"-"`
}

type ItemValue struct {
	Type             *int64            `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	NumVal           *string           `protobuf:"bytes,2,opt,name=numVal" json:"numVal,omitempty"`
	BoolVal          *bool             `protobuf:"varint,3,opt,name=boolVal" json:"boolVal,omitempty"`
	StrVal           *string           `protobuf:"bytes,4,opt,name=strVal" json:"strVal,omitempty"`
	MapVal           map[string]string `protobuf:"bytes,5,rep,name=mapVal" json:"mapVal,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ListVal          [][]byte          `protobuf:"bytes,6,rep,name=listVal" json:"listVal,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

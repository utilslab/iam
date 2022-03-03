package iam

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
)

const (
	Post   = "POST"
	Get    = "GET"
	Delete = "DELETE"
	Put    = "PUT"
)

const (
	Read  ActionType = "read"
	Write ActionType = "write"
	List  ActionType = "list"
)

type ActionType string
type Html string
type Text string

type Register interface {
	DELETE(path string, handles ...gin.HandlerFunc) gin.IRoutes
	GET(path string, handles ...gin.HandlerFunc) gin.IRoutes
	HEAD(path string, handles ...gin.HandlerFunc) gin.IRoutes
	OPTIONS(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PATCH(path string, handles ...gin.HandlerFunc) gin.IRoutes
	POST(path string, handles ...gin.HandlerFunc) gin.IRoutes
	PUT(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Any(path string, handles ...gin.HandlerFunc) gin.IRoutes
	Group(prefix string, middleware ...gin.HandlerFunc) *gin.RouterGroup
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
}

type Route struct {
	Prefix      string
	Middlewares []gin.HandlerFunc `json:"-"`
	Groups      []*Group
}

type Router interface {
	Routes() []*Route
}

type Module struct {
}

type Driver interface {
	Register(api *API) error
	Start(addr string) error
}

type HandlerInfo struct {
	Name     string
	Location string
}

func (p HandlerInfo) ParsePath() string {
	return fmt.Sprintf("/%s", p.Name)
}

type RouteTable struct {
	rows []RouteRow
}

func (p *RouteTable) AddRow(method, path string, info HandlerInfo) {
	p.rows = append(p.rows, RouteRow{
		Method:   method,
		Path:     path,
		Source:   info.Name,
		Location: info.Location,
	})
}

func (p RouteTable) Print() {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"", "接口地址", "方法名", "位置"})
	var data [][]string
	for _, v := range p.rows {
		data = append(data, []string{v.Method, v.Path, v.Source, v.Location})
	}
	t.AppendBulk(data)
	t.Render()
}

type RouteRow struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Source   string `json:"source"`
	Location string `json:"location"`
}

type Value struct {
	value   interface{}
	headers map[string]string
}

type ContextWrapper func(ctx *gin.Context) (context.Context, error)

type Resource struct {
	Name        string
	Ident       string
	Description string
	optional    bool
}

func (r Resource) Optional() Resource {
	n := r
	n.optional = true
	return n
}

type Group struct {
	Name        string
	Prefix      string
	Middlewares []gin.HandlerFunc
	Actions     []*Action
}

type Action struct {
	Type        ActionType
	Description string
	Resources   []Resource
	Codes       []Code
	Handler     interface{} `json:"-"`
	handler     reflect.Value
	group       string
	method      string
	path        string
}

type Code struct {
	Status  int
	Code    string
	Message string
}

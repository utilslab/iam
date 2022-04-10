package iam

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/utilslab/iam/binding"
	"github.com/utilslab/iam/exporter"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func New() *API {
	return &API{
		routeTable: &RouteTable{},
	}
}

type API struct {
	version        string
	routers        []Router
	engine         *gin.Engine
	routeTable     *RouteTable
	exporter       *exporter.Exporter
	methods        []*exporter.Method
	basics         *exporter.BasicTypes
	models         *exporter.Fields
	errorWrapper   ErrorWrapper
	contextWrapper ContextWrapper
}

func (p *API) SetVersion(version string) {
	p.version = version
}

func (p *API) AddRouter(router ...Router) {
	p.routers = append(p.routers, router...)
}

func (p *API) SetEngine(engine *gin.Engine) {
	p.engine = engine
}

func (p *API) SetContextWrapper(contextWrapper ContextWrapper) {
	p.contextWrapper = contextWrapper
}

func (p *API) SetErrorWrapper(errorWrapper ErrorWrapper) {
	p.errorWrapper = errorWrapper
}

func (p *API) SetExporter(addr string, options *exporter.Options) {
	basicTypes := []exporter.BasicType{
		{
			Elem: decimal.Decimal{},
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
		{
			Elem: Html(""),
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
		{
			Elem: Text(""),
			Mapping: map[string]exporter.Library{
				"ts": {Type: "string"},
			},
		},
	}
	if options == nil {
		options = new(exporter.Options)
	}
	options.BasicTypes = append(basicTypes, options.BasicTypes...)
	p.exporter = exporter.NewExporter(addr, options)
}

func (p *API) Run(addr string) {
	if p.engine == nil {
		p.engine = gin.Default()
	}
	var (
		routes []*Route
		err    error
	)
	for _, router := range p.routers {
		routes = router.Routes()
		err = p.prepareRoutes(routes)
		if err != nil {
			panic(err)
		}
		err = p.registerRoutes(p.engine, routes)
		if err != nil {
			panic(err)
			return
		}
	}
	if p.exporter != nil {
		p.exporter.Init(p.version, p.methods, p.models)
		p.exporter.Run()
	}
	err = p.engine.Run(addr)
	if err != nil {
		panic(err)
		return
	}
}

// 检查参数是否为 error 类型
func (p *API) isError(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

// 检查参数是否为 context 类型
func (p *API) isContext(v reflect.Type) bool {
	if v.Name() == "Context" && v.PkgPath() == "context" {
		return true
	}
	return false
}

// 检查参数是否接受的路由 Handler 格式
func (p *API) isHandler(t reflect.Type) error {
	if t.Kind() != reflect.Func {
		return fmt.Errorf("handler expect type func")
	}
	if t.NumIn() != 1 && t.NumIn() != 2 {
		return fmt.Errorf("max 2 parameters expected")
	}
	if t.NumIn() == 2 {
		in := t.In(1)
		for {
			if in.Kind() != reflect.Ptr {
				break
			}
			in = in.Elem()
		}
		if in.Kind() != reflect.Struct {
			return fmt.Errorf("second input parameter only acept struct")
		}
	}
	if !p.isContext(t.In(0)) {
		return fmt.Errorf("first input parameter expect type context.Context")
	}
	if t.NumOut() != 1 && t.NumOut() != 2 {
		return fmt.Errorf("max 2 output parameters expected")
	}
	if !p.isError(t.Out(t.NumOut() - 1)) {
		return fmt.Errorf("last output parameter expect type error")
	}
	return nil
}

// 反射路由 Handler, 并检查是否为可接受的格式
func (p *API) parseHandler(handler interface{}) (v reflect.Value, err error) {
	v = reflect.ValueOf(handler)
	if err = p.isHandler(v.Type()); err != nil {
		err = fmt.Errorf("unexpect handler: %s", v.Type())
		return
	}
	return
}

// 预处理路由，反射路由处理器，并检查类型
func (p *API) prepareRoutes(routes []*Route) (err error) {
	// TODO Action 冲突检查
	for _, route := range routes {
		for _, group := range route.Groups {
			for _, action := range group.Actions {
				err = p.prepareAction(action)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (p *API) prepareAction(action *Action) (err error) {
	switch action.Type {
	case Read, List:
		action.method = Get
	default:
		action.method = Post
	}
	if action.Handler != nil {
		action.handler, err = p.parseHandler(action.Handler)
	} else {
		// TODO 报错，Handler 未定义
		err = fmt.Errorf("action Handle not defined")
		return
	}
	return
}

// 递归注册路由树，处理中间件前缀逻辑，代理路由处理器为 Gin 控制器
func (p *API) registerRoutes(register Register, routes []*Route) (err error) {
	for _, route := range routes {
		routeRegister := register
		if route.Prefix != "" || len(route.Middlewares) > 0 {
			routeRegister = register.Group(route.Prefix, route.Middlewares...)
		}
		for _, group := range route.Groups {
			groupRegister := routeRegister
			if group.Prefix != "" || len(group.Middlewares) > 0 {
				groupRegister = routeRegister.Group(group.Prefix, group.Middlewares...)
			}
			for _, action := range group.Actions {
				info := p.parseHandlerInfo(action.Handler)
				path := info.ParsePath()
				fullPath := strings.Join([]string{route.Prefix, group.Prefix, path}, "")
				p.addMethod(action.method, fullPath, action.Description, info, action.handler)
				switch action.method {
				case http.MethodGet:
					groupRegister.GET(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				case http.MethodPost:
					groupRegister.POST(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				case http.MethodPut:
					groupRegister.PUT(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				case http.MethodDelete:
					groupRegister.DELETE(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				case http.MethodHead:
					groupRegister.HEAD(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				case http.MethodOptions:
					groupRegister.OPTIONS(path, append([]gin.HandlerFunc{p.proxyHandler(action.handler)}, route.Middlewares...)...)
				default:
					err = fmt.Errorf("action '%s' method '%s' unsupported", info.Name, action.method)
					return
				}
			}
		}
	}
	return
}

func (p *API) proxyHandler(handler reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		var out []reflect.Value
		var ctx context.Context
		var err error
		defer func() {
			if err != nil {
				if p.errorWrapper != nil {
					p.errorWrapper(c, err)
				} else {
					c.String(http.StatusBadRequest, err.Error())
				}
			}
			return
		}()
		if p.contextWrapper == nil {
			ctx = context.Background()
		} else {
			ctx, err = p.contextWrapper(c)
			if err != nil {
				return
			}
		}
		if handler.Type().NumIn() == 2 {
			var in reflect.Value
			in, err = bind(c, handler.Type().In(1))
			if err != nil {
				return
			}
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx), in})
		} else {
			out = handler.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}

		l := len(out)
		if !out[l-1].IsNil() {
			err = out[l-1].Interface().(error)
			return
		}
		if l == 2 {
			switch out[0].Interface().(type) {
			case Html:
				c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(out[0].Interface().(Html)))
			case Text:
				c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(out[0].Interface().(Text)))
			default:
				c.JSON(http.StatusOK, out[0].Interface())
			}
			return
		} else {
			c.String(http.StatusOK, "")
			return
		}
	}
}

func realType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t
		}
		t = t.Elem()
	}
}

func bind(c *gin.Context, t reflect.Type) (reflect.Value, error) {
	ptr := t.Kind() == reflect.Ptr
	if ptr {
		t = realType(t)
	}
	in := reflect.New(t)
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.MustBindWith(in.Interface(), b)
	if err != nil {
		return in, err
	}
	if ptr {
		return in, nil
	}
	return in.Elem(), nil
}

// 解析 Handler 的信息
func (p *API) parseHandlerInfo(h interface{}) HandlerInfo {
	target := reflect.ValueOf(h).Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

// 解析 Handler 的信息
func (p *API) parseHandlerInfoValue(v reflect.Value) HandlerInfo {
	target := v.Pointer()
	pc := runtime.FuncForPC(target)
	file, line := pc.FileLine(target)
	names := strings.Split(strings.TrimSuffix(pc.Name(), "-fm"), ".")
	return HandlerInfo{
		Name:     names[len(names)-1],
		Location: fmt.Sprintf("%s:%d", file, line),
	}
}

func (p *API) addMethod(method, path, description string, info HandlerInfo, handler reflect.Value) {
	if p.exporter == nil {
		return
	}
	m := &exporter.Method{
		Name:        info.Name,
		Path:        path,
		Method:      method,
		Description: description,
	}
	if handler.Type().NumIn() > 1 {
		m.Input = p.exporter.ReflectFields("", "", "", nil, handler.Type().In(1))
	}
	if handler.Type().NumOut() > 1 {
		m.Output = p.exporter.ReflectFields("", "", "", nil, handler.Type().Out(0))
	}
	p.methods = append(p.methods, m)
}

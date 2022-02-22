package exporter

type Options struct {
	Project    string           `json:"project"`
	Envs       []Env            `json:"envs"`
	BasicTypes []BasicType      `json:"-"`
	Makers     map[string]Maker `json:"-"`
}

type BasicType struct {
	Elem    interface{}        `json:"-"`
	Type    string             `json:"type"`
	Package string             `json:"package,omitempty"`
	Mapping map[string]Library `json:"mapping,omitempty"`
}

func (p BasicType) Fork() *BasicType {
	n := new(BasicType)
	n.Elem = p.Elem
	if p.Mapping != nil {
		n.Mapping = map[string]Library{}
		for k, v := range p.Mapping {
			n.Mapping[k] = v
		}
	}
	n.Package = p.Package
	return n
}

func (p BasicType) getMapping(lang string) *Library {
	if p.Mapping == nil {
		return nil
	}
	v, ok := p.Mapping[lang]
	if !ok {
		return nil
	}
	return &v
}

type BasicTypes struct {
	list    []*BasicType
	mapping map[string]bool
}

func (p *BasicTypes) Add(item *BasicType) {
	if p.mapping == nil {
		p.mapping = map[string]bool{}
	}
	if _, ok := p.mapping[item.Type]; ok {
		return
	}
	p.mapping[item.Type] = true
	p.list = append(p.list, item)
}

func (p BasicTypes) All() []*BasicType {
	return p.list
}

type Library struct {
	Type    string   `json:"type,omitempty"`
	Package *Package `json:"package,omitempty"`
}

type Package struct {
	Import string `json:"import,omitempty"`
	From   string `json:"from,omitempty"`
}

type Env struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type Method struct {
	Name        string `json:"name,omitempty"`
	Path        string `json:"path,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Middlewares string `json:"middlewares,omitempty"`
	Input       *Field `json:"input,omitempty"`
	Output      *Field `json:"output,omitempty"`
}

func (p Method) Fork() *Method {
	n := new(Method)
	n.Name = p.Name
	n.Path = p.Path
	n.Method = p.Method
	n.Description = p.Description
	n.Middlewares = p.Middlewares
	if p.Input != nil {
		n.Input = p.Input.Fork()
	}
	if p.Output != nil {
		n.Output = p.Output.Fork()
	}
	return n
}

type Struct struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
}

type Field struct {
	Name        string     `json:"name,omitempty"`
	Param       string     `json:"param,omitempty"`
	Label       string     `json:"label,omitempty"`
	Type        string     `json:"type,omitempty"`
	Description string     `json:"description,omitempty"`
	Array       bool       `json:"array,omitempty"`
	Struct      bool       `json:"struct,omitempty"`
	Nested      bool       `json:"nested,omitempty"`
	Origin      string     `json:"origin,omitempty"`    // 原始类型
	Fields      []*Field   `json:"fields,omitempty"`    // 描述 Struct 成员变量
	Elem        *Field     `json:"elem,omitempty"`      // 描述 Slice/Array 子元素
	Validator   *Validator `json:"validator,omitempty"` // 定义校验器
	Form        string     `json:"form,omitempty"`      // 定义表单组件
	BasicType   *BasicType `json:"-"`
}

func (p Field) Fork() *Field {
	n := new(Field)
	n.Name = p.Name
	n.Param = p.Param
	n.Label = p.Label
	n.Type = p.Type
	n.Description = p.Description
	n.Array = p.Array
	n.Struct = p.Struct
	n.Nested = p.Nested
	n.Origin = p.Origin
	for _, v := range p.Fields {
		n.Fields = append(n.Fields, v.Fork())
	}
	if p.Elem != nil {
		n.Elem = p.Elem.Fork()
	}
	n.Validator = p.Validator
	n.Form = p.Form
	n.BasicType = p.BasicType
	return n
}

type Fields struct {
	list    []*Field
	mapping map[string]bool
}

func (p *Fields) Add(item *Field) {
	if p.mapping == nil {
		p.mapping = map[string]bool{}
	}
	if _, ok := p.mapping[item.Type]; ok {
		return
	}
	p.mapping[item.Type] = true
	p.list = append(p.list, item)
}

func (p Fields) All() []*Field {
	return p.list
}

type Validator struct {
	Required bool     `json:"required,omitempty"`
	Max      *uint64  `json:"max,omitempty"`
	Min      *int64   `json:"min,omitempty"`
	Enums    []string `json:"enums,omitempty"`
}

type Component struct {
	Name string
}

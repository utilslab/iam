package iam

type Service interface {
	Actions() []ActionGroup
}

// Resource 资源
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

type ActionType string

const (
	Read  ActionType = "read"
	Write ActionType = "write"
	List  ActionType = "list"
)

// ActionGroup 动作分组
type ActionGroup struct {
	Name        string
	Description string
	Actions     []Action
}

// Action 动作
type Action struct {
	Type        ActionType
	Method      string
	Description string
	Resources   []Resource
	Handle      interface{}
	Codes       []Code
}

type Code struct {
	Status  int
	Code    string
	Message string
}

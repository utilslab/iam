package exporter

type AxiosMaker struct {
}

func (a AxiosMaker) Lang() string {
	return Ts
}

func (a AxiosMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
	data := MakeRenderData(a.Lang(), methods, EmptyNamer, TsTyper)
	for _, v := range data.Structs {
		for _, vv := range v.Fields {
			if vv.Param == "" {
				vv.Param = vv.Name
			}
		}
	}
	serviceFile := new(File)
	serviceFile.Name = "service.make.ts"
	serviceFile.Content, err = Render(axiosServiceTpl, data, EmptyFormatter)
	if err != nil {
		return
	}
	files = append(files, serviceFile)
	return
}

const axiosServiceTpl = `import axios, {AxiosPromise, AxiosRequestConfig} from 'axios';
{% for method in Methods %}
{% if method.Description %}// {{ method.Description }}{% endif %}
export function {{ method.Name }}({% if method.InputType !='' %}params: {{ method.InputType }}, {% endif %}request?: AxiosRequestConfig): {% if method.OutputType !='' %}AxiosPromise<{{ method.OutputType }}>{% else %}AxiosPromise<null>{% endif %} {
	if (!request) {
		request = {}
	}
	request.method = '{{ method.Method }}';
	request.url = '{{ method.Path }}';
	{% if method.InputType !='' %}{% if  method.Method == 'GET' or method.Method == 'DELETE' %}request.params = params;
	{% else %}request.data = params;{% endif %}{% endif %}
	return axios(request)
}
{% endfor %}
{% for struct in Structs %}
export interface {{ struct.Name }} {
{% for field in struct.Fields %}    {{field.Param}}?: {{field.Type}}, {% if field.Label or field.Description %}// {{field.Label}} {{field.Description}}{% endif %}
{% endfor %}}
{% endfor %}
`

package exporter

type UmiMaker struct {
}

func (u UmiMaker) Lang() string {
	return Ts
}
func (u UmiMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
	data := MakeRenderData(u.Lang(), methods, EmptyNamer, TsTyper)
	for _, v := range data.Structs {
		for _, vv := range v.Fields {
			if vv.Param == "" {
				vv.Param = vv.Name
			}
		}
	}
	apiFile := new(File)
	apiFile.Name = "api.make.ts"
	apiFile.Content, err = Render(umiServiceTpl, data, EmptyFormatter)
	if err != nil {
		return
	}
	typingDFile := new(File)
	typingDFile.Name = "typings.d.ts"
	typingDFile.Content, err = Render(umiTypingDTpl, data, EmptyFormatter)
	if err != nil {
		return
	}
	files = append(files, apiFile, typingDFile)
	return
}

const umiServiceTpl = `
// @ts-ignore
/* eslint-disable */
import {request} from 'umi';

{% for method in Methods %}
{% if method.Description %}// {{ method.Description }}{% endif %}
export async function {{ method.Name }}({% if method.InputType !='' %}params: API.{{ method.InputType }}, {% endif %}options?: { [key: string]: any }) {
	{% if  method.Method == 'GET' or method.Method == 'DELETE' %}return request<{% if method.OutputType !='' %}API.{{ method.OutputType }}{% else %}null{% endif %}>('{{ method.Path }}', {
		method: '{{ method.Method }}',{% if method.InputType !='' %}
		params: params,{%endif%}
		...(options || {}),
	}){% else %}return request<{% if method.OutputType !='' %}API.{{ method.OutputType }}{% else %}null{% endif %}>('{{ method.Path }}', {
		method: '{{ method.Method }}',
		headers: {
			'Content-Type': 'application/json',
		},
		{% if method.InputType !='' %}data: params,{% endif %}
		...(options || {}),
	}){% endif %}
}
{% endfor %}
`

const umiTypingDTpl = `
declare namespace API{
	{% for struct in Structs %}
	export interface {{ struct.Name }} {
	{% for field in struct.Fields %}    {{field.Param}}?: {{field.Type}}, {% if field.Label or field.Description %}// {{field.Label}} {{field.Description}}{% endif %}
	{% endfor %}}
	{% endfor %}
}
`

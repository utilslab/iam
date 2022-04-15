package exporter

import "fmt"

type AngularMaker struct {
}

var TsTyper Typer = func(s string, isStruct, isArray bool) string {
	if isArray {
		if !isStruct {
			s = typescriptTypeConverter(nil, s)
		}
		return fmt.Sprintf("%s[]", s)
	}
	if !isStruct {
		s = typescriptTypeConverter(nil, s)
	}
	return s
}

func (a AngularMaker) Lang() string {
	return Ts
}
func (a AngularMaker) Make(pkg string, methods []*Method) (files []*File, err error) {
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
	serviceFile.Content, err = Render(angularServiceTpl, data, EmptyFormatter)
	if err != nil {
		return
	}
	files = append(files, serviceFile)
	return
}

const angularServiceTpl = `import {HttpClient, HttpHeaders, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';

export class APIService {

    client: HttpClient;
    host: string;

    constructor(client:HttpClient, host:string){
	  this.client = client;
      this.host = host;
    }
{% for method in Methods %}
{% if method.Description %}    // {{ method.Description }}{% endif %}
    {{ method.Name }}({% if method.InputType !='' %}params:{{ method.InputType }}, {% endif %}options?:HttpOptions):{% if method.OutputType !='' %}Observable<{{ method.OutputType }}>{% else %}Observable<null>{% endif %}{ {% if method.InputType !='' %}
	    if(!options){
           options = {};
	    }
		{% if  method.Method == 'GET' or method.Method == 'DELETE' %}  // @ts-ignore
		  options.params = params;{% else %}  options.body = params;{% endif %}{% endif %}
	    return this.client.request('{{ method.Method }}', this.host+'{{ method.Path }}', options)
    }{% endfor %}
}

export interface HttpOptions {
    body?: any;
    headers?: HttpHeaders | {
        [header: string]: string | string[];
    };
    params?: HttpParams | {
        [param: string]: string | string[];
    };
    observe?: 'body' | 'events' | 'response';
    reportProgress?: boolean;
    responseType?: 'arraybuffer' | 'blob' | 'json' | 'text';
    withCredentials?: boolean;
}

{% for struct in Structs %}
export interface {{ struct.Name }} {
{% for field in struct.Fields %}    {{field.Param}}?: {{field.Type}}, {% if field.Label or field.Description %}// {{field.Label}} {{field.Description}}{% endif %}
{% endfor %}}
{% endfor %}
`

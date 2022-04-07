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
	serviceFile.Content, err = Render(angularServiceTpl, data, EmptyFormatter)
	if err != nil {
		return
	}
	files = append(files, serviceFile)
	return
}

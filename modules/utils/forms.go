// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utils

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/astaxie/beego/utils/forms"

	"github.com/beego/wetalk/setting"
)

func init() {
	initExtraField()
}

func initExtraField() {
	// a example to create a bootstrap style checkbox creater
	forms.RegisterFieldCreater("checkbox", func(fSet *forms.FieldSet) {
		fSet.Field = func() template.HTML {
			value := false
			if b, ok := fSet.Value.(bool); ok {
				value = b
			}
			active := ""
			if value {
				active = " active"
			}
			return template.HTML(fmt.Sprintf(`<label>
            <input type="hidden" name="%s" value="%v">
            <button type="button" data-toggle="button" data-name="%s" class="btn btn-default btn-xs btn-checked%s">
            	<i class="icon icon-ok"></i>
        	</button>%s
        </label>`, fSet.Name, value, fSet.Name, active, fSet.LabelText))
		}
	})

	// a example to create a select2 box
	forms.RegisterFieldFilter("select", func(fSet *forms.FieldSet) {
		if strings.Index(fSet.Attrs, `rel="select2"`) != -1 {
			field := fSet.Field.String()
			fSet.Field = func() template.HTML {
				field = strings.Replace(field, "<option", "<option></option><option", 1)
				return template.HTML(field)
			}
		}
	})

	// a example to create a captcha field
	forms.RegisterFieldCreater("captcha", func(fSet *forms.FieldSet) {
		fSet.Label = template.HTML(fmt.Sprintf(`
					<label class="control-label" for="%s">%s</label>`, fSet.Id, fSet.LabelText))

		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(`%v
			<input id="%s" name="%s" type="text" value="" class="form-control" autocomplete="off"%s%s>`,
				setting.Captcha.CreateCaptchaHtml(), fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs))
		}
	})

	forms.RegisterFieldCreater("empty", func(fSet *forms.FieldSet) {
		fSet.Label = ""
		fSet.Field = func() template.HTML { return "" }
	})
}

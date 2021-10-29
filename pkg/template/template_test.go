package template

import (
    "html/template"
    "reflect"
    "testing"

    "github.com/configcenter/internal/mock"
    "github.com/configcenter/pkg/repository"
    "github.com/golang/mock/gomock"
)

func TestNewCtlFindTemplate(t *testing.T) {
    type args struct {
        tmplInstanceName string
    }
    tests := []struct {
        name    string
        args    args
        want    *TemplateImpl
        wantErr bool
    }{
        {
            name: "normal",
            args: args{
                tmplInstanceName: "test",
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewCtlFindTemplate(tt.args.tmplInstanceName)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewCtlFindTemplate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if _, ok := got.funcMap["CtlFind"]; !ok {
                t.Errorf("funcMap err")
            }
            if got.allTemplates == nil {
                t.Errorf("init temp err")
            }
        })
    }
}

func TestNewTemplateImpl(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(infrastructureJson), nil),
        //mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("json err"), nil),
    )
    repository.Src = mockSrc
    type args struct {
        source           repository.Storage
        globalId         string
        localId          string
        tmplInstanceName string
        version          string
        env              string
    }
    tests := []struct {
        name    string
        args    args
        want    *TemplateImpl
        wantErr bool
    }{
        {
            name: "normal",
            args: args{
                source:           mockSrc,
                version:          "0.0.1",
                env:              "00",
                globalId:         "3254",
                localId:          "0",
                tmplInstanceName: "test",
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewTemplateImpl(tt.args.source, tt.args.globalId, tt.args.localId, tt.args.tmplInstanceName, tt.args.version, tt.args.env)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewTemplateImpl() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if _, ok := got.funcMap["CtlFind"]; !ok {
                t.Errorf("funcMap ctlfind err")
            }
            if _, ok := got.funcMap["GetInfo"]; !ok {
                t.Errorf("funcMap GetInfo err")
            }
            if _, ok := got.funcMap["UnsafeGetInfo"]; !ok {
                t.Errorf("funcMap UnsafeGetInfo err")
            }
            if got.allTemplates == nil {
                t.Errorf("init temp err")
            }
        })
    }
}

func TestTemplateImpl_Fill(t1 *testing.T) {
    tmp := template.New("ins")
    tmp, err := tmp.Parse("ins")
    if err != nil {
        t1.Fatal(err)
    }
    type fields struct {
        funcMap      map[string]interface{}
        allTemplates *template.Template
    }
    type args struct {
        tmplContent []byte
        tmplName    string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []byte
        wantErr bool
    }{
        {
            name: "normal",
            fields: fields{
                funcMap:      nil,
                allTemplates: tmp,
            },
            args: args{
                tmplName:    "test",
                tmplContent: []byte("test"),
            },
            want:    []byte("test"),
            wantErr: false,
        },
        {
            name: "nil data input",
            args: args{
                tmplName:    "test",
                tmplContent: nil,
            },
            want:    nil,
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t1.Run(tt.name, func(t1 *testing.T) {
            t := &TemplateImpl{
                funcMap:      tt.fields.funcMap,
                allTemplates: tt.fields.allTemplates,
            }
            got, err := t.Fill(tt.args.tmplContent, tt.args.tmplName)
            if (err != nil) != tt.wantErr {
                t1.Errorf("Fill() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t1.Errorf("Fill() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestTemplateImpl_addTmpl(t1 *testing.T) {
    type fields struct {
        funcMap      map[string]interface{}
        allTemplates *template.Template
    }
    type args struct {
        tmplContent []byte
        tmplName    string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t1.Run(tt.name, func(t1 *testing.T) {
            t := &TemplateImpl{
                funcMap:      tt.fields.funcMap,
                allTemplates: tt.fields.allTemplates,
            }
            if err := t.addTmpl(tt.args.tmplContent, tt.args.tmplName); (err != nil) != tt.wantErr {
                t1.Errorf("addTmpl() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

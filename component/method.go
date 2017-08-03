package component

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/lonnng/nano/session"
)

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
	typeOfSession = reflect.TypeOf(session.New(nil))
)

func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// isHandlerMethod decide a method is suitable handler method
func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, *Session, []byte or pointer.
	if mt.NumIn() != 3 {
		return false
	}

	// Method needs one outs: error
	if mt.NumOut() != 1 {
		return false
	}

	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfSession {
		return false
	}

	if (mt.In(2).Kind() != reflect.Ptr && mt.In(2) != typeOfBytes) || mt.Out(0) != typeOfError {
		return false
	}
	return true
}

// suitableMethods returns suitable methods of typ
func suitableHandlerMethods(typ reflect.Type) map[string]*Handler {
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name
		if isHandlerMethod(method) {
			raw := false
			if mt.In(2) == typeOfBytes {
				raw = true
			}
			methods[mn] = &Handler{Method: method, Type: mt.In(2), IsRawArg: raw}
		}
	}
	return methods
}
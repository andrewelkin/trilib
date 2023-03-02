package dynamic_types

import (
	"fmt"
	"github.com/andrewelkin/trilib/utils"
	"github.com/andrewelkin/trilib/utils/logger"
	"golang.org/x/net/context"
	"reflect"
)

var dynTypeRegistry = make(map[string]reflect.Type)

func RegisterDynType(name string, i interface{}) {
	fmt.Printf("Registering dynamic type %s\n", name)
	t := reflect.TypeOf(i).Elem()
	dynTypeRegistry[fmt.Sprintf("%v", t)] = t
	dynTypeRegistry[name] = t
}

func MakeDynInstance(name string, ctx context.Context, cfg *utils.Config, logger logger.Logger) interface{} {
	t, ok := dynTypeRegistry[name]
	if !ok {
		utils.Throwf("error initializing dynamical class '%s' : not registered", name)
	}
	i := reflect.New(t).Interface()
	ti := reflect.ValueOf(i)
	f := ti.MethodByName("Init")
	res := f.Call([]reflect.Value{
		reflect.ValueOf(name),
		reflect.ValueOf(ctx),
		reflect.ValueOf(cfg),
		reflect.ValueOf(logger),
	})

	if len(res) > 1 { // second parameter can be error
		if err := res[1].Interface(); err != nil {
			utils.Throwf("error initializing dynamical class '%s' : %v", name, err)
		}
	}
	ret := res[0].Elem().Interface()
	return ret
}

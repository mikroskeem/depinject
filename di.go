package depinject

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrRegisterCanOnlySupplyPointer = errors.New("can only supply a pointer")
	ErrCanOnlyInjectPointer         = errors.New("can only inject into a pointer")
	ErrUnableToSetField             = errors.New("unable to set field")
)

type DI struct {
	services map[reflect.Type]interface{}
}

func (d *DI) RegisterComponent(instance interface{}) error {
	t := reflect.TypeOf(instance)
	if t.Kind() != reflect.Ptr {
		return ErrRegisterCanOnlySupplyPointer
	}

	if d.services == nil {
		d.services = map[reflect.Type]interface{}{}
	}
	d.services[t] = instance
	return nil
}

func (d *DI) MustRegisterComponent(instance interface{}) {
	if err := d.RegisterComponent(instance); err != nil {
		panic(err)
	}
}

func (d *DI) InjectAndRegister(instance interface{}) error {
	if _, err := d.Inject(instance); err != nil {
		return err
	}
	if err := d.RegisterComponent(instance); err != nil {
		return err
	}
	return nil
}

func (d *DI) MustInjectAndRegister(instance interface{}) {
	if err := d.InjectAndRegister(instance); err != nil {
		panic(err)
	}
}

// Inject takes struct pointer
func (d *DI) Inject(into interface{}) (res interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic: %v", r)
			}
		}
	}()

	t := reflect.TypeOf(into)
	if t.Kind() != reflect.Ptr {
		err = ErrCanOnlyInjectPointer
		return
	}

	instance := reflect.Indirect(reflect.ValueOf(into))
	reflType := t.Elem()

	for i := 0; i < reflType.NumField(); i++ {
		field := reflType.Field(i)

		if field.Tag.Get("depinject") == "skip" {
			continue
		}

		instField := instance.Field(i)
		if !instField.CanSet() {
			continue
		}

		if mapv, ok := d.services[field.Type]; ok {
			instField.Set(reflect.ValueOf(mapv))
		}
	}

	res = into
	return
}

func (d *DI) MustInject(into interface{}) interface{} {
	if r, err := d.Inject(into); err != nil {
		panic(err)
	} else {
		return r
	}
}

package di

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/juju/errors"
	"github.com/samber/do"
)

type NamedProvider[T any] func(*do.Injector, string) (T, error)

func providerKey[T any, TProvider do.Provider[T] | NamedProvider[T]](name *string, _ TProvider) string {
	return diKey[T](name)
}

func invokeKey[T any](name *string) string {
	return diKey[T](name)
}

func diKey[T any](tag *string) string {
	var obj T
	t := reflect.TypeOf(obj)
	isPtr := false
	switch {
	// Interface
	case t == nil:
		t = reflect.TypeOf(new(T)).Elem()
	// Pointer to Struct
	case t.Kind() == reflect.Ptr:
		t, isPtr = t.Elem(), true
	}

	pkgName := t.PkgPath()
	typeName := t.Name()
	if isPtr {
		typeName = fmt.Sprintf("*%s", typeName)
	}

	var depKey string
	switch tag {
	case nil:
		depKey = fmt.Sprintf("%s::%s", pkgName, typeName)
	default:
		depKey = fmt.Sprintf("%s::%s#%s", pkgName, typeName, *tag)
	}
	return depKey
}

func Invoke[T any](injector *do.Injector) T {
	return InvokeNamed[T](injector, nil)
}

func InvokeNamed[T any](injector *do.Injector, name *string) T {
	key := invokeKey[T](name)
	inst, err := do.InvokeNamed[T](injector, key)
	if err != nil {
		diLogger().Error(
			"DI: failed to Invoke service",
			"key", key,
			"err", err.Error(),
			"stack", errors.ErrorStack(err),
		)
		panic(errors.Annotatef(err, "Stack: %s", errors.ErrorStack(err)))
	}
	return inst
}

func Provide[T any](injector *do.Injector, provider do.Provider[T]) {
	key := providerKey[T](nil, provider)
	inst, err := provider(injector)
	if err != nil {
		diLogger().Error(
			"DI: failed to Provide service",
			"key", key,
			"err", err.Error(),
			"stack", errors.ErrorStack(err),
		)
		panic(errors.Annotatef(err, "Stack: %s", errors.ErrorStack(err)))
	}
	do.ProvideNamedValue(injector, key, inst)
}

func ProvideNamed[T any](injector *do.Injector, name string, provider NamedProvider[T]) {
	key := providerKey[T](&name, provider)
	inst, err := provider(injector, name)
	if err != nil {
		diLogger().Error(
			"DI: failed to Provide service",
			"key", key,
			"err", err.Error(),
			"stack", errors.ErrorStack(err),
		)
		panic(errors.Annotatef(err, "Stack: %s", errors.ErrorStack(err)))
	}
	do.ProvideNamedValue(injector, key, inst)
}

func InvokeOrProvide[T any](injector *do.Injector, provider do.Provider[T]) T {
	key := providerKey[T](nil, provider)
	if err := do.HealthCheckNamed(injector, key); err != nil {
		Provide(injector, provider)
	}
	return Invoke[T](injector)
}

func InvokeOrProvideNamed[T any](injector *do.Injector, name string, provider NamedProvider[T]) T {
	key := providerKey[T](&name, provider)
	if err := do.HealthCheckNamed(injector, key); err != nil {
		ProvideNamed(injector, name, provider)
	}
	return InvokeNamed[T](injector, &name)
}

func CreateInjector(logging bool) *do.Injector {
	if !logging {
		return do.New()
	}
	return do.NewWithOpts(&do.InjectorOpts{
		Logf: func(format string, args ...interface{}) {
			diLogger().Debug(
				fmt.Sprintf(format, args...),
			)
		},
	})
}

func diLogger() *slog.Logger {
	return slog.Default().With("package", "di")
}

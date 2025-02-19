package di

import (
	"fmt"
	"reflect"

	"github.com/inconshreveable/log15"
	"github.com/juju/errors"
	"github.com/samber/do"

	"github.com/khwong-c/wtcode/tooling/log"
)

var logger = log.NewLogger("di")

type NamedProvider[T any] func(*do.Injector, string) (T, error)

func providerKey[T any, TProvider do.Provider[T] | NamedProvider[T]](name *string, _ TProvider) string {
	return invokeKey[T](name)
}

func invokeKey[T any](name *string) string {
	var stub T
	outputType := reflect.TypeOf(stub)
	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	pkgName := outputType.PkgPath()
	typeName := outputType.Name()
	var depKey string
	switch name {
	case nil:
		depKey = fmt.Sprintf("%s::%s", pkgName, typeName)
	default:
		depKey = fmt.Sprintf("%s::%s#%s", pkgName, typeName, *name)
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
		logger.Crit("DI: failed to Invoke service", "key", key, "err", err, "stack", errors.ErrorStack(err))
		panic(errors.Annotatef(err, errors.ErrorStack(err)))
		//panic(errors.Trace(err))
	}
	return inst
}

func Provide[T any](injector *do.Injector, provider do.Provider[T]) {
	key := providerKey[T](nil, provider)
	inst, err := provider(injector)
	if err != nil {
		logger.Crit("DI: failed to provide service", "key", key, "err", err, "stack", errors.ErrorStack(err))
		panic(errors.Annotatef(err, errors.ErrorStack(err)))
	}
	do.ProvideNamedValue(injector, key, inst)
}

func ProvideNamed[T any](injector *do.Injector, name string, provider NamedProvider[T]) {
	key := providerKey[T](&name, provider)
	inst, err := provider(injector, name)
	if err != nil {
		logger.Crit("DI: failed to provide service", "key", key, "err", err, "stack", errors.ErrorStack(err))
		panic(errors.Annotatef(err, errors.ErrorStack(err)))
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

func CreateInjector(logging, withStack bool) *do.Injector {
	if !logging {
		return do.New()
	}

	logger := log.NewLogger("di")
	if withStack {
		handler := logger.GetHandler()
		logger.SetHandler(log15.CallerStackHandler("%+v", handler))
	}

	return do.NewWithOpts(&do.InjectorOpts{
		Logf: func(format string, args ...interface{}) {
			logger.Debug(
				fmt.Sprintf(format, args...),
			)
		},
	})
}

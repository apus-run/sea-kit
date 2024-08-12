package async

import "encoding/json"

/*
Future is an execution result of an asynchronous function
that returns immediately, without locking execution thread.
To lock execution and wait for result, use .Await() method
or async.Await() function. As an alternative you can use
a syntax similar to JavaScript Promise, using .Then()
and .Catch() methods.

Usage:

	// Let's assume we have a future object in "ftr" variable.
	// We can lock execution and wait for a result with .Await()
	res, err := ftr.Await()
	// Or, we can use async.Await()
	res, err := async.Await(ftr)
	// Or, we can avoid locking execution and provide then/catch
	// functions to handle execution results.
	ftr.Then(func(val string) {
		println(val)
	}).Catch(func(err error) {
		println(err.Error())
	})
*/
type Future[T any] struct {
	value chan T
	cache *T
	err   error

	onthen func(T)
	oncatch func(error)
}

/*
Await for a future object results.

Usage:

	// Let's assume we have a future object in "ftr" variable.
	res, err := ftr.Await()
*/
func (f *Future[T]) Await() (T, error) {
	// Return from cache, if exists
	if f.cache != nil {
		return *f.cache, f.err
	}
	// Wait for value
	v := <-f.value
	// Save to cache
	f.cache = &v
	// Return
	return v, f.err
}

/*
AwaitRuntime is a runtime version of .Await()

Usage:

	// Let's assume we have a future object in "ftr" variable.
	// Result will be stored as "any" type, so you'll need to cast it.
	res, err := ftr.AwaitRuntime()
*/
func (f *Future[T]) AwaitRuntime() (any, error) {
	return f.Await()
}

/*
Then accepts a function, that will be executed on
future work completion.

Usage:

	// Let's assume we have a future object of string in "ftr" variable.
	ftr.Then(func(v string) {
		println(v)
	})
*/
func (f *Future[T]) Then(fn func(T)) *Future[T] {
	// Await first
	f.Await() //nolint:errcheck
	// If no error, call provided function
	if f.err == nil {
		fn(*f.cache)
	}
	// Self-return
	return f
}

/*
Catch accepts a function, that will be executed on
future execution error.

Usage:

	// Let's assume we have a future object of string in "ftr" variable.
	ftr.Catch(func(err error) {
		println(err.Error())
	})
*/
func (f *Future[T]) Catch(fn func(error)) *Future[T] {
	// Await first
	f.Await() //nolint:errcheck
	// If error, call provided function
	if f.err != nil {
		fn(f.err)
	}
	// Self-return
	return f
}

/*
MarshalJSON implements future marshalling.
*/
func (f *Future[T]) MarshalJSON() ([]byte, error) {
	val, err := Await(f)
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(val)
}

/*
UnmarshalJSON implements future unmarshalling.
*/
func (f *Future[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.cache)
}

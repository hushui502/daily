## “优雅”的处理golang的错误

首先声明我不觉得golang目前的错误处理有什么大问题，try catch也挺好，各有优缺点，这里只是给出一种错误处理的方案，并非最佳。

## The classic example
```cassandraql
func copyfile(src, dst) error {
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fdst.Close()

	err := io.Copy(fdst, fsrc)

	return err
}
```
上面这种其实还好，因为类似这种文件操作的代码我们都会封装成公共库来调用，因此可以避免if err != nil的重复灾难。


但是下面这种http.Handler如果再直接判断err就显得有点问题，主要原因是我们的业务代码复杂重复。
```
func handleThing(w http.ResponseWriter, r *http.Request) {
	// Our path is something like /thing/3
	id, err := idFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusNotFound)
		return
	}

	thing, err := store.GetThingByID(id)
	if err != nil {
		// The error might be sql.NoRows, or it might be something else.
		if store.IsNotFoundErr(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Printf("Failed to get a thing: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	acc := AccountFromRequest(r)
	if acc == nil {
		// No account attached to the request's session -> permission denied.
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	has, err := thing.HasPermissionToView(acc)
	if err != nil {
		// For some reason, we failed to check permissions. Better log it.
		log.Printf("Failed to check permissions: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !has {
		// Permission denied.
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// All good, send data to the client.
	respond(w, r, decodeThing(thing))
}
```
可以用点设计模式了。

```cassandraql
func viewRecord(w http.ResponseWriter, r *http.Request) error {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        return err
    }
    return viewTemplate.Execute(w, record)
}

// NOTE: the following is my adapted version from the example's ServeHTTP to a
// middleware/wrapper

type HandlerE = func(w http.ResponseWriter, r *http.Request) error

func WithError(h HandlerE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
```

这样似乎已经满足了，但是可以选择不直接返回错误交给客户端处理，外加一些日志处理。如下：
```cassandraql
func WithError(h HandlerE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {

			if is404err(err) {
				http.Error(w, "not found", 404)
				return
			}

			if isBadRequest(err) {
				http.Error(w, "bad request", 400)
				return
			}

			// Some other special cases...
			// ...

			log.Printf("Something went wrong: %v", err)

			http.Error(w, "Internal server error", 500)
		}
	}
}
```
手写各种is404err这样导致的问题就是重构灾难，所以需要增加一个interface。

```cassandraql
type ErrorResponder interface {
    // RespondError writes an error message to w. If it doesn't know what to
    // respond, it returns false.
	RespondError(w http.ResponseWriter, r *http.Request) bool
}
```

```cassandraql
func WithError(h HandlerE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if er, ok := err.(ErrorResponder); ok {
				if er.RespondError(w, r) {
					return
				}
			}

			log.Printf("Something went wrong: %v", err)

			http.Error(w, "Internal server error", 500)
		}
	}
}
```
剩下的就是增加一些实现这个错误接口的错误类型。
```cassandraql
// BadRequest error responds with bad request status code, and optionally with
// a json body.
type BadRequestError struct {
	err  error
	body interface{}
}

func BadRequest(err error) *BadRequestError {
	return &BadRequestError{err: err}
}

func BadRequestWithBody(body interface{}) *BadRequestError {
	return &BadRequestError{body: body}
}

func (e *BadRequestError) RespondError(w http.ResponseWriter, r *http.Request) bool {
	if e.body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusBadRequest)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(e.body)

		if err != nil {
			log.Printf("Failed to encode a response: %v", err)
		}
	}

	return true
}

func (e *BadRequestError) Error() string {
	return e.err.Error()
}

// Maybe404Error responds with not found status code, if its supplied error
// is sql.ErrNoRows.
type Maybe404Error struct {
	err error
}

func Maybe404(err error) *Maybe404Error {
	return &Maybe404Error{err: err}
}

func (e *Maybe404Error) Error() string {
	return fmt.Sprintf("Maybe404: %v", e.err.Error())
}

func (e *Maybe404Error) Is404() bool {
	return errors.Is(e.err, sql.ErrNoRows)
}

func (e *Maybe404Error) RespondError(w http.ResponseWriter, r *http.Request) bool {
	if !e.Is404() {
		return false
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return true
}
```

更改handler func，利用中间件的使用方式
```cassandraql
func handleThing(w http.ResponseWriter, r *http.Request) error {
	// Our path is something like /thing/3
	id, err := idFromPath(r.URL.Path)
	if err != nil {
		// Literally bad request. We could use BadRequestWithBody to
		// respond with a fancy information for the client.
		return BadRequest(err)
	}

	thing, err := store.GetThingByID(id)
	if err != nil {
		// Likely a not found issue, but something else might have gone wrong.
		// Maybe404Error handles both cases.
		return Maybe404(err)
	}

	acc := AccountFromRequest(r)
	if acc == nil {
		// No account attached to the request. Client needs to authenticate.
		return AuthenticationRequired()
	}

	has, err := thing.HasPermissionToView(acc)
	if err != nil {
		// Something actually went wrong. Error will be logged and 500 message
		// sent to the client.
		return err
	}

	if !has {
		// Client doesn't have permission to view this resource.
		return PermissionDenied()
	}

	// All good, send data to the client.
	respond(w, r, decodeThing(thing))
}

func main() {
	...
	mux.Handle("/thing/", WithError(handleThing))
	...
}
```

经验活，并非银弹。

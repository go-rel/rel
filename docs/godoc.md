# martini
`import "github.com/codegangsta/martini"`


## Constants
``` go
const (
    Dev  string = "development"
    Prod string = "production"
    Test string = "test"
)
```
Envs


## <a name="pkg-variables">Variables</a>
``` go
var Env = Dev
```
Env is the environment that Martini is executing in. The MARTINI_ENV is read on initialization to set this variable.

``` go
var Root string
```



## <a name="BeforeFunc">type</a> [BeforeFunc](/src/target/response_writer.go?s=962:998#L29)
``` go
type BeforeFunc func(ResponseWriter)
```
BeforeFunc is a function that is called before the ResponseWriter has been written to.










## <a name="ClassicMartini">type</a> [ClassicMartini](/src/target/martini.go?s=3381:3429#L111)
``` go
type ClassicMartini struct {
    *Martini
    Router
}
```
ClassicMartini represents a Martini with some reasonable defaults. Embeds the router functions for convenience.







### <a name="Classic">func</a> [Classic](/src/target/martini.go?s=3608:3638#L118)
``` go
func Classic() *ClassicMartini
```
Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery and martini.Static.
Classic also maps martini.Routes as a service.





## <a name="Context">type</a> [Context](/src/target/martini.go?s=4293:4663#L140)
``` go
type Context interface {
    inject.Injector
    // Next is an optional function that Middleware Handlers can call to yield the until after
    // the other Handlers have been executed. This works really well for any operations that must
    // happen after an http request
    Next()
    // Written returns whether or not the response for this context has been written.
    Written() bool
}
```
Context represents a request context. Services can be mapped on the request level from this interface.










## <a name="Handler">type</a> [Handler](/src/target/martini.go?s=4011:4035#L131)
``` go
type Handler interface{}
```
Handler can be any callable function. Martini attempts to inject services into the handler's argument list.
Martini will panic if an argument could not be fullfilled via dependency injection.







### <a name="Logger">func</a> [Logger](/src/target/logger.go?s=164:185#L10)
``` go
func Logger() Handler
```
Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.


### <a name="Recovery">func</a> [Recovery](/src/target/recovery.go?s=2766:2789#L115)
``` go
func Recovery() Handler
```
Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
While Martini is in development mode, Recovery will also output the panic as HTML.


### <a name="Static">func</a> [Static](/src/target/static.go?s=1447:1512#L53)
``` go
func Static(directory string, staticOpt ...StaticOptions) Handler
```
Static returns a middleware handler that serves static files in the given directory.





## <a name="Martini">type</a> [Martini](/src/target/martini.go?s=612:712#L30)
``` go
type Martini struct {
    inject.Injector
    // contains filtered or unexported fields
}
```
Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.







### <a name="New">func</a> [New](/src/target/martini.go?s=843:862#L38)
``` go
func New() *Martini
```
New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.





### <a name="Martini.Action">func</a> (\*Martini) [Action](/src/target/martini.go?s=1498:1539#L55)
``` go
func (m *Martini) Action(handler Handler)
```
Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().




### <a name="Martini.Handlers">func</a> (\*Martini) [Handlers](/src/target/martini.go?s=1212:1259#L47)
``` go
func (m *Martini) Handlers(handlers ...Handler)
```
Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
Will panic if any of the handlers is not a callable function




### <a name="Martini.Logger">func</a> (\*Martini) [Logger](/src/target/martini.go?s=1617:1661#L61)
``` go
func (m *Martini) Logger(logger *log.Logger)
```
Logger sets the logger




### <a name="Martini.Run">func</a> (\*Martini) [Run](/src/target/martini.go?s=2830:2853#L90)
``` go
func (m *Martini) Run()
```
Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.




### <a name="Martini.RunOnAddr">func</a> (\*Martini) [RunOnAddr](/src/target/martini.go?s=2251:2291#L79)
``` go
func (m *Martini) RunOnAddr(addr string)
```
Run the http server on a given host and port.




### <a name="Martini.ServeHTTP">func</a> (\*Martini) [ServeHTTP](/src/target/martini.go?s=2092:2163#L74)
``` go
func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request)
```
ServeHTTP is the HTTP Entry point for a Martini instance. Useful if you want to control your own HTTP server.




### <a name="Martini.Use">func</a> (\*Martini) [Use](/src/target/martini.go?s=1866:1904#L67)
``` go
func (m *Martini) Use(handler Handler)
```
Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added.




## <a name="Params">type</a> [Params](/src/target/router.go?s=227:256#L13)
``` go
type Params map[string]string
```
Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.










## <a name="ResponseWriter">type</a> [ResponseWriter](/src/target/response_writer.go?s=302:870#L13)
``` go
type ResponseWriter interface {
    http.ResponseWriter
    http.Flusher
    http.Hijacker
    // Status returns the status code of the response or 0 if the response has not been written.
    Status() int
    // Written returns whether or not the ResponseWriter has been written.
    Written() bool
    // Size returns the size of the response body.
    Size() int
    // Before allows for a function to be called before the ResponseWriter has been written to. This is
    // useful for setting headers or any other operations that must happen before a response has been written.
    Before(BeforeFunc)
}
```
ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
if the functionality calls for it.







### <a name="NewResponseWriter">func</a> [NewResponseWriter](/src/target/response_writer.go?s=1080:1141#L32)
``` go
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter
```
NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter





## <a name="ReturnHandler">type</a> [ReturnHandler](/src/target/return_handler.go?s=324:373#L13)
``` go
type ReturnHandler func(Context, []reflect.Value)
```
ReturnHandler is a service that Martini provides that is called
when a route handler returns something. The ReturnHandler is
responsible for writing to the ResponseWriter based on the values
that are passed into this function.










## <a name="Route">type</a> [Route](/src/target/router.go?s=5128:5495#L189)
``` go
type Route interface {
    // URLWith returns a rendering of the Route's url with the given string params.
    URLWith([]string) string
    // Name sets a name for the route.
    Name(string)
    // GetName returns the name of the route.
    GetName() string
    // Pattern returns the pattern of the route.
    Pattern() string
    // Method returns the method of the route.
    Method() string
}
```
Route is an interface representing a Route in Martini's routing layer.










## <a name="RouteMatch">type</a> [RouteMatch](/src/target/router.go?s=6195:6214#L228)
``` go
type RouteMatch int
```

``` go
const (
    NoMatch RouteMatch = iota
    StarMatch
    OverloadMatch
    ExactMatch
)
```









### <a name="RouteMatch.BetterThan">func</a> (RouteMatch) [BetterThan](/src/target/router.go?s=6323:6372#L238)
``` go
func (r RouteMatch) BetterThan(o RouteMatch) bool
```
Higher number = better match




## <a name="Router">type</a> [Router](/src/target/router.go?s=374:1829#L16)
``` go
type Router interface {
    Routes

    // Group adds a group where related routes can be added.
    Group(string, func(Router), ...Handler)
    // Get adds a route for a HTTP GET request to the specified matching pattern.
    Get(string, ...Handler) Route
    // Patch adds a route for a HTTP PATCH request to the specified matching pattern.
    Patch(string, ...Handler) Route
    // Post adds a route for a HTTP POST request to the specified matching pattern.
    Post(string, ...Handler) Route
    // Put adds a route for a HTTP PUT request to the specified matching pattern.
    Put(string, ...Handler) Route
    // Delete adds a route for a HTTP DELETE request to the specified matching pattern.
    Delete(string, ...Handler) Route
    // Options adds a route for a HTTP OPTIONS request to the specified matching pattern.
    Options(string, ...Handler) Route
    // Head adds a route for a HTTP HEAD request to the specified matching pattern.
    Head(string, ...Handler) Route
    // Any adds a route for any HTTP method request to the specified matching pattern.
    Any(string, ...Handler) Route
    // AddRoute adds a route for a given HTTP method request to the specified matching pattern.
    AddRoute(string, string, ...Handler) Route

    // NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.
    NotFound(...Handler)

    // Handle is the entry point for routing. This is used as a martini.Handler
    Handle(http.ResponseWriter, *http.Request, Context)
}
```
Router is Martini's de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.







### <a name="NewRouter">func</a> [NewRouter](/src/target/router.go?s=2289:2312#L68)
``` go
func NewRouter() Router
```
NewRouter creates a new Router instance.
If you aren't using ClassicMartini, then you can add Routes as a
service with:


	m := martini.New()
	r := martini.NewRouter()
	m.MapTo(r, (*martini.Routes)(nil))

If you are using ClassicMartini, then this is done for you.





## <a name="Routes">type</a> [Routes](/src/target/router.go?s=8189:8567#L328)
``` go
type Routes interface {
    // URLFor returns a rendered URL for the given route. Optional params can be passed to fulfill named parameters in the route.
    URLFor(name string, params ...interface{}) string
    // MethodsFor returns an array of methods available for the path
    MethodsFor(path string) []string
    // All returns an array with all the routes in the router.
    All() []Route
}
```
Routes is a helper service for Martini's routing layer.










## <a name="StaticOptions">type</a> [StaticOptions](/src/target/static.go?s=196:888#L13)
``` go
type StaticOptions struct {
    // Prefix is the optional prefix used to serve the static directory content
    Prefix string
    // SkipLogging will disable [Static] log messages when a static file is served.
    SkipLogging bool
    // IndexFile defines which file to serve as index if it exists.
    IndexFile string
    // Expires defines which user-defined function to use for producing a HTTP Expires Header
    // https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
    Expires func() string
    // Fallback defines a default URL to serve when the requested resource was
    // not found.
    Fallback string
    // Exclude defines a pattern for URLs this handler should never process.
    Exclude string
}
```
StaticOptions is a struct for specifying configuration options for the martini.Static middleware.

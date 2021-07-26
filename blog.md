# Generate fully working Go CRUD HTTP API with Ent

## Schema as Code + Code Generation = 😻

When we say that one of the core principles of Ent is "Schema as Code", we mean by that more than "Ent's DSL for
defining entities and their edges is done using regular Go code". Ent's unique approach, compared to many other ORMs, is
to express all of the logic related to an entity, as code, directly in the schema definition.

With Ent, developers can write all authorization logic (called "Privacy" within Ent), and all of the mutation
side-effects (called "Hooks" within Ent) directly on the schema. Having everything in the same place can be very
convenient, but its true power is revealed when paired with code generation.

If schemas are defined this way, it becomes possible to generate code for fully working production-grade servers
automatically. If we move the responsibility for authorization decisions and custom side effects from the RPC layer to
the data layer, the implementation of the basic CRUD (Create, Read, Update and Delete) endpoints becomes generic to the
extent that it can be machine-generated. This is exactly the idea behind the popular GraphQL and gRPC Ent extensions.

Today, we would like to present a new Ent Extension named elk that can automatically generate fully working, RESTful API
endpoints from your Ent schemas. Elk strives to automate all of the tedious work of setting up the basic CRUD endpoints
for every entity you add to your graph, including logging, validation of the request body, eager loading relations and
serializing, all while leaving reflection out of sight and maintaining type-safety.

Buckle up, nerds, and strap yourself in for a journey filled with prodigies and magic!

## Getting Started

The final version of the code below can be found on GitHub.

Start by creating a new Go project:

```shell
mkdir elk-example
cd elk-example
go mod init elk-example
```

Invoke the ent code generator and create three schemas: User, Pet, Group:

```shell
go run -mod=mod entgo.io/ent/cmd/ent init Pet User Group
```

Your project should now look like this:

```
.
├── ent
│   ├── generate.go
│   └── schema
│       ├── group.go
│       ├── pet.go
│       └── user.go
├── go.mod
└── go.sum
```

The next step is to add the `elk` package to our project:

```shell
go get -u github.com/masseelch/elk
```

`elk` uses the
Ent [extension api](https://github.com/ent/ent/blob/a19a89a141cf1a5e1b38c93d7898f218a1f86c94/entc/entc.go#L197) to add
its templates. This requires you to use the `entc` (ent codegen) package as
described [here](https://entgo.io/docs/code-gen#use-entc-as-a-package). Follow the next three steps to enable it and to
command the generator to execute the `elk` templates:

1. Create a new Go file named `ent/entc.go` and paste the following content:

```go
// +build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"log"
)

func main() {
	ex, err := elk.NewExtension()
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

```

2. Edit the `ent/generate.go` file to execute the `ent/entc.go` file:

```go
package ent

//go:generate go run -mod=mod entc.go

```

3. `elk` uses a not yet released version of Ent. To have the dependencies up to date run the following:

```shell
go mod tidy
```

With the aforementioned steps completed, all is set up for `elk`-powered ent! To learn more about Ent, how to connect to
different types of databases, run migrations or work with entities head over to
the [Setup Tutorial](https://entgo.io/docs/tutorial-setup/).

## Generating HTTP CRUD Handlers with `elk`

To generate the fully-working handlers we need first create an Ent schema definition. Open and edit `ent/schema/pet.go`:

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("age"),
	}
}

```

We added two fields to our `Pet` entity: `name` and `age`. The `ent.Schema` is just the definition of the schema. To
generate runnable code from it, we have to run Ent's code generation tool, (which we extended by `elk` in the previous
step).

```shell
go generate ./...
```

You should see a bunch of new files generated. In addition to the files Ent is generating for you there is another
folder `ent/http` where the code for the HTTP handlers resides. Below you can find an excerpt of the generated code for
a read-operation on the Pet entity:

[comment]: <> (// @formatter:off)
```go
const (
    PetCreate Routes = 1 << iota
    PetRead
    PetUpdate
    PetDelete
    PetList
    PetRoutes = 1<<iota - 1
)

// PetHandler handles http crud operations on ent.Pet.
type PetHandler struct {
    handler

    client    *ent.Client
    log       *zap.Logger
    validator *validator.Validate
}

func NewPetHandler(c *ent.Client, l *zap.Logger, v *validator.Validate) *PetHandler {
    return &PetHandler{
        client:    c,
        log:       l.With(zap.String("handler", "PetHandler")),
        validator: v,
    }
}

// Read fetches the ent.Pet identified by a given url-parameter from the
// database and renders it to the client.
func (h *PetHandler) Read(w http.ResponseWriter, r *http.Request) {
    l := h.log.With(zap.String("method", "Read"))
    // ID is URL parameter.
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
        render.BadRequest(w, r, "id must be an integer greater zero")
        return
    }
    // Create the query to fetch the Pet
    q := h.client.Pet.Query().Where(pet.ID(id))
    e, err := q.Only(r.Context())
    if err != nil {
        switch err.(type) {
        case ent.IsNotFound(err):
            msg := h.stripEntError(err)
            l.Info(msg, zap.Int("id", id), zap.Error(err))
            render.NotFound(w, r, msg)
        case ent.IsNotSingular(err):
            msg := h.stripEntError(err)
            l.Error(msg, zap.Int("id", id), zap.Error(err))
            render.BadRequest(w, r, msg)
        default:
            l.Error("error fetching pet from db", zap.Int("id", id), zap.Error(err))
            render.InternalServerError(w, r, nil)
        }
        return
    }
    d, err := sheriff.Marshal(&sheriff.Options{
        IncludeEmptyTag: true,
        Groups:          []string{"pet"},
    }, e)
    if err != nil {
        l.Error("serialization error", zap.Int("id", id), zap.Error(err))
        render.InternalServerError(w, r, nil)
        return
    }
    l.Info("pet rendered", zap.Int("id", id))
    render.OK(w, r, d)
}
```
[comment]: <> (// @formatter:on)

It is plain simple to create a running RESTful HTTP api to manage your pet entities. Create `main.go` and add the
following content:

```go
package main

import (
	"context"
	"elk-example/ent"
	elk "elk-example/ent/http"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	// Create the ent client.
	c, err := ent.Open("sqlite3", "./ent.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()
	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Router, Logger and Validator.
	r, l, v := chi.NewRouter(), zap.NewExample(), validator.New()
	// Create the pet handler.
	r.Route("/pets", func(r chi.Router) {
		elk.NewPetHandler(c, l, v).Mount(r, elk.PetRoutes)
	})
	// Start listen to incoming requests.
	fmt.Println("Server running")
	defer fmt.Println("Server stopped")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

```

Go ahead and start the server:

```shell
go run -mod=mod main.go
```

Congratulations, you now have a working Pets API. You could ask the server for a list of all pets, but there are none
yet. We have to create one first:

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Kuro","age":3}' 'localhost:8080/pets'
```

You should get this response:

```json
{
  "age": 3,
  "id": 1,
  "name": "Kuro"
}
```

If you head over to the terminal where the server is running you can also see `elk`s built in logging:

```json
{
  "level": "info",
  "msg": "pet rendered",
  "handler": "PetHandler",
  "method": "Create",
  "id": 1
}
```

`elk` uses [zap](https://github.com/uber-go/zap) for logging. To learn more about it have a look at its documentation.

## Relations

To illustrate more of `elk`s features, let’s extend our graph: edit `ent/schema/user.go` and `ent/schema/pet.go`:

`ent/schema/pet.go`:

[comment]: <> (// @formatter:off)
```go
// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("owner", User.Type).
            Ref("pets").
            Unique(),
    }
}

```
[comment]: <> (// @formatter:on)

`ent/schema/user.go`:

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("age"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}

```

The above creates a One-To-Many relation between the Pet and User schemas: A pet belongs to a user, and a user can have
multiple pets.

Rerun the code generator:

```shell
go generate ./...
```

Do not forget to register the `UserHandler` on our router. Just add the following lines to `main.go`:

[comment]: <> (// @formatter:off)
```go
[...]
    r.Route("/pets", func(r chi.Router) {
        elk.NewPetHandler(c, l, v).Mount(r, elk.PetRoutes)
    })
+    // Create the user handler.
+    r.Route("/users", func(r chi.Router) {
+        elk.NewUserHandler(c, l, v).Mount(r, elk.UserRoutes)
+    })
    // Start listen to incoming requests.
    fmt.Println("Server running")
[...]
```
[comment]: <> (// @formatter:on)

After restarting the server we can create a `User` that owns the previously created Pet named Kuro:

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Elk","age":30,"owner":1}' 'localhost:8080/users'
```

The server returns the following response:

```json
{
  "age": 30,
  "edges": {},
  "id": 1,
  "name": "Elk"
}
```

You can see the user has been created, but the edges are empty. `elk` does not include edges in its output by default.
You can configure `elk` to render edges using serialization groups. Annotate your schemas with
the `elk.SchemaAnnotation`
and `elk.Annotation` structs. Edit `ent/schema/user.go` and add those:

[comment]: <> (// @formatter:off)
```go
// Edges of the User.
func (User) Edges() []ent.Edge {
    return []ent.Edge{
    	edge.To("pets", Pet.Type).
    		Annotations(elk.Groups("user")),
    }
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
    return []schema.Annotation{elk.ReadGroups("user")}
}

```
[comment]: <> (// @formatter:on)

The `elk.Annotation`s added to the fields and edges command elk to serialize those if the "user" group is requested.
The `elk.SchemaAnnotation` is used to make the read-operation of the `UserHandler` request "user". Note, that all fields
that do not have a serialization group attached are included by default. Edges, however, are excluded, if not stated
otherwise.

Once again regenerate the code and restart the server. You should now see the pets of a user rendered if you read a
resource:

```shell
curl 'localhost:8080/users/1'
```

```json
{
  "age": 30,
  "edges": {
    "pets": [
      {
        "id": 1,
        "name": "Kuro",
        "age": 3,
        "edges": {}
      }
    ]
  },
  "id": 1,
  "name": "Elk"
}
```

## Request validation

Our current schemas allow to set a negative age for pets or users and we can create pets without an owner (as we did
with Kuro). Ent has [builtin validation](https://entgo.io/docs/schema-fields#validators) but if you'd want to validate
requests made against your api you need validation to happen even before ent is involved. `elk` does
use [this](https://github.com/go-playground/validator) package to define validation rules and validate data. We can
create different validation rules for create and update operations using `elk.Annotation`. We want our pet schema to
only allow ages greater zero and don't allow a pet without an owner. Edit `ent/schema/pet.go`:

[comment]: <> (// @formatter:off)
```go
// Fields of the Pet.
func (Pet) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.Int("age").
            Positive().
            Annotations(elk.Validation("required,gt=0")),
    }
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("owner", User.Type).
            Ref("pets").
            Unique().
        	Required().
			Annotations(elk.Validation("required")),
    }
}
```
[comment]: <> (// @formatter:on)

After regenerating the code and restart the server try to create a pet with invalid age and without an owner:

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Bob","age":-2}' 'localhost:8080/pets'
```

`elk` returns a meaningful response telling you exactly what was wrong with the given request data:

```json
{
  "code": 400,
  "status": "Bad Request",
  "errors": {
    "Age": "This value failed validation on 'gt:0'.",
    "Owner": "This value is required."
  }
}
```

Note the uppercase field names. The validator package uses the structs field name to generate its validation errors, but
you can simply override this, as stated in
the [example](https://github.com/go-playground/validator/blob/9a5bce32538f319bf69aebb3aca90d394bc6d0cb/_examples/struct-level/main.go#L37)
.

## Upcoming Features

`elk` has pretty nice features already, but there are still things to come. The below features are just some of many new
features you can expect in the near future:

- Fully working flutter frontend to administrate your nodes
- Integration of Ent’s validation in the current request validator
- More transport formats (currently only json)

## Conclusion

This post has shown just a small part of what `elk` can do. It is a powerful extension to Ent. With `elk`-powered Ent
you and your devs can focus on much more meaningful work, and you can automate stuff that otherwise would consume
precious resources and time.

Yet `elk` is in its very early stage of development and cannot be considered stable. Any suggestion and feedback is very
welcome and if you are willing to help I'd be very glad. The [GitHub Issues](https://github.com/masseelch/elk/issues) is
a wonderful place for you to reach out for help, feedback, suggestions and contribution.
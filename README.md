Egon
===
**Note: This is a work in progress.**

Egon is a templating language for go, based on [Ego](https://github.com/benbjohnson/ego).
Egon parses .egon templates and converts them into Go source files.

## Differences with Ego

* Ego generates a single source file for every template in a package. Egon
  generates a source file per template.
* Ego includes options to name the output package and output files. Egon always
  determines these names based on the source structure.
* Ego generates a single function per template for rendering the template. Egon
  includes a second function, [Template]View that returns a function wrapping
  the typical parameters that only requires an io.Writer.
* Ego requires a full function definition, Egon only requires parameter declarations.

## Usage

To install egon:

```sh
$ go get github.com/commondream/egon/cmd/egon
```

Running the `egon` command will process all templates for the package name
given.

```sh
$ egon github.com/commondream/templates
```

All egon files found in a package are converted to .egon.go files. Each .egon.go
file defines two functions:

1. The Template function.
2. The View function.


## Language Definition

An ego template is made up of several types of blocks:

* **Code Block** - These blocks execute raw Go code: `<% var foo = "bar" %>`

* **Print Block** - These blocks print a Go expression. They use `html.EscapeString` to escape it before outputting: `<%= myVar %>`

* **Raw Print Block** - These blocks print a Go expression raw into the HTML: `<%== "<script>" %>`

* **Header Block** - These blocks allow you to import packages: `<%% import "encoding/json" %%>`

* **Parameter Block** - This block defines the function signature for your template.

A single declaration block should exist at the top of your template and accept an `w io.Writer` and return an `error`. Other arguments can be added as needed. A function receiver can also be used.

```
<%! name string %>
```


## Example

Below is an example ego template for a web page:

```ego
<%% import "strings" %%>
<%! u *User %>

<html>
  <body>
    <h1>Hello <%= strings.TrimSpace(u.FirstName) %>!</h1>

    <p>Here's a list of your favorite colors:</p>
    <ul>
      <% for _, colorName := range u.FavoriteColors { %>
        <li><%= colorName %></li>
      <% } %>
    </ul>
  </body>
</html>
```

Once this template is compiled you can call it using the parameters you specified:

```go
myUser := &User{
  FirstName: "Bob",
  FavoriteColors: []string{"blue", "green", "mauve"},
}
var buf bytes.Buffer
err := mypkg.MyTmpl(&buf, myUser)
```


## Caveats

Unlike other runtime-based templating languages, Egon does not support ad hoc
templates. All templates must be generated before compile time.

Egon does not attempt to provide any security around the templates. Just like
regular Go code, the security model is up to you.

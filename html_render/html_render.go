package HTMLRender

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ContentHTML header value for HTML data.
	ContentHTML = "text/html"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentLength header constant.
	ContentLength = "Content-Length"
	// ContentText header value for Text data.
	ContentText = "text/plain"
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentXHTML header value for XHTML data.
	ContentXHTML = "application/xhtml+xml"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// Default character encoding.
	defaultCharset = "UTF-8"
)

// Delims represents a set of Left and Right delimiters for HTML template rendering.
type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

// Options is a struct for specifying configuration options for the render.Render object.
type Options struct {
	// Directory to load templates. Default is "templates".
	Directory string
	// Asset function to use in place of directory. Defaults to nil.
	Asset func(name string) ([]byte, error)
	// AssetNames function to use in place of directory. Defaults to nil.
	AssetNames func() []string
	// Layout template name. Will not render a layout if blank (""). Defaults to blank ("").
	Layout string
	// Extensions to parse template files from. Defaults to [".tmpl"].
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs []template.FuncMap
	// Delims sets the action delimiters to the specified strings in the Delims struct.
	Delims Delims
	// Appends the given character set to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON.
	IndentJSON bool
	// Outputs human readable XML. Default is false.
	IndentXML bool
	// Prefixes the JSON output with the given bytes. Default is false.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html".
	HTMLContentType string
	// If IsDevelopment is set to true, this will recompile the templates on every request. Default is false.
	IsDevelopment bool
	// Unescape HTML characters "&<>" to their original values. Default is false.
	UnEscapeHTML bool
	// Streams JSON responses instead of marshalling prior to sending. Default is false.
	StreamingJSON bool
	// Require that all partials executed in the layout are implemented in all templates using the layout. Default is false.
	RequirePartials bool
	// Deprecated: Use the above `RequirePartials` instead of this. As of Go 1.6, blocks are built in. Default is false.
	RequireBlocks bool
	// Disables automatic rendering of http.StatusInternalServerError when an error occurs. Default is false.
	DisableHTTPErrorRendering bool
}

// HTMLOptions is a struct for overriding some rendering Options for specific HTML call.
type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

// Render is a service that provides functions for easily writing JSON, XML,
// binary data, and HTML templates out to a HTTP Response.
type Render struct {
	// Customize Secure with an Options struct.
	opt             Options
	templates       *template.Template
	compiledCharset string
}

// New constructs a new Render instance with the supplied options.
func New(options ...Options) *Render {
	var o Options
	if len(options) == 0 {
		o = Options{}
	} else {
		o = options[0]
	}

	r := Render{
		opt: o,
	}

	r.prepareOptions()
	r.compileTemplates()

	// Create a new buffer pool for writing templates into.
	if bufPool == nil {
		bufPool = NewBufferPool(64)
	}

	return &r
}

func (r *Render) prepareOptions() {
	// Fill in the defaults if need be.
	if len(r.opt.Charset) == 0 {
		r.opt.Charset = defaultCharset
	}
	r.compiledCharset = "; charset=" + r.opt.Charset

	if len(r.opt.Directory) == 0 {
		r.opt.Directory = "templates"
	}
	if len(r.opt.Extensions) == 0 {
		r.opt.Extensions = []string{".tmpl"}
	}
	if len(r.opt.HTMLContentType) == 0 {
		r.opt.HTMLContentType = ContentHTML
	}

	r.opt.IsDevelopment = (gin.Mode() != gin.ReleaseMode)
}

func (r *Render) compileTemplates() {
	if r.opt.Asset == nil || r.opt.AssetNames == nil {
		r.compileTemplatesFromDir()
		return
	}
	r.compileTemplatesFromAsset()
}

func (r *Render) compileTemplatesFromDir() {
	dir := r.opt.Directory
	r.templates = template.New(dir)
	r.templates.Delims(r.opt.Delims.Left, r.opt.Delims.Right)

	// Walk the supplied directory and compile any files that match our extension list.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Fix same-extension-dirs bug: some dir might be named to: "users.tmpl", "local.html".
		// These dirs should be excluded as they are not valid golang templates, but files under
		// them should be treat as normal.
		// If is a dir, return immediately (dir is not a valid golang template).
		if info == nil || info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := ""
		if strings.Index(rel, ".") != -1 {
			ext = filepath.Ext(rel)
		}

		for _, extension := range r.opt.Extensions {
			if ext == extension {
				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (rel[0 : len(rel)-len(ext)])
				tmpl := r.templates.New(filepath.ToSlash(name))

				// Add our funcmaps.
				for _, funcs := range r.opt.Funcs {
					tmpl.Funcs(funcs)
				}

				// Break out if this parsing fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}
		return nil
	})
}

func (r *Render) compileTemplatesFromAsset() {
	dir := r.opt.Directory
	r.templates = template.New(dir)
	r.templates.Delims(r.opt.Delims.Left, r.opt.Delims.Right)

	for _, path := range r.opt.AssetNames() {
		if !strings.HasPrefix(path, dir) {
			continue
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		ext := ""
		if strings.Index(rel, ".") != -1 {
			ext = "." + strings.Join(strings.Split(rel, ".")[1:], ".")
		}

		for _, extension := range r.opt.Extensions {
			if ext == extension {

				buf, err := r.opt.Asset(path)
				if err != nil {
					panic(err)
				}

				name := (rel[0 : len(rel)-len(ext)])
				tmpl := r.templates.New(filepath.ToSlash(name))

				// Add our funcmaps.
				for _, funcs := range r.opt.Funcs {
					tmpl.Funcs(funcs)
				}

				// Break out if this parsing fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}
	}
}

// TemplateLookup is a wrapper around template.Lookup and returns
// the template with the given name that is associated with t, or nil
// if there is no such template.
func (r *Render) TemplateLookup(t string) *template.Template {
	return r.templates.Lookup(t)
}

func (r *Render) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	return buf, r.templates.ExecuteTemplate(buf, name, binding)
}

func (r *Render) addLayoutFuncs(name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(name, binding)
			// Return safe HTML here since we are rendering our own template.
			return template.HTML(buf.String()), err
		},
		"current": func() (string, error) {
			return name, nil
		},
		"block": func(partialName string) (template.HTML, error) {
			log.Print("Render's `block` implementation is now depericated. Use `partial` as a drop in replacement.")
			fullPartialName := fmt.Sprintf("%s-%s", partialName, name)
			if r.opt.RequireBlocks || r.TemplateLookup(fullPartialName) != nil {
				buf, err := r.execute(fullPartialName, binding)
				// Return safe HTML here since we are rendering our own template.
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
		"partial": func(partialName string) (template.HTML, error) {
			fullPartialName := fmt.Sprintf("%s-%s", partialName, name)
			if r.opt.RequirePartials || r.TemplateLookup(fullPartialName) != nil {
				buf, err := r.execute(fullPartialName, binding)
				// Return safe HTML here since we are rendering our own template.
				return template.HTML(buf.String()), err
			}
			return "", nil
		},
	}
	if tpl := r.templates.Lookup(name); tpl != nil {
		tpl.Funcs(funcs)
	}
}

func (r *Render) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.opt.Layout,
	}
}

func (r *Render) Instance(name string, data interface{}) render.Render {
	// If we are in development mode, recompile the templates on every HTML request.
	if r.opt.IsDevelopment {
		r.compileTemplates()
	}

	if len(r.opt.Layout) > 0 {
		r.addLayoutFuncs(name, data)
		name = r.opt.Layout
	}

	return render.HTML{
		Template: r.templates,
		Name:     name,
		Data:     data,
	}

}

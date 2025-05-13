package handlers

import (
	"dockerregistryUI/utils"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var validColor = regexp.MustCompile(`#(?:\d|[a-f]){6}`)

// HandlerContext A context for the handlers.
type HandlerContext struct {
	initialized bool
	settings    utils.DockerRegistryUISettings
	client      *utils.RegistryHTTPClient
	cache       UITemplateCache
}

// New Initializes a new HandlerContext.
func New(settings utils.DockerRegistryUISettings, client *utils.RegistryHTTPClient) *HandlerContext {
	return &HandlerContext{
		initialized: true,
		settings:    settings,
		client:      client,
		cache:       newInMemoryUITemplateCache(),
	}
}

// IndexHandler Handles requests to the main index page.
func (context *HandlerContext) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var templateData UITemplateData
	var hasCache bool
	if templateData, hasCache = context.cache.GetCached(); !hasCache {
		templateData = InitializeUITemplateData(context.settings, context.client)
		context.cache.Cache(templateData)
	} else if up := RefreshUITemplateDataIfNecessary(context.settings, context.client, &templateData); up {
		context.cache.Flush()
		context.cache.Cache(templateData)
	}
	categoryQuery := r.URL.Query().Get("category")
	if categoryID, err := strconv.ParseUint(categoryQuery, 10, 64); err == nil {
		templateData.FilterImages(uint(categoryID))
	}
	setCommonHeaders(w)
	err := templates.ExecuteTemplate(w, "index.gohtml", templateData)
	if err != nil {
		log.Printf("Error rendering template: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateCategoryHandler Creates a new category from a form using "name" and "color".
func (context *HandlerContext) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	name := r.PostFormValue("name")
	if len(name) > 0 {
		context.cache.Flush()
	}
	context.RootRedirectHandler(w, r)
}

func escapeColor(color string) string {
	if len(color) != 7 {
		return ""
	} else if validColor.MatchString(color) {
		return color
	}
	return ""
}

// RemoveCategoryHandler Remvoes a category from a POST using "id".
func (context *HandlerContext) RemoveCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	id := r.PostFormValue("id")
	if len(id) > 0 {
		if _, err := strconv.ParseUint(id, 10, 64); err == nil {
			context.cache.Flush()
		} else {
			log.Printf("Cannot delete image category with invalid id: %s\n", id)
		}
	}
	context.RootRedirectHandler(w, r)
}

// CreateDescriptionHandler Creates a description from a form using "imageName", "description", "exampleCommand".
func (context *HandlerContext) CreateDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	imageName := r.PostFormValue("imageName")
	if len(imageName) > 0 {
		context.cache.Flush()
	}
	context.RootRedirectHandler(w, r)
}

// AddCategoryToDescriptionHandler Adds a category to a descriptio from a POST using "category" and "image".
func (context *HandlerContext) AddCategoryToDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	categoryID := r.PostFormValue("category")
	descriptionID := r.PostFormValue("image")
	if len(categoryID) > 0 && len(descriptionID) > 0 {
		_, categoryErr := strconv.ParseUint(categoryID, 10, 64)
		_, descriptionErr := strconv.ParseUint(descriptionID, 10, 64)
		if categoryErr == nil && descriptionErr == nil {
			context.cache.Flush()
		}
	}
	context.RootRedirectHandler(w, r)
}

// RemoveCategoryFromDescriptionHandler Removes a category from a description from a POST using "category" and "image".
func (context *HandlerContext) RemoveCategoryFromDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	categoryID := r.PostFormValue("category")
	descriptionID := r.PostFormValue("image")
	if len(categoryID) > 0 && len(descriptionID) > 0 {
		_, categoryErr := strconv.ParseUint(categoryID, 10, 64)
		_, descriptionErr := strconv.ParseUint(descriptionID, 10, 64)
		if categoryErr == nil && descriptionErr == nil {
			context.cache.Flush()
		}
	}
	context.RootRedirectHandler(w, r)
}

// EditHelloHandler Edit the hello message from a POST using "hello".
func (context *HandlerContext) EditHelloHandler(w http.ResponseWriter, r *http.Request) {
	if !checkForPostWithError(w, r) {
		return
	}
	r.ParseForm()
	hello := r.PostFormValue("hello")
	if len(hello) > 0 {
		context.cache.Flush()
	}
	context.RootRedirectHandler(w, r)
}

// RootRedirectHandler Redirects to the index page.
func (context *HandlerContext) RootRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, context.settings.ContextRoot+"/", 302)
}

func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
}

func checkForPostWithError(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "POST" {
		return true
	}
	http.NotFound(w, r)
	return false
}

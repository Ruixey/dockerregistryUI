package main

import (
	"dockerregistryUI/handlers"
	"dockerregistryUI/utils"
	"log"
	"net/http"
)

func main() {
	settings := utils.SettingsFromEnvironmentVariables()
	client := utils.NewRegistryHTTPClient(settings)
	context := handlers.New(settings, client)
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle(settings.URIStaticDir, http.StripPrefix(settings.URIStaticDir, fileServer))
	http.HandleFunc(settings.URICreateCategory, context.CreateCategoryHandler)
	http.HandleFunc(settings.URIRemoveCategory, context.RemoveCategoryHandler)
	http.HandleFunc(settings.URIImageDescription, context.CreateDescriptionHandler)
	http.HandleFunc(settings.URIAddCategoryToImage, context.AddCategoryToDescriptionHandler)
	http.HandleFunc(settings.URIRemoveCategoryFromImage, context.RemoveCategoryFromDescriptionHandler)
	http.HandleFunc(settings.URIHello, context.EditHelloHandler)
	http.HandleFunc(settings.ContextRoot+"/", context.IndexHandler)
	http.HandleFunc(settings.ContextRoot, context.RootRedirectHandler)
	http.HandleFunc("/", context.RootRedirectHandler)
	log.Println("Started Docker Registry UI")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

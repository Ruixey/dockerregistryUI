package handlers

import (
	"dockerregistryUI/persistence"
	"dockerregistryUI/utils"
	"html/template"
	"regexp"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

var templates = template.Must(template.ParseFiles("templates/index.gohtml"))
var autolinkRegex = regexp.MustCompile(`(?:^|\s)(?P<link>https??://\S+)(?:$|\s)`)
var markdownRenderer = markdown.New()

// UITemplateData Data passed to the HTML template.
type UITemplateData struct {
	Settings utils.DockerRegistryUISettings
	Images   []ImageData
}

// ImageData Data about an image, collected from registry and database.
type ImageData struct {
	Name          string
	Tags          []string
	FormattedTags string
}

// UITemplateCache Caches the template's data until the next database change.*/
type UITemplateCache interface {
	Cache(data UITemplateData)
	Flush()
	GetCached() (UITemplateData, bool)
}

// InMemoryUITemplateCache Simple cache that holds a template's data in local memory until flushed.
type inMemoryUITemplateCache struct {
	cached     UITemplateData
	cleanCache bool
}

// MergeAndFormatImageData Merges image data retreived from the registry with database data for use in the template.
// Also applies markdown processing to stored MD user inputs.
func MergeAndFormatImageData(image utils.RegistryImage, description *persistence.ImageDescription) ImageData {
	var data ImageData
	data.Name = image.ImageName
	data.Tags = image.ImageTags
	data.FormattedTags = formatTags(data.Tags)
	return data
}

func formatTags(tags []string) string {
	var sb strings.Builder
	for i, tag := range tags {
		sb.WriteString(tag)
		if i < len(tags)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// InitializeUITemplateData Initializes the data for a UI template, fetching it from registry and db and merging.
func InitializeUITemplateData(settings utils.DockerRegistryUISettings,
	client *utils.RegistryHTTPClient) UITemplateData {
	return *createUITemplateDataWithRegistryImages(settings, client.RetreiveRegistryImages())
}

// RefreshUITemplateDataIfNecessary Checks if newer information for the template is available.
// Updates und returns true if yes. Returns false if unchanged.
func RefreshUITemplateDataIfNecessary(settings utils.DockerRegistryUISettings,
	client *utils.RegistryHTTPClient, data *UITemplateData) bool {
	var imageMetaData []utils.DockerImageMetaData
	for _, image := range data.Images {
		imageMetaData = append(imageMetaData, &image)
	}
	if newerRegistryImages, useOldData := client.CheckUpToDateOrRetreiveRegistryImages(imageMetaData); !useOldData {
		(*data) = (*createUITemplateDataWithRegistryImages(settings, newerRegistryImages))
		return true
	}
	return false
}

func createUITemplateDataWithRegistryImages(settings utils.DockerRegistryUISettings,
	registryImages []utils.RegistryImage) *UITemplateData {
	data := UITemplateData{Settings: settings}
	for _, registryImage := range registryImages {
		var imageData ImageData
		imageData = MergeAndFormatImageData(registryImage, &persistence.ImageDescription{})
		data.Images = append(data.Images, imageData)
	}
	return &data
}

// Cache Caches the current UI template data in local memory.
func (cache *inMemoryUITemplateCache) Cache(data UITemplateData) {
	cache.cached = data
	cache.cleanCache = true
}

// Flush Flushes the current UI template data from local memory.
func (cache *inMemoryUITemplateCache) Flush() {
	cache.cleanCache = false
}

// GetCached Returns the cached element and a bool indicating if the cached element is clean (true if it is clean).
func (cache *inMemoryUITemplateCache) GetCached() (UITemplateData, bool) {
	return cache.cached, cache.cleanCache
}

func newInMemoryUITemplateCache() *inMemoryUITemplateCache {
	return &inMemoryUITemplateCache{
		cached:     UITemplateData{},
		cleanCache: false,
	}
}

// GetImageName Get the image name.
func (imageData *ImageData) GetImageName() string {
	return imageData.Name
}

// GetImageTags Get the image tags.
func (imageData *ImageData) GetImageTags() []string {
	return imageData.Tags
}

package controller

//Language configuration processor
type Language struct {
	content map[string]interface{}
	dir     string
	src     string
	status bool
}

//Initialize the language configuration processor
func (this *Language) Init(languageType string) bool {
	sep := GetPathSep()
	this.dir = "." + sep + "language" + sep
	this.src = this.dir + languageType + ".json"
	this.status = false
	if IsFile(this.src) == false {
		log.NewLog("The language configuration file does not exist.", nil)
		return false
	}
	this.content, err = LoadConfigFile(this.src)
	if err != nil {
		log.NewLog("The language configuration file could not be read properly", err)
		return false
	}
	this.status = true
	return true
}

//Get the language
func (this *Language) Get(name string) string {
	if this.status == false{
		return ""
	}
	if this.content[name] == "" {
		return ""
	}
	return this.content[name].(string)
}

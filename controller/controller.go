package controller

//glob error
var err error

//glob log
var log Log

//glob coll
var coll Coll

//Profile data
//config file : ./config/config.json
var configData map[string]interface{}

//controller class
type Controller struct {
	//database
	db Database
	//router
	router Router
	//Configure the directory
	configDir string
}

//Initialize the structure
func (this *Controller) Init() {
	//Configure the directory
	sep := GetPathSep()
	this.configDir = "." + sep + "config" + sep
	configSrc := this.configDir + "config.json"
	//Read the configuration data
	configData, err = LoadConfigFile(configSrc)
	if err != nil {
		log.SendFmtPrintln("Unable to read the configuration file, can not start the program.Config src : " + configSrc)
		return
	}
	if configData["database-type"] == nil {
		log.SendFmtPrintln("The content of the configuration file is incorrect. Please check again.")
		return
	}
	//Initialize the log
	log.init(configData["data-src"].(string), true, true, true, true, true, true)
	//Connect database
	err = this.db.Connect(configData["database-type"].(string), configData["database-dns"].(string))
	defer this.db.Close()
	if err != nil {
		log.NewLog("Unable to connect to the database.", err)
		return
	}
	//Initializes the coll object
	coll.init(&this.db,configData["data-src"].(string))
	//Start the server
	this.router.RunServer(&this.db)
}

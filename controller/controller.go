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
	if configData["server-local"] == nil || configData["language"] == nil || configData["data-src"] == nil {
		log.SendFmtPrintln("The content of the configuration file is incorrect. Please check again.")
		return
	}
	//Initialize the log
	log.init(configData["data-src"].(string) + GetPathSep() + "sys-log", true, true, true, true, true, true)
	//Connect database
	dbTemplateSrc := "config" + sep + "coll-mz-default.sqlite"
	dbDirSrc :=  configData["data-src"].(string) + sep + "database"
	dbSrc := dbDirSrc + sep + "coll-mz.sqlite"
	if IsFile(dbSrc) == false{
		err = CreateDir(dbDirSrc)
		if err != nil{
			log.NewLog("",err)
			return
		}
		b,err := CopyFile(dbTemplateSrc,dbSrc)
		if err != nil || b == false{
			log.NewLog("Unable to create the total database file.",err)
			return
		}
	}
	err = this.db.Connect("sqlite3",dbSrc)
	if err != nil {
		log.NewLog("Unable to connect to the database.", err)
		return
	}
	defer this.db.Close()
	//Initializes the coll object
	collDatabaseTemplateSrc := "config" + sep + "coll-default.sqlite"
	coll.init(&this.db,configData["data-src"].(string),collDatabaseTemplateSrc)
	go coll.AutoTask()
	//Start the server
	this.router.RunServer(&this.db)
}

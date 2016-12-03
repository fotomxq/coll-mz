package controller

//Collect local data
func (this *Coll) CollLocal() {
	//Gets the object
	thisChildren := &this.collList.local
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	//
	if thisChildren.status == false{
		return
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}

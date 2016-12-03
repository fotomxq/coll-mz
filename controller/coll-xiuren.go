package controller

//Collect xiuren data
func (this *Coll) CollXiuren() {
	//Gets the object
	thisChildren := &this.collList.xiuren
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

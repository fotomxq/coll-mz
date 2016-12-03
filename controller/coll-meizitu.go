package controller

//Collect Mzitu data
func (this *Coll) CollMeizitu() {
	//Gets the object
	thisChildren := &this.collList.meizitu
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
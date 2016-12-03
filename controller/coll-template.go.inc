package controller

//Collect `edit name`
func (this *Coll) CollName() { `edit name`
	//Gets the object
	thisChildren := &this.collList.Name `edit name`
	var collOperate CollOperate
	if this.CollStart(thisChildren,&collOperate) == false{
		return
	}
	// `edit`
	if thisChildren.status == false{
		return
	}
	//finish
	this.CollEnd(thisChildren,&collOperate)
}


export namespace models {
	
	export class ContactsList {
	    name: string;
	    uid: string;
	
	    static createFrom(source: any = {}) {
	        return new ContactsList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.uid = source["uid"];
	    }
	}

}


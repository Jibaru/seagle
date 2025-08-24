export namespace handlers {
	
	export class ConnectInput {
	    host: string;
	    port: number;
	    database: string;
	    username: string;
	    password: string;
	    sslmode: string;
	    connectionString: string;
	    useConnectionString: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConnectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.database = source["database"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.sslmode = source["sslmode"];
	        this.connectionString = source["connectionString"];
	        this.useConnectionString = source["useConnectionString"];
	    }
	}
	export class ConnectOutput {
	    success: boolean;
	    message?: string;
	    databases?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ConnectOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.databases = source["databases"];
	    }
	}
	export class GetTableColumnsInput {
	    database: string;
	    table: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTableColumnsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.database = source["database"];
	        this.table = source["table"];
	    }
	}
	export class GetTableColumnsOutput {
	    success: boolean;
	    message?: string;
	    columns?: services.TableColumn[];
	
	    static createFrom(source: any = {}) {
	        return new GetTableColumnsOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.columns = this.convertValues(source["columns"], services.TableColumn);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetTablesInput {
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTablesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.database = source["database"];
	    }
	}
	export class GetTablesOutput {
	    success: boolean;
	    message?: string;
	    tables?: string[];
	
	    static createFrom(source: any = {}) {
	        return new GetTablesOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.tables = source["tables"];
	    }
	}
	export class TestConnectionInput {
	    host: string;
	    port: number;
	    database: string;
	    username: string;
	    password: string;
	    sslmode: string;
	    connectionString: string;
	    useConnectionString: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TestConnectionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.database = source["database"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.sslmode = source["sslmode"];
	        this.connectionString = source["connectionString"];
	        this.useConnectionString = source["useConnectionString"];
	    }
	}

}

export namespace services {
	
	export class TableColumn {
	    name: string;
	    dataType: string;
	    isNullable: boolean;
	    defaultValue?: string;
	
	    static createFrom(source: any = {}) {
	        return new TableColumn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.dataType = source["dataType"];
	        this.isNullable = source["isNullable"];
	        this.defaultValue = source["defaultValue"];
	    }
	}

}


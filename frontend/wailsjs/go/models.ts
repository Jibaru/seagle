export namespace handlers {
	
	export class AnalyzeMetadataInput {
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeMetadataInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class AnalyzeMetadataOutput {
	    success: boolean;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeMetadataOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class ConnectByIDInput {
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectByIDInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class ConnectByIDOutput {
	    success: boolean;
	    message: string;
	    databases?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ConnectByIDOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.databases = source["databases"];
	    }
	}
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
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.databases = source["databases"];
	        this.id = source["id"];
	    }
	}
	export class DisconnectInput {
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new DisconnectInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class ExecuteQueryInput {
	    id: string;
	    database: string;
	    query: string;
	
	    static createFrom(source: any = {}) {
	        return new ExecuteQueryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.database = source["database"];
	        this.query = source["query"];
	    }
	}
	export class ExecuteQueryOutput {
	    success: boolean;
	    message?: string;
	    result?: types.QueryResult;
	
	    static createFrom(source: any = {}) {
	        return new ExecuteQueryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.result = this.convertValues(source["result"], types.QueryResult);
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
	export class GenerateQueryInput {
	    id: string;
	    database: string;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateQueryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.database = source["database"];
	        this.prompt = source["prompt"];
	    }
	}
	export class GenerateQueryOutput {
	    success: boolean;
	    message?: string;
	    result?: types.GenerateQueryResult;
	
	    static createFrom(source: any = {}) {
	        return new GenerateQueryOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.result = this.convertValues(source["result"], types.GenerateQueryResult);
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
	export class GetConfigOutput {
	    success: boolean;
	    message: string;
	    // Go type: struct { OpenAIAPIKey string "json:\"openAIAPIKey\"" }
	    Config: any;
	
	    static createFrom(source: any = {}) {
	        return new GetConfigOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.Config = this.convertValues(source["Config"], Object);
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
	export class GetTableColumnsInput {
	    id: string;
	    database: string;
	    table: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTableColumnsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.database = source["database"];
	        this.table = source["table"];
	    }
	}
	export class GetTableColumnsOutput {
	    success: boolean;
	    message?: string;
	    columns?: types.TableColumn[];
	
	    static createFrom(source: any = {}) {
	        return new GetTableColumnsOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.columns = this.convertValues(source["columns"], types.TableColumn);
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
	    id: string;
	    database: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTablesInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
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
	export class ListConnectionsOutput {
	    success: boolean;
	    message: string;
	    connections: types.ConnectionSummary[];
	
	    static createFrom(source: any = {}) {
	        return new ListConnectionsOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.connections = this.convertValues(source["connections"], types.ConnectionSummary);
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
	export class SetConfigInput {
	    openAIAPIKey: string;
	
	    static createFrom(source: any = {}) {
	        return new SetConfigInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.openAIAPIKey = source["openAIAPIKey"];
	    }
	}
	export class SetConfigOutput {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new SetConfigOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
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

export namespace types {
	
	export class ConnectionSummary {
	    id: string;
	    host: string;
	    port: number;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.host = source["host"];
	        this.port = source["port"];
	    }
	}
	export class GenerateQueryResult {
	    generatedQuery: string;
	    originalPrompt: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerateQueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generatedQuery = source["generatedQuery"];
	        this.originalPrompt = source["originalPrompt"];
	    }
	}
	export class QueryResult {
	    columns: string[];
	    rows: any[][];
	    rowsAffected: number;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.columns = source["columns"];
	        this.rows = source["rows"];
	        this.rowsAffected = source["rowsAffected"];
	        this.duration = source["duration"];
	    }
	}
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


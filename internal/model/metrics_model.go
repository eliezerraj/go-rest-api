package model

type ManagerInfo struct {
	App 		*ManagerInfoApp `json:"app"`
	Server     	Server     		`json:"servers"`
	Setup		Setup			`json:"setup_behaviour"`
	AwsEnv		AwsEnv			`json:"aws_env"`
	DatabaseRDS DatabaseRDS		`json:"database_rds"`
}

type ManagerInfoApp struct {
	Name 				string `json:"name"`
	Description 		string `json:"description"`
	Version 			string `json:"version"`
	OSPID				string `json:"os_pid"`
	IpAdress			string `json:"ip_adress"`
}

type ManagerHealthDiskSpace struct {
	Status    string `json:"status"`
	Total     uint64 `json:"total"`
	Free      uint64 `json:"free"`
	Threshold uint64 `json:"threshold"`
}

type ManagerHealth struct {
	Status    	string                 	`json:"status"`
	Liveness	bool					`json:"liveness"`
	Readiness	bool					`json:"readiness"`
	DB        	ManagerHealthDB        	`json:"db"`
	DiskSpace 	ManagerHealthDiskSpace 	`json:"diskSpace"`
}

type ManagerHealthDB struct {
	Status bool `json:"status"`
}

type Server struct {
	Port 			int `json:"port"`
	ReadTimeout		int `json:"readTimeout"`
	WriteTimeout	int `json:"writeTimeout"`
	IdleTimeout		int `json:"idleTimeout"`
	CtxTimeout		int `json:"ctxTimeout"`
}

type Setup struct {
    ResponseTime 		int 	`json:"response_time"`
    ResponseStatusCode  int 	`json:"response_status_code"`
	IsRandomTime		bool 	`json:"is_random_time"`
	Count				int 	`json:"count"`
	DatabaseType		string  `json:"databaseType"`
	Liveness			bool	`json:"liveness"`
	Readiness			bool	`json:"readiness"`
}

type AwsEnv struct {
    Aws_region 			string `json:"aws_region"`
    Aws_access_id  		string `json:"aws_access_id"`
	Aws_access_secret	string `json:"aws_access_secret"`
}

type DatabaseRDS struct {
    Host 				string `json:"host"`
    Port  				string `json:"port"`
	Schema				string `json:"schema"`
	DatabaseName		string `json:"databaseName"`
	User				string `json:"user"`
	Password			string `json:"password"`
	Db_timeout			int	`json:"db_timeout"`
	Postgres_Driver		string `json:"postgres_driver"`
}
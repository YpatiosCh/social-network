## How to run

### Initial Run
- [Install mkcert 🔗](#installation-instructions-for-mkcert)

- Run ```make mkcert-local``` to generate  the certificates. *You may need to install first.*
- Run postgres on docker ```make db-up```
- Initialize database ```make db-init``` (make sure docker is running prior to this)
- Populate database with mock data ```make db-populate``` (optional)   
- Run backend ```make run-backend```
- Run frontend ```make run-frontend``` (in a different terminal)
and visit https://localhost:8080


### To change the schema or reset populate
- Make sure postgress is running on docker   
```make db-up```
- Reset the database:    
```make db-reset```
- Populate database with mock data:   
```make db-populate```

### After Initial Run to run again on last state
#### Run in one go (one terminal)
- ```make run-all```
#### Or in different terminals
- Terminal 1:   
	- ```make db-up```
	- ```make db-psql``` to acces sql from command line
- Terminal 2:  
	- ```make run-backend```
- Terminal 3:  
	- ```make run-frontend```

and visit https://localhost:8080


### To access the db command line
- ```make db-psql```

## Installation instructions for mkcert:
#### Windows:
 		choco install mkcert
 		choco install nss -y 
#### MacOS:
 		brew install mkcert
 		brew install nss
#### Linux:
 		curl -LO "https://github.com/FiloSottile/mkcert/releases/latest/download/mkcert-v1.4.5-linux-amd64"
 		`chmod +x mkcert-v1.4.5-linux-amd64
 		sudo mv mkcert-v1.4.5-linux-amd64 /usr/local/bin/mkcert
 		sudo apt install libnss3-tools`
*Note that Firefox doesn't work well with the current certificate*
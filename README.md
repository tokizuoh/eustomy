# eustomy
  
![](./docs/eustomy.png)
  
## Docker
  
### Version

```bash
> docker --version
Docker version 19.03.12, build 48a66213fe

> docker-compose --version
docker-compose version 1.27.2, build 18f557f9
```
  
### Build
  
```bash
> docker-compose up --build -d
```
  
### Run
  
```bash
> docker-compose exec app go run main.go
```
## Flowchart
  
|  using                            |  update DB                     |
| --------------------------------- | ------------------------------ |
|  ![](./docs/flowchart_using.png)  |  ![](./docs/flowchart_db.png)  |
  
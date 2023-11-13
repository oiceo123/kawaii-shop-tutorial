<h1>Kawaii Shop</h1>
<p>Kawaii shop is a mini REST API e-commerce project that made by Golang</p>

<h2>🧰 Tools</h2>
<p>If you don't want to use Google Cloud, you don't need to use <strong>Gcloud CLI</strong>, <strong>Google Cloud Run</strong>, <strong>Google Cloud SQL</strong>, <strong>Google Cloud Storage</strong>, <strong>Google Cloud Container Registry</strong>, but you need to edit some code.</p>
<ul>
    <li><a href="https://www.docker.com/">Docker</a></li>
	<li><a href="https://code.visualstudio.com/">Vscode</a></li>
	<li><a href="https://dbeaver.io/">DBeaver</a></li>
    <li><a href="https://www.postman.com/">Postman</a></li>
    <li><a href="https://cloud.google.com/sdk/docs/install">Gcloud CLI</a></li>
	<li><a href="">Google Cloud Run</a></li>
	<li><a href="">Google Cloud SQL</a></li>
	<li><a href="">Google Cloud Storage</a></li>
	<li><a href="">Google Cloud Container Regitry</a></li>
</ul>

<h2>📦 Packages</h2>

```bash
go get github.com/gofiber/fiber/v2
go get github.com/joho/godotenv
go get github.com/jmoiron/sqlx
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get cloud.google.com/go/storage
```

<h2>Install the gcloud CLI</h2>
<a href="https://cloud.google.com/sdk/docs/install#windows">Download the Google Cloud CLI installer.</a>

<h2>Login to GCP</h2>

```bash
gcloud auth login
```

<h2>Set project</h2>

```bash
gcloud config set project PROJECT_ID
```

<h2>Check Configurations</h2>

```bash
gcloud config configurations list
```

<h2>Database Schema</h2>

<img alt="schema" src="./screenshots/schema.png"/>

<h2>🐋 Pull PostgreSQL from Docker</h2>

```bash
docker pull postgres:alpine
```

<h2>🐋 Start PostgreSQL on Docker</h2>

```bash
docker run --name kawaii_db_test -e POSTGRES_USER=kawaii -e POSTGRES_PASSWORD=123456 -e TZ=Asia/Bangkok -p 4444:5432 -d postgres:alpine
```

<h2>🐋 Execute a container</h2>

```bash
docker exec -it kawaii_db_test bash
psql -U kawaii
```

<h2>✏️ Create a new database</h2>

```bash
CREATE DATABASE kawaii_db_test;
\l
```

<h2>💣 Remove a database</h2>
<p>If the migration fails, delete the database and then create a new one and migrate again.</p>

```bash
DROP DATABASE kawaii_db_test;
\l
```

<h2>📝 Migration</h2>

```bash
# Migrate up
migrate -source file://path/to/migrations -database 'postgres://kawaii:123456@localhost:4444/kawaii_db_test?sslmode=disable' -verbose up

# Migrate down
migrate -source file://path/to/migrations -database 'postgres://kawaii:123456@localhost:4444/kawaii_db_test?sslmode=disable' -verbose down
```

<h2>Build a docker image</h2>

```bash
docker build -t asia.gcr.io/prject-id/container-bucket .
```

<h2>Error and Solution</h2>
<p>If an error like this occurs Follow the steps below.</p>
<p><img alt="error" src="./screenshots/error.JPG"/></p>
<p>1. open file %userprofile%\.docker\config.json</p>
<p>2. rename credsStore to credStore</p>

<h2>Enable GCR and Create service account</h2>
<a href="./doc/how-to-push-image-to-gcr.docx">Download document.</a>

<h2>Push a docker image to GCP</h2>

```bash
docker push asia.gcr.io/prject-id/container-bucket
```

<h2>Set time zone on Google Cloud SQL</h2>
<p>1. Choose database</p>
<p>2. Choose Edit</p>
<p>3. Choose Flags</p>
<p>4. ADD A DATABASE FLAG</p>
<p>5. Choose "timezone" variable</p>
<p>6. Set the variable to Asia/Bangkok.</p>

<h2>Postman</h2>
<ul>
    <li><a href="./kawaii-shop-tutorial.postman_collection.json">Collection</a></li>
    <li><a href="./kawaii-shop-tutorial-dev.postman_environment.json">Environment</a></li>
</ul>

<!-- <h2>In case you don't want to use Google Cloud storage, Please follow this step</h2>
<p><strong>***Don't forget to change a function that related along with files module in products and orders module</strong></p>

<ol>
<li>

<p>Add this to your config for IAppConfig in config.go</p>

```go
type IAppConfig interface {
	Host() string
	Port() int
    ...
}

...
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }
```
</li>
<li>

<p>Add this middleware handler in your middleware module in middlewaresHandler.go</p>

```go
type IMiddlewaresHandler interface {
	StreamingFile() fiber.Handler
    ...
}

...
func (h *middlewaresHandler) StreamingFile() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root: http.Dir("./assets/images"),
	})
}
```
</li>
<li>

<p>Declare a StrammingFile() in func (s *server) Start() {} in server.go</p>

```go
func (s *server) Start() {
    ...
	s.app.Use(middlewares.StreamingFile())
}
```
</li>
<li>

<p>Add this usecase in your files module in filesUsecase.go</p>

```go
func (u *filesUsecase) uploadToStorageWorker(ctx context.Context, jobs <-chan *files.FileReq, results chan<- *files.FileRes, errs chan<- error) {
	for job := range jobs {
		cotainer, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}
		b, err := ioutil.ReadAll(cotainer)
		if err != nil {
			errs <- err
			return
		}

		// Upload an object to storage
		dest := fmt.Sprintf("./assets/images/%s", job.Destination)
		if err := os.WriteFile(dest, b, 0777); err != nil {
			if err := os.MkdirAll("./assets/images/"+strings.Replace(job.Destination, job.FileName, "", 1), 0777); err != nil {
				errs <- fmt.Errorf("mkdir \"./assets/images/%s\" failed: %v", err, job.Destination)
				return
			}
			if err := os.WriteFile(dest, b, 0777); err != nil {
				errs <- fmt.Errorf("write file failed: %v", err)
				return
			}
		}

		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("http://%s:%d/%s", u.cfg.App().Host(), u.cfg.App().Port(), job.Destination),
			},
			destination: job.Destination,
		}

		errs <- nil
		results <- newFile.file
	}
}

func (u *filesUsecase) UploadToStorage(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *files.FileReq, len(req))
	resultsCh := make(chan *files.FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadToStorageWorker(ctx, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}
	return res, nil
}
```
</li>
<li>
<p>Change usecase function in UploadFiles handler in filesHandler.go</p>

```go
func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
    ...
	res, err := h.filesUsecase.UploadToStorage(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
}
```
</li>
<li>
<p>Add this usecase in your files module in filesUsecase.go</p>

```go
func (u *filesUsecase) deleteFromStorageFileWorkers(ctx context.Context, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		if err := os.Remove("./assets/images/" + job.Destination); err != nil {
			errs <- fmt.Errorf("remove file: %s failed: %v", job.Destination, err)
			return
		}
		errs <- nil
	}
}

func (u *filesUsecase) DeleteFileOnStorage(req []*files.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.deleteFromStorageFileWorkers(ctx, jobsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		return err
	}
	return nil
}
```
</li>
<li>

<p>Change usecase function in DeleteFile handler in filesHandler.go</p>

```go
func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	...
	if err := h.filesUsecase.DeleteFileOnStorage(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}
}
```
</li>
</ol> -->
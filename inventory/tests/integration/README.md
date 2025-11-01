# Запуск тестов локально
Если не хочется долго ждать сборки образа приложения, то можно предварительно его сбилдить
```bash
docker build \
  -f ./deploy/docker/inventory/Dockerfile \
  -t inventory-app:latest \
  .
```
> ⚠️ _Команду необходимо выполнять из корня проекта_

Далее необходимо в структуру `Config`, файл `app/app.go`, добавить поле `Image` и соответствующую _Optional-функцию_ в файле `app/opts.go`
```go
type Config struct {
	Name          string
	DockerfileDir string
	Dockerfile    string
	Image         string // <- Новое поле
	Port          string
	Env           map[string]string
	Networks      []string
	LogOutput     io.Writer
	StartupWait   wait.Strategy
	Logger        Logger
}

```

```go
// WithImage - добавляет имя образа для запуска контейнера
func WithImage(image string) Option {
	return func(c *Config) {
		c.Image = image
	}
}
```

Также в модуле `testcontainers`, файл `app/app.go` необходимо исправить функцию `New`
```go
req := testcontainers.ContainerRequest{
		Name: cfg.Name,
		Image:              cfg.Image, // <- Передаем имя образа из Config
		Networks:           cfg.Networks,
		Env:                cfg.Env,
		WaitingFor:         cfg.StartupWait,
		ExposedPorts:       []string{cfg.Port + "/tcp"},
		HostConfigModifier: DefaultHostConfig(),
	}
```


После этого в файле `setup.go` необходимо выставить передачу образа
```go
const inventoryImageName  = "inventory-app:latest"

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
        // Код до...
	appContainer, err := app.NewContainer(ctx,
            app.WithName(inventoryAppName),
            app.WithPort(grpcPort),
            app.WithImage(inventoryImageName), // <- Передаем имя нашего образа
            app.WithNetwork(generatedNetwork.Name()),
            app.WithEnv(appEnv),
            app.WithLogOutput(os.Stdout),
            app.WithStartupWait(waitStrategy),
            app.WithLogger(logger.Logger()),
        )
	// Код после...
}
```

И далее уже можно запускать командой
```bash
task integration:test:inventory
```

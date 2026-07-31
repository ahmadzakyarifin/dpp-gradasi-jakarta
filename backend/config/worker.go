package config

// WorkerConfig tetap ada walau sederhana — untuk background job email
type WorkerConfig struct {
	Concurrency int `env:"WORKER_CONCURRENCY" envDefault:"5"`
}

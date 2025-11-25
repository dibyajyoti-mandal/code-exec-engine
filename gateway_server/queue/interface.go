package queue

import "github.com/dibyajyoti-mandal/code-exec-engine/gateway/models"

//interface implemented by redis
type Publisher interface {
	Publish(job models.Job) error
}

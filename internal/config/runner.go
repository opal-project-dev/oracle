package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
)

type Runner interface {
	RunnerInterval() time.Duration
}

func NewRunnerInfo(getter kv.Getter) Runner {
	return &runner{
		getter: getter,
	}
}

type runner struct {
	getter kv.Getter
	once   comfig.Once
}

func (r *runner) RunnerInterval() time.Duration {
	var result time.Duration
	r.once.Do(func() interface{} {
		var info struct {
			Interval string `fig:"interval,required"`
		}

		err := figure.
			Out(&info).
			From(kv.MustGetStringMap(r.getter, "runner")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out runner"))
		}
		result, err = time.ParseDuration(info.Interval)
		return nil
	})
	return result
}

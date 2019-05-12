package splash

type job struct {
	f    func(...interface{}) error
	args []interface{}
}

func newJob(f func(...interface{}) error, args []interface{}) job {
	return job{f: f, args: args}
}

func (j *job) exec() error {
	return j.f(j.args...)
}

package gochan

type TaskFunc func() error

type gochan struct {
	uuid int
	taskChan chan TaskFunc
	dieChan chan struct{}
}

func newGochan(bufferNum int) *gochan {
	gc := &gochan{
		uuid:     defaultUUID(),
		taskChan: make(chan TaskFunc, bufferNum),
		dieChan:  make(chan struct{}),
	}

	return gc
}

func (gc *gochan) setUUID(uuid int) {
	gc.uuid = uuid
}

func (gc *gochan) run() {
	go gc.start()
}

func (gc *gochan) start() {
	defer func() {
		logger.Infof("gochan %d ending...", gc.uuid)
	}()

	logger.Infof("gochan %d starting...", gc.uuid)

	for {
		select {
		case task := <-gc.taskChan:
			err := task()
			if err != nil {
				logger.Errorf("task in gochan %d error: %s", gc.uuid, err.Error())
			}
		case <-gc.dieChan:
			return
		}
	}
}

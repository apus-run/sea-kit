package websocket

import "time"

func NewClientWithReconnect(opt *Option, stop chan<- struct{}) {
	defer func() {
		opt.Cancel()
		if stop != nil {
			stop <- struct{}{}
		}
	}()
	for {
		switch opt.Status {
		case OptionClosed:
			return
		case OptionActive, OptionInActive:
			client, err := NewClient(opt)
			if err != nil {
				if ok := opt.Next(); !ok {
					return
				}
				retryWaitInMillisecond(opt.RetryDuration)
				continue
			}
			opt.ChangeStatus(OptionActive)
			if err := opt.Prepare(); err != nil {
				return
			}
			select {
			case <-opt.Done():
				return
			case <-client.Context.Done():
				if ok := opt.Next(); !ok {
					return
				}
				retryWaitInMillisecond(opt.RetryDuration)
			}
		}
	}
}

func retryWaitInMillisecond(sleep int64) {
	if sleep == 0 {
		return
	}
	time.Sleep(time.Millisecond * time.Duration(sleep))
}

package gotelemetry

import (
	"time"
)

type batchStreamSubmission struct {
	submissionType BatchType
	tag            string
	data           interface{}
}

type BatchStream struct {
	C              chan batchStreamSubmission
	errorChannel   chan error
	credentials    Credentials
	channelTag     string
	control        chan bool
	updates        map[string]batchStreamSubmission
	updateInterval time.Duration
}

func NewBatchStream(credentials Credentials, channelTag string, submissionInterval time.Duration, errorChannel chan error) (*BatchStream, error) {
	if submissionInterval < time.Second {
		return nil, NewError(500, "Invalid submission interval (must be >= 1s)")
	}

	result := &BatchStream{
		C:              make(chan batchStreamSubmission, 10),
		errorChannel:   errorChannel,
		credentials:    credentials,
		channelTag:     channelTag,
		control:        make(chan bool, 0),
		updates:        map[string]batchStreamSubmission{},
		updateInterval: submissionInterval,
	}

	go result.handle()

	return result, nil
}

func (b *BatchStream) Send(f *Flow) {
	b.C <- batchStreamSubmission{
		submissionType: BatchTypePOST,
		tag:            f.Tag,
		data:           f.Data,
	}
}

func (b *BatchStream) SendData(tag string, data interface{}, submissionType BatchType) {
	b.C <- batchStreamSubmission{
		submissionType: submissionType,
		tag:            tag,
		data:           data,
	}
}

func (b *BatchStream) Stop() {
	b.control <- true
}

func (b *BatchStream) Flush() {
	b.sendUpdates()
	b.Stop()
}

func (b *BatchStream) handle() {
	t := time.After(b.updateInterval)

	for {
		select {
		case update := <-b.C:

			b.updates[update.tag] = update

		case <-b.control:

			b.sendUpdates()
			return

			return

		case <-t:

			b.sendUpdates()
			t = time.After(b.updateInterval)
		}

	}
}

func (b *BatchStream) sendUpdates() {
	if len(b.updates) == 0 {
		return
	}

	batches := map[BatchType]Batch{}

	for _, update := range b.updates {
		if _, ok := batches[update.submissionType]; !ok {
			batches[update.submissionType] = Batch{}
		}

		batches[update.submissionType].SetData(update.tag, update.data)
	}

	for submissionType, batch := range batches {
		err := batch.Publish(b.credentials, b.channelTag, submissionType)

		if err != nil && b.errorChannel != nil {
			b.errorChannel <- err
		}
	}

	b.updates = map[string]batchStreamSubmission{}
}

package gotelemetry

import (
	"net/http"
	"time"
)

type batchStreamSubmission struct {
	submissionType BatchType
	tag            string
	data           interface{}
}

// BatchStream struct
type BatchStream struct {
	C              chan batchStreamSubmission
	credentials    Credentials
	channelTag     string
	control        chan bool
	updates        map[string]batchStreamSubmission
	updateInterval time.Duration
}

// NewBatchStream function
func NewBatchStream(credentials Credentials, channelTag string, submissionInterval time.Duration, disableKeepAlives, disableCompression bool) (*BatchStream, error) {
	if submissionInterval < time.Second {
		return nil, NewError(500, "Invalid submission interval (must be >= 1s)")
	}

	// TODO: the client should belong to BatchStream, however, this is impossible
	//       in the current implementation, we need to work with one global
	//       static value.
	client = &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     disableKeepAlives,
			DisableCompression:    disableCompression,
		},
		Timeout: 30 * time.Second,
	}

	result := &BatchStream{
		C:              make(chan batchStreamSubmission, 10),
		credentials:    credentials,
		channelTag:     channelTag,
		control:        make(chan bool, 0),
		updates:        map[string]batchStreamSubmission{},
		updateInterval: submissionInterval,
	}

	go result.handle()

	return result, nil
}

// Send function
func (b *BatchStream) Send(f *Flow) {
	b.C <- batchStreamSubmission{
		submissionType: BatchTypePOST,
		tag:            f.Tag,
		data:           f.Data,
	}
}

// SendData function
func (b *BatchStream) SendData(tag string, data interface{}, submissionType BatchType) {
	b.C <- batchStreamSubmission{
		submissionType: submissionType,
		tag:            tag,
		data:           data,
	}
}

// Stop function
func (b *BatchStream) Stop() {
	b.control <- true
}

// Flush function
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
		if err != nil {
			logger.Error("failed to publish batch", "error", err.Error())
		}
	}

	b.updates = map[string]batchStreamSubmission{}
}

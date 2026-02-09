package kafka

import "context"

type ProducerMock struct {
	SendMessageFunc func(ctx context.Context, msg Message) error
}

func (m *ProducerMock) SendMessage(ctx context.Context, msg Message) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, msg)
	}
	return nil
}

func (m *ProducerMock) Close() error {
	return nil
}

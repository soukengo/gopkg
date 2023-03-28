package log

import "context"

type proxy struct {
	logger Logger
}

func newProxy(logger Logger) *proxy {
	return &proxy{logger: logger}
}

func (p *proxy) Log(level Level, paris Pairs) {
	p.logger.Log(level, paris)
}

func (p *proxy) With(paris Pairs) Logger {
	return p.logger.With(paris)
}

func (p *proxy) WithContext(ctx context.Context) Logger {
	return p.logger.WithContext(ctx)
}

func (p *proxy) Helper() *Helper {
	return p.logger.Helper()
}

func (p *proxy) HelperWith(paris Pairs) *Helper {
	return p.logger.HelperWith(paris)
}

func (p *proxy) Factory() Factory {
	return p.logger.Factory()
}

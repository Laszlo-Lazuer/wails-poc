package update

type Bindings struct {
	updater *Updater
}

func (b *Bindings) Update() error {
	return b.updater.update()
}

func (b *Bindings) Restart() error {
	return b.updater.restart()
}

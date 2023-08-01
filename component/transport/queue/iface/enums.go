package iface

type Action int

const (
	CommitMessage  Action = 1
	ReconsumeLater Action = 2
)

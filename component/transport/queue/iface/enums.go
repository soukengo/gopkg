package iface

type Action int

const (
	None           Action = 0
	CommitMessage  Action = 1
	ReconsumeLater Action = 2
)

package model

var (
	Normal   = 1
	Personal = 1
)

var AESkey = "qwertyuiopasdfghjklzxcvb"

const (
	NoDeleted = iota
	Deleted
)
const (
	NoArcheve = iota
	Archeve
)

const (
	Open = iota
	private
	custom
)
const (
	Default = "default"
	Simple  = "simple"
)

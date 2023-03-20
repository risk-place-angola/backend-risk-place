package drive

type IDrive interface {
	Connect() string
}

type DnsDrive struct {}

func NewDnsDrive() *DnsDrive {
	return &DnsDrive{}
}

func Drive(drive IDrive) string {
	return drive.Connect()
}

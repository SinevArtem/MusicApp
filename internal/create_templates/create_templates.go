package createtemplates

// 	"html/template"

type exeptionOnRegOrLog struct {
	Exeption string
}

func GetExeptionOnRegister(s string) *exeptionOnRegOrLog {

	exeptionOnRegOrLog := &exeptionOnRegOrLog{Exeption: s}
	return exeptionOnRegOrLog
}

type AddTrack struct {
	Exeption     string
	Notification string
}

func TemplAddTrack(ex, notification string) *AddTrack {

	addtrack := &AddTrack{Exeption: ex, Notification: notification}
	return addtrack
}

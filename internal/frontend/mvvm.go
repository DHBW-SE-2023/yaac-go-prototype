package yaac_frontend

type mvvm interface {
	MailFormUpdated(data EmailData)
}

type Frontend struct {
	MVVM mvvm
}

func New(mvvm mvvm) *Frontend {
	return &Frontend{
		MVVM: mvvm,
	}
}

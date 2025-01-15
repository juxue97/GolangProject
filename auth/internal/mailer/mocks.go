package mailer

type MockMailerService struct{}

func NewMockMailerClient() *Client {
	return &Client{
		MailTrapService: &MockMailerService{},
	}
}

func (m *MockMailerService) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
	// return 200 for successful send
	return 200, nil
}

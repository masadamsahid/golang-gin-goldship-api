package xenditService

import (
	"os"

	"github.com/xendit/xendit-go/v7"
)

var Client *xendit.APIClient
var XENDIT_WEBHOOK_VERIFICATION_TOKEN string

func InitXendit() {
	API_KEY := os.Getenv("XENDIT_SECRET_API_KEY")
	Client = xendit.NewClient(API_KEY)

	XENDIT_WEBHOOK_VERIFICATION_TOKEN = os.Getenv("XENDIT_WEBHOOK_VERIFICATION_TOKEN")
}

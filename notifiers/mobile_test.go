package notifiers

import (
	"github.com/nordcloud/statping-ng/database"
	"github.com/nordcloud/statping-ng/types/core"
	"github.com/nordcloud/statping-ng/types/failures"
	"github.com/nordcloud/statping-ng/types/notifications"
	"github.com/nordcloud/statping-ng/types/null"
	"github.com/nordcloud/statping-ng/types/services"
	"github.com/nordcloud/statping-ng/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	mobileToken string
)

func TestMobileNotifier(t *testing.T) {
	err := utils.InitLogs()
	require.Nil(t, err)

	t.Parallel()

	mobileToken = utils.Params.GetString("MOBILE_TOKEN")
	if mobileToken == "" {
		t.Log("Mobile notifier testing skipped, missing MOBILE_ID environment variable")
		t.SkipNow()
	}

	Mobile.Var1 = null.NewNullString(mobileToken)

	db, err := database.OpenTester()
	require.Nil(t, err)
	db.AutoMigrate(&notifications.Notification{})
	notifications.SetDB(db)
	core.Example()

	t.Run("Load Mobile", func(t *testing.T) {
		Mobile.Var1 = null.NewNullString(mobileToken)
		Mobile.Delay = time.Duration(100 * time.Millisecond)
		Mobile.Limits = 10
		Mobile.Enabled = null.NewNullBool(true)

		Add(Mobile)

		assert.Equal(t, "Hunter Long", Mobile.Author)
		assert.Equal(t, mobileToken, Mobile.Var1.String)
	})

	t.Run("Mobile Notifier Tester", func(t *testing.T) {
		assert.True(t, Mobile.CanSend())
	})

	t.Run("Mobile OnSave", func(t *testing.T) {
		_, err := Mobile.OnSave()
		assert.Nil(t, err)
	})

	t.Run("Mobile OnFailure", func(t *testing.T) {
		_, err := Mobile.OnFailure(services.Example(false), failures.Example())
		assert.Nil(t, err)
	})

	t.Run("Mobile OnSuccess", func(t *testing.T) {
		_, err := Mobile.OnSuccess(services.Example(true))
		assert.Nil(t, err)
	})

	t.Run("Mobile Test", func(t *testing.T) {
		_, err := Mobile.OnTest()
		assert.Nil(t, err)
	})

}

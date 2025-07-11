package notification

import (
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"reflect"
	"sync"
	"time"
)

// ========================================================================================
//                                 Notification Structure
// ========================================================================================

// INotifiable is the interface that defines a notification that can be sent
// through different channels like Mail, SMS, Slack, etc.
// Each notification handler must implement this interface.
type INotifiable interface {
	// Notify sends the actual notification through the specific channel
	Notify()
}

// fnHandler is a factory function type that creates a notification handler
// implementing INotifiable interface for a given notification type.
type fnHandler func(notification any) INotifiable

// handlerInfo contains a handler function and the interface type it expects
type handlerInfo struct {
	handler       fnHandler
	interfaceType reflect.Type
}

var (
	// handlerFactory handler factory to keep all type of notification type
	handlerFactory []handlerInfo
)

// Register adds a notification handler factory function to the handlerFactory slice.
// The handler function creates an INotifiable implementation for a specific notification type.
// This allows registering different notification channels (Mail, SMS, Slack, etc.).
// The interfacePtr should be a pointer to an interface (e.g., (*IMailNotification)(nil))
func Register(handler fnHandler, interfacePtr any) {
	interfaceType := reflect.TypeOf(interfacePtr).Elem()
	handlerFactory = append(handlerFactory, handlerInfo{
		handler:       handler,
		interfaceType: interfaceType,
	})
}

// ========================================================================================
//                                 Notification Orchestration
// ========================================================================================

// Send delivers notifications concurrently through multiple notification handlers (SMS|Mail|Slack|Database).
// It takes a notification object of any type and sends it through all registered handlers that implement
// the corresponding notification interfaces. If notifications are disabled via NOTIFICATION_ENABLE env var,
// it will skip sending and return nil. Returns error.NotImplemented if no valid handlers are found.
//
// Parameters:
//   - notification: Any object implementing notification interfaces for registered handlers
//     To send Mail. For example, a struct implementing IMailNotification interface for mail notifications.
//     To send SMS. For example, a struct implementing ISmsNotification interface for sms notifications.
//     To send Slack. For example, a struct implementing ISlackNotification interface for Slack notifications.
//     To send Database. For example, a struct implementing IDatabaseNotification interface for database notifications.
//
// Returns:
//   - error: error.NotImplemented if no handlers match, nil on success or disabled notifications
func Send(notification any) error {
	if !utils.Getenv("NOTIFICATION_ENABLE", false) {
		log.Warnf("[STOP] Notification at %v", time.Now())

		return nil
	}

	var notificationHandlers []INotifiable

	for _, handlerInfo := range handlerFactory {
		notifyType := reflect.TypeOf(notification)

		if handlerInfo.handler == nil {
			continue
		}

		if !notifyType.Implements(handlerInfo.interfaceType) {
			continue
		}

		log.Tracef("[RUN] Notification handler for type %v", notifyType)

		notificationHandlers = append(notificationHandlers, handlerInfo.handler(notification))
	}

	if len(notificationHandlers) == 0 {
		return errors.NotImplemented
	}

	startTime := time.Now()

	// Send notifications concurrently using goroutines
	var wg sync.WaitGroup
	for _, notificationHandler := range notificationHandlers {
		wg.Add(1)
		go func(handler INotifiable) {
			defer wg.Done()
			handler.Notify()
		}(notificationHandler)
	}
	wg.Wait()

	log.Infof("[RUN] Notification time %v", time.Since(startTime))

	return nil
}

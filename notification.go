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
	// handlerFactory keeps all registered notification handlers.
	handlerFactory []handlerInfo

	// mu guards handlerFactory so Register and Send are safe to call concurrently.
	mu sync.RWMutex
)

// Register adds a notification handler factory function to the handlerFactory slice.
// The handler function creates an INotifiable implementation for a specific notification type.
// This allows registering different notification channels (Mail, SMS, Slack, etc.).
// The interfacePtr should be a pointer to an interface (e.g., (*IMailNotification)(nil)).
//
// Register is safe for concurrent use. Passing a nil handler or a non-interface
// pointer is a programming error and will panic.
func Register(handler fnHandler, interfacePtr any) {
	if handler == nil {
		panic("notification.Register: handler must not be nil")
	}

	ptrType := reflect.TypeOf(interfacePtr)
	if ptrType == nil || ptrType.Kind() != reflect.Ptr || ptrType.Elem().Kind() != reflect.Interface {
		panic("notification.Register: interfacePtr must be a pointer to an interface, e.g. (*IMailNotification)(nil)")
	}

	interfaceType := ptrType.Elem()

	mu.Lock()
	defer mu.Unlock()

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
// it will skip sending and return nil. Returns errors.NotImplemented if no valid handlers are found.
//
// Each handler runs in its own goroutine; a panic in one handler is recovered and logged so it
// cannot bring down the process or prevent the other handlers from completing.
//
// Parameters:
//   - notification: Any object implementing notification interfaces for registered handlers.
//     To send Mail. For example, a struct implementing IMailNotification interface for mail notifications.
//     To send SMS. For example, a struct implementing ISmsNotification interface for sms notifications.
//     To send Slack. For example, a struct implementing ISlackNotification interface for Slack notifications.
//     To send Database. For example, a struct implementing IDatabaseNotification interface for database notifications.
//
// Returns:
//   - error: errors.InvalidParameter if notification is nil, errors.NotImplemented if no handlers
//     match, nil on success or when notifications are disabled.
func Send(notification any) error {
	if !utils.Getenv("NOTIFICATION_ENABLE", false) {
		log.Warnf("[STOP] Notification at %v", time.Now())

		return nil
	}

	if notification == nil {
		return errors.InvalidParameter
	}

	notifyType := reflect.TypeOf(notification)

	// Snapshot the matching handlers under a read lock so registration can happen
	// concurrently without racing on handlerFactory.
	mu.RLock()
	var notificationHandlers []INotifiable
	for _, info := range handlerFactory {
		if info.handler == nil {
			continue
		}

		if !notifyType.Implements(info.interfaceType) {
			continue
		}

		log.Tracef("[RUN] Notification handler for type %v", notifyType)

		notificationHandlers = append(notificationHandlers, info.handler(notification))
	}
	mu.RUnlock()

	if len(notificationHandlers) == 0 {
		return errors.NotImplemented
	}

	startTime := time.Now()

	// Send notifications concurrently using goroutines.
	var wg sync.WaitGroup
	for _, notificationHandler := range notificationHandlers {
		wg.Add(1)
		go func(handler INotifiable) {
			defer wg.Done()
			// Isolate handler failures so one channel cannot crash the others or the process.
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("[PANIC] Notification handler %T recovered: %v", handler, r)
				}
			}()

			handler.Notify()
		}(notificationHandler)
	}
	wg.Wait()

	log.Infof("[RUN] Notification time %v", time.Since(startTime))

	return nil
}

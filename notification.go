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

type INotifiable interface {
	Notify()
}

type fnHandler func(notification any) INotifiable

var (
	// handlerFactory handler factory to keep all type of notification type
	handlerFactory = map[reflect.Type]fnHandler{}
)

func Register(notifyType reflect.Type, handler fnHandler) {
	handlerFactory[notifyType] = handler
}

func Type(i any) reflect.Type {
	return reflect.TypeOf(i).Elem()
}

// ========================================================================================
//                                 Notification Orchestration
// ========================================================================================

// Send Deliver many notification types SMS|Mail|Slack|Database.
func Send(notification any) error {
	if !utils.Getenv("NOTIFICATION_ENABLE", false) {
		log.Warnf("[STOP] Notification at %v", time.Now())

		return nil
	}

	var notificationHandlers []INotifiable

	for notifyType, handler := range handlerFactory {
		if reflect.TypeOf(notification).Implements(notifyType) {
			notificationHandlers = append(notificationHandlers, handler(notification))
		}
	}

	if len(notificationHandlers) == 0 {
		return errors.NotYetImplemented
	}

	var wg sync.WaitGroup

	startTime := time.Now()
	// Send message to each channel
	for _, notificationHandler := range notificationHandlers {
		wg.Add(1)
		notificationHandler := notificationHandler
		go func() {
			defer wg.Done()
			notificationHandler.Notify()
		}()
	}

	wg.Wait()

	log.Infof("[RUN] Notification time %v", time.Since(startTime))

	return nil
}

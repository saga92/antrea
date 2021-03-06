// Copyright 2019 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ram

import (
	"context"

	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"

	"github.com/vmware-tanzu/antrea/pkg/apiserver/storage"
)

// storeWatcher implements watch.Interface
type storeWatcher struct {
	// input represents the channel for incoming internal events that should be processed.
	input chan storage.InternalEvent
	// result represents the channel for outgoing events that will be sent to the client.
	result chan watch.Event
	done   chan struct{}
	// selectors represent a watcher's conditions to select objects.
	selectors *storage.Selectors
	// forget is used to cleanup the watcher
	forget func()
}

func newStoreWatcher(chanSize int, selectors *storage.Selectors, forget func()) *storeWatcher {
	return &storeWatcher{
		input:     make(chan storage.InternalEvent, chanSize),
		result:    make(chan watch.Event, chanSize),
		done:      make(chan struct{}),
		selectors: selectors,
		forget:    forget,
	}
}

// add sends InternalEvent to channel input.
// It's non blocking and simply discards the event if the channel is currently full.
// Finally, we should wait for the channel to be available and terminate the
// watcher, if it times out.
func (w *storeWatcher) add(event storage.InternalEvent) bool {
	select {
	case w.input <- event:
		return true
	default:
		// TODO: handle the case the channel is full
		klog.Errorf("The input channel is full. Event %+v was discarded", event)
		return false
	}
}

// process first sends initEvents and then keeps sending events got from channel input
// if they are newer than the specified resourceVersion.
func (w *storeWatcher) process(ctx context.Context, initEvents []storage.InternalEvent, resourceVersion uint64) {
	for _, event := range initEvents {
		w.sendWatchEvent(event)
	}
	defer close(w.result)
	for {
		select {
		case event, ok := <-w.input:
			if !ok {
				klog.Error("The input channel had been closed")
				return
			}
			if event.GetResourceVersion() > resourceVersion {
				w.sendWatchEvent(event)
			}
		case <-ctx.Done():
			return
		}
	}
}

// sendWatchEvent converts an InternalEvent to watch.Event based on the watcher's selectors.
// It sends the converted event to result channel, if not nil.
func (w *storeWatcher) sendWatchEvent(event storage.InternalEvent) {
	watchEvent := event.ToWatchEvent(w.selectors)
	if watchEvent == nil {
		// Watcher is not interested in that object.
		return
	}

	select {
	case <-w.done:
		return
	default:
	}

	select {
	case w.result <- *watchEvent:
	case <-w.done:
	}
}

// ResultChan returns the channel for outgoing events to the client.
func (w *storeWatcher) ResultChan() <-chan watch.Event {
	return w.result
}

// Stop stops the store watcher.
func (w *storeWatcher) Stop() {
	w.forget()
}

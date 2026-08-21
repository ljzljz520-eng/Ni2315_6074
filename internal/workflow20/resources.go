package workflow20

import (
	"fmt"
	"sync"

	"independentweeklylog/internal/domain"
)

type CloseRecorder interface{ Record(string) }

type ResourceCloser struct {
	resource *domain.Resource
	parent   *ResourceCloser
	recorder CloseRecorder
	closed   bool
	mu       sync.Mutex
}

func NewResourceCloser(resource *domain.Resource, parent *ResourceCloser, recorder CloseRecorder) *ResourceCloser {
	return &ResourceCloser{resource: resource, parent: parent, recorder: recorder}
}

func (c *ResourceCloser) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.parent != nil && !c.parent.IsClosed() {
		c.resource.Payload = ""
	}
	c.resource.Closed = true
	if c.recorder != nil {
		c.recorder.Record(c.resource.ID)
	}
	return nil
}

func (c *ResourceCloser) IsClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *ResourceCloser) Resource() domain.Resource {
	c.mu.Lock()
	defer c.mu.Unlock()
	return *c.resource
}

type CloseLog struct{ Items []string }

func (l *CloseLog) Record(id string) { l.Items = append(l.Items, id) }

func (l *CloseLog) ValidParentFirst(parentID, childID string) bool {
	if len(l.Items) < 2 {
		return false
	}
	return l.Items[0] == parentID && l.Items[1] == childID
}

func ensureResourceTree(resources []domain.Resource) (domain.Resource, domain.Resource, error) {
	var parent, child domain.Resource
	for _, resource := range resources {
		if resource.ParentID == "" && parent.ID == "" {
			parent = resource
		}
	}
	for _, resource := range resources {
		if resource.ParentID == parent.ID && child.ID == "" {
			child = resource
		}
	}
	if parent.ID == "" || child.ID == "" {
		return parent, child, fmt.Errorf("workflow20 requires a parent and child resource")
	}
	return parent, child, nil
}

package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"repairdesk.local/internal/model"
	"sync"
	"sync/atomic"
	"time"
)

var buckets = [][]byte{[]byte("parts"), []byte("equipment"), []byte("requests"), []byte("approvals"), []byte("transactions"), []byte("notifications"), []byte("audit")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err = s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
				return e
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	e := s.db.Close()
	s.db = nil
	return e
}
func (s *Store) Path() string { return s.path }
func put[T any](s *Store, b []byte, id string, v T) error {
	raw, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(b).Put([]byte(id), raw) })
}
func get[T any](s *Store, b []byte, id string) (T, error) {
	var out T
	e := s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(b).Get([]byte(id))
		if raw == nil {
			return model.ErrNotFound
		}
		return json.Unmarshal(raw, &out)
	})
	return out, e
}
func list[T any](s *Store, b []byte) ([]T, error) {
	out := []T{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(b).ForEach(func(_, v []byte) error {
			var x T
			if err := json.Unmarshal(v, &x); err != nil {
				return err
			}
			out = append(out, x)
			return nil
		})
	})
	return out, e
}
func del(s *Store, b []byte, id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(b).Delete([]byte(id)) })
}

var idCounter uint64

func nextID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, atomic.AddUint64(&idCounter, 1))
}
func (s *Store) SavePart(v model.SparePart) error { return put(s, buckets[0], v.ID, v) }
func (s *Store) GetPart(id string) (model.SparePart, error) {
	return get[model.SparePart](s, buckets[0], id)
}
func (s *Store) ListParts() ([]model.SparePart, error) { return list[model.SparePart](s, buckets[0]) }
func (s *Store) SaveEquipment(v model.Equipment) error { return put(s, buckets[1], v.ID, v) }
func (s *Store) GetEquipment(id string) (model.Equipment, error) {
	return get[model.Equipment](s, buckets[1], id)
}
func (s *Store) ListEquipment() ([]model.Equipment, error) {
	return list[model.Equipment](s, buckets[1])
}
func (s *Store) SaveRequest(v model.SpareIssueRequest) error { return put(s, buckets[2], v.ID, v) }
func (s *Store) GetRequest(id string) (model.SpareIssueRequest, error) {
	return get[model.SpareIssueRequest](s, buckets[2], id)
}
func (s *Store) ListRequests() ([]model.SpareIssueRequest, error) {
	return list[model.SpareIssueRequest](s, buckets[2])
}
func (s *Store) DeleteRequest(id string) error             { return del(s, buckets[2], id) }
func (s *Store) SaveApproval(v model.ApprovalRecord) error { return put(s, buckets[3], v.ID, v) }
func (s *Store) ListApprovals() ([]model.ApprovalRecord, error) {
	return list[model.ApprovalRecord](s, buckets[3])
}
func (s *Store) SaveTransaction(v model.InventoryTransaction) error {
	return put(s, buckets[4], v.ID, v)
}
func (s *Store) ListTransactions() ([]model.InventoryTransaction, error) {
	return list[model.InventoryTransaction](s, buckets[4])
}
func (s *Store) SaveNotification(v model.Notification) error { return put(s, buckets[5], v.ID, v) }
func (s *Store) ListNotifications() ([]model.Notification, error) {
	return list[model.Notification](s, buckets[5])
}
func (s *Store) SaveAudit(v model.AuditEvent) error     { return put(s, buckets[6], v.ID, v) }
func (s *Store) ListAudit() ([]model.AuditEvent, error) { return list[model.AuditEvent](s, buckets[6]) }
func NewID(prefix string) string                        { return nextID(prefix) }
func EnsureDir(path string) string                      { return filepath.Dir(path) }

package identity

import (
	"encoding/json"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserAgent string
	Cookies   map[string]string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	vault    *Vault
}

func NewSessionManager(vault *Vault) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		vault:    vault,
	}
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessions[id]
	if !ok {
		data := sm.vault.Get("session_" + id)
		if data == "" {
			return nil, false
		}
		var sess Session
		if err := json.Unmarshal([]byte(data), &sess); err != nil {
			return nil, false
		}
		return &sess, true
	}
	return s, ok
}

func (sm *SessionManager) SaveSession(session *Session) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[session.ID] = session
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return sm.vault.Set("session_"+session.ID, string(data))
}

func (sm *SessionManager) InvalidateSession(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, id)
	return sm.vault.Set("session_"+id, "")
}

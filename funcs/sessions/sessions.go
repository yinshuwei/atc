package sessions

import "github.com/gorilla/sessions"

// Sessions 会话数据
type Sessions struct {
	session *sessions.Session
	written bool
}

// New 新建Session
func New(session *sessions.Session) *Sessions {
	s := Sessions{
		session: session,
		written: false,
	}
	return &s
}

// Get 获得会话数据
func (s *Sessions) Get(key interface{}) interface{} {
	if v, ok := s.session.Values[key]; ok && v != nil {
		return v
	}
	return ""
}

// Set 设置会话数据
func (s *Sessions) Set(key, value interface{}) interface{} {
	s.session.Values[key] = value
	s.written = true
	return nil
}

// Delete 删除会话数据
func (s *Sessions) Delete(key interface{}) interface{} {
	delete(s.session.Values, key)
	s.written = true
	return nil
}

// Clear 删除所有会话数据
func (s *Sessions) Clear() interface{} {
	for key := range s.session.Values {
		s.Delete(key)
	}
	return nil
}

// Exist 是否存在会话数据
func (s *Sessions) Exist(key string) bool {
	v, ok := s.session.Values[key]
	return ok && v != nil
}

// Written 是否已有写入
func (s *Sessions) Written() bool {
	return s.written
}

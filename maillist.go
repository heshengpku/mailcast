package main

import (
	"errors"
	"sync"
	"fmt"
)

// MailManager 管理邮件相关操作
// 使用互斥锁确保并发安全
type MailManager struct {
	content  Content     // 邮件配置信息
	mailList []string    // 收件人邮箱列表
	mu       sync.RWMutex // 读写锁，保护并发访问
}

// NewMailManager 创建新的邮件管理器实例
// 返回初始化后的 MailManager 指针
func NewMailManager() *MailManager {
	return &MailManager{
		mailList: make([]string, 0, 10), // 预分配10个容量，提高性能
	}
}

// LoadContent 从指定路径加载邮件配置
// path: 配置文件路径
func (m *MailManager) LoadContent(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if err := loadData(&m.content, path); err != nil {
		return fmt.Errorf("failed to load content: %w", err)
	}
	return nil
}

// SaveContent 保存邮件配置到指定路径
// path: 配置文件保存路径
func (m *MailManager) SaveContent(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if err := saveData(m.content, path); err != nil {
		return fmt.Errorf("failed to save content: %w", err)
	}
	return nil
}

// SetContent 设置邮件内容和配置信息
// 所有参数都不能为空，且邮箱格式必须正确
func (m *MailManager) SetContent(name, pwd, host, subject, body string) error {
	// 参数验证
	if name == "" || pwd == "" || host == "" {
		return errors.New("required parameters cannot be empty")
	}
	
	if !emailValid(name) {
		return errors.New("invalid email format")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.content = Content{
		Name:    name,
		Pwd:     pwd,
		Host:    host,
		Subject: subject,
		Body:    body,
		MailList: m.mailList,
	}
	return nil
}

// LoadMails 从文件加载邮件列表
// path: 邮件列表文件路径
func (m *MailManager) LoadMails(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	
	mails, err := readMailsFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to read mails: %w", err)
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.mailList = mails
	return nil
}

// GetMails 获取邮件列表的副本
// 返回一个新的切片，避免外部修改影响内部数据
func (m *MailManager) GetMails() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make([]string, len(m.mailList))
	copy(result, m.mailList)
	return result
}

// DeleteMail 从列表中删除指定邮件地址
// mail: 要删除的邮件地址
func (m *MailManager) DeleteMail(mail string) error {
	if mail == "" {
		return errors.New("mail address cannot be empty")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	originalLen := len(m.mailList)
	m.mailList = delMailFromList(m.mailList, mail)
	
	// 如果长度没变，说明没有找到要删除的邮件
	if originalLen == len(m.mailList) {
		return fmt.Errorf("mail address %s not found", mail)
	}
	
	return nil
}

// GetContent 获取当前邮件配置的副本
func (m *MailManager) GetContent() Content {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.content
}

// AddMail 添加新的邮件地址到列表
func (m *MailManager) AddMail(mail string) error {
	if !emailValid(mail) {
		return errors.New("invalid email format")
	}
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 检查是否已存在
	for _, existingMail := range m.mailList {
		if existingMail == mail {
			return fmt.Errorf("mail address %s already exists", mail)
		}
	}
	
	m.mailList = append(m.mailList, mail)
	return nil
}
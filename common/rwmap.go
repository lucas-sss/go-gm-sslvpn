/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-14 22:25:49
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-14 23:22:35
 * @FilePath: /gmvpn/common/rwmap.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package common

import (
	"sync"
)

type RWMutexMap struct {
	lock sync.RWMutex
	m    map[string]interface{}
}

// 新建一个RWMap
func NewRWMutexMap(n int) *RWMutexMap {
	return &RWMutexMap{
		m: make(map[string]interface{}, n),
	}
}

func (receiver *RWMutexMap) Get(key string) (interface{}, bool) {
	receiver.lock.RLock()
	value, ok := receiver.m[key]
	receiver.lock.RUnlock()
	return value, ok
}

func (receiver *RWMutexMap) Set(key string, value interface{}) {
	receiver.lock.Lock()
	receiver.m[key] = value
	receiver.lock.Unlock()
}

func (receiver *RWMutexMap) TrySet(key string, value interface{}) bool {
	ok := receiver.lock.TryLock()
	if !ok {
		return false
	}
	receiver.m[key] = value
	receiver.lock.Unlock()
	return true
}

func (receiver *RWMutexMap) Del(key string) {
	receiver.lock.Lock()
	delete(receiver.m, key)
	receiver.lock.Unlock()
}

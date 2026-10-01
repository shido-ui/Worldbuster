package inventory

import (
	"errors"
	"sync"
)

var ErrUnknownItem = errors.New("unknown item")
var ErrInvalidQuantity = errors.New("quantity must be positive")
var ErrCapacity = errors.New("inventory capacity reached")

type Service struct {
	mu          sync.RWMutex
	items       map[string]Item
	inventories map[string]*Inventory
}

func NewService() *Service {
	return &Service{items: map[string]Item{}, inventories: map[string]*Inventory{}}
}
func (s *Service) RegisterItem(i Item) error {
	if i.ID == "" || i.MaxStack <= 0 {
		return ErrUnknownItem
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[i.ID] = i
	return nil
}
func (s *Service) Get(characterID string) Inventory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if inv, ok := s.inventories[characterID]; ok {
		return clone(*inv)
	}
	return Inventory{CharacterID: characterID, Capacity: 20, Stacks: []Stack{}}
}
func (s *Service) Add(characterID, itemID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[itemID]
	if !ok {
		return ErrUnknownItem
	}
	inv, ok := s.inventories[characterID]
	if !ok {
		inv = &Inventory{CharacterID: characterID, Capacity: 20}
		s.inventories[characterID] = inv
	}
	for i := range inv.Stacks {
		if inv.Stacks[i].ItemID == itemID && item.Stackable {
			space := item.MaxStack - inv.Stacks[i].Quantity
			if space > 0 {
				take := quantity
				if take > space {
					take = space
				}
				inv.Stacks[i].Quantity += take
				quantity -= take
				if quantity == 0 {
					return nil
				}
			}
		}
	}
	for quantity > 0 {
		if len(inv.Stacks) >= inv.Capacity {
			return ErrCapacity
		}
		take := quantity
		if item.Stackable && take > item.MaxStack {
			take = item.MaxStack
		}
		inv.Stacks = append(inv.Stacks, Stack{ItemID: itemID, Quantity: take})
		quantity -= take
		if !item.Stackable {
			break
		}
	}
	return nil
}
func (s *Service) Remove(characterID, itemID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	total := 0
	for _, st := range s.inventories[characterID].Stacks {
		if st.ItemID == itemID {
			total += st.Quantity
		}
	}
	if total < quantity {
		return errors.New("insufficient quantity")
	}
	for i := len(s.inventories[characterID].Stacks) - 1; i >= 0 && quantity > 0; i-- {
		st := &s.inventories[characterID].Stacks[i]
		if st.ItemID != itemID {
			continue
		}
		take := quantity
		if take > st.Quantity {
			take = st.Quantity
		}
		st.Quantity -= take
		quantity -= take
		if st.Quantity == 0 {
			s.inventories[characterID].Stacks = append(s.inventories[characterID].Stacks[:i], s.inventories[characterID].Stacks[i+1:]...)
		}
	}
	return nil
}
func clone(i Inventory) Inventory { i.Stacks = append([]Stack(nil), i.Stacks...); return i }

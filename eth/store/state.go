package store

type State struct {
	stateRoot *Trie
	storagMap map[string]*Trie
}

func InitState() *State {
	return &State{
		stateRoot: InitTrie(),
		storagMap: make(map[string]*Trie),
	}
}

func (s *State) UpdateAccount(addr string, data []byte) {
	if s.storagMap[addr] == nil {
		s.storagMap[addr] = InitTrie()
	}
	s.stateRoot.Update(addr, data)
}

func (s *State) GetAccount(addr string) []byte {
	node, err := s.stateRoot.Get(addr)
	if err != nil {
		panic("获取账户错误")
	}
	return node.Value
}

func (s *State) GetStateRoot() []byte {
	return s.stateRoot.RootHash
}

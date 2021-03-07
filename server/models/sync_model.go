package models

type SyncModel struct {
	onSave           chan bool `bson:"-"`
}

func (s *SyncModel) Saved() error {
	s.onSave <- true
	return nil
}

func (s *SyncModel) WaitForSave() {
	<-s.onSave
}

func (s *SyncModel) CollectionName() string {
	return "messages"
}

func (s *SyncModel) New() {
	s.onSave = make(chan bool, 1)
}

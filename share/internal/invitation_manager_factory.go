package internal

import (
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type InvitationManagerFactory struct {
	Rr RelationRepository
	T  EncodeDecoder
	C  XTaskRequester
}

func (i *InvitationManagerFactory) NewManager(conf contract.ConfigReader) *InvitationManager {
	return &InvitationManager{conf: conf, rr: i.Rr, tokenizer: i.T, xtaskClient: i.C}
}

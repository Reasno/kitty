package handlers

import (
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	code "glab.tagtic.cn/ad_gains/kitty/pkg/invitecode"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewShareService(
	manager InvitationManager,
	ur UserRepository,
	dispatcher contract.Dispatcher,
	tokenizer *code.Tokenizer,
) *shareService {
	return &shareService{
		manager:    manager,
		ur:         ur,
		dispatcher: dispatcher,
		tokenizer:  tokenizer,
	}
}

func ProvideShareServer(service *shareService) pb.ShareServer {
	return service
}

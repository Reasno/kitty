package handlers

import pb "glab.tagtic.cn/ad_gains/kitty/proto"

// NewService returns a naïve, stateless implementation of Service.
func NewShareService(manager InvitationManager, ur UserRepository) *shareService {
	return &shareService{manager: manager, ur: ur}
}

func ProvideShareServer(service *shareService) pb.ShareServer {
	return service
}

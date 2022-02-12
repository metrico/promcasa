package unmarshal

import (
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/utils/proto/logproto"
	"google.golang.org/protobuf/proto"
)

func UnmarshalProto(body []byte) (model.PushRequest, error) {
	oRes := logproto.PushRequest{}
	err := proto.Unmarshal(body, &oRes)
	if err != nil {
		return model.PushRequest{}, err
	}
	res := model.PushRequest{
		Streams: make([]model.Stream, len(oRes.GetStreams())),
	}
	for i := range oRes.GetStreams() {
		res.Streams[i] = model.Stream{
			Labels:  oRes.GetStreams()[i].Labels,
			Entries: make([]model.Entry, len(oRes.GetStreams()[i].GetEntries())),
		}
		for j := range oRes.GetStreams()[i].GetEntries() {
			enrty := oRes.GetStreams()[i].GetEntries()[j]
			res.Streams[i].Entries[j] = model.Entry{
				Timestamp: model.FromNano(enrty.Timestamp.GetSeconds()*1e9 + int64(enrty.Timestamp.GetNanos())),
				Line:      enrty.Line,
			}
		}
	}
	return res, nil
}

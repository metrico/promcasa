package controllerv1

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/metrico/promcasa/utils/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/metrico/promcasa/config"
	"github.com/metrico/promcasa/model"
	"github.com/metrico/promcasa/service"
	"github.com/metrico/promcasa/system/webmessages"
	"github.com/metrico/promcasa/utils/heputils"
	httpresponse "github.com/metrico/promcasa/utils/httpresponse"
	"github.com/metrico/promcasa/utils/labels"
	"github.com/metrico/promcasa/utils/logger"
	"github.com/metrico/promcasa/utils/unmarshal"
	unmarshal_legacy "github.com/metrico/promcasa/utils/unmarshal/legacy"
	"github.com/patrickmn/go-cache"
	"github.com/valyala/bytebufferpool"
)

type PromController struct {
	Controller
	PromService *service.PromService
}

// swagger:route GET /api/v1/prom/remote/write Data WriteData
//
// Returns data from server in array
//
// ---
//     Consumes:
//     - application/json
//
// 	   Produces:
// 	   - application/json
//
//	   Security:
//	   - JWT
//     - ApiKeyAuth
//
//
// SecurityDefinitions:
// JWT:
//      type: apiKey
//      name: Authorization
//      in: header
// ApiKeyAuth:
//      type: apiKey
//      in: header
//      name: Auth-Token
///
//  Responses:
//    201: body:TableUserList
//    400: body:FailureResponse

func (uc *PromController) WriteStream(ctx *fiber.Ctx) error {
	var req model.PushRequest
	var err error
	buf, err := helpers.GetRawBody(ctx)
	if err != nil {
		return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, "Request body too long")
	}
	defer buf.Release()
	if err != nil {
		return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, "Content-Length is incorrect")
	}
	if ctx.Get("content-type", "") == "application/x-protobuf" {
		req, err = unmarshal.UnmarshalProto(buf.Bytes())
		if err != nil {
			logger.Error(err.Error())
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}
	} else if heputils.GetVersion(string(ctx.Request().RequestURI())) == heputils.VersionV1 {
		if req, err = unmarshal.DecodePushRequestString(buf.Bytes()); err != nil {
			logger.Error(err.Error())
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}
	} else {
		if req, err = unmarshal_legacy.DecodePushRequestString(buf.Bytes()); err != nil {
			logger.Error(err.Error())
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}
	}

	logger.Debug("Data: ", req)

	tsReq := make([]*model.TableTimeSeries, 0, len(req.Streams))
	splReq := make([]*model.TableSample, 0, 1000)

	for _, stream := range req.Streams {
		var lbs []model.Label = nil
		if stream.Labels != "" {
			logger.Debug("Stream: ", stream)

			labels, err := labels.ParseLabels(stream.Labels)
			if err != nil {
				logger.Error(err.Error())
				return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
			}

			lbs = make([]model.Label, len(labels))
			labelKey := make([]string, len(labels))

			//fix the format = / :
			labelValue := make(map[string][]string)
			i := 0
			for k, l := range labels {
				keyJson := k
				valueJson := l
				labelKey[i] = keyJson
				labelValue[keyJson] = append(labelValue[keyJson], valueJson)
				lbs[i] = model.Label{
					Key:   keyJson,
					Value: valueJson,
				}
				i++
			}

			// lets insert only the unique values for key
			for k, v := range labelValue {
				if keys, exist := uc.PromService.GoCache.Get(k); exist {
					uc.PromService.GoCache.Replace(k, heputils.AppendTwoSlices(keys.([]string), heputils.UniqueSlice(v)), 0)
				} else {
					uc.PromService.GoCache.Add(k, heputils.UniqueSlice(v), 0)
				}
			}

		} else if stream.Stream != nil {
			lbs = make([]model.Label, len(stream.Stream))
			i := 0
			for k, v := range stream.Stream {
				lbs[i].Key = k
				lbs[i].Value = v
				i++
			}
		} else {
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}
		sort.Slice(lbs[:], func(i, j int) bool {
			return lbs[i].Key < lbs[j].Key
		})
		var fingerPrint uint64
		fingerByte, err := unmarshal.MarshalLabelsPushRequestString(lbs)
		if err != nil {
			logger.Error("bad unmarshaling of fingerprinting : ", err.Error())
			continue
		}

		logger.Debug("Label string: ", string(fingerByte))
		logger.Debug("fingertype : ", config.Setting.FingerPrintType)

		switch config.Setting.FingerPrintType {
		case config.FINGERPRINT_CityHash:
			fingerPrint = heputils.FingerprintLabelsCityHash(fingerByte)
		case config.FINGERPRINT_Bernstein:
			fingerPrint = uint64(heputils.FingerprintLabelsDJBHashPrometheus(fingerByte))
		}

		// if fingerprint was not found, lets insert into time_series
		if _, found := uc.PromService.GoCache.Get(fmt.Sprint(fingerPrint)); !found {

			b := bytebufferpool.Get()

			uc.PromService.GoCache.Set(fmt.Sprint(fingerPrint), true, cache.DefaultExpiration)
			tsReq = append(tsReq, &model.TableTimeSeries{
				Date:        time.Now(),
				FingerPrint: fingerPrint,
				Labels:      heputils.MakeJson(lbs, b),
				Name:        "",
			})
		}
		if stream.Entries != nil && len(stream.Entries) > 0 {
			for _, entries := range stream.Entries {
				splReq = append(splReq, &model.TableSample{
					FingerPrint: fingerPrint,
					TimestampMS: entries.Timestamp.GetNanos() / 1e6,
					Value:       0,
					String:      entries.Line,
				})

			}

		} else if stream.Values != nil && len(stream.Values) > 0 {
			for _, val := range stream.Values {
				iVal, err := strconv.ParseInt(val[0], 10, 64)
				if err != nil {
					logger.Error(err.Error())
					return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
				}
				splReq = append(splReq, &model.TableSample{
					FingerPrint: fingerPrint,
					TimestampMS: iVal / 1000000,
					Value:       0,
					String:      val[1],
				})

			}
		} else {
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}

		if err != nil {
			return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
		}
	}
	res := make([]chan error, 0, 2)
	res = append(res, uc.PromService.InsertTableSamples(splReq))
	if len(tsReq) > 0 {
		res = append(res, uc.PromService.InsertTimeSeriesRequest(tsReq))
	}
	err = nil
	for _, r := range res {
		_err := <-r
		if _err != nil {
			err = _err
		}
	}
	if err != nil {
		return httpresponse.CreateBadResponse(ctx, http.StatusBadRequest, webmessages.UserRequestFormatIncorrect)
	}

	return httpresponse.CreateSuccessResponseWTBody(ctx, http.StatusNoContent)
}

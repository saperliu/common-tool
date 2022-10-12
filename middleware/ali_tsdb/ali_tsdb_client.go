package ali_tsdb

import (
	"common-tool/common"
	"common-tool/logger"
	"common-tool/vo"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type AliTSDBClient struct {
	Address  string //地址
	Port     string //地址
	Api      string //地址
	User     string //用户
	Password string //密码
	DbClient *common.HttpClient
}

var loggerKairosdb = logger.CreateLogger("tsdb.log")

func NewAliTSDBClient(url string) *AliTSDBClient {
	httpClient := common.HttpClient{}
	httpClient.Init()
	return &AliTSDBClient{
		Address:  url,
		DbClient: &httpClient}
}

/**
 * 保存数据
 */
func (tSDBClient *AliTSDBClient) SaveHisPoint(storeHisData vo.StoreHisData) (error error) {

	mb := tSDBClient.builderMetric(storeHisData)
	loggerKairosdb.Info().Msgf(" SaveHisPoint  mb : %v ", mb)
	mjson, _ := json.Marshal(mb)
	mString := string(mjson)
	response, err := tSDBClient.DbClient.HttpPostJson("", mString)

	if err != nil {
		return err
	} else {
		code := response
		loggerKairosdb.Info().Msgf(" save point  code %v ", code)
	}
	return nil
}

func (tSDBClient *AliTSDBClient) builderMetric(storeHisData vo.StoreHisData) AliMetric {

	//增加数据
	metricBuilder := AliMetric{}
	timeTemp, err := strconv.ParseInt(storeHisData.Timestamp, 10, 64)
	if err != nil {
		logger.Error("builderPointMetric storeHisData.Timestamp  %v      %v ", storeHisData.Timestamp, err)
	}
	dateTemp := time.Unix(0, timeTemp*1e6)
	metricBuilder.Metric = strings.ToUpper(strings.TrimSpace(storeHisData.Tagid))
	v, err := strconv.ParseFloat(storeHisData.Value, 64)
	if err != nil {
		//v,err:=strconv.ParseFloat(storePointVo.Value,64)
		loggerKairosdb.Info().Msgf("----storePointVo %v    err  %v     ", v, err)
		metricBuilder.Value = v
	} else {
		//metricBuilder.Value = storeHisData.Value
	}
	metricBuilder.Timestamp = timeTemp
	tags := make(map[string]string)
	tags["status"] = common.NumberToString(storeHisData.Status)
	tags["year"] = common.NumberToString(dateTemp.Year())
	tags["month"] = dateTemp.Format("01")
	tags["day"] = dateTemp.Format("02")
	tags["hour"] = dateTemp.Format("15")
	tags["min"] = dateTemp.Format("04")
	tags["sec"] = dateTemp.Format("05")

	metricBuilder.Tags = tags

	return metricBuilder
}

/**
 * 保存数据
 */
func (tSDBClient *AliTSDBClient) SaveStorePoint(storePointVo vo.StorePoint) (error error) {

	mb := tSDBClient.builderPointMetric(storePointVo)
	loggerKairosdb.Info().Msgf(" SaveStorePoint  storePointVo : %v ", storePointVo)
	mjson, _ := json.Marshal(mb)
	mString := string(mjson)
	response, err := tSDBClient.DbClient.HttpPostJson("", mString)
	if err != nil {
		return err
	} else {
		code := response
		loggerKairosdb.Info().Msgf(" SaveStorePoint  code %v ", code)
	}
	return nil
}

func (tSDBClient *AliTSDBClient) builderPointMetric(storePointVo vo.StorePoint) AliMetric {

	//增加数据
	metricBuilder := AliMetric{}

	timeTemp, err := strconv.ParseInt(storePointVo.Timestamp, 10, 64)
	if err != nil {
		logger.Error("builderPointMetric storeHisData.Timestamp  %v      %v ", storePointVo.Timestamp, err)
	}
	dateTemp := time.Unix(0, timeTemp*1e6)
	metricBuilder.Metric = strings.ToUpper(strings.TrimSpace(storePointVo.Tagid))
	v, err := strconv.ParseFloat(storePointVo.Value, 64)
	if err != nil {
		loggerKairosdb.Info().Msgf("----storePointVo %v    err  %v     ", v, err)
		metricBuilder.Value = v
	} else {
		//metricBuilder.Value = storePointVo.Value
	}
	metricBuilder.Timestamp = timeTemp
	tags := make(map[string]string)
	tags["status"] = common.NumberToString(storePointVo.Status)
	tags["year"] = common.NumberToString(dateTemp.Year())
	tags["month"] = dateTemp.Format("01")
	tags["day"] = dateTemp.Format("02")
	tags["hour"] = dateTemp.Format("15")
	tags["min"] = dateTemp.Format("04")
	tags["sec"] = dateTemp.Format("05")

	metricBuilder.Tags = tags
	return metricBuilder
}

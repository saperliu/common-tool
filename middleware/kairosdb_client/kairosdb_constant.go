package kairosdb_client

const TIME_FORMAT = "2006-01-02 15:04:05"

const KAIROSDB_TAG_ORG = "org"

const KAIROSDB_TAG_SITE = "site"

const KAIROSDB_TAG_YEAR = "year"

const KAIROSDB_TAG_MONTH = "month"

const KAIROSDB_TAG_DAY = "day"

const KAIROSDB_TAG_HOUR = "hour"

const KAIROSDB_TAG_MIN = "min"

const KAIROSDB_TAG_SEC = "sec"

const KAIROSDB_TAG_WEEK = "week"

const KAIROSDB_TAG_REMARK = "remark"              //自定义的标签
const KAIROSDB_TAG_POSITION_TYPE = "positionType" //自定义的标签
const KAIROSDB_TAG_DATA_TYPE = "dataType"         //自定义的标签
const KAIROSDB_TAG_DATA_DESC = "dataDesc"         //自定义的标签
const KAIROSDB_TAG_STATUS = "status"              //原始数据值，此值做为原始数据的备份，在修复时使用。
const KAIROSDB_TAG_RSTATUS = "rstatus"            //修复后的数据值，查询时只查询此值，

const KAIROSDB_DATA_TYPE_GPS = "gps"         //自定义的数据
const KAIROSDB_DATA_TYPE_PERFECT = "perfect" //自定义的带数据描述的数据

const KAIROSDB_DATA_TYPE_NUMBER = "number" //自带的数字类型
const KAIROSDB_DATA_TYPE_TEXT = "text"     //自带的文本类型

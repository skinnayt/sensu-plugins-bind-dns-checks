package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type bindJsonStats struct {
	JsonStatsVersion string    `json:"json-stats-version"`
	BootTime         time.Time `json:"boot-time"`
	ConfigTime       time.Time `json:"config-time"`
	CurrentTime      time.Time `json:"current-time"`
	Version          string    `json:"version"`
	OpCodes          struct {
		Query      int `json:"QUERY"`
		IQuery     int `json:"IQUERY"`
		Status     int `json:"STATUS"`
		Reserved3  int `json:"RESERVED3"`
		Notify     int `json:"NOTIFY"`
		Update     int `json:"UPDATE"`
		Reserved6  int `json:"RESERVED6"`
		Reserved7  int `json:"RESERVED7"`
		Reserved8  int `json:"RESERVED8"`
		Reserved9  int `json:"RESERVED9"`
		Reserved10 int `json:"RESERVED10"`
		Reserved11 int `json:"RESERVED11"`
		Reserved12 int `json:"RESERVED12"`
		Reserved13 int `json:"RESERVED13"`
		Reserved14 int `json:"RESERVED14"`
		Reserved15 int `json:"RESERVED15"`
	} `json:"opcodes"`
	RCodes  RCode  `json:"rcodes"`
	QTypes  QTypes `json:"qtypes"`
	NSStats struct {
		Requestv4        int `json:"Requestv4"`
		Requestv6        int `json:"Requestv6"`
		ReqEdns0         int `json:"ReqEdns0"`
		ReqTCP           int `json:"ReqTCP"`
		TCPConnHighWater int `json:"TCPConnHighWater"`
		AuthQryRej       int `json:"AuthQryRej"`
		RecQryRej        int `json:"RecQryRej"`
		Response         int `json:"Response"`
		TruncatedResp    int `json:"TruncatedResp"`
		RespEDNS0        int `json:"RespEDNS0"`
		QrySuccess       int `json:"QrySuccess"`
		QryAuthAns       int `json:"QryAuthAns"`
		QryNoauthAns     int `json:"QryNoauthAns"`
		QryReferral      int `json:"QryReferral"`
		QryNxrrset       int `json:"QryNxrrset"`
		QryNXDOMAIN      int `json:"QryNXDOMAIN"`
		QryFailure       int `json:"QryFailure"`
		QryUDP           int `json:"QryUDP"`
		QryTCP           int `json:"QryTCP"`
		CookieIn         int `json:"CookieIn"`
		CookieNew        int `json:"CookieNew"`
		CookieMatch      int `json:"CookieMatch"`
		ECSOpt           int `json:"ECSOpt"`
	} `json:"nsstats"`
	ZoneStats struct {
		NotifyInv4 int `json:"NotifyInv4"`
		SOAOutv4   int `json:"SOAOutv4"`
		AXFRReqv4  int `json:"AXFRReqv4"`
		IXFRReqv4  int `json:"IXFRReqv4"`
		XfrSuccess int `json:"XfrSuccess"`
	} `json:"zonestats"`
	Views struct {
		Default BindView `json:"_default"`
		Bind    BindView `json:"_bind"`
	} `json:"views"`
	SocketStats struct {
		UDP4Open    int `json:"UDP4Open"`
		UDP6Open    int `json:"UDP6Open"`
		TCP4Open    int `json:"TCP4Open"`
		TCP6Open    int `json:"TCP6Open"`
		RawOpen     int `json:"RawOpen"`
		UDP4Close   int `json:"UDP4Close"`
		UDP6Close   int `json:"UDP6Close"`
		TCP4Close   int `json:"TCP4Close"`
		TCP6Close   int `json:"TCP6Close"`
		UDP6Conn    int `json:"UDP6Conn"`
		TCP4Conn    int `json:"TCP4Conn"`
		TCP6Conn    int `json:"TCP6Conn"`
		TCP4Accept  int `json:"TCP4Accept"`
		TCP6Accept  int `json:"TCP6Accept"`
		TCP4RecvErr int `json:"TCP4RecvErr"`
		UDP4Active  int `json:"UDP4Active"`
		UDP6Active  int `json:"UDP6Active"`
		TCP4Active  int `json:"TCP4Active"`
		TCP6Active  int `json:"TCP6Active"`
		RawActive   int `json:"RawActive"`
	} `json:"socketstats"`
	SocketMgr struct {
		Sockets []SocketMgrSocket `json:"sockets"`
	} `json:"socketmgr"`
	TaskMgr struct {
		ThreadModel    string         `json:"thread-model"`
		DefaultQuantum int            `json:"default-quantum"`
		Tasks          []*TaskMgrTask `json:"tasks"`
	} `json:"taskmgr"`
	Memory struct {
		TotalUse    int        `json:"TotalUse"`
		InUse       int        `json:"InUse"`
		Malloced    int        `json:"Malloced"`
		BlockSize   int        `json:"BlockSize"`
		ContextSize int        `json:"ContextSize"`
		Lost        int        `json:"Lost"`
		Contexts    []*Context `json:"Contexts"`
	} `json:"memory"`
	Traffic struct {
		DnsUDPRequestsSizesReceivedIPv4 struct {
			U0_15    int `json:"0-15,omitempty"`
			U16_31   int `json:"16-31,omitempty"`
			U32_47   int `json:"32-47,omitempty"`
			U48_63   int `json:"48-63,omitempty"`
			U64_79   int `json:"64-79,omitempty"`
			U80_95   int `json:"80-95,omitempty"`
			U96_111  int `json:"96-111,omitempty"`
			U112_127 int `json:"112-127,omitempty"`
			U128_143 int `json:"128-143,omitempty"`
			U144_159 int `json:"144-159,omitempty"`
		} `json:"dns-udp-requests-sizes-received-ipv4"`
		DnsUDPResponsesSizesSentIPv4 struct {
			U0_15      int `json:"0-15,omitempty"`
			U16_31     int `json:"16-31,omitempty"`
			U32_47     int `json:"32-47,omitempty"`
			U48_63     int `json:"48-63,omitempty"`
			U64_79     int `json:"64-79,omitempty"`
			U80_95     int `json:"80-95,omitempty"`
			U96_111    int `json:"96-111,omitempty"`
			U112_127   int `json:"112-127,omitempty"`
			U128_143   int `json:"128-143,omitempty"`
			U144_159   int `json:"144-159,omitempty"`
			U160_175   int `json:"160-175,omitempty"`
			U176_191   int `json:"176-191,omitempty"`
			U192_207   int `json:"192-207,omitempty"`
			U208_223   int `json:"208-223,omitempty"`
			U224_239   int `json:"224-239,omitempty"`
			U240_255   int `json:"240-255,omitempty"`
			U256_271   int `json:"256-271,omitempty"`
			U272_287   int `json:"272-287,omitempty"`
			U288_303   int `json:"288-303,omitempty"`
			U304_319   int `json:"304-319,omitempty"`
			U320_335   int `json:"320-335,omitempty"`
			U336_351   int `json:"336-351,omitempty"`
			U352_367   int `json:"352-367,omitempty"`
			U368_383   int `json:"368-383,omitempty"`
			U384_399   int `json:"384-399,omitempty"`
			U400_415   int `json:"400-415,omitempty"`
			U416_431   int `json:"416-431,omitempty"`
			U432_447   int `json:"432-447,omitempty"`
			U448_463   int `json:"448-463,omitempty"`
			U464_479   int `json:"464-479,omitempty"`
			U480_495   int `json:"480-495,omitempty"`
			U496_511   int `json:"496-511,omitempty"`
			U512_527   int `json:"512-527,omitempty"`
			U528_543   int `json:"528-543,omitempty"`
			U544_559   int `json:"544-559,omitempty"`
			U560_575   int `json:"560-575,omitempty"`
			U576_591   int `json:"576-591,omitempty"`
			U592_607   int `json:"592-607,omitempty"`
			U608_623   int `json:"608-623,omitempty"`
			U624_639   int `json:"624-639,omitempty"`
			U640_655   int `json:"640-655,omitempty"`
			U656_671   int `json:"656-671,omitempty"`
			U672_687   int `json:"672-687,omitempty"`
			U688_703   int `json:"688-703,omitempty"`
			U704_719   int `json:"704-719,omitempty"`
			U720_735   int `json:"720-735,omitempty"`
			U736_751   int `json:"736-751,omitempty"`
			U752_767   int `json:"752-767,omitempty"`
			U768_783   int `json:"768-783,omitempty"`
			U784_799   int `json:"784-799,omitempty"`
			U800_815   int `json:"800-815,omitempty"`
			U816_831   int `json:"816-831,omitempty"`
			U832_847   int `json:"832-847,omitempty"`
			U848_863   int `json:"848-863,omitempty"`
			U864_879   int `json:"864-879,omitempty"`
			U880_895   int `json:"880-895,omitempty"`
			U896_911   int `json:"896-911,omitempty"`
			U912_927   int `json:"912-927,omitempty"`
			U928_943   int `json:"928-943,omitempty"`
			U944_959   int `json:"944-959,omitempty"`
			U960_975   int `json:"960-975,omitempty"`
			U976_991   int `json:"976-991,omitempty"`
			U992_1007  int `json:"992-1007,omitempty"`
			U1008_1023 int `json:"1008-1023,omitempty"`
			U1024_1039 int `json:"1024-1039,omitempty"`
			U1040_1055 int `json:"1040-1055,omitempty"`
			U1056_1071 int `json:"1056-1071,omitempty"`
			U1072_1087 int `json:"1072-1087,omitempty"`
			U1088_1103 int `json:"1088-1103,omitempty"`
			U1104_1119 int `json:"1104-1119,omitempty"`
			U1120_1135 int `json:"1120-1135,omitempty"`
			U1136_1151 int `json:"1136-1151,omitempty"`
			U1152_1167 int `json:"1152-1167,omitempty"`
			U1168_1183 int `json:"1168-1183,omitempty"`
			U1184_1199 int `json:"1184-1199,omitempty"`
			U1200_1215 int `json:"1200-1215,omitempty"`
			U1216_1231 int `json:"1216-1231,omitempty"`
		} `json:"dns-udp-responses-sizes-sent-ipv4"`
		DnsTCPRequestsSizesReceivedIPv4 struct {
			T0_15  int `json:"0-15,omitempty"`
			T16_31 int `json:"16-31,omitempty"`
			T32_47 int `json:"32-47,omitempty"`
			T48_63 int `json:"48-63,omitempty"`
			T64_79 int `json:"64-79,omitempty"`
			T80_95 int `json:"80-95,omitempty"`
		} `json:"dns-tcp-requests-sizes-received-ipv4"`
		DnsTCPResponsesSizesSentIPv4 struct {
			T0_15      int `json:"0-15,omitempty"`
			T16_31     int `json:"16-31,omitempty"`
			T32_47     int `json:"32-47,omitempty"`
			T48_63     int `json:"48-63,omitempty"`
			T64_79     int `json:"64-79,omitempty"`
			T80_95     int `json:"80-95,omitempty"`
			T96_111    int `json:"96-111,omitempty"`
			T112_127   int `json:"112-127,omitempty"`
			T128_143   int `json:"128-143,omitempty"`
			T144_159   int `json:"144-159,omitempty"`
			T160_175   int `json:"160-175,omitempty"`
			T176_191   int `json:"176-191,omitempty"`
			T192_207   int `json:"192-207,omitempty"`
			T208_223   int `json:"208-223,omitempty"`
			T224_239   int `json:"224-239,omitempty"`
			T240_255   int `json:"240-255,omitempty"`
			T256_271   int `json:"256-271,omitempty"`
			T272_287   int `json:"272-287,omitempty"`
			T288_303   int `json:"288-303,omitempty"`
			T304_319   int `json:"304-319,omitempty"`
			T320_335   int `json:"320-335,omitempty"`
			T336_351   int `json:"336-351,omitempty"`
			T352_367   int `json:"352-367,omitempty"`
			T368_383   int `json:"368-383,omitempty"`
			T384_399   int `json:"384-399,omitempty"`
			T400_415   int `json:"400-415,omitempty"`
			T416_431   int `json:"416-431,omitempty"`
			T432_447   int `json:"432-447,omitempty"`
			T448_463   int `json:"448-463,omitempty"`
			T464_479   int `json:"464-479,omitempty"`
			T480_495   int `json:"480-495,omitempty"`
			T496_511   int `json:"496-511,omitempty"`
			T512_527   int `json:"512-527,omitempty"`
			T528_543   int `json:"528-543,omitempty"`
			T544_559   int `json:"544-559,omitempty"`
			T560_575   int `json:"560-575,omitempty"`
			T576_591   int `json:"576-591,omitempty"`
			T592_607   int `json:"592-607,omitempty"`
			T608_623   int `json:"608-623,omitempty"`
			T624_639   int `json:"624-639,omitempty"`
			T640_655   int `json:"640-655,omitempty"`
			T656_671   int `json:"656-671,omitempty"`
			T672_687   int `json:"672-687,omitempty"`
			T688_703   int `json:"688-703,omitempty"`
			T704_719   int `json:"704-719,omitempty"`
			T720_735   int `json:"720-735,omitempty"`
			T736_751   int `json:"736-751,omitempty"`
			T752_767   int `json:"752-767,omitempty"`
			T768_783   int `json:"768-783,omitempty"`
			T784_799   int `json:"784-799,omitempty"`
			T800_815   int `json:"800-815,omitempty"`
			T816_831   int `json:"816-831,omitempty"`
			T832_847   int `json:"832-847,omitempty"`
			T848_863   int `json:"848-863,omitempty"`
			T864_879   int `json:"864-879,omitempty"`
			T880_895   int `json:"880-895,omitempty"`
			T896_911   int `json:"896-911,omitempty"`
			T912_927   int `json:"912-927,omitempty"`
			T928_943   int `json:"928-943,omitempty"`
			T944_959   int `json:"944-959,omitempty"`
			T960_975   int `json:"960-975,omitempty"`
			T976_991   int `json:"976-991,omitempty"`
			T992_1007  int `json:"992-1007,omitempty"`
			T1008_1023 int `json:"1008-1023,omitempty"`
			T1024_1039 int `json:"1024-1039,omitempty"`
			T1040_1055 int `json:"1040-1055,omitempty"`
			T1056_1071 int `json:"1056-1071,omitempty"`
			T1072_1087 int `json:"1072-1087,omitempty"`
			T1088_1103 int `json:"1088-1103,omitempty"`
			T1104_1119 int `json:"1104-1119,omitempty"`
			T1120_1135 int `json:"1120-1135,omitempty"`
			T1136_1151 int `json:"1136-1151,omitempty"`
			T1152_1167 int `json:"1152-1167,omitempty"`
			T1168_1183 int `json:"1168-1183,omitempty"`
			T1184_1199 int `json:"1184-1199,omitempty"`
			T1200_1215 int `json:"1200-1215,omitempty"`
			T1216_1231 int `json:"1216-1231,omitempty"`
			T1232_1247 int `json:"1232-1247,omitempty"`
			T1248_1263 int `json:"1248-1263,omitempty"`
			T1264_1279 int `json:"1264-1279,omitempty"`
			T1280_1295 int `json:"1280-1295,omitempty"`
			T1296_1311 int `json:"1296-1311,omitempty"`
			T1312_1327 int `json:"1312-1327,omitempty"`
			T1328_1343 int `json:"1328-1343,omitempty"`
			T1344_1359 int `json:"1344-1359,omitempty"`
			T1360_1375 int `json:"1360-1375,omitempty"`
			T1376_1391 int `json:"1376-1391,omitempty"`
			T1392_1407 int `json:"1392-1407,omitempty"`
			T1408_1423 int `json:"1408-1423,omitempty"`
			T1424_1439 int `json:"1424-1439,omitempty"`
			T1440_1455 int `json:"1440-1455,omitempty"`
			T1456_1471 int `json:"1456-1471,omitempty"`
			T1472_1487 int `json:"1472-1487,omitempty"`
			T1488_1503 int `json:"1488-1503,omitempty"`
			T1504_1519 int `json:"1504-1519,omitempty"`
			T1520_1535 int `json:"1520-1535,omitempty"`
			T1536_1551 int `json:"1536-1551,omitempty"`
			T1552_1567 int `json:"1552-1567,omitempty"`
			T1568_1583 int `json:"1568-1583,omitempty"`
			T1584_1599 int `json:"1584-1599,omitempty"`
			T1600_1615 int `json:"1600-1615,omitempty"`
			T1616_1631 int `json:"1616-1631,omitempty"`
			T1632_1647 int `json:"1632-1647,omitempty"`
			T1648_1663 int `json:"1648-1663,omitempty"`
			T1664_1679 int `json:"1664-1679,omitempty"`
			T1680_1695 int `json:"1680-1695,omitempty"`
			T1696_1711 int `json:"1696-1711,omitempty"`
			T1712_1727 int `json:"1712-1727,omitempty"`
			T1728_1743 int `json:"1728-1743,omitempty"`
			T1744_1759 int `json:"1744-1759,omitempty"`
			T1760_1775 int `json:"1760-1775,omitempty"`
			T1776_1791 int `json:"1776-1791,omitempty"`
			T1792_1807 int `json:"1792-1807,omitempty"`
			T1808_1823 int `json:"1808-1823,omitempty"`
			T1824_1839 int `json:"1824-1839,omitempty"`
			T1840_1855 int `json:"1840-1855,omitempty"`
			T1856_1871 int `json:"1856-1871,omitempty"`
			T1872_1887 int `json:"1872-1887,omitempty"`
			T1888_1903 int `json:"1888-1903,omitempty"`
			T1904_1919 int `json:"1904-1919,omitempty"`
			T1920_1935 int `json:"1920-1935,omitempty"`
			T1936_1951 int `json:"1936-1951,omitempty"`
			T1952_1967 int `json:"1952-1967,omitempty"`
			T1968_1983 int `json:"1968-1983,omitempty"`
			T1984_1999 int `json:"1984-1999,omitempty"`
			T2000_2015 int `json:"2000-2015,omitempty"`
			T2016_2031 int `json:"2016-2031,omitempty"`
			T2032_2047 int `json:"2032-2047,omitempty"`
			T2048_2063 int `json:"2048-2063,omitempty"`
			T2064_2079 int `json:"2064-2079,omitempty"`
			T2080_2095 int `json:"2080-2095,omitempty"`
			T2096_2111 int `json:"2096-2111,omitempty"`
			T2112_2127 int `json:"2112-2127,omitempty"`
			T2128_2143 int `json:"2128-2143,omitempty"`
			T2144_2159 int `json:"2144-2159,omitempty"`
			T2160_2175 int `json:"2160-2175,omitempty"`
			T2176_2191 int `json:"2176-2191,omitempty"`
			T2192_2207 int `json:"2192-2207,omitempty"`
			T2208_2223 int `json:"2208-2223,omitempty"`
			T2224_2239 int `json:"2224-2239,omitempty"`
			T2240_2255 int `json:"2240-2255,omitempty"`
			T2256_2271 int `json:"2256-2271,omitempty"`
			T2272_2287 int `json:"2272-2287,omitempty"`
			T2288_2303 int `json:"2288-2303,omitempty"`
			T2304_2319 int `json:"2304-2319,omitempty"`
			T2320_2335 int `json:"2320-2335,omitempty"`
		} `json:"dns-tcp-responses-sizes-sent-ipv4"`
		DnsUDPRequestsSizesReceivedIPv6 struct {
			U0_15    int `json:"0-15,omitempty"`
			U16_31   int `json:"16-31,omitempty"`
			U32_47   int `json:"32-47,omitempty"`
			U48_63   int `json:"48-63,omitempty"`
			U64_79   int `json:"64-79,omitempty"`
			U80_95   int `json:"80-95,omitempty"`
			U96_111  int `json:"96-111,omitempty"`
			U112_127 int `json:"112-127,omitempty"`
			U128_143 int `json:"128-143,omitempty"`
			U144_159 int `json:"144-159,omitempty"`
		} `json:"dns-udp-requests-sizes-received-ipv6"`
		DnsUDPResponsesSizesSentIPv6 struct {
			U0_15      int `json:"0-15,omitempty"`
			U16_31     int `json:"16-31,omitempty"`
			U32_47     int `json:"32-47,omitempty"`
			U48_63     int `json:"48-63,omitempty"`
			U64_79     int `json:"64-79,omitempty"`
			U80_95     int `json:"80-95,omitempty"`
			U96_111    int `json:"96-111,omitempty"`
			U112_127   int `json:"112-127,omitempty"`
			U128_143   int `json:"128-143,omitempty"`
			U144_159   int `json:"144-159,omitempty"`
			U160_175   int `json:"160-175,omitempty"`
			U176_191   int `json:"176-191,omitempty"`
			U192_207   int `json:"192-207,omitempty"`
			U208_223   int `json:"208-223,omitempty"`
			U224_239   int `json:"224-239,omitempty"`
			U240_255   int `json:"240-255,omitempty"`
			U256_271   int `json:"256-271,omitempty"`
			U272_287   int `json:"272-287,omitempty"`
			U288_303   int `json:"288-303,omitempty"`
			U304_319   int `json:"304-319,omitempty"`
			U320_335   int `json:"320-335,omitempty"`
			U336_351   int `json:"336-351,omitempty"`
			U352_367   int `json:"352-367,omitempty"`
			U368_383   int `json:"368-383,omitempty"`
			U384_399   int `json:"384-399,omitempty"`
			U400_415   int `json:"400-415,omitempty"`
			U416_431   int `json:"416-431,omitempty"`
			U432_447   int `json:"432-447,omitempty"`
			U448_463   int `json:"448-463,omitempty"`
			U464_479   int `json:"464-479,omitempty"`
			U480_495   int `json:"480-495,omitempty"`
			U496_511   int `json:"496-511,omitempty"`
			U512_527   int `json:"512-527,omitempty"`
			U528_543   int `json:"528-543,omitempty"`
			U544_559   int `json:"544-559,omitempty"`
			U560_575   int `json:"560-575,omitempty"`
			U576_591   int `json:"576-591,omitempty"`
			U592_607   int `json:"592-607,omitempty"`
			U608_623   int `json:"608-623,omitempty"`
			U624_639   int `json:"624-639,omitempty"`
			U640_655   int `json:"640-655,omitempty"`
			U656_671   int `json:"656-671,omitempty"`
			U672_687   int `json:"672-687,omitempty"`
			U688_703   int `json:"688-703,omitempty"`
			U704_719   int `json:"704-719,omitempty"`
			U720_735   int `json:"720-735,omitempty"`
			U736_751   int `json:"736-751,omitempty"`
			U752_767   int `json:"752-767,omitempty"`
			U768_783   int `json:"768-783,omitempty"`
			U784_799   int `json:"784-799,omitempty"`
			U800_815   int `json:"800-815,omitempty"`
			U816_831   int `json:"816-831,omitempty"`
			U832_847   int `json:"832-847,omitempty"`
			U848_863   int `json:"848-863,omitempty"`
			U864_879   int `json:"864-879,omitempty"`
			U880_895   int `json:"880-895,omitempty"`
			U896_911   int `json:"896-911,omitempty"`
			U912_927   int `json:"912-927,omitempty"`
			U928_943   int `json:"928-943,omitempty"`
			U944_959   int `json:"944-959,omitempty"`
			U960_975   int `json:"960-975,omitempty"`
			U976_991   int `json:"976-991,omitempty"`
			U992_1007  int `json:"992-1007,omitempty"`
			U1008_1023 int `json:"1008-1023,omitempty"`
			U1024_1039 int `json:"1024-1039,omitempty"`
			U1040_1055 int `json:"1040-1055,omitempty"`
			U1056_1071 int `json:"1056-1071,omitempty"`
			U1072_1087 int `json:"1072-1087,omitempty"`
			U1088_1103 int `json:"1088-1103,omitempty"`
			U1104_1119 int `json:"1104-1119,omitempty"`
			U1120_1135 int `json:"1120-1135,omitempty"`
			U1136_1151 int `json:"1136-1151,omitempty"`
			U1152_1167 int `json:"1152-1167,omitempty"`
			U1168_1183 int `json:"1168-1183,omitempty"`
			U1184_1199 int `json:"1184-1199,omitempty"`
			U1200_1215 int `json:"1200-1215,omitempty"`
			U1216_1231 int `json:"1216-1231,omitempty"`
		} `json:"dns-udp-responses-sizes-sent-ipv6"`
		DnsTCPRequestsSizesReceivedIPv6 struct {
			T0_15  int `json:"0-15,omitempty"`
			T16_31 int `json:"16-31,omitempty"`
			T32_47 int `json:"32-47,omitempty"`
			T48_63 int `json:"48-63,omitempty"`
			T64_79 int `json:"64-79,omitempty"`
			T80_95 int `json:"80-95,omitempty"`
		} `json:"dns-tcp-requests-sizes-received-ipv6"`
		DnsTCPResponsesSizesSentIPv6 struct {
			T0_15      int `json:"0-15,omitempty"`
			T16_31     int `json:"16-31,omitempty"`
			T32_47     int `json:"32-47,omitempty"`
			T48_63     int `json:"48-63,omitempty"`
			T64_79     int `json:"64-79,omitempty"`
			T80_95     int `json:"80-95,omitempty"`
			T96_111    int `json:"96-111,omitempty"`
			T112_127   int `json:"112-127,omitempty"`
			T128_143   int `json:"128-143,omitempty"`
			T144_159   int `json:"144-159,omitempty"`
			T160_175   int `json:"160-175,omitempty"`
			T176_191   int `json:"176-191,omitempty"`
			T192_207   int `json:"192-207,omitempty"`
			T208_223   int `json:"208-223,omitempty"`
			T224_239   int `json:"224-239,omitempty"`
			T240_255   int `json:"240-255,omitempty"`
			T256_271   int `json:"256-271,omitempty"`
			T272_287   int `json:"272-287,omitempty"`
			T288_303   int `json:"288-303,omitempty"`
			T304_319   int `json:"304-319,omitempty"`
			T320_335   int `json:"320-335,omitempty"`
			T336_351   int `json:"336-351,omitempty"`
			T352_367   int `json:"352-367,omitempty"`
			T368_383   int `json:"368-383,omitempty"`
			T384_399   int `json:"384-399,omitempty"`
			T400_415   int `json:"400-415,omitempty"`
			T416_431   int `json:"416-431,omitempty"`
			T432_447   int `json:"432-447,omitempty"`
			T448_463   int `json:"448-463,omitempty"`
			T464_479   int `json:"464-479,omitempty"`
			T480_495   int `json:"480-495,omitempty"`
			T496_511   int `json:"496-511,omitempty"`
			T512_527   int `json:"512-527,omitempty"`
			T528_543   int `json:"528-543,omitempty"`
			T544_559   int `json:"544-559,omitempty"`
			T560_575   int `json:"560-575,omitempty"`
			T576_591   int `json:"576-591,omitempty"`
			T592_607   int `json:"592-607,omitempty"`
			T608_623   int `json:"608-623,omitempty"`
			T624_639   int `json:"624-639,omitempty"`
			T640_655   int `json:"640-655,omitempty"`
			T656_671   int `json:"656-671,omitempty"`
			T672_687   int `json:"672-687,omitempty"`
			T688_703   int `json:"688-703,omitempty"`
			T704_719   int `json:"704-719,omitempty"`
			T720_735   int `json:"720-735,omitempty"`
			T736_751   int `json:"736-751,omitempty"`
			T752_767   int `json:"752-767,omitempty"`
			T768_783   int `json:"768-783,omitempty"`
			T784_799   int `json:"784-799,omitempty"`
			T800_815   int `json:"800-815,omitempty"`
			T816_831   int `json:"816-831,omitempty"`
			T832_847   int `json:"832-847,omitempty"`
			T848_863   int `json:"848-863,omitempty"`
			T864_879   int `json:"864-879,omitempty"`
			T880_895   int `json:"880-895,omitempty"`
			T896_911   int `json:"896-911,omitempty"`
			T912_927   int `json:"912-927,omitempty"`
			T928_943   int `json:"928-943,omitempty"`
			T944_959   int `json:"944-959,omitempty"`
			T960_975   int `json:"960-975,omitempty"`
			T976_991   int `json:"976-991,omitempty"`
			T992_1007  int `json:"992-1007,omitempty"`
			T1008_1023 int `json:"1008-1023,omitempty"`
			T1024_1039 int `json:"1024-1039,omitempty"`
			T1040_1055 int `json:"1040-1055,omitempty"`
			T1056_1071 int `json:"1056-1071,omitempty"`
			T1072_1087 int `json:"1072-1087,omitempty"`
			T1088_1103 int `json:"1088-1103,omitempty"`
			T1104_1119 int `json:"1104-1119,omitempty"`
			T1120_1135 int `json:"1120-1135,omitempty"`
			T1136_1151 int `json:"1136-1151,omitempty"`
			T1152_1167 int `json:"1152-1167,omitempty"`
			T1168_1183 int `json:"1168-1183,omitempty"`
			T1184_1199 int `json:"1184-1199,omitempty"`
			T1200_1215 int `json:"1200-1215,omitempty"`
			T1216_1231 int `json:"1216-1231,omitempty"`
			T1232_1247 int `json:"1232-1247,omitempty"`
			T1248_1263 int `json:"1248-1263,omitempty"`
			T1264_1279 int `json:"1264-1279,omitempty"`
			T1280_1295 int `json:"1280-1295,omitempty"`
			T1296_1311 int `json:"1296-1311,omitempty"`
			T1312_1327 int `json:"1312-1327,omitempty"`
			T1328_1343 int `json:"1328-1343,omitempty"`
			T1344_1359 int `json:"1344-1359,omitempty"`
			T1360_1375 int `json:"1360-1375,omitempty"`
			T1376_1391 int `json:"1376-1391,omitempty"`
			T1392_1407 int `json:"1392-1407,omitempty"`
			T1408_1423 int `json:"1408-1423,omitempty"`
			T1424_1439 int `json:"1424-1439,omitempty"`
			T1440_1455 int `json:"1440-1455,omitempty"`
			T1456_1471 int `json:"1456-1471,omitempty"`
			T1472_1487 int `json:"1472-1487,omitempty"`
			T1488_1503 int `json:"1488-1503,omitempty"`
			T1504_1519 int `json:"1504-1519,omitempty"`
			T1520_1535 int `json:"1520-1535,omitempty"`
			T1536_1551 int `json:"1536-1551,omitempty"`
			T1552_1567 int `json:"1552-1567,omitempty"`
			T1568_1583 int `json:"1568-1583,omitempty"`
			T1584_1599 int `json:"1584-1599,omitempty"`
			T1600_1615 int `json:"1600-1615,omitempty"`
			T1616_1631 int `json:"1616-1631,omitempty"`
			T1632_1647 int `json:"1632-1647,omitempty"`
			T1648_1663 int `json:"1648-1663,omitempty"`
			T1664_1679 int `json:"1664-1679,omitempty"`
			T1680_1695 int `json:"1680-1695,omitempty"`
			T1696_1711 int `json:"1696-1711,omitempty"`
			T1712_1727 int `json:"1712-1727,omitempty"`
			T1728_1743 int `json:"1728-1743,omitempty"`
			T1744_1759 int `json:"1744-1759,omitempty"`
			T1760_1775 int `json:"1760-1775,omitempty"`
			T1776_1791 int `json:"1776-1791,omitempty"`
			T1792_1807 int `json:"1792-1807,omitempty"`
			T1808_1823 int `json:"1808-1823,omitempty"`
			T1824_1839 int `json:"1824-1839,omitempty"`
			T1840_1855 int `json:"1840-1855,omitempty"`
			T1856_1871 int `json:"1856-1871,omitempty"`
			T1872_1887 int `json:"1872-1887,omitempty"`
			T1888_1903 int `json:"1888-1903,omitempty"`
			T1904_1919 int `json:"1904-1919,omitempty"`
			T1920_1935 int `json:"1920-1935,omitempty"`
			T1936_1951 int `json:"1936-1951,omitempty"`
			T1952_1967 int `json:"1952-1967,omitempty"`
			T1968_1983 int `json:"1968-1983,omitempty"`
			T1984_1999 int `json:"1984-1999,omitempty"`
			T2000_2015 int `json:"2000-2015,omitempty"`
			T2016_2031 int `json:"2016-2031,omitempty"`
			T2032_2047 int `json:"2032-2047,omitempty"`
			T2048_2063 int `json:"2048-2063,omitempty"`
			T2064_2079 int `json:"2064-2079,omitempty"`
			T2080_2095 int `json:"2080-2095,omitempty"`
			T2096_2111 int `json:"2096-2111,omitempty"`
			T2112_2127 int `json:"2112-2127,omitempty"`
			T2128_2143 int `json:"2128-2143,omitempty"`
			T2144_2159 int `json:"2144-2159,omitempty"`
			T2160_2175 int `json:"2160-2175,omitempty"`
			T2176_2191 int `json:"2176-2191,omitempty"`
			T2192_2207 int `json:"2192-2207,omitempty"`
			T2208_2223 int `json:"2208-2223,omitempty"`
			T2224_2239 int `json:"2224-2239,omitempty"`
			T2240_2255 int `json:"2240-2255,omitempty"`
			T2256_2271 int `json:"2256-2271,omitempty"`
			T2272_2287 int `json:"2272-2287,omitempty"`
			T2288_2303 int `json:"2288-2303,omitempty"`
			T2304_2319 int `json:"2304-2319,omitempty"`
			T2320_2335 int `json:"2320-2335,omitempty"`
		} `json:"dns-tcp-responses-sizes-sent-ipv6"`
	} `json:"traffic"`
}

type QTypes struct {
	Others     int `json:"Others,omitempty"`
	A          int `json:"A,omitempty"`
	Ns         int `json:"NS,omitempty"`
	Cname      int `json:"CNAME,omitempty"`
	Soa        int `json:"SOA,omitempty"`
	Ptr        int `json:"PTR,omitempty"`
	Mx         int `json:"MX,omitempty"`
	Txt        int `json:"TXT,omitempty"`
	Afsdb      int `json:"AFSDB,omitempty"`
	Aaaa       int `json:"AAAA,omitempty"`
	Srv        int `json:"SRV,omitempty"`
	Naptr      int `json:"NAPTR,omitempty"`
	Dname      int `json:"DNAME,omitempty"`
	Ds         int `json:"DS,omitempty"`
	Rrsig      int `json:"RRSIG,omitempty"`
	Dnskey     int `json:"DNSKEY,omitempty"`
	Nsec3param int `json:"NSEC3PARAM,omitempty"`
	Tlsa       int `json:"TLSA,omitempty"`
	Cds        int `json:"CDS,omitempty"`
	Cdnskey    int `json:"CDNSKEY,omitempty"`
	Zonemd     int `json:"ZONEMD,omitempty"`
	Svcb       int `json:"SVCB,omitempty"`
	Https      int `json:"HTTPS,omitempty"`
	Spf        int `json:"SPF,omitempty"`
	Any        int `json:"ANY,omitempty"`
}

func (q *QTypes) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	metrics = append(metrics, &Metric{
		Name:      "Others",
		Value:     int64(q.Others),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "A",
		Value:     int64(q.A),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NS",
		Value:     int64(q.Ns),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "CNAME",
		Value:     int64(q.Cname),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "SOA",
		Value:     int64(q.Soa),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "PTR",
		Value:     int64(q.Ptr),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "MX",
		Value:     int64(q.Mx),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "TXT",
		Value:     int64(q.Txt),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "AFSDB",
		Value:     int64(q.Afsdb),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "AAAA",
		Value:     int64(q.Aaaa),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "SRV",
		Value:     int64(q.Srv),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NAPTR",
		Value:     int64(q.Naptr),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "DNAME",
		Value:     int64(q.Dname),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "DS",
		Value:     int64(q.Ds),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RRSIG",
		Value:     int64(q.Rrsig),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "DNSKEY",
		Value:     int64(q.Dnskey),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NSEC3PARAM",
		Value:     int64(q.Nsec3param),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "TLSA",
		Value:     int64(q.Tlsa),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "CDS",
		Value:     int64(q.Cds),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "CDNSKEY",
		Value:     int64(q.Cdnskey),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "ZONEMD",
		Value:     int64(q.Zonemd),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "SVCB",
		Value:     int64(q.Svcb),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "HTTPS",
		Value:     int64(q.Https),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "SPF",
		Value:     int64(q.Spf),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "ANY",
		Value:     int64(q.Any),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})

	return metrics
}

type RCode struct {
	Noerror    int `json:"NOERROR,omitempty"`
	Formerr    int `json:"FORMERR,omitempty"`
	Servfail   int `json:"SERVFAIL,omitempty"`
	Nxdomain   int `json:"NXDOMAIN,omitempty"`
	Notimp     int `json:"NOTIMP,omitempty"`
	Refused    int `json:"REFUSED,omitempty"`
	Yxdomain   int `json:"YXDOMAIN,omitempty"`
	Yxrrset    int `json:"YXRRSET,omitempty"`
	Nxrrset    int `json:"NXRRSET,omitempty"`
	Notauth    int `json:"NOTAUTH,omitempty"`
	Notzone    int `json:"NOTZONE,omitempty"`
	Reserved11 int `json:"RESERVED11,omitempty"`
	Reserved12 int `json:"RESERVED12,omitempty"`
	Reserved13 int `json:"RESERVED13,omitempty"`
	Reserved14 int `json:"RESERVED14,omitempty"`
	Reserved15 int `json:"RESERVED15,omitempty"`
	Badvers    int `json:"BADVERS,omitempty"`
	R17        int `json:"17,omitempty"`
	R18        int `json:"18,omitempty"`
	R19        int `json:"19,omitempty"`
	R20        int `json:"20,omitempty"`
	R21        int `json:"21,omitempty"`
	R22        int `json:"22,omitempty"`
	Badcookie  int `json:"BADCOOKIE,omitempty"`
}

func (r *RCode) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	metrics = append(metrics, &Metric{
		Name:      "NOERROR",
		Value:     int64(r.Noerror),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "FORMERR",
		Value:     int64(r.Formerr),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "SERVFAIL",
		Value:     int64(r.Servfail),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NXDOMAIN",
		Value:     int64(r.Nxdomain),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NOTIMP",
		Value:     int64(r.Notimp),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "REFUSED",
		Value:     int64(r.Refused),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "YXDOMAIN",
		Value:     int64(r.Yxdomain),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "YXRRSET",
		Value:     int64(r.Yxrrset),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NXRRSET",
		Value:     int64(r.Nxrrset),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NOTAUTH",
		Value:     int64(r.Notauth),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "NOTZONE",
		Value:     int64(r.Notzone),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RESERVED11",
		Value:     int64(r.Reserved11),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RESERVED12",
		Value:     int64(r.Reserved12),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RESERVED13",
		Value:     int64(r.Reserved13),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RESERVED14",
		Value:     int64(r.Reserved14),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RESERVED15",
		Value:     int64(r.Reserved15),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "BADVERS",
		Value:     int64(r.Badvers),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "17",
		Value:     int64(r.R17),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "18",
		Value:     int64(r.R18),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "19",
		Value:     int64(r.R19),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "20",
		Value:     int64(r.R20),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "21",
		Value:     int64(r.R21),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "22",
		Value:     int64(r.R22),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "BADCOOKIE",
		Value:     int64(r.Badcookie),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})

	return metrics
}

type BindView struct {
	Zones    []*ZoneView `json:"zones"`
	Resolver struct {
		Stats struct {
			Queryv6         int `json:"Queryv6"`
			Responsev6      int `json:"Responsev6"`
			NXDOMAIN        int `json:"NXDOMAIN"`
			Truncated       int `json:"Truncated"`
			Retry           int `json:"Retry"`
			ValAttempt      int `json:"ValAttempt"`
			ValOk           int `json:"ValOk"`
			ValNegOk        int `json:"ValNegOk"`
			QryRTT100       int `json:"QryRTT100"`
			QryRTT500       int `json:"QryRTT500"`
			BucketSize      int `json:"BucketSize"`
			ClientCookieOut int `json:"ClientCookieOut"`
			ServerCookieOut int `json:"ServerCookieOut"`
			CookieIn        int `json:"CookieIn"`
			CookieClientOk  int `json:"CookieClientOk"`
			Priming         int `json:"Priming"`
		} `json:"stats"`
		QTypes     QTypes `json:"qtypes"`
		Cache      QTypes `json:"cache"`
		CacheStats struct {
			CacheHits    int `json:"CacheHits"`
			CacheMisses  int `json:"CacheMisses"`
			QueryHits    int `json:"QueryHits"`
			QueryMisses  int `json:"QueryMisses"`
			DeleteLRU    int `json:"DeleteLRU"`
			DeleteTTL    int `json:"DeleteTTL"`
			CacheNodes   int `json:"CacheNodes"`
			CacheBuckets int `json:"CacheBuckets"`
			TreeMemTotal int `json:"TreeMemTotal"`
			TreeMemInUse int `json:"TreeMemInUse"`
			TreeMemMax   int `json:"TreeMemMax"`
			HeapMemTotal int `json:"HeapMemTotal"`
			HeapMemInUse int `json:"HeapMemInUse"`
			HeapMemMax   int `json:"HeapMemMax"`
		} `json:"cache-stats"`
		Adb struct {
			Nentries   int `json:"nentries"`
			Entriescnt int `json:"entriescnt"`
			Nnames     int `json:"nnames"`
			Namescnt   int `json:"namescnt"`
		} `json:"adb"`
	} `json:"resolver"`
}

func (bv *BindView) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	zone_metrics := make([]*Metric, 0)
	for _, zone := range bv.Zones {
		zone_metrics = append(zone_metrics, zone.toMetrics(metric_time)...)
	}
	metrics = append(metrics, zone_metrics...)
	stats_tag := &MetricTag{"resolver", "stats"}
	resolver_stats_metrics := make([]*Metric, 0)
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "Queryv6",
		Value:     int64(bv.Resolver.Stats.Queryv6),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "Responsev6",
		Value:     int64(bv.Resolver.Stats.Responsev6),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "NXDOMAIN",
		Value:     int64(bv.Resolver.Stats.NXDOMAIN),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "Truncated",
		Value:     int64(bv.Resolver.Stats.Truncated),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "Retry",
		Value:     int64(bv.Resolver.Stats.Retry),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "ValAttempt",
		Value:     int64(bv.Resolver.Stats.ValAttempt),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "ValOk",
		Value:     int64(bv.Resolver.Stats.ValOk),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "ValNegOk",
		Value:     int64(bv.Resolver.Stats.ValNegOk),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "QryRTT100",
		Value:     int64(bv.Resolver.Stats.QryRTT100),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "QryRTT500",
		Value:     int64(bv.Resolver.Stats.QryRTT500),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "BucketSize",
		Value:     int64(bv.Resolver.Stats.BucketSize),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "ClientCookieOut",
		Value:     int64(bv.Resolver.Stats.ClientCookieOut),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "ServerCookieOut",
		Value:     int64(bv.Resolver.Stats.ServerCookieOut),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "CookieIn",
		Value:     int64(bv.Resolver.Stats.CookieIn),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "CookieClientOk",
		Value:     int64(bv.Resolver.Stats.CookieClientOk),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	resolver_stats_metrics = append(resolver_stats_metrics, &Metric{
		Name:      "Priming",
		Value:     int64(bv.Resolver.Stats.Priming),
		Timestamp: metric_time,
		Tags:      []*MetricTag{stats_tag},
	})
	metrics = append(metrics, resolver_stats_metrics...)

	qtypes_tag := &MetricTag{"resolver", "qtypes"}
	resolver_qtypes_metrics := bv.Resolver.QTypes.toMetrics(metric_time)
	for _, metric := range resolver_qtypes_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, qtypes_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	metrics = append(metrics, resolver_qtypes_metrics...)
	cache_tag := &MetricTag{"resolver", "cache"}
	resolver_cache_metrics := bv.Resolver.Cache.toMetrics(metric_time)
	for _, metric := range resolver_cache_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, cache_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	metrics = append(metrics, resolver_cache_metrics...)

	cachestats_tag := &MetricTag{"resolver", "cachestats"}
	resolver_cache_stats_metrics := make([]*Metric, 0)
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "CacheHits",
		Value:     int64(bv.Resolver.CacheStats.CacheHits),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "CacheMisses",
		Value:     int64(bv.Resolver.CacheStats.CacheMisses),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "QueryHits",
		Value:     int64(bv.Resolver.CacheStats.QueryHits),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "QueryMisses",
		Value:     int64(bv.Resolver.CacheStats.QueryMisses),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "DeleteLRU",
		Value:     int64(bv.Resolver.CacheStats.DeleteLRU),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "DeleteTTL",
		Value:     int64(bv.Resolver.CacheStats.DeleteTTL),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "CacheNodes",
		Value:     int64(bv.Resolver.CacheStats.CacheNodes),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "CacheBuckets",
		Value:     int64(bv.Resolver.CacheStats.CacheBuckets),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "TreeMemTotal",
		Value:     int64(bv.Resolver.CacheStats.TreeMemTotal),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "TreeMemInUse",
		Value:     int64(bv.Resolver.CacheStats.TreeMemInUse),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "TreeMemMax",
		Value:     int64(bv.Resolver.CacheStats.TreeMemMax),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "HeapMemTotal",
		Value:     int64(bv.Resolver.CacheStats.HeapMemTotal),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "HeapMemInUse",
		Value:     int64(bv.Resolver.CacheStats.HeapMemInUse),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	resolver_cache_stats_metrics = append(resolver_cache_stats_metrics, &Metric{
		Name:      "HeapMemMax",
		Value:     int64(bv.Resolver.CacheStats.HeapMemMax),
		Timestamp: metric_time,
		Tags:      []*MetricTag{cachestats_tag},
	})
	metrics = append(metrics, resolver_cache_stats_metrics...)
	adb_tag := &MetricTag{"resolver", "adb"}
	resolver_adb_metrics := make([]*Metric, 0)
	resolver_adb_metrics = append(resolver_adb_metrics, &Metric{
		Name:      "nentries",
		Value:     int64(bv.Resolver.Adb.Nentries),
		Timestamp: metric_time,
		Tags:      []*MetricTag{adb_tag},
	})
	resolver_adb_metrics = append(resolver_adb_metrics, &Metric{
		Name:      "entriescnt",
		Value:     int64(bv.Resolver.Adb.Entriescnt),
		Timestamp: metric_time,
		Tags:      []*MetricTag{adb_tag},
	})
	resolver_adb_metrics = append(resolver_adb_metrics, &Metric{
		Name:      "nnames",
		Value:     int64(bv.Resolver.Adb.Nnames),
		Timestamp: metric_time,
		Tags:      []*MetricTag{adb_tag},
	})
	resolver_adb_metrics = append(resolver_adb_metrics, &Metric{
		Name:      "namescnt",
		Value:     int64(bv.Resolver.Adb.Namescnt),
		Timestamp: metric_time,
		Tags:      []*MetricTag{adb_tag},
	})
	metrics = append(metrics, resolver_adb_metrics...)

	return metrics
}

type SocketMgrSocket struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	References   int      `json:"references"`
	Type         string   `json:"type"`
	PeerAddress  string   `json:"peer-address,omitempty"`
	LocalAddress string   `json:"local-address,omitempty"`
	States       []string `json:"states"`
}

func (s *SocketMgrSocket) toMetric(metric_time time.Time) *Metric {
	socket_metric := &Metric{
		Name:      s.Name,
		Value:     int64(s.References),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	}
	if s.LocalAddress != "" {
		socket_metric.Tags = append(
			socket_metric.Tags,
			&MetricTag{
				"local-address",
				strings.Replace(s.LocalAddress, ".", "_", -1),
			},
		)
	}
	if s.PeerAddress != "" {
		socket_metric.Tags = append(
			socket_metric.Tags,
			&MetricTag{
				"peer-address",
				strings.Replace(s.PeerAddress, ".", "_", -1),
			},
		)
	}
	socket_metric.Tags = append(socket_metric.Tags, &MetricTag{"type", s.Type})
	return socket_metric
}

type TaskMgrTask struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	References int    `json:"references"`
	State      string `json:"state"`
	Quantum    int    `json:"quantum"`
	Events     int    `json:"events"`
}

func (t *TaskMgrTask) toMetric(metric_time time.Time) *Metric {
	task_metric := &Metric{
		Name:      t.Id,
		Value:     int64(t.Events),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	}
	if t.Name != "" {
		task_metric.Tags = append(task_metric.Tags, &MetricTag{"name", t.Name})
	}
	return task_metric
}

type Context struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	References  int    `json:"references"`
	Total       int    `json:"total"`
	Inuse       int    `json:"inuse"`
	Maxinuse    int    `json:"maxinuse"`
	Malloced    int    `json:"malloced"`
	Maxmalloced int    `json:"maxmalloced"`
	Blocksize   int    `json:"blocksize"`
	Pools       int    `json:"pools"`
	Hiwater     int    `json:"hiwater"`
	Lowater     int    `json:"lowater"`
}

func (c *Context) toMetric(metric_time time.Time) []*Metric {
	context_name_tag := &MetricTag{"context", c.Name}
	context_id_tag := &MetricTag{"context-id", c.Id}

	context_metrics := make([]*Metric, 0)
	context_metrics = append(context_metrics, &Metric{
		Name:      "References",
		Value:     int64(c.References),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Total",
		Value:     int64(c.Total),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Inuse",
		Value:     int64(c.Inuse),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Maxinuse",
		Value:     int64(c.Maxinuse),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Malloced",
		Value:     int64(c.Malloced),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Maxmalloced",
		Value:     int64(c.Maxmalloced),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Blocksize",
		Value:     int64(c.Blocksize),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Pools",
		Value:     int64(c.Pools),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Hiwater",
		Value:     int64(c.Hiwater),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})
	context_metrics = append(context_metrics, &Metric{
		Name:      "Lowater",
		Value:     int64(c.Lowater),
		Timestamp: metric_time,
		Tags:      []*MetricTag{context_name_tag, context_id_tag},
	})

	return context_metrics
}

type ZoneView struct {
	Name    string    `json:"name"`
	Class   string    `json:"class"`
	Serial  int       `json:"serial"`
	Type    string    `json:"type"`
	Loaded  time.Time `json:"loaded"`
	Expires time.Time `json:"expires,omitempty"`
	Refresh time.Time `json:"refresh,omitempty"`
	RCodes  RCode     `json:"rcodes,omitempty"`
	QTypes  QTypes    `json:"qtypes"`
}

func (z *ZoneView) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	zone_name_tag := &MetricTag{"zone", strings.Replace(z.Name, ".", "_", -1)}
	zone_class_tag := &MetricTag{"class", z.Class}
	zone_type_tag := &MetricTag{"type", z.Type}

	rcode_metrics := z.RCodes.toMetrics(metric_time)
	for _, rcode_metric := range rcode_metrics {
		rcode_metric_tags := make([]*MetricTag, 0, len(rcode_metric.Tags)+3)
		rcode_metric_tags = append(rcode_metric_tags, zone_name_tag)
		rcode_metric_tags = append(rcode_metric_tags, zone_class_tag)
		rcode_metric_tags = append(rcode_metric_tags, zone_type_tag)
		rcode_metric_tags = append(rcode_metric_tags, rcode_metric.Tags...)
		rcode_metric.Tags = rcode_metric_tags
	}
	metrics = append(metrics, rcode_metrics...)

	qtype_metrics := z.QTypes.toMetrics(metric_time)
	for _, qtype_metric := range qtype_metrics {
		qtype_metric_tags := make([]*MetricTag, 0, len(qtype_metric.Tags)+3)
		qtype_metric_tags = append(qtype_metric_tags, zone_name_tag)
		qtype_metric_tags = append(qtype_metric_tags, zone_class_tag)
		qtype_metric_tags = append(qtype_metric_tags, zone_type_tag)
		qtype_metric.Tags = qtype_metric_tags
	}
	metrics = append(metrics, qtype_metrics...)

	return metrics
}

func ReadJsonStats(statsData []byte) error {
	// Read the JSON statistics
	var jsonStats bindJsonStats

	err := json.Unmarshal(statsData, &jsonStats)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return err
	}

	return_metrics := make([]*Metric, 0)

	opcodes_tag := &MetricTag{"server", "opcodes"}
	return_metrics = append(return_metrics, &Metric{
		Name:      "Query",
		Value:     int64(jsonStats.OpCodes.Query),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "IQuery",
		Value:     int64(jsonStats.OpCodes.IQuery),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Status",
		Value:     int64(jsonStats.OpCodes.Status),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Notify",
		Value:     int64(jsonStats.OpCodes.Notify),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Update",
		Value:     int64(jsonStats.OpCodes.Update),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved6",
		Value:     int64(jsonStats.OpCodes.Reserved6),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved7",
		Value:     int64(jsonStats.OpCodes.Reserved7),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved8",
		Value:     int64(jsonStats.OpCodes.Reserved8),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved9",
		Value:     int64(jsonStats.OpCodes.Reserved9),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved10",
		Value:     int64(jsonStats.OpCodes.Reserved10),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved11",
		Value:     int64(jsonStats.OpCodes.Reserved11),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved12",
		Value:     int64(jsonStats.OpCodes.Reserved12),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved13",
		Value:     int64(jsonStats.OpCodes.Reserved13),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved14",
		Value:     int64(jsonStats.OpCodes.Reserved14),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Reserved15",
		Value:     int64(jsonStats.OpCodes.Reserved15),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})

	rcodes_tag := &MetricTag{"server", "rcodes"}
	json_rcodes := jsonStats.RCodes.toMetrics(jsonStats.CurrentTime)
	for _, rcode_metric := range json_rcodes {
		metric_tags := make([]*MetricTag, 0, len(rcode_metric.Tags)+1)
		metric_tags = append(metric_tags, rcodes_tag)
		metric_tags = append(metric_tags, rcode_metric.Tags...)
		rcode_metric.Tags = metric_tags
	}
	return_metrics = append(return_metrics, json_rcodes...)

	qtypes_tag := &MetricTag{"server", "qtypes"}
	json_qtypes := jsonStats.QTypes.toMetrics(jsonStats.CurrentTime)
	for _, qtype_metric := range json_qtypes {
		metric_tags := make([]*MetricTag, 0, len(qtype_metric.Tags)+1)
		metric_tags = append(metric_tags, qtypes_tag)
		metric_tags = append(metric_tags, qtype_metric.Tags...)
		qtype_metric.Tags = metric_tags
	}
	return_metrics = append(return_metrics, json_qtypes...)

	nsstat_tag := &MetricTag{"server", "nsstat"}
	return_metrics = append(return_metrics, &Metric{
		Name:      "AuthQryRej",
		Value:     int64(jsonStats.NSStats.AuthQryRej),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "CookieIn",
		Value:     int64(jsonStats.NSStats.CookieIn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "CookieMatch",
		Value:     int64(jsonStats.NSStats.CookieMatch),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "CookieNew",
		Value:     int64(jsonStats.NSStats.CookieNew),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "ECSOpt",
		Value:     int64(jsonStats.NSStats.ECSOpt),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryAuthAns",
		Value:     int64(jsonStats.NSStats.QryAuthAns),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryFailure",
		Value:     int64(jsonStats.NSStats.QryFailure),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryNXDOMAIN",
		Value:     int64(jsonStats.NSStats.QryNXDOMAIN),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryNoauthAns",
		Value:     int64(jsonStats.NSStats.QryNoauthAns),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryNxrrset",
		Value:     int64(jsonStats.NSStats.QryNxrrset),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryReferral",
		Value:     int64(jsonStats.NSStats.QryReferral),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QrySuccess",
		Value:     int64(jsonStats.NSStats.QrySuccess),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryTCP",
		Value:     int64(jsonStats.NSStats.QryTCP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "QryUDP",
		Value:     int64(jsonStats.NSStats.QryUDP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "RecQryRej",
		Value:     int64(jsonStats.NSStats.RecQryRej),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "ReqEdns0",
		Value:     int64(jsonStats.NSStats.ReqEdns0),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "ReqTCP",
		Value:     int64(jsonStats.NSStats.ReqTCP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Requestv4",
		Value:     int64(jsonStats.NSStats.Requestv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Requestv6",
		Value:     int64(jsonStats.NSStats.Requestv6),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "RespEDNS0",
		Value:     int64(jsonStats.NSStats.RespEDNS0),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Response",
		Value:     int64(jsonStats.NSStats.Response),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCPConnHighWater",
		Value:     int64(jsonStats.NSStats.TCPConnHighWater),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TruncatedResp",
		Value:     int64(jsonStats.NSStats.TruncatedResp),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})

	zone_tag := &MetricTag{"server", "zonestats"}
	return_metrics = append(return_metrics, &Metric{
		Name:      "AXFRReqv4",
		Value:     int64(jsonStats.ZoneStats.AXFRReqv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "IXFRReqv4",
		Value:     int64(jsonStats.ZoneStats.IXFRReqv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "NotifyInv4",
		Value:     int64(jsonStats.ZoneStats.NotifyInv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "SOAOutv4",
		Value:     int64(jsonStats.ZoneStats.SOAOutv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "XfrSuccess",
		Value:     int64(jsonStats.ZoneStats.XfrSuccess),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})

	bind_view_bind_tag := &MetricTag{"view", "_bind"}
	bind_view_bind_metrics := jsonStats.Views.Bind.toMetrics(jsonStats.CurrentTime)
	for _, bind_view_bind_metric := range bind_view_bind_metrics {
		view_metrics := make([]*MetricTag, 0, len(bind_view_bind_metric.Tags)+1)
		view_metrics = append(view_metrics, bind_view_bind_tag)
		view_metrics = append(view_metrics, bind_view_bind_metric.Tags...)
		bind_view_bind_metric.Tags = view_metrics
	}
	return_metrics = append(return_metrics, bind_view_bind_metrics...)

	bind_view_default_tag := &MetricTag{"view", "_default"}
	bind_view_default_metrics := jsonStats.Views.Default.toMetrics(jsonStats.CurrentTime)
	for _, bind_view_default_metric := range bind_view_default_metrics {
		view_metrics := make([]*MetricTag, 0, len(bind_view_default_metric.Tags)+1)
		view_metrics = append(view_metrics, bind_view_default_tag)
		view_metrics = append(view_metrics, bind_view_default_metric.Tags...)
		bind_view_default_metric.Tags = view_metrics
	}
	return_metrics = append(return_metrics, bind_view_default_metrics...)

	sockstats_tag := &MetricTag{"server", "sockstats"}
	return_metrics = append(return_metrics, &Metric{
		Name:      "RawActive",
		Value:     int64(jsonStats.SocketStats.RawActive),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "RawOpen",
		Value:     int64(jsonStats.SocketStats.RawOpen),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4Accept",
		Value:     int64(jsonStats.SocketStats.TCP4Accept),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4Active",
		Value:     int64(jsonStats.SocketStats.TCP4Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4Close",
		Value:     int64(jsonStats.SocketStats.TCP4Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4Conn",
		Value:     int64(jsonStats.SocketStats.TCP4Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4Open",
		Value:     int64(jsonStats.SocketStats.TCP4Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP4RecvErr",
		Value:     int64(jsonStats.SocketStats.TCP4RecvErr),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP6Accept",
		Value:     int64(jsonStats.SocketStats.TCP6Accept),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP6Active",
		Value:     int64(jsonStats.SocketStats.TCP6Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP6Close",
		Value:     int64(jsonStats.SocketStats.TCP6Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP6Conn",
		Value:     int64(jsonStats.SocketStats.TCP6Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TCP6Open",
		Value:     int64(jsonStats.SocketStats.TCP6Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP4Active",
		Value:     int64(jsonStats.SocketStats.UDP4Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP4Close",
		Value:     int64(jsonStats.SocketStats.UDP4Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP4Open",
		Value:     int64(jsonStats.SocketStats.UDP4Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP6Active",
		Value:     int64(jsonStats.SocketStats.UDP6Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP6Close",
		Value:     int64(jsonStats.SocketStats.UDP6Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP6Conn",
		Value:     int64(jsonStats.SocketStats.UDP6Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "UDP6Open",
		Value:     int64(jsonStats.SocketStats.UDP6Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})

	socket_mgr_metrics := make([]*Metric, 0)
	socket_mgr_tag := &MetricTag{"server", "socketmgr"}
	for _, socket := range jsonStats.SocketMgr.Sockets {
		socket_metric := socket.toMetric(jsonStats.CurrentTime)
		if socket_metric.Name != "" {
			socket_tags := make([]*MetricTag, 0, len(socket_metric.Tags)+1)
			socket_tags = append(socket_tags, socket_mgr_tag)
			socket_tags = append(socket_tags, socket_metric.Tags...)
			socket_metric.Tags = socket_tags
			socket_mgr_metrics = append(socket_mgr_metrics, socket_metric)
		}
	}
	return_metrics = append(return_metrics, socket_mgr_metrics...)

	task_mgr_metrics := make([]*Metric, 0)
	for _, task := range jsonStats.TaskMgr.Tasks {
		task_metric := task.toMetric(jsonStats.CurrentTime)
		if task_metric.Name != "" {
			task_mgr_metrics = append(task_mgr_metrics, task_metric)
		}
	}
	return_metrics = append(return_metrics, task_mgr_metrics...)

	memory_tag := &MetricTag{"server", "memory"}
	return_metrics = append(return_metrics, &Metric{
		Name:      "BlockSize",
		Value:     int64(jsonStats.Memory.BlockSize),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "ContextSize",
		Value:     int64(jsonStats.Memory.ContextSize),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "InUse",
		Value:     int64(jsonStats.Memory.InUse),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Lost",
		Value:     int64(jsonStats.Memory.Lost),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "Malloced",
		Value:     int64(jsonStats.Memory.Malloced),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	return_metrics = append(return_metrics, &Metric{
		Name:      "TotalUse",
		Value:     int64(jsonStats.Memory.TotalUse),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})

	context_tag := &MetricTag{"server", "context"}
	context_metrics := make([]*Metric, 0)
	for _, context := range jsonStats.Memory.Contexts {
		context_metric := context.toMetric(jsonStats.CurrentTime)
		for _, context_metric := range context_metric {
			context_metric_tags := make([]*MetricTag, 0, len(context_metric.Tags)+1)
			context_metric_tags = append(context_metric_tags, context_tag)
			context_metric_tags = append(context_metric_tags, context_metric.Tags...)
			context_metric.Tags = context_metric_tags
		}
		context_metrics = append(context_metrics, context_metric...)
	}
	return_metrics = append(return_metrics, context_metrics...)

	plugin.returnMetrics = return_metrics
	return nil
}

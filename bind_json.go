package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
	Traffic Traffic `json:"traffic"`
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
	R17         int `json:"17,omitempty"`
	R18         int `json:"18,omitempty"`
	R19         int `json:"19,omitempty"`
	R20         int `json:"20,omitempty"`
	R21         int `json:"21,omitempty"`
	R22         int `json:"22,omitempty"`
	AuthQryRej  int `json:"AuthQryRej,omitempty"`
	Badcookie   int `json:"BADCOOKIE,omitempty"`
	Badvers     int `json:"BADVERS,omitempty"`
	Formerr     int `json:"FORMERR,omitempty"`
	Noerror     int `json:"NOERROR,omitempty"`
	Notauth     int `json:"NOTAUTH,omitempty"`
	Notimp      int `json:"NOTIMP,omitempty"`
	Notzone     int `json:"NOTZONE,omitempty"`
	Nxdomain    int `json:"NXDOMAIN,omitempty"`
	Nxrrset     int `json:"NXRRSET,omitempty"`
	QryAuthAns  int `json:"QryAuthAns,omitempty"`
	QryNXDOMAIN int `json:"QryNXDOMAIN,omitempty"`
	QryNxrrset  int `json:"QryNxrrset,omitempty"`
	QrySuccess  int `json:"QrySuccess,omitempty"`
	QryTCP      int `json:"QryTCP,omitempty"`
	QryUDP      int `json:"QryUDP,omitempty"`
	Refused     int `json:"REFUSED,omitempty"`
	Reserved11  int `json:"RESERVED11,omitempty"`
	Reserved12  int `json:"RESERVED12,omitempty"`
	Reserved13  int `json:"RESERVED13,omitempty"`
	Reserved14  int `json:"RESERVED14,omitempty"`
	Reserved15  int `json:"RESERVED15,omitempty"`
	RecQryRej   int `json:"RecQryRej,omitempty"`
	Servfail    int `json:"SERVFAIL,omitempty"`
	UpdateDone  int `json:"UpdateDone,omitempty"`
	XfrRej      int `json:"XfrRej,omitempty"`
	XfrReqDone  int `json:"XfrReqDone,omitempty"`
	Yxdomain    int `json:"YXDOMAIN,omitempty"`
	Yxrrset     int `json:"YXRRSET,omitempty"`
}

func (r *RCode) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
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
		Name:      "AuthQryRej",
		Value:     int64(r.AuthQryRej),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Badcookie",
		Value:     int64(r.Badcookie),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Badvers",
		Value:     int64(r.Badvers),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Formerr",
		Value:     int64(r.Formerr),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Noerror",
		Value:     int64(r.Noerror),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Notauth",
		Value:     int64(r.Notauth),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Notimp",
		Value:     int64(r.Notimp),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Notzone",
		Value:     int64(r.Notzone),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Nxdomain",
		Value:     int64(r.Nxdomain),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Nxrrset",
		Value:     int64(r.Nxrrset),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QryAuthAns",
		Value:     int64(r.QryAuthAns),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QryNXDOMAIN",
		Value:     int64(r.QryNXDOMAIN),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QryNxrrset",
		Value:     int64(r.QryNxrrset),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QrySuccess",
		Value:     int64(r.QrySuccess),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QryTCP",
		Value:     int64(r.QryTCP),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "QryUDP",
		Value:     int64(r.QryUDP),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Refused",
		Value:     int64(r.Refused),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Reserved11",
		Value:     int64(r.Reserved11),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Reserved12",
		Value:     int64(r.Reserved12),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Reserved13",
		Value:     int64(r.Reserved13),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Reserved14",
		Value:     int64(r.Reserved14),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Reserved15",
		Value:     int64(r.Reserved15),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "RecQryRej",
		Value:     int64(r.RecQryRej),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Servfail",
		Value:     int64(r.Servfail),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "UpdateDone",
		Value:     int64(r.UpdateDone),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "XfrRej",
		Value:     int64(r.XfrRej),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "XfrReqDone",
		Value:     int64(r.XfrReqDone),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Yxdomain",
		Value:     int64(r.Yxdomain),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	})
	metrics = append(metrics, &Metric{
		Name:      "Yxrrset",
		Value:     int64(r.Yxrrset),
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
	stats_tag := &MetricTag{"type", "resstats"}
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

	qtypes_tag := &MetricTag{"type", "resqtype"}
	resolver_qtypes_metrics := bv.Resolver.QTypes.toMetrics(metric_time)
	for _, metric := range resolver_qtypes_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, qtypes_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	metrics = append(metrics, resolver_qtypes_metrics...)
	cache_tag := &MetricTag{"type", "cache"}
	resolver_cache_metrics := bv.Resolver.Cache.toMetrics(metric_time)
	for _, metric := range resolver_cache_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, cache_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	metrics = append(metrics, resolver_cache_metrics...)

	cachestats_tag := &MetricTag{"type", "cachestats"}
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
	adb_tag := &MetricTag{"type", "adbstat"}
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
		Name:      t.Name,
		Value:     int64(t.Events),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	}
	task_metric.Tags = append(task_metric.Tags, &MetricTag{"task_id", t.Id})
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
	context_id_tag := &MetricTag{"context_id", c.Id}

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
		Name:      "InUse",
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
	Name          string    `json:"name"`
	Class         string    `json:"class"`
	Serial        int       `json:"serial"`
	Type          string    `json:"type"`
	Loaded        time.Time `json:"loaded"`
	Expires       time.Time `json:"expires,omitempty"`
	Refresh       time.Time `json:"refresh,omitempty"`
	RCodes        RCode     `json:"rcodes,omitempty"`
	QTypes        QTypes    `json:"qtypes"`
	DnsSecSign    DnsSec    `json:"dnssec-sign,omitempty"`
	DnsSecRefresh DnsSec    `json:"dnssec-refresh,omitempty"`
}

func (z *ZoneView) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	zone_name_tag := &MetricTag{"zone", strings.Replace(z.Name, ".", "_", -1)}
	zone_class_tag := &MetricTag{"class", z.Class}
	zone_type_tag := &MetricTag{"type", z.Type}

	rcode_tag := &MetricTag{"type", "rcode"}
	rcode_metrics := z.RCodes.toMetrics(metric_time)
	for _, rcode_metric := range rcode_metrics {
		if rcode_metric.Value != 0 {
			rcode_metric_tags := make([]*MetricTag, 0, len(rcode_metric.Tags)+3)
			rcode_metric_tags = append(rcode_metric_tags, zone_name_tag)
			rcode_metric_tags = append(rcode_metric_tags, zone_class_tag)
			rcode_metric_tags = append(rcode_metric_tags, zone_type_tag)
			rcode_metric_tags = append(rcode_metric_tags, rcode_tag)
			rcode_metric_tags = append(rcode_metric_tags, rcode_metric.Tags...)
			rcode_metric.Tags = rcode_metric_tags
			metrics = append(metrics, rcode_metric)
		}
	}

	qtype_tag := &MetricTag{"type", "qtype"}
	qtype_metrics := z.QTypes.toMetrics(metric_time)
	for _, qtype_metric := range qtype_metrics {
		if qtype_metric.Value != 0 {
			qtype_metric_tags := make([]*MetricTag, 0, len(qtype_metric.Tags)+3)
			qtype_metric_tags = append(qtype_metric_tags, zone_name_tag)
			qtype_metric_tags = append(qtype_metric_tags, zone_class_tag)
			qtype_metric_tags = append(qtype_metric_tags, zone_type_tag)
			qtype_metric_tags = append(qtype_metric_tags, qtype_tag)
			qtype_metric_tags = append(qtype_metric_tags, qtype_metric.Tags...)
			qtype_metric.Tags = qtype_metric_tags
			metrics = append(metrics, qtype_metric)
		}
	}

	dnssecsign_tag := &MetricTag{"type", "dnssec-sign"}
	dnssecsign_metrics := z.DnsSecSign.toMetrics(metric_time)
	for _, dnssecsign_metric := range dnssecsign_metrics {
		if dnssecsign_metric.Value != 0 {
			dnssecsign_metric_tags := make([]*MetricTag, 0, len(dnssecsign_metric.Tags)+4)
			dnssecsign_metric_tags = append(dnssecsign_metric_tags, zone_name_tag)
			dnssecsign_metric_tags = append(dnssecsign_metric_tags, zone_class_tag)
			dnssecsign_metric_tags = append(dnssecsign_metric_tags, zone_type_tag)
			dnssecsign_metric_tags = append(dnssecsign_metric_tags, dnssecsign_tag)
			dnssecsign_metric_tags = append(dnssecsign_metric_tags, dnssecsign_metric.Tags...)
			dnssecsign_metric.Tags = dnssecsign_metric_tags
			metrics = append(metrics, dnssecsign_metric)
		}
	}
	dnsssecrefresh_tag := &MetricTag{"type", "dnssec-refresh"}
	dnsssecrefresh_metrics := z.DnsSecRefresh.toMetrics(metric_time)
	for _, dnsssecrefresh_metric := range dnsssecrefresh_metrics {
		if dnsssecrefresh_metric.Value != 0 {
			dnsssecrefresh_metric_tags := make([]*MetricTag, 0, len(dnsssecrefresh_metric.Tags)+4)
			dnsssecrefresh_metric_tags = append(dnsssecrefresh_metric_tags, zone_name_tag)
			dnsssecrefresh_metric_tags = append(dnsssecrefresh_metric_tags, zone_class_tag)
			dnsssecrefresh_metric_tags = append(dnsssecrefresh_metric_tags, zone_type_tag)
			dnsssecrefresh_metric_tags = append(dnsssecrefresh_metric_tags, dnsssecrefresh_tag)
			dnsssecrefresh_metric_tags = append(dnsssecrefresh_metric_tags, dnsssecrefresh_metric.Tags...)
			dnsssecrefresh_metric.Tags = dnsssecrefresh_metric_tags
			metrics = append(metrics, dnsssecrefresh_metric)
		}
	}

	return metrics
}

type DnsSec struct {
	DnsSecTypes []struct {
		Name  string
		Value int64
	}
}

func (d *DnsSec) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	for _, dnssec_type := range d.DnsSecTypes {
		metrics = append(metrics, &Metric{
			Name:      dnssec_type.Name,
			Value:     dnssec_type.Value,
			Timestamp: metric_time,
			Tags:      []*MetricTag{},
		})
	}
	return metrics
}

func (d *DnsSec) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	dnssec := string(data)

	// String the beginning and ending curly braces
	dnssec = dnssec[1 : len(dnssec)-1]

	// Strip beginning and ending whitespace
	dnssec = strings.TrimSpace(dnssec)

	// Split the string into key value pairs
	d.DnsSecTypes = make([]struct {
		Name  string
		Value int64
	}, 0)
	dnssec_pairs := strings.Split(dnssec, ",")
	for _, pair := range dnssec_pairs {
		// Split the key value pair
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			continue
		}
		k := strings.Replace(strings.TrimSpace(kv[0]), `"`, "", -1)
		v, _ := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)

		d.DnsSecTypes = append(d.DnsSecTypes, struct {
			Name  string
			Value int64
		}{
			Name:  k,
			Value: v,
		})
	}

	return nil
}

type Traffic struct {
	TrafficTypes []struct {
		Protocol string
		Type     string
		IPVer    string
		Name     string
		Value    int64
	}
}

func (t *Traffic) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0)
	for _, traffic_type := range t.TrafficTypes {
		protocol_tag := &MetricTag{"protocol", traffic_type.Protocol}
		type_tag := &MetricTag{"type", traffic_type.Type}
		ipver_tag := &MetricTag{"ipver", traffic_type.IPVer}
		metrics = append(metrics, &Metric{
			Name:      traffic_type.Name,
			Value:     traffic_type.Value,
			Timestamp: metric_time,
			Tags:      []*MetricTag{ipver_tag, protocol_tag, type_tag},
		})
	}
	return metrics
}

func (t *Traffic) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	traffic := string(data)

	// String the beginning and ending curly braces
	traffic = traffic[1 : len(traffic)-1]

	// Strip beginning and ending whitespace
	traffic = strings.TrimSpace(traffic)

	traffic_attribute, _ := regexp.Compile("dns-(?P<protocol>udp|tcp|tcp6|udp6)-(?P<type>(?:requests|responses)-sizes)-(?:received|sent)-(?P<ipver>ipv4|ipv6)")
	for {
		if len(traffic) == 0 {
			break
		}

		// Parse the attribute name from the string
		attribute_name := traffic[:strings.Index(traffic, ":")+1]
		traffic = strings.Replace(traffic, attribute_name, "", 1)

		// Parse the attribute name into pieces
		attribute_name = strings.Replace(attribute_name, `"`, "", -1)
		attribute_name = strings.Replace(attribute_name, ":", "", -1)
		attribute_pieces := traffic_attribute.FindStringSubmatch(attribute_name)
		protocol := attribute_pieces[1]
		traffic_type := attribute_pieces[2]
		traffic_type = strings.Replace(traffic_type, "s-sizes", "-size", -1)
		ipver := attribute_pieces[3]

		// Parse the attribute value from the string
		attribute_value := traffic[:strings.Index(traffic, "}")+1]

		// Strip the attribute value from the string
		traffic = strings.Replace(traffic, attribute_value, "", 1)

		// Strip the beginning and ending curly braces
		attribute_value = attribute_value[1 : len(attribute_value)-1]

		// Strip beginning and ending whitespace
		attribute_value = strings.TrimSpace(attribute_value)

		// Parse the attribute value into pieces
		attribute_value_pieces := strings.Split(attribute_value, ",")

		if len(attribute_value_pieces) > 0 {
			for _, piece := range attribute_value_pieces {
				// Split the key value pair
				kv := strings.Split(piece, ":")
				if len(kv) == 2 {
					k := strings.Replace(strings.TrimSpace(kv[0]), `"`, "", -1)
					v, _ := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)

					traffic_type := struct {
						Protocol string
						Type     string
						IPVer    string
						Name     string
						Value    int64
					}{
						Protocol: protocol,
						Type:     traffic_type,
						IPVer:    ipver,
						Name:     k,
						Value:    v,
					}

					t.TrafficTypes = append(t.TrafficTypes, traffic_type)
				}
			}
		}

		// Strip whitespace from the string
		traffic = strings.TrimSpace(traffic)

		// Strip the comma from the string
		if strings.Index(traffic, ",") == 0 {
			traffic = traffic[1:]
		}
	}

	return nil
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

	opscodes_metrics := make([]*Metric, 0)
	opcodes_tag := &MetricTag{"server", "opcodes"}
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Query",
		Value:     int64(jsonStats.OpCodes.Query),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "IQuery",
		Value:     int64(jsonStats.OpCodes.IQuery),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Status",
		Value:     int64(jsonStats.OpCodes.Status),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Notify",
		Value:     int64(jsonStats.OpCodes.Notify),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Update",
		Value:     int64(jsonStats.OpCodes.Update),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved6",
		Value:     int64(jsonStats.OpCodes.Reserved6),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved7",
		Value:     int64(jsonStats.OpCodes.Reserved7),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved8",
		Value:     int64(jsonStats.OpCodes.Reserved8),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved9",
		Value:     int64(jsonStats.OpCodes.Reserved9),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved10",
		Value:     int64(jsonStats.OpCodes.Reserved10),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved11",
		Value:     int64(jsonStats.OpCodes.Reserved11),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved12",
		Value:     int64(jsonStats.OpCodes.Reserved12),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved13",
		Value:     int64(jsonStats.OpCodes.Reserved13),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved14",
		Value:     int64(jsonStats.OpCodes.Reserved14),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	opscodes_metrics = append(opscodes_metrics, &Metric{
		Name:      "Reserved15",
		Value:     int64(jsonStats.OpCodes.Reserved15),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{opcodes_tag},
	})
	for _, opcode_metric := range opscodes_metrics {
		if opcode_metric.Value != 0 {
			return_metrics = append(return_metrics, opcode_metric)
		}
	}

	rcodes_tag := &MetricTag{"server", "rcodes"}
	json_rcodes := jsonStats.RCodes.toMetrics(jsonStats.CurrentTime)
	for _, rcode_metric := range json_rcodes {
		if rcode_metric.Value != 0 {
			metric_tags := make([]*MetricTag, 0, len(rcode_metric.Tags)+1)
			metric_tags = append(metric_tags, rcodes_tag)
			metric_tags = append(metric_tags, rcode_metric.Tags...)
			rcode_metric.Tags = metric_tags
			return_metrics = append(return_metrics, rcode_metric)
		}
	}

	qtypes_tag := &MetricTag{"server", "qtypes"}
	json_qtypes := jsonStats.QTypes.toMetrics(jsonStats.CurrentTime)
	for _, qtype_metric := range json_qtypes {
		if qtype_metric.Value != 0 {
			metric_tags := make([]*MetricTag, 0, len(qtype_metric.Tags)+1)
			metric_tags = append(metric_tags, qtypes_tag)
			metric_tags = append(metric_tags, qtype_metric.Tags...)
			qtype_metric.Tags = metric_tags
			return_metrics = append(return_metrics, qtype_metric)
		}
	}

	nsstat_tag := &MetricTag{"server", "nsstat"}
	nsstat_metrics := make([]*Metric, 0)
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "AuthQryRej",
		Value:     int64(jsonStats.NSStats.AuthQryRej),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "CookieIn",
		Value:     int64(jsonStats.NSStats.CookieIn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "CookieMatch",
		Value:     int64(jsonStats.NSStats.CookieMatch),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "CookieNew",
		Value:     int64(jsonStats.NSStats.CookieNew),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "ECSOpt",
		Value:     int64(jsonStats.NSStats.ECSOpt),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryAuthAns",
		Value:     int64(jsonStats.NSStats.QryAuthAns),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryFailure",
		Value:     int64(jsonStats.NSStats.QryFailure),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryNXDOMAIN",
		Value:     int64(jsonStats.NSStats.QryNXDOMAIN),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryNoauthAns",
		Value:     int64(jsonStats.NSStats.QryNoauthAns),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryNxrrset",
		Value:     int64(jsonStats.NSStats.QryNxrrset),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryReferral",
		Value:     int64(jsonStats.NSStats.QryReferral),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QrySuccess",
		Value:     int64(jsonStats.NSStats.QrySuccess),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryTCP",
		Value:     int64(jsonStats.NSStats.QryTCP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "QryUDP",
		Value:     int64(jsonStats.NSStats.QryUDP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "RecQryRej",
		Value:     int64(jsonStats.NSStats.RecQryRej),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "ReqEdns0",
		Value:     int64(jsonStats.NSStats.ReqEdns0),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "ReqTCP",
		Value:     int64(jsonStats.NSStats.ReqTCP),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "Requestv4",
		Value:     int64(jsonStats.NSStats.Requestv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "Requestv6",
		Value:     int64(jsonStats.NSStats.Requestv6),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "RespEDNS0",
		Value:     int64(jsonStats.NSStats.RespEDNS0),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "Response",
		Value:     int64(jsonStats.NSStats.Response),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "TCPConnHighWater",
		Value:     int64(jsonStats.NSStats.TCPConnHighWater),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	nsstat_metrics = append(nsstat_metrics, &Metric{
		Name:      "TruncatedResp",
		Value:     int64(jsonStats.NSStats.TruncatedResp),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{nsstat_tag},
	})
	for _, nsstat_metric := range nsstat_metrics {
		if nsstat_metric.Value != 0 {
			return_metrics = append(return_metrics, nsstat_metric)
		}
	}

	zone_tag := &MetricTag{"server", "zonestats"}
	zonestats_metrics := make([]*Metric, 0)
	zonestats_metrics = append(zonestats_metrics, &Metric{
		Name:      "AXFRReqv4",
		Value:     int64(jsonStats.ZoneStats.AXFRReqv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	zonestats_metrics = append(zonestats_metrics, &Metric{
		Name:      "IXFRReqv4",
		Value:     int64(jsonStats.ZoneStats.IXFRReqv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	zonestats_metrics = append(zonestats_metrics, &Metric{
		Name:      "NotifyInv4",
		Value:     int64(jsonStats.ZoneStats.NotifyInv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	zonestats_metrics = append(zonestats_metrics, &Metric{
		Name:      "SOAOutv4",
		Value:     int64(jsonStats.ZoneStats.SOAOutv4),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	zonestats_metrics = append(zonestats_metrics, &Metric{
		Name:      "XfrSuccess",
		Value:     int64(jsonStats.ZoneStats.XfrSuccess),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{zone_tag},
	})
	for _, zonestats_metric := range zonestats_metrics {
		if zonestats_metric.Value != 0 {
			return_metrics = append(return_metrics, zonestats_metric)
		}
	}

	bind_view_bind_tag := &MetricTag{"view", "_bind"}
	bind_view_bind_metrics := jsonStats.Views.Bind.toMetrics(jsonStats.CurrentTime)
	for _, bind_view_bind_metric := range bind_view_bind_metrics {
		if bind_view_bind_metric.Value != 0 {
			view_metrics := make([]*MetricTag, 0, len(bind_view_bind_metric.Tags)+1)
			view_metrics = append(view_metrics, bind_view_bind_tag)
			view_metrics = append(view_metrics, bind_view_bind_metric.Tags...)
			bind_view_bind_metric.Tags = view_metrics
			return_metrics = append(return_metrics, bind_view_bind_metric)
		}
	}

	bind_view_default_tag := &MetricTag{"view", "_default"}
	bind_view_default_metrics := jsonStats.Views.Default.toMetrics(jsonStats.CurrentTime)
	for _, bind_view_default_metric := range bind_view_default_metrics {
		if bind_view_default_metric.Value != 0 {
			view_metrics := make([]*MetricTag, 0, len(bind_view_default_metric.Tags)+1)
			view_metrics = append(view_metrics, bind_view_default_tag)
			view_metrics = append(view_metrics, bind_view_default_metric.Tags...)
			bind_view_default_metric.Tags = view_metrics
			return_metrics = append(return_metrics, bind_view_default_metric)
		}
	}

	sockstats_tag := &MetricTag{"server", "sockstats"}
	sockstats_metrics := make([]*Metric, 0)
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "RawActive",
		Value:     int64(jsonStats.SocketStats.RawActive),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "RawOpen",
		Value:     int64(jsonStats.SocketStats.RawOpen),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4Accept",
		Value:     int64(jsonStats.SocketStats.TCP4Accept),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4Active",
		Value:     int64(jsonStats.SocketStats.TCP4Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4Close",
		Value:     int64(jsonStats.SocketStats.TCP4Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4Conn",
		Value:     int64(jsonStats.SocketStats.TCP4Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4Open",
		Value:     int64(jsonStats.SocketStats.TCP4Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP4RecvErr",
		Value:     int64(jsonStats.SocketStats.TCP4RecvErr),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP6Accept",
		Value:     int64(jsonStats.SocketStats.TCP6Accept),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP6Active",
		Value:     int64(jsonStats.SocketStats.TCP6Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP6Close",
		Value:     int64(jsonStats.SocketStats.TCP6Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP6Conn",
		Value:     int64(jsonStats.SocketStats.TCP6Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "TCP6Open",
		Value:     int64(jsonStats.SocketStats.TCP6Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP4Active",
		Value:     int64(jsonStats.SocketStats.UDP4Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP4Close",
		Value:     int64(jsonStats.SocketStats.UDP4Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP4Open",
		Value:     int64(jsonStats.SocketStats.UDP4Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP6Active",
		Value:     int64(jsonStats.SocketStats.UDP6Active),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP6Close",
		Value:     int64(jsonStats.SocketStats.UDP6Close),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP6Conn",
		Value:     int64(jsonStats.SocketStats.UDP6Conn),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	sockstats_metrics = append(sockstats_metrics, &Metric{
		Name:      "UDP6Open",
		Value:     int64(jsonStats.SocketStats.UDP6Open),
		Timestamp: jsonStats.CurrentTime,
		Tags:      []*MetricTag{sockstats_tag},
	})
	for _, sockstats_metric := range sockstats_metrics {
		if sockstats_metric.Value != 0 {
			return_metrics = append(return_metrics, sockstats_metric)
		}
	}

	socket_mgr_metrics := make([]*Metric, 0)
	socket_mgr_tag := &MetricTag{"server", "socketmgr"}
	for _, socket := range jsonStats.SocketMgr.Sockets {
		socket_metric := socket.toMetric(jsonStats.CurrentTime)
		if socket_metric.Name != "" {
			socket_tags := make([]*MetricTag, 0, len(socket_metric.Tags)+1)
			socket_tags = append(socket_tags, socket_mgr_tag)
			socket_tags = append(socket_tags, socket_metric.Tags...)
			socket_metric.Tags = socket_tags
			if socket_metric.Value != 0 {
				socket_mgr_metrics = append(socket_mgr_metrics, socket_metric)
			}
		}
	}
	return_metrics = append(return_metrics, socket_mgr_metrics...)

	task_mgr_tag := &MetricTag{"server", "taskmgr"}
	task_mgr_metrics := make([]*Metric, 0)
	for _, task := range jsonStats.TaskMgr.Tasks {
		task_metric := task.toMetric(jsonStats.CurrentTime)
		if task_metric.Name != "" {
			task_tags := make([]*MetricTag, 0, len(task_metric.Tags)+1)
			task_tags = append(task_tags, task_mgr_tag)
			task_tags = append(task_tags, task_metric.Tags...)
			task_metric.Tags = task_tags
			if task_metric.Value != 0 {
				task_mgr_metrics = append(task_mgr_metrics, task_metric)
			}
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
	for _, context := range jsonStats.Memory.Contexts {
		context_metrics := context.toMetric(jsonStats.CurrentTime)
		for _, context_metric := range context_metrics {
			context_metric_tags := make([]*MetricTag, 0, len(context_metric.Tags)+1)
			context_metric_tags = append(context_metric_tags, context_tag)
			context_metric_tags = append(context_metric_tags, context_metric.Tags...)
			context_metric.Tags = context_metric_tags
			if context_metric.Value != 0 {
				return_metrics = append(return_metrics, context_metric)
			}
		}
	}
	// traffic_tag := &MetricTag{"server", "traffic"}
	traffic_metrics := jsonStats.Traffic.toMetrics(jsonStats.CurrentTime)
	for _, traffic_metric := range traffic_metrics {
		traffic_metric_tags := make([]*MetricTag, 0, len(traffic_metric.Tags))
		// traffic_metric_tags = append(traffic_metric_tags, traffic_tag)
		traffic_metric_tags = append(traffic_metric_tags, traffic_metric.Tags...)
		traffic_metric.Tags = traffic_metric_tags
		if traffic_metric.Value != 0 {
			return_metrics = append(return_metrics, traffic_metric)
		}
	}

	plugin.returnMetrics = return_metrics
	return nil
}

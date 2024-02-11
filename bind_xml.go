package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type bindXmlStats struct {
	Memory struct {
		Contexts struct {
			Context []struct {
				Blocksize   string `xml:"blocksize"`
				Hiwater     int    `xml:"hiwater"`
				ID          string `xml:"id"`
				Inuse       int    `xml:"inuse"`
				Lowater     int    `xml:"lowater"`
				Malloced    int    `xml:"malloced"`
				Maxinuse    int    `xml:"maxinuse"`
				Maxmalloced int    `xml:"maxmalloced"`
				Name        string `xml:"name"`
				Pools       int    `xml:"pools"`
				References  int    `xml:"references"`
				Total       int    `xml:"total"`
			} `xml:"context"`
		} `xml:"contexts"`
		Summary struct {
			BlockSize   int `xml:"BlockSize"`
			ContextSize int `xml:"ContextSize"`
			InUse       int `xml:"InUse"`
			Lost        int `xml:"Lost"`
			Malloced    int `xml:"Malloced"`
			TotalUse    int `xml:"TotalUse"`
		} `xml:"summary"`
	} `xml:"memory"`
	Server struct {
		BootTime    time.Time      `xml:"boot-time"`
		ConfigTime  time.Time      `xml:"config-time"`
		Counters    []*XmlCounters `xml:"counters"`
		CurrentTime time.Time      `xml:"current-time"`
		Version     string         `xml:"version"`
	} `xml:"server"`
	Socketmgr struct {
		Sockets struct {
			Socket []struct {
				ID           string  `xml:"id"`
				LocalAddress *string `xml:"local-address"`
				Name         *string `xml:"name"`
				PeerAddress  string  `xml:"peer-address"`
				References   int     `xml:"references"`
				States       struct {
					State []string `xml:"state"`
				} `xml:"states"`
				Type string `xml:"type"`
			} `xml:"socket"`
		} `xml:"sockets"`
	} `xml:"socketmgr"`
	Taskmgr struct {
		Tasks struct {
			Task []struct {
				Events     int     `xml:"events"`
				ID         string  `xml:"id"`
				Name       *string `xml:"name"`
				Quantum    int     `xml:"quantum"`
				References int     `xml:"references"`
				State      string  `xml:"state"`
			} `xml:"task"`
		} `xml:"tasks"`
		ThreadModel struct {
			DefaultQuantum int    `xml:"default-quantum"`
			Type           string `xml:"type"`
		} `xml:"thread-model"`
	} `xml:"taskmgr"`
	Traffic struct {
		Ipv4 *XmlIp `xml:"ipv4"`
		Ipv6 *XmlIp `xml:"ipv6"`
	} `xml:"traffic"`
	Views struct {
		View []struct {
			Name  string `xml:"name,attr"`
			Cache struct {
				Name  string `xml:"name,attr"`
				Rrset []struct {
					Counter int    `xml:"counter"`
					Name    string `xml:"name"`
				} `xml:"rrset"`
			} `xml:"cache"`
			Counters []*XmlCounters `xml:"counters"`
			Zones    struct {
				Zone []struct {
					Name       string         `xml:"name,attr"`
					Rdataclass string         `xml:"rdataclass,attr"`
					Counters   []*XmlCounters `xml:"counters"`
					Expires    *time.Time     `xml:"expires"`
					Loaded     time.Time      `xml:"loaded"`
					Refresh    *time.Time     `xml:"refresh"`
					Serial     int            `xml:"serial"`
					Type       string         `xml:"type"`
				} `xml:"zone"`
			} `xml:"zones"`
		} `xml:"view"`
	} `xml:"views"`
}

type XmlCounter struct {
	Name     string `xml:"name,attr"`
	CharData string `xml:",chardata"`
}

func (xc *XmlCounter) toMetric(metric_time time.Time) *Metric {
	metric_value, _ := strconv.Atoi(xc.CharData)
	return &Metric{
		Name:      xc.Name,
		Value:     int64(metric_value),
		Timestamp: metric_time,
		Tags:      []*MetricTag{},
	}
}

type XmlCounters struct {
	Type    string        `xml:"type,attr"`
	Counter []*XmlCounter `xml:"counter"`
}

func (xc *XmlCounters) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0, len(xc.Counter))
	for _, counter := range xc.Counter {
		metric := counter.toMetric(metric_time)
		metric_tags := []*MetricTag{}
		metric_tags = append(metric_tags, &MetricTag{"type", xc.Type})
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
		metrics = append(metrics, metric)
	}
	return metrics
}

type XmlIp struct {
	Tcp struct {
		Counters []*XmlCounters `xml:"counters"`
	} `xml:"tcp"`
	Udp struct {
		Counters []*XmlCounters `xml:"counters"`
	} `xml:"udp"`
}

func (xi *XmlIp) toMetrics(metric_time time.Time) []*Metric {
	metrics := make([]*Metric, 0, len(xi.Tcp.Counters)+len(xi.Udp.Counters))
	tcp_metrics := make([]*Metric, 0, len(xi.Tcp.Counters))
	for _, counter := range xi.Tcp.Counters {
		tcp_metric := counter.toMetrics(metric_time)
		for _, metric := range tcp_metric {
			metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
			metric_tag := &MetricTag{"protocol", "tcp"}
			metric_tags = append(metric_tags, metric_tag)
			metric_tags = append(metric_tags, metric.Tags...)
			metric.Tags = metric_tags
			tcp_metrics = append(tcp_metrics, metric)
		}
	}
	udp_metrics := make([]*Metric, 0, len(xi.Udp.Counters))
	for _, counter := range xi.Udp.Counters {
		udp_metric := counter.toMetrics(metric_time)
		for _, metric := range udp_metric {
			metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
			metric_tag := &MetricTag{"protocol", "udp"}
			metric_tags = append(metric_tags, metric_tag)
			metric_tags = append(metric_tags, metric.Tags...)
			metric.Tags = metric_tags
			udp_metrics = append(udp_metrics, metric)
		}
	}
	metrics = append(metrics, tcp_metrics...)
	metrics = append(metrics, udp_metrics...)
	return metrics
}

func ReadXmlStats(statsData []byte) error {
	var xmlStats bindXmlStats

	// Parse the XML statistics
	err := xml.Unmarshal(statsData, &xmlStats)
	if err != nil {
		fmt.Printf("Error parsing XML: %s\n", err)
		return err
	}

	returnMetrics := make([]*Metric, 0, 100)

	// Process the memory context statistics
	context_tag := &MetricTag{"server", "context"}
	contextMetrics := make([]*Metric, 0, 10)
	for _, context := range xmlStats.Memory.Contexts.Context {
		context_name_tag := &MetricTag{"context", context.Name}
		context_id_tag := &MetricTag{"context_id", context.ID}
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "References",
			Value:     int64(context.References),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		blocksize, err := strconv.ParseInt(context.Blocksize, 10, 64)
		if err != nil {
			blocksize = 0
		}
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Blocksize",
			Value:     int64(blocksize),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})

		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Hiwater",
			Value:     int64(context.Hiwater),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "InUse",
			Value:     int64(context.Inuse),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Lowater",
			Value:     int64(context.Lowater),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Malloced",
			Value:     int64(context.Malloced),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Maxinuse",
			Value:     int64(context.Maxinuse),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Maxmalloced",
			Value:     int64(context.Maxmalloced),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Pools",
			Value:     int64(context.Pools),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
		contextMetrics = append(contextMetrics, &Metric{
			Name:      "Total",
			Value:     int64(context.Total),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{context_tag, context_name_tag, context_id_tag},
		})
	}
	for _, context_metric := range contextMetrics {
		if context_metric.Value != 0 {
			returnMetrics = append(returnMetrics, context_metric)
		}
	}

	// Process the memory statistics
	memory_tag := &MetricTag{"server", "memory"}
	memoryMetrics := make([]*Metric, 0, 10)
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "BlockSize",
		Value:     int64(xmlStats.Memory.Summary.BlockSize),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "ContextSize",
		Value:     int64(xmlStats.Memory.Summary.ContextSize),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "InUse",
		Value:     int64(xmlStats.Memory.Summary.InUse),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "Lost",
		Value:     int64(xmlStats.Memory.Summary.Lost),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "Malloced",
		Value:     int64(xmlStats.Memory.Summary.Malloced),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "TotalUse",
		Value:     int64(xmlStats.Memory.Summary.TotalUse),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{memory_tag},
	})
	returnMetrics = append(returnMetrics, memoryMetrics...)

	// Process the socketmgr statistics
	socketMetrics := make([]*Metric, 0, 10)
	socket_mgr_tag := &MetricTag{"server", "socketmgr"}
	for _, socket := range xmlStats.Socketmgr.Sockets.Socket {
		if socket.Name != nil {
			socket_metric := &Metric{
				Name:      *socket.Name,
				Value:     int64(socket.References),
				Timestamp: xmlStats.Server.CurrentTime,
				Tags:      []*MetricTag{socket_mgr_tag},
			}
			if socket.LocalAddress != nil {
				socket_metric.Tags = append(
					socket_metric.Tags,
					&MetricTag{
						"local-address",
						strings.Replace(*socket.LocalAddress, ".", "_", -1),
					},
				)
			}
			if socket.PeerAddress != "" {
				socket_metric.Tags = append(
					socket_metric.Tags,
					&MetricTag{
						"peer-address",
						strings.Replace(socket.PeerAddress, ".", "_", -1),
					},
				)
			}
			socket_metric.Tags = append(socket_metric.Tags, &MetricTag{"type", socket.Type})
			if socket_metric.Value != 0 {
				socketMetrics = append(socketMetrics, socket_metric)
			}
		}
	}
	returnMetrics = append(returnMetrics, socketMetrics...)

	// Process the taskmgr statistics
	taskMetrics := make([]*Metric, 0, 10)
	taskmgr_tag := &MetricTag{"server", "taskmgr"}
	for _, task := range xmlStats.Taskmgr.Tasks.Task {
		if task.Name != nil {
			task_metric := &Metric{
				Name:      *task.Name,
				Value:     int64(task.Events),
				Timestamp: xmlStats.Server.CurrentTime,
				Tags:      []*MetricTag{taskmgr_tag},
			}
			task_metric.Tags = append(task_metric.Tags, &MetricTag{"task_id", task.ID})
			if task_metric.Value != 0 {
				taskMetrics = append(taskMetrics, task_metric)
			}
		}
	}
	returnMetrics = append(returnMetrics, taskMetrics...)

	// Process the traffic statistics
	trafficMetrics := make([]*Metric, 0, 10)
	ipv4_metrics := xmlStats.Traffic.Ipv4.toMetrics(xmlStats.Server.CurrentTime)
	ipv4_tag := &MetricTag{"ipver", "ipv4"}
	for _, metric := range ipv4_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, ipv4_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	trafficMetrics = append(trafficMetrics, ipv4_metrics...)
	ipv6_metrics := xmlStats.Traffic.Ipv6.toMetrics(xmlStats.Server.CurrentTime)
	ipv6_tag := &MetricTag{"ipver", "ipv6"}
	for _, metric := range ipv6_metrics {
		metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
		metric_tags = append(metric_tags, ipv6_tag)
		metric_tags = append(metric_tags, metric.Tags...)
		metric.Tags = metric_tags
	}
	trafficMetrics = append(trafficMetrics, ipv6_metrics...)
	returnMetrics = append(returnMetrics, trafficMetrics...)

	// Process the view statistics
	viewMetrics := make([]*Metric, 0, 10)
	for _, view := range xmlStats.Views.View {
		// Zones
		view_tag := &MetricTag{"view", view.Name}
		cache_tag := &MetricTag{"type", "cache"}
		for _, cache_counter := range view.Cache.Rrset {
			cache_metric := &Metric{
				Name:      cache_counter.Name,
				Value:     int64(cache_counter.Counter),
				Timestamp: xmlStats.Server.CurrentTime,
				Tags:      []*MetricTag{view_tag, cache_tag},
			}
			if cache_metric.Value != 0 {
				viewMetrics = append(viewMetrics, cache_metric)
			}
		}

		for _, view_counter := range view.Counters {
			view_counters := view_counter.toMetrics(xmlStats.Server.CurrentTime)
			for _, metric := range view_counters {
				if metric.Value != 0 {
					view_counter_metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
					view_counter_metric_tags = append(view_counter_metric_tags, view_tag)
					view_counter_metric_tags = append(view_counter_metric_tags, metric.Tags...)
					metric.Tags = view_counter_metric_tags
					viewMetrics = append(viewMetrics, metric)
				}
			}
		}

		for _, zone := range view.Zones.Zone {
			zone_tags := make([]*MetricTag, 0, 10)
			zone_tags = append(zone_tags, view_tag)
			// zone_tags = append(zone_tags, cache_tag)
			zone_tags = append(zone_tags, &MetricTag{"zone", strings.Replace(zone.Name, ".", "_", -1)})
			zone_tags = append(zone_tags, &MetricTag{"class", zone.Rdataclass})
			zone_tags = append(zone_tags, &MetricTag{"type", zone.Type})
			for _, zone_counter := range zone.Counters {
				zone_counters := zone_counter.toMetrics(xmlStats.Server.CurrentTime)
				for _, metric := range zone_counters {
					if metric.Value != 0 {
						zone_counter_metric_tags := make([]*MetricTag, 0, len(metric.Tags)+4)
						zone_counter_metric_tags = append(zone_counter_metric_tags, zone_tags...)
						zone_counter_metric_tags = append(zone_counter_metric_tags, metric.Tags...)
						metric.Tags = zone_counter_metric_tags
						viewMetrics = append(viewMetrics, metric)
					}
				}
			}
		}
	}
	returnMetrics = append(returnMetrics, viewMetrics...)

	plugin.returnMetrics = returnMetrics

	return nil
}

package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
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
	fmt.Printf("Read %d bytes of XML\n", len(statsData))

	var xmlStats bindXmlStats

	// Parse the XML statistics
	err := xml.Unmarshal(statsData, &xmlStats)
	if err != nil {
		fmt.Printf("Error parsing XML: %s\n", err)
		return err
	}

	returnMetrics := make([]*Metric, 0, 100)

	// Process the memory context statistics
	contextMetrics := make([]*Metric, 0, 10)
	for _, context := range xmlStats.Memory.Contexts.Context {
		context_metric := &Metric{
			Name:      context.ID,
			Value:     int64(context.References),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{},
		}
		if context.Name != "" {
			context_metric.Tags = append(context_metric.Tags, &MetricTag{"name", context.Name})
		}
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"blocksize", context.Blocksize})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"hiwater", strconv.Itoa(context.Hiwater)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"inuse", strconv.Itoa(context.Inuse)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"lowater", strconv.Itoa(context.Lowater)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"malloced", strconv.Itoa(context.Malloced)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"maxinuse", strconv.Itoa(context.Maxinuse)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"maxmalloced", strconv.Itoa(context.Maxmalloced)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"pools", strconv.Itoa(context.Pools)})
		context_metric.Tags = append(context_metric.Tags, &MetricTag{"total", strconv.Itoa(context.Total)})
		contextMetrics = append(contextMetrics, context_metric)
	}
	returnMetrics = append(returnMetrics, contextMetrics...)

	// Process the memory statistics
	memoryMetrics := make([]*Metric, 0, 10)
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.BlockSize",
		Value:     int64(xmlStats.Memory.Summary.BlockSize),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.ContextSize",
		Value:     int64(xmlStats.Memory.Summary.ContextSize),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.InUse",
		Value:     int64(xmlStats.Memory.Summary.InUse),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.Lost",
		Value:     int64(xmlStats.Memory.Summary.Lost),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.Malloced",
		Value:     int64(xmlStats.Memory.Summary.Malloced),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	memoryMetrics = append(memoryMetrics, &Metric{
		Name:      "memory.summary.TotalUse",
		Value:     int64(xmlStats.Memory.Summary.TotalUse),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	returnMetrics = append(returnMetrics, memoryMetrics...)

	// Process the server statistics
	serverMetrics := make([]*Metric, 0, 10)
	serverMetrics = append(serverMetrics, &Metric{
		Name:      "server.BootTime",
		Value:     int64(xmlStats.Server.BootTime.Unix()),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	serverMetrics = append(serverMetrics, &Metric{
		Name:      "server.ConfigTime",
		Value:     int64(xmlStats.Server.ConfigTime.Unix()),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	})
	returnMetrics = append(returnMetrics, serverMetrics...)

	// Process the socketmgr statistics
	socketMetrics := make([]*Metric, 0, 10)
	for _, socket := range xmlStats.Socketmgr.Sockets.Socket {
		socket_metric := &Metric{
			Name:      socket.ID,
			Value:     int64(socket.References),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{},
		}
		if socket.Name != nil {
			socket_metric.Name = *socket.Name
		}
		if socket.LocalAddress != nil {
			socket_metric.Tags = append(socket_metric.Tags, &MetricTag{"local-address", *socket.LocalAddress})
		}
		socket_metric.Tags = append(socket_metric.Tags, &MetricTag{"peer-address", socket.PeerAddress})
		socket_metric.Tags = append(socket_metric.Tags, &MetricTag{"type", socket.Type})
		socketMetrics = append(socketMetrics, socket_metric)
	}
	returnMetrics = append(returnMetrics, socketMetrics...)

	// Process the taskmgr statistics
	taskMetrics := make([]*Metric, 0, 10)
	for _, task := range xmlStats.Taskmgr.Tasks.Task {
		task_metric := &Metric{
			Name:      task.ID,
			Value:     int64(task.Events),
			Timestamp: xmlStats.Server.CurrentTime,
			Tags:      []*MetricTag{},
		}
		if task.Name != nil {
			task_metric.Tags = append(task_metric.Tags, &MetricTag{"name", *task.Name})
		}
		task_metric.Tags = append(task_metric.Tags, &MetricTag{"state", task.State})
		taskMetrics = append(taskMetrics, task_metric)
	}
	returnMetrics = append(returnMetrics, taskMetrics...)

	// Process the task thread model statistics
	thread_model_metric := &Metric{
		Name:      "taskmgr.thread-model",
		Value:     int64(xmlStats.Taskmgr.ThreadModel.DefaultQuantum),
		Timestamp: xmlStats.Server.CurrentTime,
		Tags:      []*MetricTag{},
	}
	thread_model_metric.Tags = append(thread_model_metric.Tags, &MetricTag{"type", xmlStats.Taskmgr.ThreadModel.Type})
	returnMetrics = append(returnMetrics, thread_model_metric)

	// Process the traffic statistics
	trafficMetrics := make([]*Metric, 0, 10)
	trafficMetrics = append(trafficMetrics, xmlStats.Traffic.Ipv4.toMetrics(xmlStats.Server.CurrentTime)...)
	trafficMetrics = append(trafficMetrics, xmlStats.Traffic.Ipv6.toMetrics(xmlStats.Server.CurrentTime)...)
	returnMetrics = append(returnMetrics, trafficMetrics...)

	// Process the view statistics
	viewMetrics := make([]*Metric, 0, 10)
	for _, view := range xmlStats.Views.View {
		// Zones
		view_tag := &MetricTag{"name", view.Name}
		cache_tag := &MetricTag{"cache", view.Cache.Name}
		for _, cache_counter := range view.Cache.Rrset {
			cache_metric := &Metric{
				Name:      cache_counter.Name,
				Value:     int64(cache_counter.Counter),
				Timestamp: xmlStats.Server.CurrentTime,
				Tags:      []*MetricTag{},
			}
			cache_metric.Tags = append(cache_metric.Tags, view_tag, cache_tag)
			viewMetrics = append(viewMetrics, cache_metric)
		}

		for _, view_counter := range view.Counters {
			view_counters := view_counter.toMetrics(xmlStats.Server.CurrentTime)
			for _, metric := range view_counters {
				view_counter_metric_tags := make([]*MetricTag, 0, len(metric.Tags)+1)
				view_counter_metric_tags = append(view_counter_metric_tags, view_tag)
				view_counter_metric_tags = append(view_counter_metric_tags, metric.Tags...)
				metric.Tags = view_counter_metric_tags
			}
			viewMetrics = append(viewMetrics, view_counters...)
		}

		for _, zone := range view.Zones.Zone {
			zone_tags := make([]*MetricTag, 0, 10)
			zone_tags = append(zone_tags, view_tag)
			zone_tags = append(zone_tags, &MetricTag{"name", zone.Name})
			zone_tags = append(zone_tags, &MetricTag{"rdataclass", zone.Rdataclass})
			zone_tags = append(zone_tags, &MetricTag{"type", zone.Type})
			for _, zone_counter := range zone.Counters {
				zone_counters := zone_counter.toMetrics(xmlStats.Server.CurrentTime)
				for _, metric := range zone_counters {
					zone_counter_metric_tags := make([]*MetricTag, 0, len(metric.Tags)+4)
					zone_counter_metric_tags = append(zone_counter_metric_tags, zone_tags...)
					zone_counter_metric_tags = append(zone_counter_metric_tags, metric.Tags...)
					metric.Tags = zone_counter_metric_tags
				}
				viewMetrics = append(viewMetrics, zone_counters...)
			}
		}
	}
	returnMetrics = append(returnMetrics, viewMetrics...)

	plugin.returnMetrics = returnMetrics

	return nil
}

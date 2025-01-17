package metric

import (
	"github.com/kubernetes-sigs/alibaba-cloud-csi-driver/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

var (
	fsClientPathPrefix = "/host/var/run/"
	fsClientTypeArray  = []string{"ossfs"}
	podInfo            = "pod_info"
	mountPointInfo     = "mount_point_info"
	counterTypeArray   = []string{"capacity_counter", "inodes_counter", "throughput_counter", "iops_counter", "latency_counter", "posix_counter", "oss_object_counter"}
	hotSpotArray       = []string{"hot_spot_read_file_top", "hot_spot_write_file_top", "hot_spot_head_file_top"}
)

var (
	usFsStatLabelNames = []string{"client_name", "backend_storage", "bucket_name", "namespace", "pod", "pv", "mount_point", "file_name"}
)

var (
	capacityBytesUsedCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "capacity_bytes_used_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	capacityBytesAvailableCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "capacity_bytes_available_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	capacityBytesTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "capacity_bytes_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	inodeBytesUsedCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "inode_bytes_used_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	inodeBytesAvailableCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "inode_bytes_available_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	inodeBytesTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "inode_bytes_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	readBytesTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "read_bytes_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	writeBytesTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "write_bytes_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	readCompletedTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "read_completed_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	writeCompletedTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "write_completed_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	readTimeMillisecondsTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "read_time_milliseconds_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	writeTimeMillisecondsTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "write_time_milliseconds_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixMkdirTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_mkdir_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixRmdirTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_rmdir_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixOpendirTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_opendir_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixReaddirTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_readdir_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixWriteTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_write_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixFlushTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_flush_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixFsyncTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_fsync_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixReleaseTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_release_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixReadTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_read_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixCreateTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_create_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixOpenTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_open_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixAccessTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_access_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixRenameTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_rename_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixChownTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_chown_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixChmodTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_chmod_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	posixTruncateTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "posix_truncate_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	ossPutObjectTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "oss_put_object_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	ossGetObjectTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "oss_get_object_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	ossHeadObjectTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "oss_head_object_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	ossDeleteObjectTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "oss_delete_object_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	ossPostObjectTotalCounterDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "oss_post_object_total_counter"),
		".",
		usFsStatLabelNames, nil,
	)
	hotSpotReadFileTopDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "hot_spot_read_file_top"),
		".",
		usFsStatLabelNames, nil,
	)
	hotSpotWriteFileTopDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "hot_spot_write_file_top"),
		".",
		usFsStatLabelNames, nil,
	)
	hotSpotHeadFileTopDesc = prometheus.NewDesc(
		prometheus.BuildFQName(nodeNamespace, volumeSubSystem, "hot_spot_head_file_top"),
		".",
		usFsStatLabelNames, nil,
	)
)

type fuseInfo struct {
	ClientName     string
	BackendStorage string
	BucketName     string
	Namespace      string
	PodName        string
	PodUID         string
	PvName         string
	MountPoint     string
}

type usFsStatCollector struct {
	hotSpotReadFileTop  *prometheus.Desc
	hotSpotWriteFileTop *prometheus.Desc
	hotSpotHeadFileTop  *prometheus.Desc
	descs               []typedFactorDesc
}

func init() {
	registerCollector("fuse_oss_stat", NewFuseOssStatCollector)
}

// NewUsFsStatCollector returns a new Collector exposing user space fs stats.
func NewFuseOssStatCollector() (Collector, error) {
	return &usFsStatCollector{
		hotSpotReadFileTop:  hotSpotReadFileTopDesc,
		hotSpotWriteFileTop: hotSpotWriteFileTopDesc,
		hotSpotHeadFileTop:  hotSpotHeadFileTopDesc,
		descs: []typedFactorDesc{
			//0-2
			{desc: capacityBytesUsedCounterDesc, valueType: prometheus.CounterValue},
			{desc: capacityBytesAvailableCounterDesc, valueType: prometheus.CounterValue},
			{desc: capacityBytesTotalCounterDesc, valueType: prometheus.CounterValue},
			//3-5
			{desc: inodeBytesUsedCounterDesc, valueType: prometheus.CounterValue},
			{desc: inodeBytesAvailableCounterDesc, valueType: prometheus.CounterValue},
			{desc: inodeBytesTotalCounterDesc, valueType: prometheus.CounterValue},
			//6-7
			{desc: readBytesTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: writeBytesTotalCounterDesc, valueType: prometheus.CounterValue},
			//8-9
			{desc: readCompletedTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: writeCompletedTotalCounterDesc, valueType: prometheus.CounterValue},
			//10-11
			{desc: readTimeMillisecondsTotalCounterDesc, valueType: prometheus.CounterValue, factor: .001},
			{desc: writeTimeMillisecondsTotalCounterDesc, valueType: prometheus.CounterValue, factor: .001},
			//12-27
			{desc: posixMkdirTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixRmdirTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixOpendirTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixReaddirTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixWriteTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixFlushTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixFsyncTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixReleaseTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixReadTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixCreateTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixOpenTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixAccessTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixRenameTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixChownTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixChmodTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: posixTruncateTotalCounterDesc, valueType: prometheus.CounterValue},
			//28-32
			{desc: ossPutObjectTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: ossGetObjectTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: ossHeadObjectTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: ossDeleteObjectTotalCounterDesc, valueType: prometheus.CounterValue},
			{desc: ossPostObjectTotalCounterDesc, valueType: prometheus.CounterValue},
		},
	}, nil
}

func getPodUID(fsClientPathPrefix string, fsClientType string) ([]string, error) {
	fsClientPath := fsClientPathPrefix + fsClientType
	if !utils.IsFileExisting(fsClientPath) {
		_ = utils.MkdirAll(fsClientPath, os.FileMode(0755))
	}
	return listDirectory(fsClientPathPrefix + fsClientType)
}

func setCounterStat(start int, end int, stat *[35]string, metricsArray []string) {
	if len(metricsArray) == 0 {
		return
	}
	if len(metricsArray) < end-start+1 {
		return
	}
	for i := 0; i < end-start+1; i++ {
		(*stat)[i+start] = metricsArray[i]
	}
}

func (p *usFsStatCollector) Update(ch chan<- prometheus.Metric) error {
	initFsClientFlag := false
	var stat = [35]string{}
	fsClientInfo := new(fuseInfo)
	// foreach fuse client type
	for _, fsClientType := range fsClientTypeArray {
		// get pod uid
		podUIDArray, err := getPodUID(fsClientPathPrefix, fsClientType)
		if err != nil {
			continue
		}
		//foreach pod uid
		for _, podUID := range podUIDArray {
			//get pod info
			podInfoArray, err := readFirstLines(fsClientPathPrefix + fsClientType + "/" + podUID + "/" + podInfo)
			if err != nil {
				continue
			}
			if len(podInfoArray) < 4 {
				continue
			}
			fsClientInfo.Namespace = podInfoArray[0]
			fsClientInfo.PodName = podInfoArray[1]
			fsClientInfo.PodUID = podInfoArray[2]
			// list volume from pod
			volumeArray, err := listDirectory(fsClientPathPrefix + fsClientType + "/" + podUID + "/")
			if err != nil {
				continue
			}
			// foreach volume
			for _, volume := range volumeArray {
				mountPointInfoArray, err := readFirstLines(fsClientPathPrefix + fsClientType + "/" + podUID + "/" + volume + "/" + mountPointInfo)
				if err != nil {
					continue
				}
				if len(mountPointInfoArray) < 5 {
					continue
				}
				fsClientInfo.ClientName = mountPointInfoArray[0]
				fsClientInfo.BackendStorage = mountPointInfoArray[1]
				fsClientInfo.BucketName = mountPointInfoArray[2]
				fsClientInfo.PvName = mountPointInfoArray[3]
				fsClientInfo.MountPoint = mountPointInfoArray[4]
				initFsClientFlag = true
				// foreach counter metrics
				for _, counterType := range counterTypeArray {
					metricsArray, err := readFirstLines(fsClientPathPrefix + fsClientType + "/" + podUID + "/" + volume + "/" + counterType)
					if err != nil {
						continue
					}
					switch counterType {
					case "capacity_counter":
						setCounterStat(0, 2, &stat, metricsArray)
					case "inodes_counter":
						setCounterStat(3, 5, &stat, metricsArray)
					case "throughput_counter":
						setCounterStat(6, 7, &stat, metricsArray)
					case "iops_counter":
						setCounterStat(8, 9, &stat, metricsArray)
					case "latency_counter":
						setCounterStat(10, 11, &stat, metricsArray)
					case "posix_counter":
						setCounterStat(12, 27, &stat, metricsArray)
					case "oss_object_counter":
						setCounterStat(28, 32, &stat, metricsArray)
					default:
						log.Errorf("Unknow counterType:%s", counterType)
					}
				}
				if initFsClientFlag {
					p.setCounterMetrics(fsClientInfo, stat, ch)
				}
				for _, hotSpotType := range hotSpotArray {
					metricsArray, err := readFirstLines(fsClientPathPrefix + fsClientType + "/" + podUID + "/" + volume + "/" + hotSpotType)
					if err != nil {
						continue
					}
					for _, metricsValue := range metricsArray {
						start := strings.LastIndex(metricsValue, ":")
						if start == -1 {
							continue
						}
						fileName := metricsValue[0:start]
						value := metricsValue[start+1:]
						valueFloat64, err := strconv.ParseFloat(value, 64)
						if err != nil {
							continue
						}
						switch hotSpotType {
						case "hot_spot_read_file_top":
							ch <- prometheus.MustNewConstMetric(p.hotSpotReadFileTop, prometheus.GaugeValue, valueFloat64, fsClientInfo.ClientName, fsClientInfo.BackendStorage, fsClientInfo.BucketName, fsClientInfo.Namespace, fsClientInfo.PodName, fsClientInfo.PvName, fsClientInfo.MountPoint, fileName)
						case "hot_spot_write_file_top":
							ch <- prometheus.MustNewConstMetric(p.hotSpotWriteFileTop, prometheus.GaugeValue, valueFloat64, fsClientInfo.ClientName, fsClientInfo.BackendStorage, fsClientInfo.BucketName, fsClientInfo.Namespace, fsClientInfo.PodName, fsClientInfo.PvName, fsClientInfo.MountPoint, fileName)
						case "hot_spot_head_file_top":
							ch <- prometheus.MustNewConstMetric(p.hotSpotHeadFileTop, prometheus.GaugeValue, valueFloat64, fsClientInfo.ClientName, fsClientInfo.BackendStorage, fsClientInfo.BucketName, fsClientInfo.Namespace, fsClientInfo.PodName, fsClientInfo.PvName, fsClientInfo.MountPoint, fileName)

						default:
							log.Errorf("Unknow hotSpotType:%s", hotSpotType)
						}
					}
				}
			}
		}
	}
	return nil
}

func (p *usFsStatCollector) setCounterMetrics(fsClientInfo *fuseInfo, stats [35]string, ch chan<- prometheus.Metric) {
	for i, value := range stats {
		if i >= len(p.descs) {
			return
		}
		if len(strings.TrimSpace(value)) == 0 {
			continue
		}
		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Errorf("Convert value %s to float64 is failed, err:%s, stat:%+v", value, err, stats)
			continue
		}

		ch <- p.descs[i].mustNewConstMetric(valueFloat64, fsClientInfo.ClientName, fsClientInfo.BackendStorage, fsClientInfo.BucketName, fsClientInfo.Namespace, fsClientInfo.PodName, fsClientInfo.PvName, fsClientInfo.MountPoint, "")
	}
}

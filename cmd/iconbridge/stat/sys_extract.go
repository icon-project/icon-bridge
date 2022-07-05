package stat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type SystemMetrics struct {
	LoadAverage *LoadAverage `json:"LoadAverage"`
	MemoryUsage *MemoryUsage `json:"MemoryUsage"`
	DiskUsage   *DiskUsage   `json:"DiskUsage"`
}

type LoadAverage struct {
	LoadAvg1       float64 `json:"LoadAvg1"`
	LoadAvg5       float64 `json:"LoadAvg5"`
	LoadAvg15      float64 `json:"LoadAvg15"`
	RunningThreads float64 `json:"RunningThreads"`
	TotalThreads   float64 `json:"TotalThreads"`
}

type MemoryUsage struct {
	MemFree      float64 `json:"MemFree"`
	MemTotal     float64 `json:"MemTotal"`
	MemAvailable float64 `json:"MemAvailable"`
	UsedPercent  float64 `json:"UsedPercent"`
	SwapFree     float64 `json:"SwapFree"`
	SwapTotal    float64 `json:"SwapTotal"`
	SwapCached   float64 `json:"SwapCached"`
}

type DiskUsage struct {
	DiskUsagePerMount []*DiskUsagePerMount `json:"DiskUsagePerMount"`
}

type DiskUsagePerMount struct {
	FileSystem  string  `json:"FileSystem"`
	MountedOn   string  `json:"MountedOn"`
	Used        string  `json:"Used"`
	Available   string  `json:"Available"`
	UsedPercent float64 `json:"UsedPercent"`
}

func getSystemMetrics() (sysMetrics *SystemMetrics, err error) {
	sysMetrics = &SystemMetrics{
		LoadAverage: &LoadAverage{},
		MemoryUsage: &MemoryUsage{},
		DiskUsage:   &DiskUsage{},
	}
	if sysMetrics.LoadAverage, err = getLoadAverage(); err != nil {
		return
	}
	if sysMetrics.MemoryUsage, err = getMemoryUsage(); err != nil {
		return
	}
	if sysMetrics.DiskUsage, err = getDiskUsage(); err != nil {
		return
	}
	return
}

func getLoadAverage() (lavg *LoadAverage, err error) {
	const NUM_FIELDS = 5
	var f *os.File
	if f, err = os.Open("/proc/loadavg"); err != nil {
		err = errors.Wrap(err, "getLoadAverageFunc; Open; Err: ")
		return
	}
	defer f.Close()

	lavg = &LoadAverage{}
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		err = errors.New("getLoadAverageFunc; failed to scan /proc/loadavg")
		return
	}
	items := strings.Fields((scanner.Text()))
	if len(items) != NUM_FIELDS {
		err = errors.New("getLoadAverageFunc; Number of fields different than expected")
		return
	}
	lavg.LoadAvg1, _ = strconv.ParseFloat(items[0], 64)
	lavg.LoadAvg5, _ = strconv.ParseFloat(items[1], 64)
	lavg.LoadAvg15, _ = strconv.ParseFloat(items[2], 64)
	if threadItem := strings.Split(items[3], "/"); len(threadItem) == 2 {
		lavg.RunningThreads, _ = strconv.ParseFloat(threadItem[0], 64)
		lavg.TotalThreads, _ = strconv.ParseFloat(threadItem[1], 64)
	}
	return
}

func getMemoryUsage() (memUsage *MemoryUsage, err error) {
	var f *os.File
	if f, err = os.Open("/proc/meminfo"); err != nil {
		err = errors.Wrap(err, "getMemoryUsageFunc; Open; Err: ")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	memUsage = &MemoryUsage{}
	refVal := reflect.ValueOf(memUsage).Elem()
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		fld := line[:i]
		sf := refVal.FieldByName(fld)
		if sf.IsValid() {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				sf.SetFloat(v / (1024 * 1024))
			}
		}
	}
	if memUsage != nil {
		memUsage.UsedPercent = 100 - ((100 * memUsage.MemAvailable) / memUsage.MemTotal)
	}
	return
}

func getDiskUsage() (diskUsage *DiskUsage, err error) {
	const NUM_FIELDS = 6
	cmd := exec.Command("bash", "-c", `df -h`)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	diskUsage = &DiskUsage{DiskUsagePerMount: []*DiskUsagePerMount{}}
	for _, line := range strings.Split(out.String(), "\n")[1:] {
		flds := strings.Fields(line)
		if len(flds) != NUM_FIELDS {
			continue
		}
		perc, _ := strconv.ParseFloat(strings.TrimSuffix(flds[4], "%"), 64)
		diskUsage.DiskUsagePerMount = append(diskUsage.DiskUsagePerMount,
			&DiskUsagePerMount{FileSystem: flds[0], Used: flds[2], Available: flds[3], UsedPercent: perc, MountedOn: flds[5]})
	}
	return
}

func (l *LoadAverage) filter(c *Trigger, verbose bool) (resMap map[string]interface{}, err error) {
	key := string(LOADAVERAGE)
	resMap = map[string]interface{}{}
	resBytes, err := json.Marshal(l)
	if err != nil {
		err = errors.Wrap(err, "LoadAverage.filterFunc; JSON Marshal; Err: ")
		return
	}
	json.Unmarshal(resBytes, &resMap)
	if v, ok := resMap[c.Field]; ok && v != nil {
		vfloat, ok := v.(float64)
		if !ok {
			err = errors.New("LoadAverage.filterFunc; JSON Marsha; Err: Value of field " + c.Field + " is not float64 ")
			return
		}
		if compare(vfloat, c.Sign, c.Value) {
			if !verbose {
				n := map[string]interface{}{c.Field: vfloat}
				return map[string]interface{}{key: n}, nil
			}
			return map[string]interface{}{key: resMap}, nil
		}
	}
	return nil, nil
}

func (m *MemoryUsage) filter(c *Trigger, verbose bool) (resMap map[string]interface{}, err error) {
	key := string(MEMORYUSAGE)
	resMap = map[string]interface{}{}
	resBytes, err := json.Marshal(m)
	if err != nil {
		err = errors.Wrap(err, "MemoryUsage.filterFunc; JSON Marshal; Err: ")
		return
	}
	json.Unmarshal(resBytes, &resMap)
	if v, ok := resMap[c.Field]; ok && v != nil {
		vfloat, ok := v.(float64)
		if !ok {
			err = errors.New("MemoryUsage.filterFunc; JSON Marsha; Err: Value of field " + c.Field + " is not float64 ")
			return
		}
		if compare(vfloat, c.Sign, c.Value) {
			if !verbose {
				n := map[string]interface{}{c.Field: vfloat}
				return map[string]interface{}{key: n}, nil
			}
			return map[string]interface{}{key: resMap}, nil
		}
	}
	return nil, nil
}

func (d *DiskUsage) filter(c *Trigger, verbose bool) (res map[string]interface{}, err error) {
	NECESSARY_FIELD := "MountedOn"
	resMapArr := []map[string]interface{}{}
	key := string(DISKUSAGE)
	for _, dm := range d.DiskUsagePerMount {
		if dm == nil {
			continue
		}
		resMap := map[string]interface{}{}
		resBytes := []byte{}
		resBytes, err = json.Marshal(dm)
		if err != nil {
			err = errors.Wrap(err, "DiskUsage.filterFunc; JSON Marshal; Err: ")
			return
		}
		json.Unmarshal(resBytes, &resMap)
		if v, ok := resMap[c.Field]; ok && v != nil {
			vfloat, ok := v.(float64)
			if !ok {
				err = errors.New("DiskUsage.filterFunc; JSON Marsha; Err: Value of field " + c.Field + " is not float64 ")
				return
			}
			if compare(vfloat, c.Sign, c.Value) {
				if !verbose {
					n := map[string]interface{}{c.Field: vfloat}
					if nf, nfok := resMap[NECESSARY_FIELD]; nfok && nf != nil {
						n[NECESSARY_FIELD] = nf
					}
					resMapArr = append(resMapArr, n)
				} else {
					resMapArr = append(resMapArr, resMap)
				}
			}
		}
	}
	if len(resMapArr) > 0 {
		return map[string]interface{}{key: resMapArr}, nil
	}
	return nil, nil
}

func compare(left float64, s string, right float64) bool {
	sign := strings.TrimSpace(s)
	if sign == "==" {
		return left == right
	} else if s == ">" {
		return left > right
	} else if s == "<" {
		return left < right
	} else if s == ">=" {
		return left >= right
	} else if s == "<=" {
		return left <= right
	} else if s == "!=" {
		return left != right
	}
	return false
}

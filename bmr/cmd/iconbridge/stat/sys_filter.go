package stat

import "errors"

func getFilteredMetrics(criterion []*Trigger, verbose bool) (mets map[string]interface{}, err error) {
	if len(criterion) == 0 {
		return nil, nil
	}
	errorMessage := ""
	mets = map[string]interface{}{}
	for _, c := range criterion {
		if c.Measurement == LOADAVERAGE {
			if l, err := getLoadAverage(); l != nil && err == nil {
				if r, err := l.filter(c, verbose); err == nil && r != nil {
					for k, v := range r {
						mets[k] = v
					}
				} else if err != nil {
					errorMessage += "getLoadAverageFunc; filterFunc; Err: " + err.Error() + "\n"
				}
			} else if err != nil {
				errorMessage += "getLoadAverageFunc; Err: " + err.Error() + "\n"
			}
		} else if c.Measurement == MEMORYUSAGE {
			if m, err := getMemoryUsage(); m != nil && err == nil {
				if r, err := m.filter(c, verbose); err == nil && r != nil {
					for k, v := range r {
						mets[k] = v
					}
				} else if err != nil {
					errorMessage += "getMemoryUsage; filterFunc; Err: " + err.Error() + "\n"
				}
			} else if err != nil {
				errorMessage += "getMemoryUsage; Err: " + err.Error() + "\n"
			}
		} else if c.Measurement == DISKUSAGE {
			if d, err := getDiskUsage(); d != nil && err == nil {
				if r, err := d.filter(c, verbose); err == nil && r != nil {
					for k, v := range r {
						mets[k] = v
					}
				} else if err != nil {
					errorMessage += "getDiskUsage; filterFunc; Err: " + err.Error() + "\n"
				}
			} else if err != nil {
				errorMessage += "getDiskUsage; Err: " + err.Error() + "\n"
			}
		}
	}
	if len(errorMessage) > 0 {
		err = errors.New(errorMessage)
	}
	return
}

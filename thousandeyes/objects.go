package thousandeyes

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

const (
	apiURLAlerts          = "https://api.thousandeyes.com/v6/alerts?format=json"
	apiURLTests           = "https://api.thousandeyes.com/v6/tests.json"
	apiURLTestBGB         = "https://api.thousandeyes.com/v6/net/bgp-metrics/%d.json"
	apiURLTestHTTP        = "https://api.thousandeyes.com/v6/web/http-server/%d.json"
	apiURLTestHTTPMetrics = "https://api.thousandeyes.com/v6/net/metrics/%d.json"
	// Added
	apiURLNetworkPathVisualizationResults = "https://api.thousandeyes.com/v6/net/path-vis/%d.json"
)

// ThousandeyesRequest the request struct
type Request struct {
	URL            string
	ResponseCode   int
	ResponseObject interface{}
	Error          error
}

// ThousandAlerts describes the JSON returned by a request active alerts to ThousandEyes
type ThousandAlerts struct {
	From  string `json:"from"`
	Alert []struct {
		Active    int        `json:"active"`
		AlertID   int        `json:"alertId"`
		DateEnd   string     `json:"dateEnd,omitempty"`
		DateStart string     `json:"dateStart"`
		Monitors  []struct { //array of monitors where the alert has at some point been active since the point that the alert was triggered. Only shown on BGP alerts.
			Active         int    `json:"active"`
			MetricsAtStart string `json:"metricsAtStart"`
			MetricsAtEnd   string `json:"metricsAtEnd"`
			MonitorID      int    `json:"monitorId"`
			MonitorName    string `json:"monitorName"`
			PrefixID       int    `json:"prefixId"`
			Prefix         string `json:"prefix"`
			DateStart      string `json:"dateStart"`
			DateEnd        string `json:"dateEnd"`
			Permalink      string `json:"permalink"`
			Network        string `json:"network"`
		} `json:"monitors,omitempty"`
		Permalink      string `json:"permalink"`
		RuleExpression string `json:"ruleExpression"`
		RuleID         int    `json:"ruleId"`
		RuleName       string `json:"ruleName"`
		TestID         int    `json:"testId"`
		TestName       string `json:"testName"`
		ViolationCount int    `json:"violationCount"`
		Type           string `json:"type"`
		APILinks       []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"apiLinks,omitempty"`
		Agents []struct { //array of monitors where the alert has at some point been active since the point that the alert was triggered. Not shown on BGP alerts.
			Active         int    `json:"active"`
			MetricsAtStart string `json:"metricsAtStart"`
			MetricsAtEnd   string `json:"metricsAtEnd"`
			AgentID        int    `json:"agentId"`
			AgentName      string `json:"agentName"`
			DateStart      string `json:"dateStart"`
			DateEnd        string `json:"dateEnd"`
			Permalink      string `json:"permalink"`
		} `json:"agents,omitempty"`
	} `json:"alert"`
	Pages struct {
		Current int `json:"current"`
	} `json:"pages"`
}

//ThousandTests describes needed Fields from the JSON returned by a request  to ThousandEyes
type ThousandTests struct {
	Tests []ThousandTest `json:"test"`
}

//ThousandTest in detail
type ThousandTest struct {
	TestID   int    `json:"testId"`
	TestName string `json:"testName"`
	Type     string `json:"type"`
	Prefix   string `json:"prefix"`
	Interval int    `json:"interval"`
	URL      string `json:"url"`
}

//https://api.thousandeyes.com/v6/net/bgp-metrics/557962.json

// BGPTestResults BGP Test details
type BGPTestResults struct {
	Net struct {
		Test       ThousandTest `json:"test"`
		BgpMetrics []struct {
			CountryID    string  `json:"countryId"`
			Prefix       string  `json:"prefix"`
			MonitorName  string  `json:"monitorName"`
			Reachability float32 `json:"reachability"`
			Updates      float32 `json:"updates"`
			PathChanges  float32 `json:"pathChanges"`
		} `json:"bgpMetrics"`
	} `json:"net"`
}

// https://api.thousandeyes.com/v6/net/metrics/612434.json

// HTTPTestMetricResults HTTP Test details on network metrics
type HTTPTestMetricResults struct {
	Net struct {
		Test        ThousandTest `json:"test"`
		HTTPMetrics []struct {
			AvgLatency float32 `json:"avgLatency"`
			Loss       float32 `json:"loss"`
			MaxLatency float32 `json:"maxLatency"`
			Jitter     float32 `json:"jitter"`
			MinLatency float32 `json:"minLatency"`
			ServerIP   string  `json:"serverIp"`
			AgentName  string  `json:"agentName"`
			CountryID  string  `json:"countryId"`
		} `json:"metrics"`
	} `json:"net"`
}

// Test Type agent-to-agent
// Network - Path Visualization
// https://developer.thousandeyes.com/v6/test_data/#/path-vis
type NetworkPathVisualizationResults struct {
	Net struct {
		Test    ThousandTest `json:"test"`
		PathVis []struct {
			AgentName string `json:"agentName"`
			AgentID   int    `json:"agentId"`
			CountryID string `json:"countryId"`
			ServerIp  string `json:"serverIp"`
			SourceIp  string `json:"sourceIp"`
			EndPoints []struct {
				NumberOfHops int    `json:"numberOfHops"`
				IpAddress    string `json:"ipAddress"`
				ResponseTime int    `json:"responseTime"`
				PathId       string `json:"pathId"`
			} `json:"endPoints"`
		} `json:"pathVis"`
	} `json:"net"`
}

// HTTPTestWebServerResults HTTP Test details on Server Response
type HTTPTestWebServerResults struct {
	Web struct {
		Test       ThousandTest `json:"test"`
		HTTPServer []struct {
			ConnectTime  int    `json:"connectTime"`
			DNSTime      int    `json:"dnsTime"`
			ErrorType    string `json:"errorType"`
			NumRedirects int    `json:"numRedirects"`
			ReceiveTime  int    `json:"receiveTime"`
			ResponseCode int    `json:"responseCode"`
			ResponseTime int    `json:"responseTime"`
			TotalTime    int    `json:"totalTime"`
			WaitTime     int    `json:"waitTime"`
			WireSize     int    `json:"wireSize"`
			AgentName    string `json:"agentName"`
			CountryID    string `json:"countryId"`
			Date         string `json:"date"`
			AgentID      int    `json:"agentId"`
			RoundID      int    `json:"roundId"`
		} `json:"httpServer"`
	} `json:"web"`
}

func thousandEyesDateTime() string {
	// Go back a bit to have some alerts to parse
	t := time.Now().UTC().Add(-*RetrospectionPeriod)
	// 2006-01-02T15:04:05 is a magic date to format dates using example based layouts
	f := t.Format("2006-01-02T15:04:05")
	return string(f)
}

func (t *Collector) GetAlerts() (ThousandAlerts, bool, bool) {

	r := Request{
		URL:            apiURLAlerts,
		ResponseObject: new(ThousandAlerts),
	}

	bHitAPILimit, bError := CallSingle(t.Token, t.User, t.IsBasicAuth, &r)

	return *r.ResponseObject.(*ThousandAlerts), bHitAPILimit, bError
}

func (t *Collector) GetTests() (bgpMs []BGPTestResults,
	httpMs []HTTPTestMetricResults,
	httpWs []HTTPTestWebServerResults,
	netPathVis []NetworkPathVisualizationResults,
	bHitAPILimit,
	bError bool) {

	rTests := Request{
		URL:            apiURLTests,
		ResponseObject: new(ThousandTests),
	}
	bHitAPILimit, bError = CallSingle(t.Token, t.User, t.IsBasicAuth, &rTests)
	if rTests.Error != nil {
		return bgpMs, httpMs, httpWs, netPathVis, bHitAPILimit, bError
	}

	te := rTests.ResponseObject.(*ThousandTests)

	var testRequests []Request

	log.Println(fmt.Sprintf("INFO: ThousandEyes Test Count: %d", len(te.Tests)))

	for i := range te.Tests {
		log.Println(fmt.Sprintf("INFO: Test list: %s %d", te.Tests[i].Type, te.Tests[i].TestID))
		switch te.Tests[i].Type {
		case "http-server":

			if t.IsCollectHttp {
				testRequests = append(testRequests, Request{
					URL:            fmt.Sprintf(apiURLTestHTTP, te.Tests[i].TestID),
					ResponseObject: new(HTTPTestWebServerResults),
				})
			}
			if t.IsCollectHttpMetrics {
				testRequests = append(testRequests, Request{
					URL:            fmt.Sprintf(apiURLTestHTTPMetrics, te.Tests[i].TestID),
					ResponseObject: new(HTTPTestMetricResults),
				})
			}

		case "bgp":

			if t.IsCollectBgp {
				testRequests = append(testRequests, Request{
					URL:            fmt.Sprintf(apiURLTestBGB, te.Tests[i].TestID),
					ResponseObject: new(BGPTestResults),
				})
			}

		case "agent-to-agent":
			if t.IsCollectNetPathViz {
				testRequests = append(testRequests, Request{
					URL:            fmt.Sprintf(apiURLNetworkPathVisualizationResults, te.Tests[i].TestID),
					ResponseObject: new(NetworkPathVisualizationResults),
				})
				//log.Println(fmt.Sprintf("INFO: testResults [%s]",testRequests ))
			}

		case "agent-to-server":
			break

		case "web-transactions":
			break

		case "dns-server":
			break

		case "dns-trace":
			break

		case "page-load":
			break

		case "voice":
			break

		default:
			log.Println(fmt.Sprintf("WARN: Not a handled test type: %s. Bug. Fix Code.", te.Tests[i].Type))
		}
	}

	//CallSequence(t.token, testRequests)
	bHitAPILimit, bError = CallParallel(t.Token, t.User, t.IsBasicAuth, testRequests)

	for c, o := range testRequests {

		//switch v:=o.ResponseObject.(type) {
		//v := reflect.TypeOf(o.ResponseObject)
		switch o.ResponseObject.(type) {
		case *BGPTestResults:
			bgpMs = append(bgpMs, *testRequests[c].ResponseObject.(*BGPTestResults))
		case *HTTPTestMetricResults:
			httpMs = append(httpMs, *testRequests[c].ResponseObject.(*HTTPTestMetricResults))
		case *HTTPTestWebServerResults:
			httpWs = append(httpWs, *testRequests[c].ResponseObject.(*HTTPTestWebServerResults))
		case *NetworkPathVisualizationResults:
			netPathVis = append(netPathVis, *testRequests[c].ResponseObject.(*NetworkPathVisualizationResults))
		default:
			log.Println(fmt.Sprintf("ERROR: Not a handled test type %s (%d of %d). Bug. Fix Code.", reflect.TypeOf(o.ResponseObject), c, len(testRequests)))
		}
	}

	return bgpMs, httpMs, httpWs, netPathVis, bHitAPILimit, bError
}

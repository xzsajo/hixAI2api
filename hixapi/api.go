package hixapi

import (
	"encoding/json"
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"github.com/gin-gonic/gin"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"
)

const (
	baseURL            = "https://hix.ai"
	chatEndpoint       = baseURL + "/api/hix/chat"
	createChatEndpoint = baseURL + "/api/trpc/hixChat.createChat?batch=1"
	subUsageEndpoint   = baseURL + "/api/trpc/subUsage.getSubUsage?batch=1"
	deleteEndpoint     = baseURL + "/api/project/delete?project_id=%s"
	uploadEndpoint     = baseURL + "/api/get_upload_personal_image_url"
	chatType           = "COPILOT_MOA_CHAT"
	imageType          = "COPILOT_MOA_IMAGE"
)

type CreateChatResponse struct {
	Result struct {
		Data struct {
			JSON struct {
				ID string `json:"id"`
			} `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

func MakeCreateChatRequest(client cycletls.CycleTLS, cookie string, modelId int) (string, error) {
	createChatBody := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"title": "Untitled",
				"botId": modelId,
			},
		},
	}
	bytes, err := json.Marshal(createChatBody)
	if err != nil {
		return "", err
	}
	accept := "application/json"

	response, err := client.Do(fmt.Sprintf(createChatEndpoint), cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Method:  "POST",
		Body:    string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "POST")

	if err != nil {
		return "", err
	}

	var responses []CreateChatResponse
	err = json.Unmarshal([]byte(response.Body), &responses)
	if err != nil {
		return "", err
	}

	// 检查数组是否非空并提取ID
	if len(responses) > 0 {
		id := responses[0].Result.Data.JSON.ID
		return id, nil
	} else {
		return "", fmt.Errorf("MakeCreateChatRequest err")
	}
}

type SubUsageResponse struct {
	Result struct {
		Data struct {
			JSON struct {
				UsageList []struct {
					ID             int    `json:"id"`
					SubscriptionID *int   `json:"subscription_id"`
					TotalCount     int    `json:"total_count"`
					UseCount       int    `json:"use_count"`
					Status         string `json:"status"`
					UsageType      string `json:"usage_type"`
					AppName        string `json:"app_name"`
					DateStart      string `json:"date_start"`
					DateEnd        string `json:"date_end"`
					PriceID        *int   `json:"price_id"`
				} `json:"usageList"`
				IsActiveSub bool `json:"isActiveSub"`
			} `json:"json"`
			Meta struct {
				Values struct {
					DateStart []string `json:"usageList.0.date_start"`
					DateEnd   []string `json:"usageList.0.date_end"`
				} `json:"values"`
			} `json:"meta"`
		} `json:"data"`
	} `json:"result"`
}

func MakeSubUsageRequest(client cycletls.CycleTLS, cookie string) (int, error) {
	subUsageReqParam := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"appName": "HIXChat",
			},
		},
		"1": map[string]interface{}{
			"json": map[string]interface{}{
				"appName": "HIXChat",
			},
		},
	}
	bytes, err := json.Marshal(subUsageReqParam)
	if err != nil {
		return 0, err
	}
	accept := "application/json"

	response, err := client.Do(fmt.Sprintf(subUsageEndpoint+"&input=%s", string(bytes)), cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Method:  "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "GET")

	if err != nil {
		return 0, err
	}

	var responses []SubUsageResponse
	err = json.Unmarshal([]byte(response.Body), &responses)
	if err != nil {
		return 0, err
	}

	// 检查数组是否非空并提取ID
	if len(responses) > 0 {
		if len(responses[0].Result.Data.JSON.UsageList) > 0 {
			totalCount := responses[0].Result.Data.JSON.UsageList[0].TotalCount
			useCount := responses[0].Result.Data.JSON.UsageList[0].UseCount
			return totalCount - useCount, nil
		} else {
			return 0, fmt.Errorf("MakeSubUsageRequest err")

		}

	} else {
		return 0, fmt.Errorf("MakeSubUsageRequest err")
	}
}

// makeRequest 发送HTTP请求
func makeChatRequest(client cycletls.CycleTLS, jsonData []byte, cookie string, isStream bool) (cycletls.Response, error) {
	accept := "application/json"
	if isStream {
		accept = "text/event-stream"
	}

	return client.Do(chatEndpoint, cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Body:    string(jsonData),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "POST")
}

func MakeStreamChatRequest(c *gin.Context, client cycletls.CycleTLS, jsonData []byte, cookie string) (<-chan cycletls.SSEResponse, error) {

	options := cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Body:    string(jsonData),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "text/event-stream",
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       "user_group=146; first-visit-url=https%3A%2F%2Fhix.ai%2Fzh; device-id=3ab3c15d54dda3a82aef9461f306e6fb; _gcl_au=1.1.1766983635.1741332968; _fbp=fb.1.1741332968166.96038000529332917; _tt_enable_cookie=1; _ttp=01JNQRGWQPYNM40F73N6A7DNWX_.tt.1; _ga=GA1.1.83506002.1741332969; _clck=1sjmshv%7C2%7Cfu0%7C0%7C1892; FPID=FPID2.2.cveO53H73pFObELS1NXF0joKa3smasV1Fd7c6cbjXEs%3D.1741332969; FPLC=l83hJaBgVHVKJy2BIZZgfnjYg29Gcav793BsFjFIFboOeJBP5RH6oJt8l4IZcdGDMTTz%2Bn5G34rmTKPUY1NZQ%2FzMIUyFQV9P%2F1Y0bvf7l5lOeBLjodjZ9JgcCB5lVw%3D%3D; g_state={\"i_l\":0}; __gsas=ID=0cca1082a003af27:T=1741333199:RT=1741333199:S=ALNI_MaEH-B7EYhLTgL2xoSVgNf0IhQJsw; __Secure-next-auth.callback-url=https%3A%2F%2Fhix.ai; __Host-next-auth.csrf-token=e4bbb2c59d9e74642b5d7926d5c479b82bb3d59319d93a703456e4e671769c5e%7C436ace319f2234970a324ea6255857f09db7a4ddacbf7891125959289e759c4b; cf_clearance=2Z1OzrAJBYIeCRKFdZlJHWmqumQe73mTFibHm8hbs6w-1741338686-1.2.1.1-cL2qKRVDMWEVsnj07nHh2lQM9FezrqGgDWnp5N.INALKEgFlCnCxvzrh0ZzUFM7_iUi.1cCI.xeEv35zFz8joF6D4Di7RusnEeYYxohIwzbt7aUi10I0n4oxicsDWM0T9aVLiQeDGe0Oo.iIVANLdQ.KZSstJcg6xEu4XDYvr1swZ7ui5xxZP_xg.INEbpjgTo1hdyFeRtbB0yqPwH_OfD23NnSVNKsIdFS14sdm5VgKlR0y0z68HfAKLeNMfmL17XJ_AVhongj2aBQwGhG4Yu3lsXOTSbA0wnyoSpZFzL1ksfjAJfiYtpnBnSgMGxxTywtcOSocTP9d3zj1pmDPTMQ9RoIuyFUafK4TFTaJRw0mmKifvbNkrs47g.9Rg5pbHm7PVA3AIsZ6sYxAMwJQ13JRRSjxBd1HFHOMEfnBezs; user-info=%7B%22id%22%3A%22cm7ygtgt9005rczp0k9klmjo2%22%2C%22name%22%3A%22Dean%20Xv%22%2C%22email%22%3A%22alanxv1024%40gmail.com%22%2C%22firstName%22%3A%22Dean%22%2C%22lastName%22%3A%22Xv%22%2C%22image%22%3A%22https%3A%2F%2Flh3.googleusercontent.com%2Fa%2FACg8ocJho_JMw-_dArGBTNa24ipeEGZ8dzv9ttGFTO1A7VWlc7p7Gw%3Ds96-c%22%7D; ph_phc_9LvbXawaTFdrUSVPTOjwbVv7bZWE1iOQDhF8U7dPa0E_posthog=%7B%22distinct_id%22%3A%22cm7ygtgt9005rczp0k9klmjo2%22%2C%22%24sesid%22%3A%5B1741339012823%2C%2201956fba-2879-714d-8252-1bf3b9ef56f8%22%2C1741336225913%5D%7D; _uetsid=d5317470fb2611efa28ab1f1448f428f; _uetvid=d531a870fb2611ef8e2ccd156937b0a5; _ga_JS0YXNKF22=GS1.1.1741332968.1.1.1741339491.60.0.0; _ga_NSMTS5QKXE=GS1.1.1741332968.1.1.1741339491.0.0.0; __Secure-next-auth.session-token=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..Gf5_CCHz-w2ZTDdR.o4GTsCOt8qS0toEWuCJs8eblG9WYO7OhILgj7_qjaWqiVGObYVh-SKZL_9xvn8m7-uA0MbyhLFs9HAAVx8nrsaQToWDrQRATkg7ZZqoSnV23Fddb-KAWUxVBXP_OikXN3LWvpeLRT22Llxu3a6BBj5q_SpIqHzci64GGk5R3aUiSlzirJrYwtDa8PG2QaDHlH2_h9nysUYp4sRDPqXy-jzDJ2iPqEsjeXObOjkKRO_jP-DXlOVn5TtSeh82IqkXb1jvBKm82D7XQDEARPNiM_YTWkNj2-lZExMg86Vj22xcu4t0UJqdIE69YbVddwNT6DMvxTo1MFKi8P3Kz1DGOHbh5AVEH_hQw-sdeCBWHzC2YMA4pHYjpEdVhxJ2EvL-GC7MKenb89OhnjTMPjX3eaOFA8GbhKexFnIphSqKM-wlMrj9e_iUcoL_wQt9ymyx6skqcAs5X1Uc7uJHXVvBYxzCuoSTwGlKFFvH790Q3HN1nHGRTRvoyBHtitm5qD0R2xm38RS77A9o.-3frk1mcOin2AMMHzdsHlA",
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}

	logger.Debug(c.Request.Context(), fmt.Sprintf("cookie: %v", cookie))

	sseChan, err := client.DoSSE(chatEndpoint, options, "POST")
	if err != nil {
		logger.Errorf(c, "Failed to make stream request: %v", err)
		return nil, fmt.Errorf("Failed to make stream request: %v", err)
	}
	return sseChan, nil
}

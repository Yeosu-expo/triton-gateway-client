package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ahr-i/triton-client-gateway/models"
	"github.com/ahr-i/triton-client-gateway/setting"
	"github.com/ahr-i/triton-client-gateway/src/errController"
	"github.com/gorilla/mux"
)

/* Response Struct */
type TritonResponse struct {
	Outputs []struct {
		Data []float32 `json:"data"`
	} `json:"outputs"`
}

type ResponseData struct {
	Image string `json:"image"`
}

/* Request Struct */
type RequestData struct {
	Prompt string `json:"prompt"`
}

/* Inference Handler: Triton Server에 Inference Request 및 Image 전달 */
func (h *Handler) inferHandler(w http.ResponseWriter, r *http.Request) {
	_, fp, _, _ := runtime.Caller(1)

	// Request Decode
	var request RequestData
	err_ := json.NewDecoder(r.Body).Decode(&request)
	errController.ErrorCheck(err_, "REQUEST JSON DECODE ERROR", fp)
	defer r.Body.Close()

	vars := mux.Vars(r)
	model := vars["name"]

	// log.Println(request.Prompt)
	if request.Prompt == "" || request.Prompt == " " {
		rend.JSON(w, http.StatusBadRequest, nil)
	}

	// Model, Version Check And Setting
	modelMap := models.GetModelList()
	version, err := modelMap[model]
	if !err {
		rend.JSON(w, http.StatusNotFound, nil)
	}

	// Triton Inference Request
	rand.Seed(time.Now().UnixNano())

	seed := rand.Intn(10001)
	url := "http://" + setting.TritonUrl + "/v2/models/" + model + "/versions/" + version + "/infer"
	requestData := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"name":     "PROMPT",
				"datatype": "BYTES",
				"shape":    []int{1},
				"data":     []string{request.Prompt},
			},
			{
				"name":     "SAMPLES",
				"datatype": "INT32",
				"shape":    []int{1},
				"data":     []int{1},
			},
			{
				"name":     "STEPS",
				"datatype": "INT32",
				"shape":    []int{1},
				"data":     []int{45},
			},
			{
				"name":     "GUIDANCE_SCALE",
				"datatype": "FP32",
				"shape":    []int{1},
				"data":     []float32{7.5},
			},
			{
				"name":     "SEED",
				"datatype": "INT64",
				"shape":    []int{1},
				"data":     []int{seed},
			},
		},
		"outputs": []map[string]string{
			{
				"name": "IMAGES",
			},
		},
	}

	requestJSON, err_ := json.Marshal(requestData)
	errController.ErrorCheck(err_, "JSON MARSHAL ERROR", fp)

	req, err_ := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	errController.ErrorCheck(err_, "HTTP REQUEST ERROR", fp)
	req.Header.Set("Content-Type", "application/json")

	// Triton Server Response
	client := &http.Client{}
	resp, err_ := client.Do(req)
	errController.ErrorCheck(err_, "HTTP RESPONSE ERROR", fp)
	defer resp.Body.Close()

	// Response Decode
	body, err_ := ioutil.ReadAll(resp.Body)
	errController.ErrorCheck(err_, "HTTP BODY READ ERROR", fp)

	var tritonResponse TritonResponse
	if err := json.Unmarshal(body, &tritonResponse); err != nil {
		log.Fatalf("RESPONSE JSON PARSE ERROR: %v", err)
	}

	// Uint8 Array To Image
	if len(tritonResponse.Outputs) > 0 && len(tritonResponse.Outputs[0].Data) > 0 {
		imgData := tritonResponse.Outputs[0].Data

		// Image의 크기 가정
		width, height := 512, 512
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// ImgData에서 픽셀 값 추출 및 Image 생성
		for i := 0; i < len(imgData); i += 3 {
			x := (i / 3) % width
			y := (i / 3) / width
			r := uint8(imgData[i] * 255)
			g := uint8(imgData[i+1] * 255)
			b := uint8(imgData[i+2] * 255)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}

		// Image Local 저장
		currentTime := time.Now().Format("20060102-150405.999")
		fileName := "result-" + currentTime + ".png"
		file, err := os.Create("./result/" + fileName)
		if err != nil {
			log.Fatalf("이미지 파일 생성 실패: %v", err)
		}
		defer file.Close()

		if err := png.Encode(file, img); err != nil {
			log.Fatalf("이미지 저장 실패: %v", err)
		}

		// Image Base64 Encoding
		var buffer bytes.Buffer

		if err := png.Encode(&buffer, img); err != nil {
			log.Println("BASE ENCODE FAIL")

			os.Exit(1)
		}
		imgBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

		// Response에 Image 추가
		w.Header().Set("Content-Type", "application/json")
		responseData := ResponseData{Image: imgBase64}

		// Response
		rend.JSON(w, http.StatusOK, responseData)
		return
	}

	// Inference Fail
	rend.JSON(w, http.StatusBadRequest, nil)
}
